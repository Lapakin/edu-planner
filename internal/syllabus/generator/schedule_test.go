package generator

import (
	"testing"
)

func makeTestLesson(groupID, subjectID uint64, teacherID uint64) *lesson {
	return &lesson{
		groupID:   groupID,
		subjectID: subjectID,
		format:    formatUnited,
		lType:     lessonTypeRegular,
		subLessons: []*internalSubLesson{
			{groupID: groupID, teacherID: teacherID},
		},
	}
}

func TestSchedule_Add_SetsIsScheduled(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}
	l := makeTestLesson(1, 10, 42)

	s.add(d, 1, l)

	if !l.isScheduled {
		t.Error("expected isScheduled=true after add")
	}
}

func TestSchedule_Add_RetrievableViaGetLessons(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}
	l := makeTestLesson(1, 10, 42)

	s.add(d, 1, l)
	lessons := s.getLessons(d, 1)

	if len(lessons) != 1 {
		t.Fatalf("expected 1 lesson, got %d", len(lessons))
	}
	if lessons[0] != l {
		t.Error("retrieved lesson is not the same as added")
	}
}

func TestSchedule_Remove_SetsIsScheduledFalse(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}
	l := makeTestLesson(1, 10, 42)

	s.add(d, 1, l)
	s.remove(d, 1, l)

	if l.isScheduled {
		t.Error("expected isScheduled=false after remove")
	}
}

func TestSchedule_Remove_LessonNotInSchedule(_ *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}
	l := makeTestLesson(1, 10, 42)

	// Remove without adding should not panic
	s.remove(d, 1, l)
}

func TestSchedule_GetLessons_EmptySlot(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	lessons := s.getLessons(d, 1)
	if lessons != nil {
		t.Errorf("expected nil for empty slot, got %v", lessons)
	}
}

func TestSchedule_GetDaySchedule_NilForEmptyDate(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	day := s.getDaySchedule(d)
	if day != nil {
		t.Error("expected nil day schedule for empty date")
	}
}

func TestSchedule_GetDaySchedule_ReturnsAllSlotsOnDate(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l1 := makeTestLesson(1, 10, 42)
	l2 := makeTestLesson(1, 11, 42)
	s.add(d, 1, l1)
	s.add(d, 2, l2)

	day := s.getDaySchedule(d)
	if day == nil {
		t.Fatal("expected non-nil day schedule")
	}
	if len(day) != 2 {
		t.Errorf("expected 2 lesson slots, got %d", len(day))
	}
}

func TestSchedule_GetGroupDayLessons(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l1 := makeTestLesson(1, 10, 42)
	l2 := makeTestLesson(2, 11, 43) // different group
	s.add(d, 1, l1)
	s.add(d, 1, l2)

	group1Lessons := s.getGroupDayLessons(d, 1)
	if len(group1Lessons) != 1 {
		t.Errorf("expected 1 lesson for group 1, got %d", len(group1Lessons))
	}
	if group1Lessons[0].groupID != 1 {
		t.Error("expected lesson for group 1")
	}
}

func TestSchedule_GetGroupDayLessonCount(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	s.add(d, 1, makeTestLesson(1, 10, 42))
	s.add(d, 2, makeTestLesson(1, 11, 42))
	s.add(d, 3, makeTestLesson(2, 12, 43)) // different group

	count := s.getGroupDayLessonCount(d, 1)
	if count != 2 {
		t.Errorf("expected count=2 for group 1, got %d", count)
	}
}

func TestSchedule_GetTeacherDayLessons(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l1 := makeTestLesson(1, 10, 42)
	l2 := makeTestLesson(2, 11, 43) // different teacher
	s.add(d, 1, l1)
	s.add(d, 2, l2)

	teacher42Lessons := s.getTeacherDayLessons(d, 42)
	if len(teacher42Lessons) != 1 {
		t.Errorf("expected 1 lesson for teacher 42, got %d", len(teacher42Lessons))
	}
}

func TestSchedule_GetGroupWeekLessonCount(t *testing.T) {
	s := newSchedule()
	dates := []date{
		{year: 2024, month: 9, day: 2},
		{year: 2024, month: 9, day: 3},
	}

	s.add(dates[0], 1, makeTestLesson(1, 10, 42))
	s.add(dates[1], 1, makeTestLesson(1, 11, 42))

	count := s.getGroupWeekLessonCount(dates, 1)
	if count != 2 {
		t.Errorf("expected week count=2 for group 1, got %d", count)
	}
}

func TestSchedule_GetTeacherWeekLessonCount(t *testing.T) {
	s := newSchedule()
	dates := []date{
		{year: 2024, month: 9, day: 2},
		{year: 2024, month: 9, day: 3},
	}

	s.add(dates[0], 1, makeTestLesson(1, 10, 42))
	s.add(dates[1], 1, makeTestLesson(1, 11, 42))

	count := s.getTeacherWeekLessonCount(dates, 42)
	if count != 2 {
		t.Errorf("expected week count=2 for teacher 42, got %d", count)
	}
}

func TestSchedule_GetScheduledLessonNumbers(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	s.add(d, 3, makeTestLesson(1, 10, 42))
	s.add(d, 1, makeTestLesson(1, 11, 42))

	nums := s.getScheduledLessonNumbers(d, 1)
	if len(nums) != 2 {
		t.Fatalf("expected 2 lesson numbers, got %v", nums)
	}
	// Should be sorted
	if nums[0] != 1 || nums[1] != 3 {
		t.Errorf("expected [1, 3], got %v", nums)
	}
}

func TestSchedule_GetScheduledLessonNumbers_OtherGroupNotIncluded(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	s.add(d, 1, makeTestLesson(1, 10, 42))
	s.add(d, 2, makeTestLesson(2, 11, 43)) // group 2

	nums := s.getScheduledLessonNumbers(d, 1)
	if len(nums) != 1 || nums[0] != 1 {
		t.Errorf("expected only slot 1 for group 1, got %v", nums)
	}
}

func TestSchedule_HasConflict_GroupConflict(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l := makeTestLesson(1, 10, 42)
	s.add(d, 1, l)

	// Same group at same slot
	if !s.hasConflict(d, 1, 1, []uint64{99}) {
		t.Error("expected group conflict")
	}
}

func TestSchedule_HasConflict_TeacherConflict(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l := makeTestLesson(1, 10, 42) // teacher 42
	s.add(d, 1, l)

	// Different group, same teacher
	if !s.hasConflict(d, 1, 2, []uint64{42}) {
		t.Error("expected teacher conflict")
	}
}

func TestSchedule_HasConflict_NoConflict(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l := makeTestLesson(1, 10, 42)
	s.add(d, 1, l)

	// Different group, different teacher, different slot
	if s.hasConflict(d, 2, 2, []uint64{99}) {
		t.Error("expected no conflict for different slot")
	}

	// Different group, different teacher, same slot
	if s.hasConflict(d, 1, 2, []uint64{99}) {
		t.Error("expected no conflict for different group and teacher")
	}
}

func TestSchedule_FindConflictingLesson_ReturnsConflict(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l := makeTestLesson(1, 10, 42)
	s.add(d, 1, l)

	// Conflict by group ID
	conflicting := s.findConflictingLesson(d, 1, 1, []uint64{99})
	if conflicting != l {
		t.Error("expected to find conflicting lesson by group ID")
	}
}

func TestSchedule_FindConflictingLesson_ReturnsNilWhenNone(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l := makeTestLesson(1, 10, 42)
	s.add(d, 1, l)

	// No conflict
	conflicting := s.findConflictingLesson(d, 1, 2, []uint64{99})
	if conflicting != nil {
		t.Error("expected nil for no conflict")
	}
}

func TestSchedule_AllDates_ReturnsAllDates(t *testing.T) {
	s := newSchedule()
	d1 := date{year: 2024, month: 9, day: 2}
	d2 := date{year: 2024, month: 9, day: 3}
	d3 := date{year: 2024, month: 9, day: 5}

	s.add(d1, 1, makeTestLesson(1, 10, 42))
	s.add(d3, 1, makeTestLesson(1, 11, 42))
	s.add(d2, 1, makeTestLesson(1, 12, 42))

	dates := s.allDates()
	if len(dates) != 3 {
		t.Fatalf("expected 3 dates, got %d", len(dates))
	}
	// Should be sorted
	if !dates[0].equal(d1) || !dates[1].equal(d2) || !dates[2].equal(d3) {
		t.Errorf("dates not sorted correctly: %v", dates)
	}
}

func TestSchedule_AllDates_EmptySchedule(t *testing.T) {
	s := newSchedule()
	dates := s.allDates()
	if len(dates) != 0 {
		t.Errorf("expected empty dates for empty schedule, got %v", dates)
	}
}

func TestSchedule_Remove_OnlyRemovesSpecificLesson(t *testing.T) {
	s := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l1 := makeTestLesson(1, 10, 42)
	l2 := makeTestLesson(2, 11, 43)
	s.add(d, 1, l1)
	s.add(d, 1, l2)

	s.remove(d, 1, l1)

	lessons := s.getLessons(d, 1)
	if len(lessons) != 1 {
		t.Fatalf("expected 1 lesson after removing l1, got %d", len(lessons))
	}
	if lessons[0] != l2 {
		t.Error("expected l2 to remain after removing l1")
	}
	if l1.isScheduled {
		t.Error("expected l1.isScheduled=false after remove")
	}
	if !l2.isScheduled {
		t.Error("expected l2.isScheduled=true (not removed)")
	}
}
