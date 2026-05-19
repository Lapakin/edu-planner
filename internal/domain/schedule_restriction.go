package domain

import "time"

type ScheduleRestriction struct {
	ID                      uint64     `json:"id" db:"id"`
	MinGroupLessonsPerDay   int        `json:"min_group_lessons_per_day" db:"min_group_lessons_per_day"`
	MaxGroupLessonsPerDay   int        `json:"max_group_lessons_per_day" db:"max_group_lessons_per_day"`
	MaxTeacherLessonsPerDay int        `json:"max_teacher_lessons_per_day" db:"max_teacher_lessons_per_day"`
	NoGapsInGroupSchedule   bool       `json:"no_gaps_in_group_schedule" db:"no_gaps_in_group_schedule"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt              *time.Time `json:"modified_at" db:"modified_at"`
}

type ScheduleRestrictions []*ScheduleRestriction

func (s ScheduleRestrictions) SetCreatedAt(t time.Time) {
	for _, r := range s {
		r.CreatedAt = t
	}
}

func (s ScheduleRestrictions) SetModifiedAt(t time.Time) {
	for _, r := range s {
		r.ModifiedAt = &t
	}
}

func DefaultScheduleRestriction() *ScheduleRestriction {
	return &ScheduleRestriction{
		ID:                      0,
		MinGroupLessonsPerDay:   2,
		MaxGroupLessonsPerDay:   4,
		MaxTeacherLessonsPerDay: 5,
		NoGapsInGroupSchedule:   true,
		CreatedAt:               time.Time{},
		ModifiedAt:              nil,
	}
}
