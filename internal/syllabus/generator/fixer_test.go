package generator

import (
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func makeFixerCoordinator() (*coordinator, *settings) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(4) // slots 1..4
	r := &domain.ScheduleRestriction{
		MinGroupLessonsPerDay:   0,
		MaxGroupLessonsPerDay:   4,
		MaxTeacherLessonsPerDay: 5,
		NoGapsInGroupSchedule:   false,
	}
	cfg := newSettings(ts, bells, r, nil, nil)
	coord := newCoordinator(cfg, []uint64{101})
	return coord, cfg
}

func TestFixer_FixDayForGroup_MovesLessonToFillGap(t *testing.T) {
	coord, _ := makeFixerCoordinator()
	d := date{year: 2024, month: 9, day: 2}
	groupID := uint64(1)
	teacherID := uint64(42)

	// Place lesson at slot 1
	lesson1 := &lesson{
		groupID:   groupID,
		subjectID: 10,
		subLessons: []*internalSubLesson{
			{groupID: groupID, teacherID: teacherID},
		},
	}
	coord.schedule.add(d, 1, lesson1)
	coord.availability.markBusy(d, groupID, 1)
	coord.availability.markBusy(d, teacherID, 1)

	// Place lesson at slot 3 (gap at 2)
	lesson3 := &lesson{
		groupID:   groupID,
		subjectID: 11,
		subLessons: []*internalSubLesson{
			{groupID: groupID, teacherID: teacherID},
		},
	}
	coord.schedule.add(d, 3, lesson3)
	coord.availability.markBusy(d, groupID, 3)
	coord.availability.markBusy(d, teacherID, 3)

	// Fix the gap
	coord.fixer.fixDayForGroup(d, groupID)

	// After fixing, slot 3 lesson should have moved to slot 2
	scheduledNumbers := coord.schedule.getScheduledLessonNumbers(d, groupID)
	for _, n := range scheduledNumbers {
		if n == 3 {
			// If still at 3, then we need both 2 and NOT a gap
			// Check if slot 2 is now occupied
			lessonsAt2 := coord.schedule.getLessons(d, 2)
			if len(lessonsAt2) == 0 {
				t.Error("expected lesson at slot 2 after fixing gap (moved from slot 3)")
			}
		}
	}

	// Verify no gap: consecutive from slot 1
	nums := coord.schedule.getScheduledLessonNumbers(d, groupID)
	if len(nums) >= 2 {
		for i := 1; i < len(nums); i++ {
			if nums[i]-nums[i-1] > 1 {
				t.Errorf("gap still exists at indices %d-%d: slots %v", i-1, i, nums)
			}
		}
	}
}

func TestFixer_FixDayForGroup_SingleLesson_NoChange(t *testing.T) {
	coord, _ := makeFixerCoordinator()
	d := date{year: 2024, month: 9, day: 2}
	groupID := uint64(1)

	l := &lesson{groupID: groupID, subjectID: 10}
	coord.schedule.add(d, 3, l)
	coord.availability.markBusy(d, groupID, 3)

	// Single lesson — no action should be taken
	coord.fixer.fixDayForGroup(d, groupID)

	nums := coord.schedule.getScheduledLessonNumbers(d, groupID)
	if len(nums) != 1 || nums[0] != 3 {
		t.Errorf("expected lesson still at slot 3, got %v", nums)
	}
}

func TestFixer_FixDayForGroup_NoGap_NoChange(t *testing.T) {
	coord, _ := makeFixerCoordinator()
	d := date{year: 2024, month: 9, day: 2}
	groupID := uint64(1)

	l1 := &lesson{groupID: groupID, subjectID: 10}
	l2 := &lesson{groupID: groupID, subjectID: 11}
	coord.schedule.add(d, 1, l1)
	coord.schedule.add(d, 2, l2)
	coord.availability.markBusy(d, groupID, 1)
	coord.availability.markBusy(d, groupID, 2)

	coord.fixer.fixDayForGroup(d, groupID)

	nums := coord.schedule.getScheduledLessonNumbers(d, groupID)
	if len(nums) != 2 || nums[0] != 1 || nums[1] != 2 {
		t.Errorf("expected slots [1,2] unchanged, got %v", nums)
	}
}

func TestFixer_FixWeekSchedule_ProcessesAllDatesAndGroups(_ *testing.T) {
	coord, _ := makeFixerCoordinator()
	dates := []date{
		{year: 2024, month: 9, day: 2},
		{year: 2024, month: 9, day: 3},
	}
	groupIDs := []uint64{1, 2}

	// Add lessons with gaps for group 1 on both dates
	for _, d := range dates {
		l1 := &lesson{groupID: 1, subjectID: 10}
		l3 := &lesson{groupID: 1, subjectID: 11}
		coord.schedule.add(d, 1, l1)
		coord.schedule.add(d, 3, l3)
		coord.availability.markBusy(d, 1, 1)
		coord.availability.markBusy(d, 1, 3)
	}

	// Should not panic
	coord.fixer.fixWeekSchedule(dates, groupIDs)
}
