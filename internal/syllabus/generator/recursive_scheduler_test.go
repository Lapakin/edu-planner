package generator

import (
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func makeRecursiveSchedulerSetup() (*settings, *coordinator) {
	ts := makeTemplateSetting(2.0, 8, 40)
	bells := makeBellSchedules(4) // slots 1..4
	r := &domain.ScheduleRestriction{
		MinGroupLessonsPerDay:   0,
		MaxGroupLessonsPerDay:   4,
		MaxTeacherLessonsPerDay: 5,
		NoGapsInGroupSchedule:   false,
	}
	cfg := newSettings(ts, bells, r)
	coord := newCoordinator(cfg, []uint64{101})
	return cfg, coord
}

func TestRecursiveScheduler_Schedule_PlacesSingleLesson(t *testing.T) {
	cfg, coord := makeRecursiveSchedulerSetup()
	dates := []date{
		{year: 2024, month: 9, day: 2},
		{year: 2024, month: 9, day: 3},
		{year: 2024, month: 9, day: 4},
		{year: 2024, month: 9, day: 5},
		{year: 2024, month: 9, day: 6},
	}

	l := &lesson{
		groupID:   1,
		subjectID: 10,
		format:    formatUnited,
		lType:     lessonTypeRegular,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 42},
		},
	}

	rs := newRecursiveScheduler(cfg, coord)
	err := rs.schedule([]*lesson{l}, dates, dates)
	if err != nil {
		t.Fatalf("schedule() returned error: %v", err)
	}
	if !l.isScheduled {
		t.Error("expected lesson to be scheduled")
	}
}

func TestRecursiveScheduler_Schedule_ErrorWhenNoDates(t *testing.T) {
	cfg, coord := makeRecursiveSchedulerSetup()

	l := &lesson{
		groupID:   1,
		subjectID: 10,
		format:    formatUnited,
		lType:     lessonTypeRegular,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 42},
		},
	}

	rs := newRecursiveScheduler(cfg, coord)
	// No dates means impossible to schedule
	err := rs.schedule([]*lesson{l}, []date{}, []date{})
	if err == nil {
		t.Error("expected error when no dates are available")
	}
}

func TestRecursiveScheduler_Schedule_SkipsAlreadyScheduled(t *testing.T) {
	cfg, coord := makeRecursiveSchedulerSetup()
	dates := []date{{year: 2024, month: 9, day: 2}}

	l := &lesson{
		groupID:     1,
		subjectID:   10,
		isScheduled: true,
		subLessons:  []*internalSubLesson{{groupID: 1, teacherID: 42}},
	}

	rs := newRecursiveScheduler(cfg, coord)
	err := rs.schedule([]*lesson{l}, dates, dates)
	if err != nil {
		t.Fatalf("expected no error for already-scheduled lesson, got: %v", err)
	}
}

func TestRecursiveScheduler_Schedule_MultipleLessons(t *testing.T) {
	cfg, coord := makeRecursiveSchedulerSetup()
	dates := []date{
		{year: 2024, month: 9, day: 2},
		{year: 2024, month: 9, day: 3},
		{year: 2024, month: 9, day: 4},
	}

	lessons := []*lesson{
		{
			groupID:    1,
			subjectID:  10,
			format:     formatUnited,
			lType:      lessonTypeRegular,
			subLessons: []*internalSubLesson{{groupID: 1, teacherID: 42}},
		},
		{
			groupID:    1,
			subjectID:  11,
			format:     formatUnited,
			lType:      lessonTypeRegular,
			subLessons: []*internalSubLesson{{groupID: 1, teacherID: 42}},
		},
	}

	rs := newRecursiveScheduler(cfg, coord)
	err := rs.schedule(lessons, dates, dates)
	if err != nil {
		t.Fatalf("schedule() returned error for multiple lessons: %v", err)
	}

	for i, l := range lessons {
		if !l.isScheduled {
			t.Errorf("lesson %d should be scheduled", i)
		}
	}
}

func TestRecursiveScheduler_Schedule_UpdatesScheduleAndAvailability(t *testing.T) {
	cfg, coord := makeRecursiveSchedulerSetup()
	dates := []date{{year: 2024, month: 9, day: 2}}

	l := &lesson{
		groupID:   1,
		subjectID: 10,
		format:    formatUnited,
		lType:     lessonTypeRegular,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 42},
		},
	}

	rs := newRecursiveScheduler(cfg, coord)
	err := rs.schedule([]*lesson{l}, dates, dates)
	if err != nil {
		t.Fatalf("schedule() error: %v", err)
	}

	d := dates[0]
	nums := coord.schedule.getScheduledLessonNumbers(d, 1)
	if len(nums) == 0 {
		t.Error("expected lesson in schedule")
	} else {
		if coord.availability.isFree(d, 1, nums[0]) {
			t.Error("expected group 1 to be marked busy")
		}
		if coord.availability.isFree(d, 42, nums[0]) {
			t.Error("expected teacher 42 to be marked busy")
		}
	}
}
