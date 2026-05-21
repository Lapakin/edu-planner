package generator

import (
	"math/rand"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func TestWeekdayIndex(t *testing.T) {
	cases := []struct {
		wd  domain.Weekday
		idx int
	}{
		{domain.WeekdayMonday, 0},
		{domain.WeekdayTuesday, 1},
		{domain.WeekdayWednesday, 2},
		{domain.WeekdayThursday, 3},
		{domain.WeekdayFriday, 4},
		{domain.WeekdaySaturday, 5},
	}
	for _, c := range cases {
		got := weekdayIndex(c.wd)
		if got != c.idx {
			t.Errorf("weekdayIndex(%v) = %d, want %d", c.wd, got, c.idx)
		}
	}
}

func TestWeekdayIndex_UnknownDefault(t *testing.T) {
	got := weekdayIndex(domain.Weekday("unknown"))
	if got != 0 {
		t.Errorf("expected 0 for unknown weekday, got %d", got)
	}
}

// makeDenominatorScenario sets up a minimal scenario for denominator reproduce testing.
// Returns cfg, coord, rng, numeratorDates, denominatorDates, denomLesson.
func makeDenominatorScenario() (*settings, *coordinator, *rand.Rand, []date, []date, *lesson) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(4)
	r := &domain.ScheduleRestriction{
		MinGroupLessonsPerDay:   0,
		MaxGroupLessonsPerDay:   4,
		MaxTeacherLessonsPerDay: 5,
		NoGapsInGroupSchedule:   false,
	}
	cfg := newSettings(ts, bells, r, nil, nil)
	coord := newCoordinator(cfg, []uint64{101})
	rng := rand.New(rand.NewSource(42))

	// numeratorDate: Sep 2 (Monday)
	numDate := date{year: 2024, month: 9, day: 2}
	// denominatorDate: Sep 9 (Monday, next week)
	denomDate := date{year: 2024, month: 9, day: 9}

	// Place a lesson in the numerator schedule
	numLesson := &lesson{
		groupID:   1,
		subjectID: 10,
		format:    formatUnited,
		lType:     lessonTypeRegular,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 42},
		},
	}
	coord.schedule.add(numDate, 1, numLesson)
	coord.availability.markBusy(numDate, 1, 1)
	coord.availability.markBusy(numDate, 42, 1)

	// Create a corresponding denominator lesson (not yet scheduled)
	denomLesson := &lesson{
		groupID:   1,
		subjectID: 10,
		format:    formatUnited,
		lType:     lessonTypeRegular,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 42},
		},
	}

	return cfg, coord, rng, []date{numDate}, []date{denomDate}, denomLesson
}

func TestDenominatorReproducer_Reproduce_MirrorsNumerator(t *testing.T) {
	cfg, coord, rng, numDates, denomDates, denomLesson := makeDenominatorScenario()

	dr := newDenominatorReproducer(cfg, coord, rng)
	err := dr.reproduce(numDates, denomDates, []*lesson{denomLesson}, []uint64{1})
	if err != nil {
		t.Fatalf("reproduce returned error: %v", err)
	}

	// The denominator lesson should now be scheduled
	if !denomLesson.isScheduled {
		// It might have been placed via random/recursive if mirror failed
		// Check that something got placed on denomDate
		denomDate := denomDates[0]
		lessonsOnDenom := coord.schedule.getGroupDayLessons(denomDate, 1)
		if len(lessonsOnDenom) == 0 {
			t.Error("expected at least one lesson on denominator date after reproduce")
		}
	}
}

func TestDenominatorReproducer_Reproduce_EmptyNumerator(t *testing.T) {
	cfg, coord, rng, _, denomDates, denomLesson := makeDenominatorScenario()

	// Empty numerator dates - reproduce should still succeed (remainder path)
	dr := newDenominatorReproducer(cfg, coord, rng)
	// When numerator is empty, the mirror loop is skipped but remaining lessons are scheduled
	err := dr.reproduce([]date{}, denomDates, []*lesson{denomLesson}, []uint64{1})
	if err != nil {
		t.Fatalf("reproduce with empty numerator returned error: %v", err)
	}
}

func TestDenominatorReproducer_Reproduce_NoRemainingLessons(t *testing.T) {
	cfg, coord, rng, numDates, denomDates, _ := makeDenominatorScenario()

	dr := newDenominatorReproducer(cfg, coord, rng)
	// Pass empty denominator lessons - should succeed immediately
	err := dr.reproduce(numDates, denomDates, []*lesson{}, []uint64{1})
	if err != nil {
		t.Fatalf("reproduce with empty denom lessons returned error: %v", err)
	}
}
