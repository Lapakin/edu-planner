package generator

import (
	"math/rand"
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

		// Find available lesson numbers
		busySlots := rs.coord.schedule.getScheduledLessonNumbers(d, l.groupID)
		availableSlots := rs.coord.finder.findUnscheduledNumbers(busySlots)

		if len(availableSlots) == 0 {
			continue
		}

		// Try each available slot
		slotIdx := rs.rng.Intn(len(availableSlots))
		lessonNumber := availableSlots[slotIdx]

		// Check for conflicts
		if rs.coord.schedule.hasConflict(d, lessonNumber, l.groupID, l.teacherIDs()) {
			continue
		}

		// Check availability
		allParticipants := append([]uint64{l.groupID}, l.teacherIDs()...)
		if !rs.coord.availability.areFree(d, allParticipants, lessonNumber) {
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
