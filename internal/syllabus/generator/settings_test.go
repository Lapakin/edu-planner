package generator

import (
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func makeTemplateSetting(hoursPerLesson float64, maxStudyHoursPerDay, maxTeacherHoursPerWeek int) *domain.ScheduleTemplateSetting {
	return &domain.ScheduleTemplateSetting{
		HoursPerLesson:            hoursPerLesson,
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
	ts := makeTemplateSetting(2.0, 8, 40)
	bells := makeBellSchedules(7)
	s := newSettings(ts, bells, nil)

	// hours-based: 8/2 = 4 group lessons/day
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
}

func TestNewSettings_RestrictionsOverrideHoursBasedMax(t *testing.T) {
	ts := makeTemplateSetting(2.0, 8, 40) // hours-based: 8/2 = 4 max/day
	bells := makeBellSchedules(7)
	r := &domain.ScheduleRestriction{
		MinGroupLessonsPerDay:   1,
		MaxGroupLessonsPerDay:   3, // more restrictive than 4
		MaxTeacherLessonsPerDay: 4,
		NoGapsInGroupSchedule:   false,
	}
	s := newSettings(ts, bells, r)

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
}

func TestNewSettings_RestrictionDoesNotIncreaseHoursBasedMax(t *testing.T) {
	ts := makeTemplateSetting(2.0, 6, 40) // hours-based: 6/2 = 3 max/day
	bells := makeBellSchedules(7)
	r := &domain.ScheduleRestriction{
		MaxGroupLessonsPerDay:   5, // higher than hours-based, so hours-based wins
		MaxTeacherLessonsPerDay: 5,
	}
	s := newSettings(ts, bells, r)

	if s.maxGroupLessonsPerDay != 3 {
		t.Errorf("expected maxGroupLessonsPerDay=3 (hours cap), got %d", s.maxGroupLessonsPerDay)
	}
}

func TestNewSettings_TeacherWeeklyMax(t *testing.T) {
	ts := makeTemplateSetting(2.0, 8, 36) // 36/2 = 18 teacher lessons/week
	bells := makeBellSchedules(7)
	s := newSettings(ts, bells, nil)

	if s.maxTeacherLessonsPerWeek() != 18 {
		t.Errorf("expected maxTeacherLessonsPerWeek=18, got %d", s.maxTeacherLessonsPerWeek())
	}
}

func TestNewSettings_LessonNumbers(t *testing.T) {
	ts := makeTemplateSetting(2.0, 8, 40)
	bells := makeBellSchedules(6)
	s := newSettings(ts, bells, nil)

	nums := s.lessonNumbers()
	if len(nums) != 6 {
		t.Errorf("expected 6 lesson numbers, got %d", len(nums))
	}
	if nums[0] != 1 || nums[5] != 6 {
		t.Errorf("expected lessons 1-6, got %v", nums)
	}
}
