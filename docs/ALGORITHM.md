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
flow.go               — merges flow-eligible lectures into shared flow lessons; flow room feasibility, fallback, diagnostics
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
                               → flow-eligible lectures merged into shared flow lessons
                                 (mergeFlowLessons, per week period) — see "Flow Lessons"
5.  Extract first-week dates   firstWeekDates(numeratorDates)
                               (template covers one repeating week, not the full semester)
5a. Apply forbidden slots      generation.applyTeacherForbiddenSlots(availability, allDates)
                               → pre-marks forbidden teacher slots as unavailable before validation
6.  Pre-generation validation  validator.validate(numLessons, firstWeek, groupIDs)
                               validator.validate(denomLessons, denomFirstWeek, groupIDs)
6b. Flow pre-pass              partition flow lessons; drop flows no room can hold
                               (flowRoomFeasible) → fallback; place remaining flows first;
                               flows that find no shared slot fall back to per-group lessons
7.  Random scheduling          randomScheduler.scheduleRandom(regularLessons, firstWeek)
                               → forbidden slots skipped; preferred slots weighted higher
                               → unscheduled (lessons that didn't fit randomly)
8.  Gap compression            fixer.fixWeekSchedule(firstWeek, groupIDs)
9.  Recursive scheduling       recursiveScheduler.schedule(unscheduled, firstWeek)
10. Post-placement validation  validatePlacement(schedule, firstWeek, groupIDs, settings)
11. Denominator reproduction   denominatorReproducer.reproduce(numDates, denomDates, denomLessons, diag)
    11a. Drop denominator flows no room can hold (flowRoomFeasible) → fallback
    11b. Mirror numerator pattern to corresponding denominator days (flows mirror as a unit)
    11c. Flow pre-pass for remaining flows, then gap compression for denominator
    11d. Post-placement validation for denominator
12. Room assignment            roomSelector.assignRooms(schedule)
                               (a flow lesson gets ONE shared room ≥ Σ group sizes)
13. Convert to ScheduleData    group first-week dates by weekday → WeekSchedule
                               flow diagnostics written to ScheduleData.Metadata (JSON)
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
| No double-booking a group | `schedule.hasConflict` (checks **every** group of a lesson — a flow spans several) |
| No double-booking a teacher | `schedule.hasConflict` |
| Flow shared slot | `schedule.getScheduledLessonNumbersForGroups` + `availability.areFree` — a flow lesson is placed only at a slot free for **all** its groups (intersection of free intervals) |
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
| `LessonsPerClass` | Number of academic hours in one lesson slot (e.g., 1, 2, or 3) |
| `MaxStudyHoursPerDay` | Maximum study hours per day → `MaxStudyHoursPerDay / LessonsPerClass` = max group lessons/day |
| `MaxTeacherHoursPerWeek` | Teacher weekly ceiling in hours → `MaxTeacherHoursPerWeek / LessonsPerClass` = max teacher lessons/week |
| `MaxGroupLessonHoursPerWeek` | Optional group weekly ceiling in hours → `/ LessonsPerClass` = max group lessons/week |
| `MaxIdenticalLessonsPerDay` | Max times the same pattern (group+teacher+subject) may appear on one day |

### `ScheduleRestriction`

Applies additional hard caps on top of `ScheduleTemplateSetting`. The latest record from the DB is used; if none exists, defaults apply.

| Field | Default | Role |
|---|---|---|
| `MinGroupLessonsPerDay` | 2 | Minimum lessons per day when a group is scheduled |
| `MaxGroupLessonsPerDay` | 4 | Hard ceiling; takes effect only when `< LessonsPerClass-based max` |
| `MaxTeacherLessonsPerDay` | 5 | Hard ceiling for teachers (independent from groups) |
| `MaxConsecutiveTeacherLessons` | 4 | Maximum consecutive lesson slots for a teacher without a break |
| `NoGapsInGroupSchedule` | true | Require consecutive lesson slots for each group each day |
| `TimePriority` | none | Preference for scheduling lessons in the morning or afternoon |
| `AllowFlowLessons` | true | Global switch for flow lessons. When `false`, no merging happens and every lecture is scheduled per group, exactly as before the feature. Per-discipline opt-in is `Discipline.IsFlow`; both must be true to form a flow |

**Priority rule**: `effective max = min(LessonsPerClass-based max, restriction max)`. A restriction can only make limits more strict, never more lenient than the LessonsPerClass-based calculation.

### `TeacherSlotPreference`

Per-teacher time slot preferences. Stored in the database and loaded at generation time for the relevant academic year.

| Field | Role |
|---|---|
| `TeacherID` | The teacher these preferences apply to |
| `Weekday` | Day of the week (monday, tuesday, …, saturday) |
| `LessonNumber` | The lesson slot number on that day |
| `SlotType` | `preferred` or `forbidden` |

**Forbidden slots** are pre-marked as unavailable in the availability tracker before any placement begins. The random and recursive schedulers skip all slots that are forbidden for any teacher involved in a lesson.

**Preferred slots** act as a soft bias in the random scheduler: when selecting among available slots for a teacher, slots that are preferred for all involved teachers are weighted with a higher probability (`prioritySlotWeight`). If no preferred slot is available, the scheduler falls back to any free slot.

### `CycleCommitteeLabRoom`

Specifies which rooms are designated lab rooms for a given cycle committee. Used to direct lab-type lessons to the appropriate rooms.

| Field | Role |
|---|---|
| `CycleCommitteeID` | The cycle committee whose disciplines should use these rooms |
| `RoomID` | A room designated for lab work of that committee |

When the room selector assigns rooms after all lessons are placed, it checks each lesson's `cycleCommitteeID`. If the lesson belongs to a committee that has designated lab rooms, the selector preferentially picks from those rooms. If none of the designated rooms are free in that slot, it falls back to the general room pool.

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
totalLessons = ceil(totalHours / lessonsPerClass)
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

### Flow (`isFlow`)

Several groups attend one lecture together (потокова пара) — one teacher, one room, one slot. Represented as a single `lesson` with one `internalSubLesson` per group, all sharing the same teacher. The output `ScheduleLesson` has `Type = flow`, `IsFlow = true`, and a `FlowID` identifying the stream. See the dedicated **Flow Lessons** section below.

---

## Flow Lessons (Потокові пари)

A **flow lesson** is one lecture delivered to several groups simultaneously — one teacher, one room, one time slot — the typical case being a stream lecture shared by parallel groups. It is gated by two switches that **both** must be true: the per-discipline `Discipline.IsFlow` and the global `ScheduleRestriction.AllowFlowLessons`. When either is off, nothing is merged and generation behaves exactly as it did before the feature.

### Eligibility and merging (`flow.go`)

During workload distribution (stage 4), `mergeFlowLessons` runs once per week period:

1. A lesson is a **candidate** when it is `flowEligible` (a flow discipline's non-lab lecturer lesson), `formatUnited`, with at least one sub-lesson. `flowEligible` is set in `workload.go` from `Discipline.IsFlow && !isLab`. Eligibility alone is **not** a flow — only the merged result is.
2. Candidates are grouped by `(subjectID, teacherID)`. Within each stream, while **two or more** distinct groups still have an unconsumed lesson, one lesson is pulled from each and merged into a single flow `lesson` carrying one `internalSubLesson` per group. `flowOrigin` retains the constituents for fallback; `flowID = "flow-<subject>-<teacher>"`.
3. Surplus lessons (a group with more lectures than its partners) and lonely candidates (no second group) stay as ordinary single-group lessons.

Because merging happens per period independently, the numerator and denominator each form their own flow lessons.

### Multi-group placement

The placement engine is **group-set aware**. A flow lesson exposes all its groups via `lesson.groupIDs()` / `hasGroup()` / `hasAnyGroup()`, and every component that used a single `groupID` now considers the whole set:

- `schedule.hasConflict` / `findConflictingLesson` — conflict if **any** group (or the teacher) is already busy.
- `schedule.getScheduledLessonNumbersForGroups` — the union of busy slots across all groups, so a flow only lands where **all** groups are free (the intersection of their free intervals).
- `limiter.canPlaceLesson` — daily/weekly group ceilings checked for every group.
- `validator.validate` / `validatePlacement` — a flow counts once toward each participating group.
- The fixer treats flow lessons as **immovable anchors** (relocating one for a single group would disturb the others); non-flow lessons still compress around them, and any residual gap is resolved by a retry with a new seed.

Flow lessons are placed **first** (stage 6b) because they are the most constrained.

### Rooms

A flow needs one room large enough for everyone. `flowRoomFeasible` checks that some room's capacity covers the **sum** of all participating groups' sizes (a capacity of 0 = unlimited; unknown sizes impose no constraint). At room assignment, `assignFlowRoom` gives every sub-lesson of the flow the **same** room.

### Fallback and diagnostics

A flow falls back to separate per-group lessons (`lesson.unmerge()`) when:

- no room can hold the combined size (`flowRoomFeasible` is false), or
- no shared slot is found during the flow pre-pass.

The unmerged constituents re-enter the normal random + recursive pipeline, so **a flow never causes generation to fail** — at worst it degrades to ordinary lessons. Each attempt records flow outcomes in a `flowDiagnostics`, serialized into `ScheduleData.Metadata` as JSON:

```json
{"flow_lessons": 2, "flow_fallbacks": 1, "fallback_flow_ids": ["flow-20-11"]}
```

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
2. **Flow lessons** are handled as a unit (`assignFlowRoom`): one room sized to the **combined** group size (`Σ group sizes`) is chosen and assigned to every sub-lesson, so all groups of the stream sit together.
3. For each remaining `internalSubLesson` without a room:
   a. If the lesson's `cycleCommitteeID` has designated lab rooms (`CycleCommitteeLabRoom` entries), prefer those rooms.
   b. Otherwise, pick the **least-used room** from the general pool that is not already occupied in this slot — restricted to rooms whose capacity fits the group (split sub-lesson ≈ half a group; united = whole group).
4. If all preferred or available rooms are occupied in the slot, fall back to the globally least-used fitting room (double-booking is allowed as a last resort).

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
  LessonsPerClass:           2       (2 academic hours per lesson slot)
  MaxStudyHoursPerDay:       8       → max 4 lessons/day per group (8/2)
  MaxTeacherHoursPerWeek:    40      → max 20 lessons/week per teacher (40/2)
  MaxIdenticalLessonsPerDay: 2

ScheduleRestriction:
  MinGroupLessonsPerDay:   2
  MaxGroupLessonsPerDay:   4
  MaxTeacherLessonsPerDay: 5
  NoGapsInGroupSchedule:   true

BellSchedules: 4 entries (lesson numbers 1–4)

Semester: 2 weeks, Mon–Fri (10 education days)

WorkloadDistribution:
  Group 1, Study Plan A (Discipline: Mathematics, CycleCommittee: Math Committee)
WorkloadAssignment:
  Teacher 42, AssignedHours: 8
  → totalLessons = ceil(8/2) = 4
  → numeratorLessons = 2, denominatorLessons = 2
  → each placed once per week in each period

TeacherSlotPreference (optional):
  Teacher 42, monday, lesson 1: forbidden  → slot 1 on Monday skipped for Teacher 42
  Teacher 42, wednesday, lesson 2: preferred → slot 2 on Wednesday weighted higher

CycleCommitteeLabRoom (optional):
  Math Committee → Room 301  → lab lessons for Math Committee prefer Room 301

Result:
  Numerator:  { monday: {1: [...], 2: [...]}, ... }
  Denominator: { monday: {1: [...], 2: [...]}, ... }
```
