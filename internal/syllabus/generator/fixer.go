package generator

import "sort"

// fixer resolves scheduling gaps — holes between consecutive lessons for a group.
type fixer struct {
	cfg   *settings
	coord *coordinator
}

// newFixer creates a new fixer.
func newFixer(cfg *settings, coord *coordinator) *fixer {
	return &fixer{
		cfg:   cfg,
		coord: coord,
	}
}

// fixWeekSchedule attempts to compress the schedule for each date,
// minimizing gaps between consecutive lessons for each group.
func (f *fixer) fixWeekSchedule(dates []date, groupIDs []uint64) {
	for _, d := range dates {
		for _, groupID := range groupIDs {
			f.fixDayForGroup(d, groupID)
		}
	}
}

// fixDayForGroup compresses the schedule for a specific group on a specific day.
func (f *fixer) fixDayForGroup(d date, groupID uint64) {
	scheduledNumbers := f.coord.schedule.getScheduledLessonNumbers(d, groupID)
	if len(scheduledNumbers) <= 1 {
		return
	}

	sort.Ints(scheduledNumbers)

	// Check for gaps
	for i := 1; i < len(scheduledNumbers); i++ {
		gap := scheduledNumbers[i] - scheduledNumbers[i-1]
		if gap <= 1 {
			continue
		}

		// There's a gap — try to move lessons closer
		f.tryCompressGap(d, groupID, scheduledNumbers, i)
	}
}

// tryCompressGap attempts to move a lesson to fill a gap.
func (f *fixer) tryCompressGap(d date, groupID uint64, scheduledNumbers []int, gapIndex int) {
	// Try to move the lesson at scheduledNumbers[gapIndex] closer to scheduledNumbers[gapIndex-1]
	targetSlot := scheduledNumbers[gapIndex-1] + 1
	sourceSlot := scheduledNumbers[gapIndex]

	if targetSlot == sourceSlot {
		return
	}

	// Find the lesson to move
	lessons := f.coord.schedule.getLessons(d, sourceSlot)
	var lessonToMove *lesson
	for _, l := range lessons {
		if l.groupID == groupID {
			lessonToMove = l
			break
		}
	}

	if lessonToMove == nil {
		return
	}

	// Check if target slot is free for all participants
	allParticipants := append([]uint64{groupID}, lessonToMove.teacherIDs()...)
	if !f.coord.availability.areFree(d, allParticipants, targetSlot) {
		return
	}

	// Check for conflicts at target
	if f.coord.schedule.hasConflict(d, targetSlot, groupID, lessonToMove.teacherIDs()) {
		return
	}

	// Move the lesson
	f.coord.schedule.remove(d, sourceSlot, lessonToMove)
	f.coord.schedule.add(d, targetSlot, lessonToMove)

	// Update availability
	for _, pid := range allParticipants {
		f.coord.availability.markFree(d, pid, sourceSlot)
		f.coord.availability.markBusy(d, pid, targetSlot)
	}
}
