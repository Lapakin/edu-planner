package generator

import (
	"fmt"
	"testing"
)

func TestNewPattern(t *testing.T) {
	p := newPattern(1, 42, 10)
	expected := pattern(fmt.Sprintf("%d:%d:%d", 1, 42, 10))
	if p != expected {
		t.Errorf("expected %q, got %q", expected, p)
	}
}

func TestPatternController_SetDayLimit_ShouldSkipWhenAtLimit(t *testing.T) {
	pc := newPatternController()
	d := date{year: 2024, month: 9, day: 2}
	p := newPattern(1, 42, 10)

	pc.setDayLimit(p, 1)

	// First occurrence should not skip
	if pc.shouldSkip(d, p) {
		t.Error("expected shouldSkip=false before any occurrences")
	}

	// Add one occurrence
	pc.increment(d, p)

	// Now at limit — should skip
	if !pc.shouldSkip(d, p) {
		t.Error("expected shouldSkip=true when at day limit")
	}
}

func TestPatternController_SetDayLimit_BelowLimit(t *testing.T) {
	pc := newPatternController()
	d := date{year: 2024, month: 9, day: 2}
	p := newPattern(1, 42, 10)

	pc.setDayLimit(p, 2)
	pc.increment(d, p)

	// 1 occurrence, limit is 2 — should not skip
	if pc.shouldSkip(d, p) {
		t.Error("expected shouldSkip=false when below day limit")
	}
}

func TestPatternController_Increment_And_Decrement(t *testing.T) {
	pc := newPatternController()
	d := date{year: 2024, month: 9, day: 2}
	p := newPattern(1, 42, 10)

	pc.setDayLimit(p, 1)
	pc.increment(d, p)

	if !pc.shouldSkip(d, p) {
		t.Error("expected shouldSkip=true after increment to limit")
	}

	pc.decrement(d, p)

	if pc.shouldSkip(d, p) {
		t.Error("expected shouldSkip=false after decrement")
	}
}

func TestPatternController_Decrement_BelowZero_Noop(_ *testing.T) {
	pc := newPatternController()
	d := date{year: 2024, month: 9, day: 2}
	p := newPattern(1, 42, 10)

	// Decrement without increment should not panic
	pc.decrement(d, p)
}

func TestPatternController_SetDayBetween_ShouldSkipWhenAdjacentDay(t *testing.T) {
	pc := newPatternController()
	p := newPattern(1, 42, 10)

	d1 := date{year: 2024, month: 9, day: 2}
	d2 := d1.addDays(1) // Sep 3

	pc.setDayBetween(p, 1)

	// Put occurrence on d1
	pc.increment(d1, p)

	// d2 is 1 day away — should skip
	if !pc.shouldSkip(d2, p) {
		t.Error("expected shouldSkip=true for adjacent day with dayBetween=1")
	}
}

func TestPatternController_SetDayBetween_ShouldNotSkipWhenFarEnough(t *testing.T) {
	pc := newPatternController()
	p := newPattern(1, 42, 10)

	d1 := date{year: 2024, month: 9, day: 2}
	d2 := d1.addDays(2) // Sep 4, 2 days away

	pc.setDayBetween(p, 1)
	pc.increment(d1, p)

	// d2 is 2 days away, dayBetween=1 — should NOT skip
	if pc.shouldSkip(d2, p) {
		t.Error("expected shouldSkip=false when day is far enough from last occurrence")
	}
}

func TestPatternController_ShouldSkip_NeitherLimitNorBetween(t *testing.T) {
	pc := newPatternController()
	d := date{year: 2024, month: 9, day: 2}
	p := newPattern(1, 42, 10)

	// No limits set — should never skip
	if pc.shouldSkip(d, p) {
		t.Error("expected shouldSkip=false when no limits are set")
	}

	pc.increment(d, p)

	if pc.shouldSkip(d, p) {
		t.Error("expected shouldSkip=false when no limits are set even after increment")
	}
}

func TestPatternController_ShouldSkipSubLessons_AllFree(t *testing.T) {
	pc := newPatternController()
	d := date{year: 2024, month: 9, day: 2}
	groupID := uint64(1)
	subjectID := uint64(10)

	subLessons := []*internalSubLesson{
		{teacherID: 42},
		{teacherID: 43},
	}

	if pc.shouldSkipSubLessons(d, subLessons, subjectID, groupID) {
		t.Error("expected shouldSkipSubLessons=false when no limits set")
	}
}

func TestPatternController_ShouldSkipSubLessons_OneAtLimit(t *testing.T) {
	pc := newPatternController()
	d := date{year: 2024, month: 9, day: 2}
	groupID := uint64(1)
	subjectID := uint64(10)

	p42 := newPattern(groupID, 42, subjectID)
	pc.setDayLimit(p42, 1)
	pc.increment(d, p42)

	subLessons := []*internalSubLesson{
		{teacherID: 42},
		{teacherID: 43},
	}

	// Teacher 42 is at limit — should skip
	if !pc.shouldSkipSubLessons(d, subLessons, subjectID, groupID) {
		t.Error("expected shouldSkipSubLessons=true when one teacher is at limit")
	}
}
