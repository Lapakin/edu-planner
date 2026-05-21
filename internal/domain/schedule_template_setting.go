package domain

import (
	"time"
)

// StudyDaysMask bitmask constants: Monday = 1, Tuesday = 2, ..., Sunday = 64.
// Default (Mon–Fri) = 31.
const (
	StudyDayMon StudyDaysMask = 1 << iota
	StudyDayTue
	StudyDayWed
	StudyDayThu
	StudyDayFri
	StudyDaySat
	StudyDaySun

	StudyDaysMaskDefault StudyDaysMask = StudyDayMon | StudyDayTue | StudyDayWed | StudyDayThu | StudyDayFri
)

type StudyDaysMask int

// ToWeekdays converts a StudyDaysMask to the ordered list of active Weekdays.
func (m StudyDaysMask) ToWeekdays() []Weekday {
	mapping := []struct {
		bit StudyDaysMask
		day Weekday
	}{
		{StudyDayMon, WeekdayMonday},
		{StudyDayTue, WeekdayTuesday},
		{StudyDayWed, WeekdayWednesday},
		{StudyDayThu, WeekdayThursday},
		{StudyDayFri, WeekdayFriday},
		{StudyDaySat, WeekdaySaturday},
		{StudyDaySun, WeekdaySunday},
	}
	days := make([]Weekday, 0, 7)
	for _, e := range mapping {
		if m&e.bit != 0 {
			days = append(days, e.day)
		}
	}
	if len(days) == 0 {
		return DefaultEducationWeek
	}
	return days
}

type ScheduleTemplateSetting struct {
	ID                         uint64        `json:"id" db:"id"`
	AcademicYearID             *uint64       `json:"academic_year_id" db:"academic_year_id"`
	LessonsPerClass            int           `json:"lessons_per_class" db:"lessons_per_class"`
	StudyDaysMask              StudyDaysMask `json:"study_days_mask" db:"study_days_mask"`
	MaxIdenticalLessonsPerDay  int           `json:"max_identical_lessons_per_day" db:"max_identical_lessons_per_day"`
	MaxStudyHoursPerDay        int           `json:"max_study_hours_per_day" db:"max_study_hours_per_day"`
	MaxTeacherHoursPerWeek     int           `json:"max_teacher_hours_per_week" db:"max_teacher_hours_per_week"`
	MaxGroupLessonHoursPerWeek *int          `json:"max_group_lesson_hours_per_week" db:"max_group_lesson_hours_per_week"`
	CreatedAt                  time.Time     `json:"created_at" db:"created_at"`
	ModifiedAt                 *time.Time    `json:"modified_at" db:"modified_at"`
}

type ScheduleTemplateSettings []*ScheduleTemplateSetting

func (s ScheduleTemplateSettings) SetCreatedAt(t time.Time) {
	for _, sts := range s {
		sts.CreatedAt = t
	}
}

func (s ScheduleTemplateSettings) SetModifiedAt(t time.Time) {
	for _, sts := range s {
		sts.ModifiedAt = &t
	}
}
