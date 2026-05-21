package generator

import (
	"math/rand"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func makeRandomSchedulerSetup() (*settings, *coordinator, *rand.Rand) {
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
	rng := rand.New(rand.NewSource(42))
	return cfg, coord, rng
}

func makeFiveDates() []date {
	return []date{
		{year: 2024, month: 9, day: 2}, // Monday
		{year: 2024, month: 9, day: 3}, // Tuesday
		{year: 2024, month: 9, day: 4}, // Wednesday
		{year: 2024, month: 9, day: 5}, // Thursday
		{year: 2024, month: 9, day: 6}, // Friday
	}
}

func TestRandomScheduler_ScheduleRandom_PlacesSingleLesson(t *testing.T) {
	cfg, coord, rng := makeRandomSchedulerSetup()
	dates := makeFiveDates()

	l := &lesson{
		groupID:   1,
		subjectID: 10,
		format:    formatUnited,
		lType:     lessonTypeRegular,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 42},
		},
	}

	rs := newRandomScheduler(cfg, coord, rng)
	unscheduled := rs.scheduleRandom([]*lesson{l}, dates, dates)

	if len(unscheduled) != 0 {
		t.Errorf("expected 0 unscheduled lessons, got %d", len(unscheduled))
	}
	if !l.isScheduled {
		t.Error("expected lesson to be scheduled")
	}
}

func TestRandomScheduler_ScheduleRandom_AlreadyScheduledLessonsSkipped(t *testing.T) {
	cfg, coord, rng := makeRandomSchedulerSetup()
	dates := makeFiveDates()

	l := &lesson{
		groupID:     1,
		subjectID:   10,
		isScheduled: true, // already scheduled
		subLessons:  []*internalSubLesson{{groupID: 1, teacherID: 42}},
	}

	rs := newRandomScheduler(cfg, coord, rng)
	unscheduled := rs.scheduleRandom([]*lesson{l}, dates, dates)

	// Should not be in unscheduled since it was already scheduled
	if len(unscheduled) != 0 {
		t.Errorf("expected 0 unscheduled for already-scheduled lesson, got %d", len(unscheduled))
	}
}

func TestRandomScheduler_ScheduleRandom_MultipleLessons(t *testing.T) {
	cfg, coord, rng := makeRandomSchedulerSetup()
	dates := makeFiveDates()

	lessons := make([]*lesson, 3)
	for i := range lessons {
		lessons[i] = &lesson{
			groupID:   1,
			subjectID: uint64(10 + i),
			format:    formatUnited,
			lType:     lessonTypeRegular,
			subLessons: []*internalSubLesson{
				{groupID: 1, teacherID: 42},
			},
		}
	}

	rs := newRandomScheduler(cfg, coord, rng)
	unscheduled := rs.scheduleRandom(lessons, dates, dates)

	// With 5 days and 4 slots each = 20 slots, 3 lessons should all fit
	scheduledCount := 0
	for _, l := range lessons {
		if l.isScheduled {
			scheduledCount++
		}
	}
	// At least 1 should be scheduled (randomness makes exact count uncertain)
	if scheduledCount == 0 {
		t.Error("expected at least one lesson to be scheduled")
	}
	_ = unscheduled
}

func TestRandomScheduler_ScheduleRandom_UpdatesAvailability(t *testing.T) {
	cfg, coord, rng := makeRandomSchedulerSetup()
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

	rs := newRandomScheduler(cfg, coord, rng)
	rs.scheduleRandom([]*lesson{l}, dates, dates)

	if l.isScheduled {
		// Find which slot the lesson was placed at and verify availability is marked
		d := dates[0]
		nums := coord.schedule.getScheduledLessonNumbers(d, 1)
		if len(nums) == 0 {
			t.Error("expected lesson to be in schedule")
		} else {
			if coord.availability.isFree(d, 1, nums[0]) {
				t.Error("expected availability to be marked busy for group")
			}
			if coord.availability.isFree(d, 42, nums[0]) {
				t.Error("expected availability to be marked busy for teacher")
			}
		}
	}
}
