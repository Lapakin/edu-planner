package domain

import (
	"time"
)

type ScheduleTemplateSetting struct {
	ID                         uint64     `json:"id" db:"id"`
	HoursPerLesson             float64    `json:"hours_per_lesson" db:"hours_per_lesson"`
	MaxIdenticalLessonsPerDay  int        `json:"max_identical_lessons_per_day" db:"max_identical_lessons_per_day"`
	MaxStudyHoursPerDay        int        `json:"max_study_hours_per_day" db:"max_study_hours_per_day"`
	MaxTeacherHoursPerWeek     int        `json:"max_teacher_hours_per_week" db:"max_teacher_hours_per_week"`
	MaxGroupLessonHoursPerWeek *int       `json:"max_group_lesson_hours_per_week" db:"max_group_lesson_hours_per_week"`
	CreatedAt                  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt                 *time.Time `json:"modified_at" db:"modified_at"`
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
