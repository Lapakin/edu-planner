package generator

import (
	"testing"
	"time"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func TestNewDate(t *testing.T) {
	t0 := time.Date(2024, 9, 2, 15, 30, 0, 0, time.UTC)
	d := newDate(t0)
	if d.year != 2024 || d.month != time.September || d.day != 2 {
		t.Errorf("unexpected date: %+v", d)
	}
}

func TestDate_ToTime(t *testing.T) {
	d := date{year: 2024, month: time.September, day: 2}
	tt := d.toTime()
	if tt.Year() != 2024 || tt.Month() != time.September || tt.Day() != 2 {
		t.Errorf("unexpected time: %v", tt)
	}
	if tt.Hour() != 0 || tt.Minute() != 0 || tt.Second() != 0 {
		t.Errorf("expected zero time component, got %v", tt)
	}
}

func TestDate_Weekday_Monday(t *testing.T) {
	// 2024-09-02 is a Monday
	d := date{year: 2024, month: time.September, day: 2}
	if d.weekday() != domain.WeekdayMonday {
		t.Errorf("expected WeekdayMonday, got %v", d.weekday())
	}
}

func TestDate_Weekday_Friday(t *testing.T) {
	// 2024-09-06 is a Friday
	d := date{year: 2024, month: time.September, day: 6}
	if d.weekday() != domain.WeekdayFriday {
		t.Errorf("expected WeekdayFriday, got %v", d.weekday())
	}
}

func TestDate_AddDays(t *testing.T) {
	d := date{year: 2024, month: time.September, day: 2}
	d2 := d.addDays(5)
	if d2.year != 2024 || d2.month != time.September || d2.day != 7 {
		t.Errorf("expected 2024-09-07, got %+v", d2)
	}
}

func TestDate_AddDays_NegativeGoesBack(t *testing.T) {
	d := date{year: 2024, month: time.September, day: 5}
	d2 := d.addDays(-4)
	if d2.year != 2024 || d2.month != time.September || d2.day != 1 {
		t.Errorf("expected 2024-09-01, got %+v", d2)
	}
}

func TestDate_Before(t *testing.T) {
	d1 := date{year: 2024, month: time.September, day: 2}
	d2 := date{year: 2024, month: time.September, day: 5}
	if !d1.before(d2) {
		t.Error("expected d1.before(d2)")
	}
	if d2.before(d1) {
		t.Error("expected d2 not before d1")
	}
}

func TestDate_Equal(t *testing.T) {
	d1 := date{year: 2024, month: time.September, day: 2}
	d2 := date{year: 2024, month: time.September, day: 2}
	d3 := date{year: 2024, month: time.September, day: 3}
	if !d1.equal(d2) {
		t.Error("expected d1 == d2")
	}
	if d1.equal(d3) {
		t.Error("expected d1 != d3")
	}
}

func TestTimeWeekdayToWeekday(t *testing.T) {
	cases := []struct {
		in  time.Weekday
		out domain.Weekday
	}{
		{time.Monday, domain.WeekdayMonday},
		{time.Tuesday, domain.WeekdayTuesday},
		{time.Wednesday, domain.WeekdayWednesday},
		{time.Thursday, domain.WeekdayThursday},
		{time.Friday, domain.WeekdayFriday},
		{time.Saturday, domain.WeekdaySaturday},
		{time.Sunday, domain.WeekdayMonday}, // default
	}
	for _, c := range cases {
		got := timeWeekdayToWeekday(c.in)
		if got != c.out {
			t.Errorf("timeWeekdayToWeekday(%v) = %v, want %v", c.in, got, c.out)
		}
	}
}

func TestWeekdayToTimeWeekday(t *testing.T) {
	cases := []struct {
		in  domain.Weekday
		out time.Weekday
	}{
		{domain.WeekdayMonday, time.Monday},
		{domain.WeekdayTuesday, time.Tuesday},
		{domain.WeekdayWednesday, time.Wednesday},
		{domain.WeekdayThursday, time.Thursday},
		{domain.WeekdayFriday, time.Friday},
		{domain.WeekdaySaturday, time.Saturday},
	}
	for _, c := range cases {
		got := weekdayToTimeWeekday(c.in)
		if got != c.out {
			t.Errorf("weekdayToTimeWeekday(%v) = %v, want %v", c.in, got, c.out)
		}
	}
}

func TestCollectSemesterDates_TwoWeekMonFri(t *testing.T) {
	// 2024-09-02 (Mon) to 2024-09-13 (Fri) = 2 weeks
	start := time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC)

	numDates, denomDates := collectSemesterDates(start, end, domain.DefaultEducationWeek)

	// Week 0 is numerator (Mon-Fri week 1), week 1 is denominator (Mon-Fri week 2)
	if len(numDates) == 0 {
		t.Error("expected non-empty numerator dates")
	}
	if len(denomDates) == 0 {
		t.Error("expected non-empty denominator dates")
	}
	// Should have 5 + 5 = 10 total education days
	total := len(numDates) + len(denomDates)
	if total != 10 {
		t.Errorf("expected 10 total education days, got %d (num=%d, denom=%d)", total, len(numDates), len(denomDates))
	}
}

func TestCollectSemesterDates_SingleWeek(t *testing.T) {
	// 2024-09-02 (Mon) to 2024-09-06 (Fri) = 1 week, all numerator
	start := time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 9, 6, 0, 0, 0, 0, time.UTC)

	numDates, denomDates := collectSemesterDates(start, end, domain.DefaultEducationWeek)

	if len(numDates) != 5 {
		t.Errorf("expected 5 numerator dates for single week, got %d", len(numDates))
	}
	if len(denomDates) != 0 {
		t.Errorf("expected 0 denominator dates for single week, got %d", len(denomDates))
	}
}

func TestGroupDatesByWeekday(t *testing.T) {
	dates := []date{
		{year: 2024, month: time.September, day: 2},  // Monday
		{year: 2024, month: time.September, day: 3},  // Tuesday
		{year: 2024, month: time.September, day: 9},  // Monday (next week)
		{year: 2024, month: time.September, day: 10}, // Tuesday (next week)
	}
	byWeekday := groupDatesByWeekday(dates)

	monDates := byWeekday[domain.WeekdayMonday]
	if len(monDates) != 2 {
		t.Errorf("expected 2 Mondays, got %d", len(monDates))
	}
	tueDates := byWeekday[domain.WeekdayTuesday]
	if len(tueDates) != 2 {
		t.Errorf("expected 2 Tuesdays, got %d", len(tueDates))
	}
}

func TestFirstWeekDates_ReturnsEarliestPerWeekday(t *testing.T) {
	dates := []date{
		{year: 2024, month: time.September, day: 2},  // Monday week 1
		{year: 2024, month: time.September, day: 3},  // Tuesday week 1
		{year: 2024, month: time.September, day: 9},  // Monday week 2
		{year: 2024, month: time.September, day: 10}, // Tuesday week 2
	}
	educationWeek := []domain.Weekday{domain.WeekdayMonday, domain.WeekdayTuesday}
	firstWeek := firstWeekDates(dates, educationWeek)

	if len(firstWeek) != 2 {
		t.Fatalf("expected 2 first-week dates, got %d", len(firstWeek))
	}
	// First occurrence of Monday should be Sep 2
	if firstWeek[0].day != 2 {
		t.Errorf("expected first Monday = Sep 2, got day %d", firstWeek[0].day)
	}
	// First occurrence of Tuesday should be Sep 3
	if firstWeek[1].day != 3 {
		t.Errorf("expected first Tuesday = Sep 3, got day %d", firstWeek[1].day)
	}
}

func TestCountWeeks_TwoWeeks(t *testing.T) {
	start := time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC)
	weeks := countWeeks(start, end)
	if weeks != 1 {
		// 11 days / 7 = 1 week
		t.Errorf("expected 1 week, got %d", weeks)
	}
}

func TestCountWeeks_MinimumOne(t *testing.T) {
	// Same day → less than 7 days → should return 1
	start := time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC)
	weeks := countWeeks(start, end)
	if weeks < 1 {
		t.Errorf("expected at least 1 week, got %d", weeks)
	}
}

func TestCountWeeks_FourWeeks(t *testing.T) {
	start := time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 9, 29, 0, 0, 0, 0, time.UTC) // 27 days / 7 = 3 weeks
	weeks := countWeeks(start, end)
	if weeks != 3 {
		t.Errorf("expected 3 weeks for 27-day range, got %d", weeks)
	}
}
