# Schedule Generation Algorithm

## Overview

The generator produces a weekly schedule template for a university semester. Given groups, teachers, subjects, and workload assignments, it finds a valid timetable that respects all hard constraints. The result is a `ScheduleData` value containing two week patterns — **numerator** (odd weeks) and **denominator** (even weeks) — each mapping weekdays to lesson slots.

The algorithm runs **N goroutines in parallel** (default 4), each attempting an independent generation with a different random seed. The first goroutine to succeed returns its result; all others are cancelled. A configurable timeout (default 120 s) caps the total wall time.

---

## Architecture

The package is split into focused single-responsibility components. Every generation attempt creates its own isolated set of these components, so parallel goroutines do not share state.

```
generator.go          — public API: parallel orchestration, timeout, first-result wins
generation.go         — single attempt: calls each stage in order
settings.go           — aggregates ScheduleTemplateSetting + ScheduleRestriction → settings
workload.go           — converts domain assignments to lessons; distributes hours across weeks
validator.go          — pre-generation feasibility check; post-placement hard constraint check
coordinator.go        — wires together all sub-components for one attempt
availability.go       — per-participant, per-slot free/busy tracking
limiter.go            — enforces daily and weekly lesson count ceilings
pattern.go            — enforces max identical lessons per day; min-gap-between constraints
finder.go             — finds candidate lesson slots (adjacent to existing, cached)
random_scheduler.go   — first-pass: randomly places as many lessons as possible
fixer.go              — gap compression: shifts lessons closer together to remove holes
recursive_scheduler.go— second-pass: backtracking placement for remaining lessons
denominator.go        — mirrors numerator pattern to denominator; schedules denom-only lessons
room_selector.go      — assigns rooms to sub-lessons after all lessons are placed
schedule.go           — internal mutable timetable: date → slot → []*lesson
date.go               — lightweight date type; semester date collection; week utilities
lesson.go             — internal mutable lesson model; sub-lesson per teacher assignment
```

---

## Generation Pipeline

Each goroutine executes these stages in order. A failure at any stage causes the goroutine to return an error; the orchestrator retries by starting a new goroutine with a different seed (up to the configured number of goroutines, then waits for timeout).

```
1.  Build coordinator          newCoordinator(settings, roomIDs)
2.  Set pattern limits         patternController.setDayLimit per workload
3.  Collect semester dates     collectSemesterDates(start, end, educationWeek)
                               → numeratorDates (odd weeks), denominatorDates (even weeks)
4.  Distribute workloads       workloadDistributor.distributeWorkloads(workloads, numWeeks)
                               → numLessons, denomLessons (per-week count for each workload)
5.  Extract first-week dates   firstWeekDates(numeratorDates)
                               (template covers one repeating week, not the full semester)
6.  Pre-generation validation  validator.validate(numLessons, firstWeek, groupIDs)
                               validator.validate(denomLessons, denomFirstWeek, groupIDs)
7.  Random scheduling          randomScheduler.scheduleRandom(numLessons, firstWeek)
                               → unscheduled (lessons that didn't fit randomly)
8.  Gap compression            fixer.fixWeekSchedule(firstWeek, groupIDs)
9.  Recursive scheduling       recursiveScheduler.schedule(unscheduled, firstWeek)
10. Post-placement validation  validatePlacement(schedule, firstWeek, groupIDs, settings)
11. Denominator reproduction   denominatorReproducer.reproduce(numDates, denomDates, denomLessons)
    11a. Mirror numerator pattern to corresponding denominator days
    11b. Gap compression for denominator
    11c. Post-placement validation for denominator
12. Room assignment            roomSelector.assignRooms(schedule)
13. Convert to ScheduleData    group first-week dates by weekday → WeekSchedule
```

---

## Constraint System

### Pre-Generation (Fast-Fail)

Checked in **stage 6** before any placement is attempted. If violated, the attempt fails immediately with `ErrValidationFailed`.

| Constraint | Check |
|---|---|
| Available slots | `required lessons ≤ days × maxGroupLessonsPerDay` for every group |
| Weekly ceiling | `required lessons ≤ maxGroupLessonsPerWeek` for every group |
| Minimum daily feasibility | If a group has any lessons and `required < minGroupLessonsPerDay`, it is impossible to schedule them without violating the daily minimum (e.g., 1 lesson/week with min=2 means any day that group has a lesson would have only 1, which violates the constraint) |

### During Placement (Hard — reject slot, try next)

Applied by `limiter`, `availability`, `patternController`, and `schedule.hasConflict` for every candidate slot.

| Constraint | Component |
|---|---|
| No double-booking a group | `schedule.hasConflict` |
| No double-booking a teacher | `schedule.hasConflict` |
| Daily group lesson ceiling | `limiter.canAddLessonForGroup` |
| Daily teacher lesson ceiling | `limiter.canAddLessonForTeacher` |
| Weekly group lesson ceiling | `limiter.canAddWeekLessonForGroup` |
| Weekly teacher lesson ceiling | `limiter.canAddWeekLessonForTeacher` |
| Max identical pattern per day | `patternController.shouldSkip` (pattern = group+teacher+subject) |
| Participant availability | `availability.areFree` |

### Post-Placement (Hard — failed attempt triggers retry with new seed)

Checked in **stage 10** and **11c** after the fixer has run. Failure returns `ErrUnableToSchedule`.

| Constraint | Setting |
|---|---|
| No gaps in group schedule | `noGapsRequired = true` (default): a group's lessons on a given day must occupy consecutive slots with no empty slots between them |
| Minimum lessons per day | `minGroupLessonsPerDay > 0` (default 2): any day a group has at least one lesson must have at least this many |

---

## Configuration

### `ScheduleTemplateSetting`

Primary source of capacity parameters.

| Field | Role |
|---|---|
| `HoursPerLesson` | Duration of one lesson in hours (e.g., 1.5 or 2.0) |
| `MaxStudyHoursPerDay` | Maximum study hours per day → `floor(MaxStudyHoursPerDay / HoursPerLesson)` = max group lessons/day |
| `MaxTeacherHoursPerWeek` | Teacher weekly ceiling in hours → `floor(.../ HoursPerLesson)` = max teacher lessons/week |
| `MaxGroupLessonHoursPerWeek` | Optional group weekly ceiling in hours |
| `MaxIdenticalLessonsPerDay` | Max times the same pattern (group+teacher+subject) may appear on one day |

### `ScheduleRestriction`

Applies additional hard caps on top of `ScheduleTemplateSetting`. The latest record from the DB is used; if none exists, defaults apply.

| Field | Default | Role |
|---|---|---|
| `MinGroupLessonsPerDay` | 2 | Minimum lessons per day when a group is scheduled |
| `MaxGroupLessonsPerDay` | 4 | Hard ceiling; takes effect only when `< hours-based max` |
| `MaxTeacherLessonsPerDay` | 5 | Hard ceiling for teachers (independent from groups) |
| `NoGapsInGroupSchedule` | true | Require consecutive lesson slots for each group each day |

**Priority rule**: `effective max = min(hours-based max, restriction max)`. A restriction can only make limits more strict, never more lenient than the hours-based calculation.

### `GenerationConfig`

Runtime parameters for the parallel orchestrator.

| Field | Default | Role |
|---|---|---|
| `Timeout` | 120 s | Wall-clock limit; `ErrTimeout` if no goroutine succeeds before deadline |
| `NumberOfGoroutines` | 4 | Number of parallel attempts with independent random seeds |

---

## Workload Distribution

Each `WorkloadDistribution` + its `WorkloadAssignment` entries represent one teacher's assignment for a group/subject pair. The distributor converts total classroom hours into a per-week lesson count.

```
totalLessons = ceil(totalHours / hoursPerLesson)
numeratorLessons  = ceil(totalLessons / 2)   + remainder if totalLessons is odd
denominatorLessons = floor(totalLessons / 2)

perWeek = ceil(periodLessons / weeksPerPeriod)
```

Each resulting `perWeek` value becomes the count of lesson instances placed in the first week of that period (the template repeats weekly).

---

## Lesson Types

### United (`formatUnited`)

A standard lesson for a whole group taught by one teacher. One `lesson` with one `internalSubLesson`.

### Split (`formatSplit`)

The group is divided into sub-groups. Each sub-group has its own `internalSubLesson` (independent teacher, independent room). The `lesson` groups all sub-lessons so their slots can be coordinated. A group is split when `WorkloadDistribution.IsSplitting = true` and a `subGroupNumber` is present.

### Flow (not yet fully implemented)

Multiple groups share one room with one teacher. Represented as a single `lesson` with multiple `internalSubLesson` entries across different groups.

---

## Numerator / Denominator Week Splitting

Ukrainian university scheduling traditionally alternates two weekly patterns:

- **Numerator** — odd weeks (week 1, 3, 5, …)
- **Denominator** — even weeks (week 2, 4, 6, …)

`collectSemesterDates` walks from semester start to end, classifying each education day into one of the two pools based on the week number (0-indexed from the Monday of the starting week).

Each pool is then handled independently:
1. Numerator is scheduled in **stages 7–10**.
2. Denominator is handled in **stage 11** by the `denominatorReproducer`, which first mirrors matching subjects from the numerator pattern (same slot, same weekday), then uses random + recursive scheduling for the remaining denominator-only lessons.

The final `ScheduleData` contains two `WeekSchedule` maps — `Numerator` and `Denominator` — keyed by `domain.Weekday`.

---

## Slot Finding Strategy

The `finder` prefers **adjacent slots** to minimise gaps before the fixer runs:

1. If the group has no lessons yet on a day → return only the **first lesson number** (slot 1 by default).
2. If the group already has lessons → return all slots **immediately before or after** an occupied slot (deduplicated, sorted).
3. If no adjacent free slot exists → fall back to all free slots on that day.

Results are **cached** keyed on the sorted busy-slot list, since the same combination appears repeatedly during a generation attempt.

---

## Gap Compression (Fixer)

After random scheduling, a group's lessons may sit at non-consecutive slots (e.g., slots 1 and 3). The fixer iterates over each group's day schedule and attempts to shift a lesson toward the previous occupied slot:

```
for each gap between scheduledNumbers[i-1] and scheduledNumbers[i]:
    targetSlot = scheduledNumbers[i-1] + 1
    if all participants are free at targetSlot and no conflict exists:
        move lesson from scheduledNumbers[i] → targetSlot
        update availability
```

This is a single-pass greedy compression — it does not guarantee a globally optimal arrangement, but it removes the most obvious gaps efficiently before the post-placement validator runs.

---

## Recursive Backtracking (Recursive Scheduler)

Lessons that the random scheduler could not place are handed to the recursive scheduler. It performs DFS with limited depth (`maxRecursionDepth = 100`) and limited retries per lesson (`maxRetryPerLesson = 20`).

```
tryScheduleLesson(lesson):
  1. Scan all dates for a free slot with no conflicts → place directly
  2. Scan all dates for a slot occupied by a conflicting lesson:
       a. Remove the conflicting lesson from that slot
       b. Place the current lesson there
       c. Recursively call tryScheduleLesson(conflicting lesson)
       d. If recursion fails → backtrack: remove current, restore conflicting
```

Backtracking is bounded by `maxRecursionDepth` to prevent stack overflow on pathological inputs.

---

## Room Assignment

Room assignment happens **after** all lessons are placed (stage 12), as a separate pass that does not affect the constraint system.

The `roomSelector` selects rooms greedily:
1. For each lesson slot (date + lesson number), collect rooms already assigned in that slot.
2. For each `internalSubLesson` without a room, pick the **least-used room** that is not already occupied in this slot.
3. If all rooms are already in use in the slot (more sub-lessons than rooms), fall back to the globally least-used room (double-booking rooms is allowed as a last resort).

Usage counts persist across the full schedule so that room load is balanced over the week.

---

## Data Structures

### `schedule`

```
dateToLessons: map[date] → map[lessonNumber] → []*lesson
```

The central mutable timetable. Supports add/remove by (date, lessonNumber, lesson) and queries for group or teacher occupancy. Multiple lessons can share the same (date, lessonNumber) slot — used for split/flow lessons where different sub-groups occupy the same time slot independently.

### `availability`

```
data: map[date] → map[participantID] → map[lessonNumber] → status
```

Tracks whether each participant (group or teacher, identified by their domain ID) is free, busy, unavailable, or on a methodical day at each (date, lessonNumber). Updated in lock-step with `schedule` so that the two never diverge.

Statuses:
- `statusFree` — default for any unseen key
- `statusHaveLesson` — set when a lesson is placed
- `statusUnavailable` — reserved for external constraints
- `statusMethodicalDay` — blocks an entire day for a participant

### `patternController`

```
occurrences: map[date] → map[pattern] → count
dayLimits:   map[pattern] → maxPerDay
dayBetween:  map[pattern] → minDayGap
```

A **pattern** is the string key `"groupID:teacherID:subjectID"`. The controller enforces two independent constraints:
- `dayLimits` — how many times the same pattern may appear on one day (driven by `MaxIdenticalLessonsPerDay`)
- `dayBetween` — minimum number of calendar days between any two placements of the same pattern (not set by current configuration but the mechanism exists)

---

## Error Reference

| Error | Meaning |
|---|---|
| `ErrTimeout` | No goroutine succeeded before the configured timeout |
| `ErrValidationFailed` | Pre-generation check failed — the input is provably infeasible |
| `ErrUnableToSchedule` | Post-placement validation failed or recursive scheduler exhausted all options |
| `ErrWorkloadDistribution` | Workload hours cannot be distributed (e.g., zero weeks) |
| `ErrNoSettings` | No `ScheduleTemplateSetting` found in the database |
| `ErrNoBellSchedule` | No bell schedules configured |
| `ErrNoGroupSemesters` | No group semesters found for the requested semester |
| `ErrNoWorkloads` | No workload assignments found for the selected groups |

---

## Example: Minimal Feasible Input

```
ScheduleTemplateSetting:
  HoursPerLesson:            2.0
  MaxStudyHoursPerDay:       8       → max 4 lessons/day per group
  MaxTeacherHoursPerWeek:    40      → max 20 lessons/week per teacher
  MaxIdenticalLessonsPerDay: 2

ScheduleRestriction:
  MinGroupLessonsPerDay:   2
  MaxGroupLessonsPerDay:   4
  MaxTeacherLessonsPerDay: 5
  NoGapsInGroupSchedule:   true

BellSchedules: 4 entries (lesson numbers 1–4)

Semester: 2 weeks, Mon–Fri (10 education days)

WorkloadDistribution:
  Group 1, Study Plan A (Discipline: Mathematics)
WorkloadAssignment:
  Teacher 42, AssignedHours: 8
  → totalLessons = ceil(8/2) = 4
  → numeratorLessons = 2, denominatorLessons = 2
  → each placed once per week in each period

Result:
  Numerator:  { monday: {1: [...], 2: [...]}, ... }
  Denominator: { monday: {1: [...], 2: [...]}, ... }
```
