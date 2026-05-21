package generator

import (
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func makeTemplateSetting(lessonsPerClass int, maxStudyHoursPerDay, maxTeacherHoursPerWeek int) *domain.ScheduleTemplateSetting {
	return &domain.ScheduleTemplateSetting{
		LessonsPerClass:           lessonsPerClass,
		StudyDaysMask:             domain.StudyDaysMaskDefault,
		MaxIdenticalLessonsPerDay: 2,
		MaxStudyHoursPerDay:       maxStudyHoursPerDay,
		MaxTeacherHoursPerWeek:    maxTeacherHoursPerWeek,
	}
}

func makeBellSchedules(count int) domain.BellSchedules {
	bs := make(domain.BellSchedules, count)
	for i := range count {
		bs[i] = &domain.BellSchedule{LessonNumber: i + 1}
	}
	return bs
}

func TestNewSettings_DefaultRestriction(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(7)
	s := newSettings(ts, bells, nil, nil, nil)

	// lessonsPerClass-based: 8/2 = 4 group lessons/day
	// restriction default: max=4, so min(4,4)=4
	if s.maxGroupLessonsPerDay != 4 {
		t.Errorf("expected maxGroupLessonsPerDay=4, got %d", s.maxGroupLessonsPerDay)
	}
	// restriction default teacher max = 5
	if s.maxTeacherLessonsPerDay != 5 {
		t.Errorf("expected maxTeacherLessonsPerDay=5, got %d", s.maxTeacherLessonsPerDay)
	}
	// restriction default min = 2
	if s.minGroupLessonsPerDay != 2 {
		t.Errorf("expected minGroupLessonsPerDay=2, got %d", s.minGroupLessonsPerDay)
	}
	// restriction default no-gaps = true
	if !s.noGapsRequired {
		t.Error("expected noGapsRequired=true")
	}
	// restriction default consecutive = 4
	if s.maxConsecutiveTeacherLessons != 4 {
		t.Errorf("expected maxConsecutiveTeacherLessons=4, got %d", s.maxConsecutiveTeacherLessons)
	}
	// restriction default time priority = none
	if s.timePriority != domain.TimePriorityNone {
		t.Errorf("expected timePriority=none, got %s", s.timePriority)
	}
	// restriction default allow flow = true
	if !s.allowFlowLessons {
		t.Error("expected allowFlowLessons=true")
	}
}

func TestNewSettings_RestrictionsOverrideHoursBasedMax(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40) // lessonsPerClass-based: 8/2 = 4 max/day
	bells := makeBellSchedules(7)
	r := &domain.ScheduleRestriction{
		MinGroupLessonsPerDay:        1,
		MaxGroupLessonsPerDay:        3, // more restrictive than 4
		MaxTeacherLessonsPerDay:      4,
		MaxConsecutiveTeacherLessons: 2,
		TimePriority:                 domain.TimePriorityMorning,
		AllowFlowLessons:             false,
		NoGapsInGroupSchedule:        false,
	}
	s := newSettings(ts, bells, r, nil, nil)

	if s.maxGroupLessonsPerDay != 3 {
		t.Errorf("expected maxGroupLessonsPerDay=3 (restriction cap), got %d", s.maxGroupLessonsPerDay)
	}
	if s.maxTeacherLessonsPerDay != 4 {
		t.Errorf("expected maxTeacherLessonsPerDay=4, got %d", s.maxTeacherLessonsPerDay)
	}
	if s.minGroupLessonsPerDay != 1 {
		t.Errorf("expected minGroupLessonsPerDay=1, got %d", s.minGroupLessonsPerDay)
	}
	if s.noGapsRequired {
		t.Error("expected noGapsRequired=false")
	}
	if s.maxConsecutiveTeacherLessons != 2 {
		t.Errorf("expected maxConsecutiveTeacherLessons=2, got %d", s.maxConsecutiveTeacherLessons)
	}
	if s.timePriority != domain.TimePriorityMorning {
		t.Errorf("expected timePriority=morning, got %s", s.timePriority)
	}
	if s.allowFlowLessons {
		t.Error("expected allowFlowLessons=false")
	}
}

func TestNewSettings_RestrictionDoesNotIncreaseHoursBasedMax(t *testing.T) {
	ts := makeTemplateSetting(2, 6, 40) // lessonsPerClass-based: 6/2 = 3 max/day
	bells := makeBellSchedules(7)
	r := &domain.ScheduleRestriction{
		MaxGroupLessonsPerDay:        5, // higher than lessonsPerClass-based, so lessonsPerClass-based wins
		MaxTeacherLessonsPerDay:      5,
		MaxConsecutiveTeacherLessons: 4,
	}
	s := newSettings(ts, bells, r, nil, nil)

	if s.maxGroupLessonsPerDay != 3 {
		t.Errorf("expected maxGroupLessonsPerDay=3 (lessonsPerClass cap), got %d", s.maxGroupLessonsPerDay)
	}
}

func TestNewSettings_TeacherWeeklyMax(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 36) // 36/2 = 18 teacher lessons/week
	bells := makeBellSchedules(7)
	s := newSettings(ts, bells, nil, nil, nil)

	if s.maxTeacherLessonsPerWeek() != 18 {
		t.Errorf("expected maxTeacherLessonsPerWeek=18, got %d", s.maxTeacherLessonsPerWeek())
	}
}

func TestNewSettings_LessonNumbers(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(6)
	s := newSettings(ts, bells, nil, nil, nil)

	nums := s.lessonNumbers()
	if len(nums) != 6 {
		t.Errorf("expected 6 lesson numbers, got %d", len(nums))
	}
	if nums[0] != 1 || nums[5] != 6 {
		t.Errorf("expected lessons 1-6, got %v", nums)
	}
}

func TestNewSettings_StudyDaysMask_DefaultIsMonFri(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(5)
	s := newSettings(ts, bells, nil, nil, nil)

	if len(s.educationWeek) != 5 {
		t.Errorf("expected 5 education days (Mon-Fri), got %d", len(s.educationWeek))
	}
	if s.educationWeek[0] != domain.WeekdayMonday {
		t.Errorf("expected first day=Monday, got %s", s.educationWeek[0])
	}
	if s.educationWeek[4] != domain.WeekdayFriday {
		t.Errorf("expected last day=Friday, got %s", s.educationWeek[4])
	}
}

func TestNewSettings_StudyDaysMask_ThreeDays(t *testing.T) {
	ts := &domain.ScheduleTemplateSetting{
		LessonsPerClass:           2,
		StudyDaysMask:             domain.StudyDayMon | domain.StudyDayWed | domain.StudyDayFri,
		MaxIdenticalLessonsPerDay: 2,
		MaxStudyHoursPerDay:       8,
		MaxTeacherHoursPerWeek:    40,
	}
	bells := makeBellSchedules(4)
	s := newSettings(ts, bells, nil, nil, nil)

	if len(s.educationWeek) != 3 {
		t.Errorf("expected 3 education days, got %d: %v", len(s.educationWeek), s.educationWeek)
	}
	if s.educationWeek[0] != domain.WeekdayMonday {
		t.Errorf("expected Mon, got %s", s.educationWeek[0])
	}
	if s.educationWeek[1] != domain.WeekdayWednesday {
		t.Errorf("expected Wed, got %s", s.educationWeek[1])
	}
	if s.educationWeek[2] != domain.WeekdayFriday {
		t.Errorf("expected Fri, got %s", s.educationWeek[2])
	}
}

func TestNewSettings_StudyDaysMask_EmptyFallsBackToDefault(t *testing.T) {
	ts := &domain.ScheduleTemplateSetting{
		LessonsPerClass:           2,
		StudyDaysMask:             0, // empty mask → fallback to default Mon-Fri
		MaxIdenticalLessonsPerDay: 2,
		MaxStudyHoursPerDay:       8,
		MaxTeacherHoursPerWeek:    40,
	}
	bells := makeBellSchedules(4)
	s := newSettings(ts, bells, nil, nil, nil)

	if len(s.educationWeek) != 5 {
		t.Errorf("expected fallback to 5 days, got %d", len(s.educationWeek))
	}
}

func TestNewSettings_TeacherForbiddenSlots(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(5)
	prefs := domain.TeacherSlotPreferences{
		{
			TeacherID:    42,
			Weekday:      domain.WeekdayMonday,
			LessonNumber: 1,
			SlotType:     domain.SlotTypeForbidden,
		},
	}
	s := newSettings(ts, bells, nil, prefs, nil)

	if s.teacherForbiddenSlots[42] == nil {
		t.Fatal("expected teacherForbiddenSlots to have entry for teacher 42")
	}
	if !s.teacherForbiddenSlots[42][domain.WeekdayMonday][1] {
		t.Error("expected lesson 1 on Monday to be forbidden for teacher 42")
	}
}

func TestNewSettings_TeacherPreferredSlots(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(5)
	prefs := domain.TeacherSlotPreferences{
		{
			TeacherID:    7,
			Weekday:      domain.WeekdayWednesday,
			LessonNumber: 3,
			SlotType:     domain.SlotTypePreferred,
		},
	}
	s := newSettings(ts, bells, nil, prefs, nil)

	if s.teacherPreferredSlots[7] == nil {
		t.Fatal("expected teacherPreferredSlots to have entry for teacher 7")
	}
	if !s.teacherPreferredSlots[7][domain.WeekdayWednesday][3] {
		t.Error("expected lesson 3 on Wednesday to be preferred for teacher 7")
	}
}

func TestNewSettings_NilPreferences(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(5)
	// nil preferences must not panic
	s := newSettings(ts, bells, nil, nil, nil)

	if s.teacherForbiddenSlots == nil {
		t.Error("expected teacherForbiddenSlots to be an empty map, not nil")
	}
	if s.teacherPreferredSlots == nil {
		t.Error("expected teacherPreferredSlots to be an empty map, not nil")
	}
}

func TestAvailability_WouldExceedConsecutive(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}

	// Teacher 1 has lessons at slots 1 and 2
	a.markBusy(d, 1, 1)
	a.markBusy(d, 1, 2)

	// Placing at slot 3 would create a run of 3 → exceeds maxConsecutive=2
	if !a.wouldExceedConsecutive(d, 1, 3, 2) {
		t.Error("expected consecutive violation for slot 3 with max=2 when slots 1,2 are busy")
	}

	// Placing at slot 4 (no continuity with 1,2) → run of 1 → OK
	if a.wouldExceedConsecutive(d, 1, 4, 2) {
		t.Error("expected no consecutive violation for slot 4 when slots 1,2 are busy")
	}

	// maxConsecutive=0 means no limit
	if a.wouldExceedConsecutive(d, 1, 3, 0) {
		t.Error("expected no violation when maxConsecutive=0")
	}
}

func TestAvailability_WouldExceedConsecutive_SingleLesson(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 3}

	// No existing lessons — placing anything creates a run of 1
	if a.wouldExceedConsecutive(d, 2, 1, 1) {
		t.Error("expected no violation: run=1 equals maxConsecutive=1")
	}
	// But a run of 1 with max=0 (no limit) is fine
	if a.wouldExceedConsecutive(d, 2, 1, 0) {
		t.Error("expected no violation with maxConsecutive=0")
	}
}
