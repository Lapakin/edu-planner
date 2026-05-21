package generator

import (
	"time"

	"github.com/Lapakin/edu-planner/internal/domain"
)

// date is a lightweight date representation (no time component).
type date struct {
	year  int
	month time.Month
	day   int
}

func newDate(t time.Time) date {
	y, m, d := t.Date()
	return date{year: y, month: m, day: d}
}

func (d date) toTime() time.Time {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, time.UTC)
}

func (d date) weekday() domain.Weekday {
	return timeWeekdayToWeekday(d.toTime().Weekday())
}

func (d date) addDays(n int) date {
	return newDate(d.toTime().AddDate(0, 0, n))
}

func (d date) before(other date) bool {
	return d.toTime().Before(other.toTime())
}

func (d date) equal(other date) bool {
	return d.year == other.year && d.month == other.month && d.day == other.day
}

// timeWeekdayToWeekday converts Go's time.Weekday to domain.Weekday.
func timeWeekdayToWeekday(wd time.Weekday) domain.Weekday {
	switch wd {
	case time.Monday:
		return domain.WeekdayMonday
	case time.Tuesday:
		return domain.WeekdayTuesday
	case time.Wednesday:
		return domain.WeekdayWednesday
	case time.Thursday:
		return domain.WeekdayThursday
	case time.Friday:
		return domain.WeekdayFriday
	case time.Saturday:
		return domain.WeekdaySaturday
	case time.Sunday:
		return domain.WeekdayMonday
	}
	return domain.WeekdayMonday
}

// weekdayToTimeWeekday converts domain.Weekday to time.Weekday.
func weekdayToTimeWeekday(wd domain.Weekday) time.Weekday {
	switch wd {
	case domain.WeekdayMonday:
		return time.Monday
	case domain.WeekdayTuesday:
		return time.Tuesday
	case domain.WeekdayWednesday:
		return time.Wednesday
	case domain.WeekdayThursday:
		return time.Thursday
	case domain.WeekdayFriday:
		return time.Friday
	case domain.WeekdaySaturday:
		return time.Saturday
	case domain.WeekdaySunday:
		return time.Sunday
	default:
		return time.Monday
	}
}

// collectSemesterDates gathers all education dates from a semester period,
// split into numerator and denominator weeks.
func collectSemesterDates(start, end time.Time, educationWeek []domain.Weekday) (numeratorDates, denominatorDates []date) {
	educationDays := make(map[time.Weekday]bool, len(educationWeek))
	for _, wd := range educationWeek {
		educationDays[weekdayToTimeWeekday(wd)] = true
	}

	startDate := newDate(start)
	endDate := newDate(end)

	weekNumber := 0
	currentWeekStart := startDate

	// Align to Monday of the first week
	for currentWeekStart.toTime().Weekday() != time.Monday {
		currentWeekStart = currentWeekStart.addDays(-1)
	}

	current := startDate
	for !current.before(startDate) && (current.before(endDate) || current.equal(endDate)) {
		wd := current.toTime().Weekday()
		if educationDays[wd] {
			// Determine which week we're in (0-based from currentWeekStart)
			daysSinceStart := int(current.toTime().Sub(currentWeekStart.toTime()).Hours() / 24)
			weekNumber = daysSinceStart / 7

			if weekNumber%2 == 0 {
				numeratorDates = append(numeratorDates, current)
			} else {
				denominatorDates = append(denominatorDates, current)
			}
		}
		current = current.addDays(1)
	}

	return numeratorDates, denominatorDates
}

// groupDatesByWeekday groups dates by their weekday.
func groupDatesByWeekday(dates []date) map[domain.Weekday][]date {
	result := make(map[domain.Weekday][]date)
	for _, d := range dates {
		wd := d.weekday()
		result[wd] = append(result[wd], d)
	}
	return result
}

// firstWeekDates returns the earliest occurrence of each education weekday from dates.
// This is used to build a one-week template: all lessons are placed on these dates
// so that buildWeekSchedule (which reads datesForDay[0]) sees them correctly.
func firstWeekDates(dates []date, educationWeek []domain.Weekday) []date {
	byWeekday := groupDatesByWeekday(dates)
	result := make([]date, 0, len(educationWeek))
	for _, wd := range educationWeek {
		if days, ok := byWeekday[wd]; ok && len(days) > 0 {
			result = append(result, days[0])
		}
	}
	return result
}

// countWeeks calculates the number of full education weeks in a period.
func countWeeks(start, end time.Time) int {
	days := int(end.Sub(start).Hours() / 24)
	weeks := max(days/7, 1)
	return weeks
}
