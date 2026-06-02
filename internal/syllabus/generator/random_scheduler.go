package generator

import (
	"math/rand"
	"sort"

	"github.com/Lapakin/edu-planner/internal/domain"
)

const (
	maxRandomAttempts  = 15
	prioritySlotWeight = 3
)

// randomScheduler performs the initial random placement of lessons.
type randomScheduler struct {
	cfg   *settings
	coord *coordinator
	rng   *rand.Rand
}

// newRandomScheduler creates a new random scheduler.
func newRandomScheduler(cfg *settings, coord *coordinator, rng *rand.Rand) *randomScheduler {
	return &randomScheduler{
		cfg:   cfg,
		coord: coord,
		rng:   rng,
	}
}

// scheduleRandom attempts to randomly place all unscheduled lessons.
// Returns the list of lessons that could not be scheduled.
func (rs *randomScheduler) scheduleRandom(lessons []*lesson, dates []date, weekDates []date) []*lesson {
	// Shuffle lessons for randomness
	rs.rng.Shuffle(len(lessons), func(i, j int) {
		lessons[i], lessons[j] = lessons[j], lessons[i]
	})

	var unscheduled []*lesson

	for _, l := range lessons {
		if l.isScheduled {
			continue
		}

		placed := rs.tryPlace(l, dates, weekDates)
		if !placed {
			unscheduled = append(unscheduled, l)
		}
	}

	return unscheduled
}

// tryPlace attempts to place a single lesson at a random valid slot.
func (rs *randomScheduler) tryPlace(l *lesson, dates []date, weekDates []date) bool {
	for range maxRandomAttempts {
		// Pick a random date
		dateIdx := rs.rng.Intn(len(dates))
		d := dates[dateIdx]

		// Check pattern constraints
		if rs.coord.patternCtrl.shouldSkipSubLessons(d, l.subLessons, l.subjectID, l.groupID) {
			continue
		}

		// Check limiter
		if !rs.coord.limiter.canPlaceLesson(rs.coord.schedule, d, weekDates, l) {
			continue
		}

		// Find available lesson numbers (free for every group of the lesson)
		busySlots := rs.coord.schedule.getScheduledLessonNumbersForGroups(d, l.groupIDs())
		availableSlots := rs.coord.finder.findUnscheduledNumbers(busySlots)

		if len(availableSlots) == 0 {
			continue
		}

		// Filter out forbidden teacher slots
		filtered := rs.filterForbiddenSlots(availableSlots, d, l.teacherIDs())
		if len(filtered) == 0 {
			continue
		}
		availableSlots = filtered

		// Order slots according to time priority, biasing toward preferred teacher slots
		lessonNumber, ok := rs.selectSlotForTeachers(availableSlots, d, l.teacherIDs())
		if !ok {
			continue
		}

		// Check for conflicts
		if rs.coord.schedule.hasConflict(d, lessonNumber, l.groupIDs(), l.teacherIDs()) {
			continue
		}

		// Check availability
		allParticipants := append(l.groupIDs(), l.teacherIDs()...)
		if !rs.coord.availability.areFree(d, allParticipants, lessonNumber) {
			continue
		}

		// Check consecutive teacher lesson constraint
		if rs.violatesConsecutive(d, l.teacherIDs(), lessonNumber) {
			continue
		}

		// Place the lesson
		rs.coord.schedule.add(d, lessonNumber, l)

		// Update availability
		for _, pid := range allParticipants {
			rs.coord.availability.markBusy(d, pid, lessonNumber)
		}

		// Update pattern controller
		rs.coord.patternCtrl.increment(d, l.pattern())

		return true
	}

	return false
}

// selectSlot picks a lesson number from availableSlots according to the configured
// time priority. Returns the selected lesson number and true, or 0 and false if
// no valid slot could be selected.
func (rs *randomScheduler) selectSlot(slots []int) (int, bool) {
	if len(slots) == 0 {
		return 0, false
	}

	switch rs.cfg.timePriority {
	case domain.TimePriorityMorning:
		sorted := make([]int, len(slots))
		copy(sorted, slots)
		sort.Ints(sorted)
		cutoff := max(1, len(sorted)/2)
		if rs.rng.Intn(prioritySlotWeight) > 0 {
			return sorted[rs.rng.Intn(cutoff)], true
		}
		return sorted[rs.rng.Intn(len(sorted))], true

	case domain.TimePriorityAfternoon:
		sorted := make([]int, len(slots))
		copy(sorted, slots)
		sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
		cutoff := max(1, len(sorted)/2)
		if rs.rng.Intn(prioritySlotWeight) > 0 {
			return sorted[rs.rng.Intn(cutoff)], true
		}
		return sorted[rs.rng.Intn(len(sorted))], true

	case domain.TimePriorityNone:
		return slots[rs.rng.Intn(len(slots))], true

	default:
		return slots[rs.rng.Intn(len(slots))], true
	}
}

// filterForbiddenSlots removes slots that are forbidden for any teacher on a given date.
func (rs *randomScheduler) filterForbiddenSlots(slots []int, d date, teacherIDs []uint64) []int {
	wd := d.weekday()
	filtered := make([]int, 0, len(slots))
	for _, ln := range slots {
		forbidden := false
		for _, tid := range teacherIDs {
			if rs.cfg.teacherForbiddenSlots[tid] != nil && rs.cfg.teacherForbiddenSlots[tid][wd][ln] {
				forbidden = true
				break
			}
		}
		if !forbidden {
			filtered = append(filtered, ln)
		}
	}
	return filtered
}

// selectSlotForTeachers picks a slot, biasing toward preferred teacher slots if any exist.
func (rs *randomScheduler) selectSlotForTeachers(slots []int, d date, teacherIDs []uint64) (int, bool) {
	if len(slots) == 0 {
		return 0, false
	}
	wd := d.weekday()
	preferred := make([]int, 0)
	for _, ln := range slots {
		for _, tid := range teacherIDs {
			if rs.cfg.teacherPreferredSlots[tid] != nil && rs.cfg.teacherPreferredSlots[tid][wd][ln] {
				preferred = append(preferred, ln)
				break
			}
		}
	}
	if len(preferred) > 0 && rs.rng.Intn(prioritySlotWeight) > 0 {
		return preferred[rs.rng.Intn(len(preferred))], true
	}
	return rs.selectSlot(slots)
}

// violatesConsecutive returns true if placing a lesson for any of the given
// teachers at lessonNumber on date d would exceed the consecutive lesson limit.
func (rs *randomScheduler) violatesConsecutive(d date, teacherIDs []uint64, lessonNumber int) bool {
	maxConsecutiveTeacherLessons := rs.cfg.maxConsecutiveTeacherLessons
	if maxConsecutiveTeacherLessons <= 0 {
		return false
	}
	for _, tid := range teacherIDs {
		if rs.coord.availability.wouldExceedConsecutive(d, tid, lessonNumber, maxConsecutiveTeacherLessons) {
			return true
		}
	}
	return false
}
