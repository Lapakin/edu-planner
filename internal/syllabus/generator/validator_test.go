package generator

import (
	"errors"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func makeSettingsForValidator(maxPerDay, minPerDay int, noGaps bool) *settings {
	ts := makeTemplateSetting(2.0, maxPerDay*2, 40)
	bells := makeBellSchedules(maxPerDay + 2)
	r := &domain.ScheduleRestriction{
		MinGroupLessonsPerDay:   minPerDay,
		MaxGroupLessonsPerDay:   maxPerDay,
		MaxTeacherLessonsPerDay: 5,
		NoGapsInGroupSchedule:   noGaps,
	}
	return newSettings(ts, bells, r)
}

func makeDates(n int) []date {
	dates := make([]date, n)
	for i := range n {
		dates[i] = date{year: 2024, month: 9, day: i + 2}
	}
	return dates
}

func makeLesson(groupID, subjectID uint64) *lesson {
	return &lesson{groupID: groupID, subjectID: subjectID}
}

// --- Pre-generation validator tests ---

func TestValidator_ImpossibleWhenGroupNeedsMoreThanAvailableSlots(t *testing.T) {
	cfg := makeSettingsForValidator(4, 0, false)
	v := newValidator(cfg)

	dates := makeDates(1) // 1 day, max 4 slots → 4 total
	groupIDs := []uint64{1}
	lessons := make([]*lesson, 5) // 5 lessons > 4 slots
	for i := range 5 {
		lessons[i] = makeLesson(1, uint64(i+10))
	}

	err := v.validate(lessons, dates, groupIDs)
	if !errors.Is(err, ErrValidationFailed) {
		t.Errorf("expected ErrValidationFailed, got %v", err)
	}
}

func TestValidator_ImpossibleWhenGroupHasFewerThanMinLessonsPerDay(t *testing.T) {
	cfg := makeSettingsForValidator(4, 2, false)
	v := newValidator(cfg)

	dates := makeDates(5) // 5 days
	groupIDs := []uint64{1}
	lessons := []*lesson{makeLesson(1, 10)} // only 1 lesson, but min=2 → impossible

	err := v.validate(lessons, dates, groupIDs)
	if !errors.Is(err, ErrValidationFailed) {
		t.Errorf("expected ErrValidationFailed for 1 lesson with min=2, got %v", err)
	}
}

func TestValidator_ValidWhenGroupLessonsEqualMin(t *testing.T) {
	cfg := makeSettingsForValidator(4, 2, false)
	v := newValidator(cfg)

	dates := makeDates(5)
	groupIDs := []uint64{1}
	lessons := []*lesson{makeLesson(1, 10), makeLesson(1, 11)} // 2 lessons, min=2 → OK

	if err := v.validate(lessons, dates, groupIDs); err != nil {
		t.Errorf("expected no error for 2 lessons with min=2, got %v", err)
	}
}

func TestValidator_ValidWhenNoMinConstraint(t *testing.T) {
	cfg := makeSettingsForValidator(4, 0, false)
	v := newValidator(cfg)

	dates := makeDates(5)
	groupIDs := []uint64{1}
	lessons := []*lesson{makeLesson(1, 10)} // 1 lesson, min=0 → OK

	if err := v.validate(lessons, dates, groupIDs); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- Post-placement validator tests ---

func TestValidatePlacement_PassesWhenNoLessons(t *testing.T) {
	cfg := makeSettingsForValidator(4, 2, true)
	sched := newSchedule()
	dates := makeDates(3)

	if err := validatePlacement(sched, dates, []uint64{1, 2}, cfg); err != nil {
		t.Errorf("expected no error for empty schedule, got %v", err)
	}
}

func TestValidatePlacement_FailsOnGap(t *testing.T) {
	cfg := makeSettingsForValidator(4, 2, true) // noGaps=true
	sched := newSchedule()
	d := makeDates(1)[0]
	groupID := uint64(1)

	// Place lessons at slots 1 and 3 (gap at slot 2)
	sched.add(d, 1, &lesson{groupID: groupID, subjectID: 10})
	sched.add(d, 3, &lesson{groupID: groupID, subjectID: 11})

	err := validatePlacement(sched, makeDates(1), []uint64{groupID}, cfg)
	if !errors.Is(err, ErrUnableToSchedule) {
		t.Errorf("expected ErrUnableToSchedule for gap, got %v", err)
	}
}

func TestValidatePlacement_PassesWithNoGap(t *testing.T) {
	cfg := makeSettingsForValidator(4, 2, true)
	sched := newSchedule()
	d := makeDates(1)[0]
	groupID := uint64(1)

	// Place lessons at consecutive slots 2 and 3
	sched.add(d, 2, &lesson{groupID: groupID, subjectID: 10})
	sched.add(d, 3, &lesson{groupID: groupID, subjectID: 11})

	if err := validatePlacement(sched, makeDates(1), []uint64{groupID}, cfg); err != nil {
		t.Errorf("expected no error for consecutive lessons, got %v", err)
	}
}

func TestValidatePlacement_FailsOnBelowMinLessons(t *testing.T) {
	cfg := makeSettingsForValidator(4, 2, false) // min=2, noGaps=false
	sched := newSchedule()
	d := makeDates(1)[0]
	groupID := uint64(1)

	// Only 1 lesson on this day — violates min=2
	sched.add(d, 1, &lesson{groupID: groupID, subjectID: 10})

	err := validatePlacement(sched, makeDates(1), []uint64{groupID}, cfg)
	if !errors.Is(err, ErrUnableToSchedule) {
		t.Errorf("expected ErrUnableToSchedule for 1 lesson with min=2, got %v", err)
	}
}

func TestValidatePlacement_PassesWithMinLessons(t *testing.T) {
	cfg := makeSettingsForValidator(4, 2, false) // min=2
	sched := newSchedule()
	d := makeDates(1)[0]
	groupID := uint64(1)

	sched.add(d, 1, &lesson{groupID: groupID, subjectID: 10})
	sched.add(d, 2, &lesson{groupID: groupID, subjectID: 11})

	if err := validatePlacement(sched, makeDates(1), []uint64{groupID}, cfg); err != nil {
		t.Errorf("expected no error for 2 lessons with min=2, got %v", err)
	}
}

func TestValidatePlacement_SkipsDaysWithNoLessons(t *testing.T) {
	cfg := makeSettingsForValidator(4, 2, true)
	sched := newSchedule()
	dates := makeDates(3) // 3 days

	// Only put lessons on day 0 (2 consecutive)
	sched.add(dates[0], 1, &lesson{groupID: 1, subjectID: 10})
	sched.add(dates[0], 2, &lesson{groupID: 1, subjectID: 11})
	// days[1] and days[2] are empty for this group

	if err := validatePlacement(sched, dates, []uint64{1}, cfg); err != nil {
		t.Errorf("expected no error when other days are empty, got %v", err)
	}
}
