package generator

import (
	"sort"

	"github.com/Lapakin/edu-planner/internal/domain"
)

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
// minimizing gaps between consecutive lessons for each group, and then slides
// each group's daily block toward the preferred half of the day (time priority).
func (f *fixer) fixWeekSchedule(dates []date, groupIDs []uint64) {
	for _, d := range dates {
		for _, groupID := range groupIDs {
			f.fixDayForGroup(d, groupID)
		}
	}
	f.alignWeekToTimePriority(dates, groupIDs)
}

// alignWeekToTimePriority slides every group's contiguous daily block toward the
// preferred half of the day — morning → earliest slots, afternoon → latest slots.
// It is a no-op when time priority is none/unset.
func (f *fixer) alignWeekToTimePriority(dates []date, groupIDs []uint64) {
	if f.cfg.timePriority != domain.TimePriorityMorning &&
		f.cfg.timePriority != domain.TimePriorityAfternoon {
		return
	}
	for _, d := range dates {
		for _, groupID := range groupIDs {
			f.alignDayForGroup(d, groupID)
		}
	}
}

// alignDayForGroup relocates a group's lessons on a day to a contiguous block
// anchored at the preferred edge of the day. The move is all-or-nothing: if any
// lesson cannot occupy its target slot (teacher busy elsewhere, forbidden slot,
// etc.) the original placement is restored, so the schedule stays valid and
// gap-free. Time priority is a soft preference, never a hard constraint.
func (f *fixer) alignDayForGroup(d date, groupID uint64) {
	nums := f.coord.schedule.getScheduledLessonNumbers(d, groupID)
	k := len(nums)
	if k == 0 {
		return
	}

	// Desired contiguous block aligned to the preferred edge.
	targetStart := f.cfg.firstLessonNumber
	if f.cfg.timePriority == domain.TimePriorityAfternoon {
		targetStart = max(f.cfg.lastLessonNumber-k+1, f.cfg.firstLessonNumber)
	}

	targets := make([]int, k)
	aligned := true
	for i := range k {
		targets[i] = targetStart + i
		if targets[i] != nums[i] {
			aligned = false
		}
	}
	if aligned {
		return
	}

	// Collect the group's lessons in slot order.
	lessons := make([]*lesson, 0, k)
	for _, ln := range nums {
		for _, l := range f.coord.schedule.getLessons(d, ln) {
			if l.groupID == groupID {
				lessons = append(lessons, l)
				break
			}
		}
	}
	if len(lessons) != k {
		return // unexpected shape; leave as-is
	}

	// Detach the whole block, freeing the group's and teachers' old slots.
	for i, ln := range nums {
		l := lessons[i]
		f.coord.schedule.remove(d, ln, l)
		f.coord.availability.markFree(d, groupID, ln)
		for _, tid := range l.teacherIDs() {
			f.coord.availability.markFree(d, tid, ln)
		}
	}

	// Verify every lesson can take its target slot before committing.
	ok := true
	for i, l := range lessons {
		ts := targets[i]
		participants := append([]uint64{groupID}, l.teacherIDs()...)
		if !f.coord.availability.areFree(d, participants, ts) ||
			f.coord.schedule.hasConflict(d, ts, groupID, l.teacherIDs()) {
			ok = false
			break
		}
	}

	placeSlots := targets
	if !ok {
		placeSlots = nums // restore original positions
	}
	for i, l := range lessons {
		s := placeSlots[i]
		f.coord.schedule.add(d, s, l)
		f.coord.availability.markBusy(d, groupID, s)
		for _, tid := range l.teacherIDs() {
			f.coord.availability.markBusy(d, tid, s)
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
