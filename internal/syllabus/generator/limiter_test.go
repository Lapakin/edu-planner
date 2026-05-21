package generator

import (
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func makeSettingsForLimiter(maxGroupPerDay, maxTeacherPerDay int) *settings {
	ts := makeTemplateSetting(2, maxGroupPerDay*2, 40)
	bells := makeBellSchedules(maxGroupPerDay + 2)
	r := &domain.ScheduleRestriction{
		MaxGroupLessonsPerDay:   maxGroupPerDay,
		MaxTeacherLessonsPerDay: maxTeacherPerDay,
		MinGroupLessonsPerDay:   0,
		NoGapsInGroupSchedule:   false,
	}
	return newSettings(ts, bells, r, nil, nil)
}

func TestLimiter_CanAddLessonForGroup(t *testing.T) {
	cfg := makeSettingsForLimiter(4, 5)
	lim := newLimiter(cfg)
	sched := newSchedule()

	d := date{year: 2024, month: 9, day: 2}
	groupID := uint64(1)

	lesson1 := &lesson{groupID: groupID, subjectID: 10}
	lesson2 := &lesson{groupID: groupID, subjectID: 11}
	lesson3 := &lesson{groupID: groupID, subjectID: 12}
	lesson4 := &lesson{groupID: groupID, subjectID: 13}
	if !lim.canAddLessonForGroup(sched, d, groupID) {
		t.Error("should be able to add first lesson")
	}
	sched.add(d, 1, lesson1)
	sched.add(d, 2, lesson2)
	sched.add(d, 3, lesson3)
	sched.add(d, 4, lesson4)

	if lim.canAddLessonForGroup(sched, d, groupID) {
		t.Error("should not be able to add 5th lesson (max=4)")
	}

	// Verify a different group can still add
	if !lim.canAddLessonForGroup(sched, d, 2) {
		t.Error("different group should still be able to add lessons")
	}
}

func TestLimiter_CanAddLessonForTeacher(t *testing.T) {
	cfg := makeSettingsForLimiter(4, 5)
	lim := newLimiter(cfg)
	sched := newSchedule()

	d := date{year: 2024, month: 9, day: 2}
	teacherID := uint64(42)

	addTeacherLesson := func(slot int) {
		l := &lesson{
			groupID:   uint64(slot),
			subjectID: 10,
			subLessons: []*internalSubLesson{
				{teacherID: teacherID, groupID: uint64(slot)},
			},
		}
		sched.add(d, slot, l)
	}

	for i := 1; i <= 4; i++ {
		addTeacherLesson(i)
		if !lim.canAddLessonForTeacher(sched, d, teacherID) {
			t.Errorf("should be able to add lesson %d for teacher (max=5)", i+1)
		}
	}

	addTeacherLesson(5)
	if lim.canAddLessonForTeacher(sched, d, teacherID) {
		t.Error("should not be able to add 6th lesson for teacher (max=5)")
	}
}

func TestLimiter_GroupAndTeacherLimitsAreSeparate(t *testing.T) {
	// Group max=4, teacher max=5 — they are separate constraints
	cfg := makeSettingsForLimiter(4, 5)
	lim := newLimiter(cfg)

	// Teacher limit is 5, not 4
	if cfg.maxTeacherLessonsPerDay != 5 {
		t.Errorf("expected maxTeacherLessonsPerDay=5, got %d", cfg.maxTeacherLessonsPerDay)
	}
	if cfg.maxGroupLessonsPerDay != 4 {
		t.Errorf("expected maxGroupLessonsPerDay=4, got %d", cfg.maxGroupLessonsPerDay)
	}
	_ = lim
}
