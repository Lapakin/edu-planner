package domain

import "time"

// TimePriority controls the preferred half-of-day for lesson placement.
type TimePriority string

const (
	TimePriorityNone      TimePriority = "none"
	TimePriorityMorning   TimePriority = "morning"
	TimePriorityAfternoon TimePriority = "afternoon"
)

type ScheduleRestriction struct {
	ID                           uint64       `json:"id" db:"id"`
	AcademicYearID               *uint64      `json:"academic_year_id" db:"academic_year_id"`
	MinGroupLessonsPerDay        int          `json:"min_group_lessons_per_day" db:"min_group_lessons_per_day"`
	MaxGroupLessonsPerDay        int          `json:"max_group_lessons_per_day" db:"max_group_lessons_per_day"`
	MaxTeacherLessonsPerDay      int          `json:"max_teacher_lessons_per_day" db:"max_teacher_lessons_per_day"`
	MaxConsecutiveTeacherLessons int          `json:"max_consecutive_teacher_lessons" db:"max_consecutive_teacher_lessons"`
	TimePriority                 TimePriority `json:"time_priority" db:"time_priority"`
	AllowFlowLessons             bool         `json:"allow_flow_lessons" db:"allow_flow_lessons"`
	NoGapsInGroupSchedule        bool         `json:"no_gaps_in_group_schedule" db:"no_gaps_in_group_schedule"`
	CreatedAt                    time.Time    `json:"created_at" db:"created_at"`
	ModifiedAt                   *time.Time   `json:"modified_at" db:"modified_at"`
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
		ID:                           0,
		AcademicYearID:               nil,
		MinGroupLessonsPerDay:        2,
		MaxGroupLessonsPerDay:        4,
		MaxTeacherLessonsPerDay:      5,
		MaxConsecutiveTeacherLessons: 4,
		TimePriority:                 TimePriorityNone,
		AllowFlowLessons:             true,
		NoGapsInGroupSchedule:        true,
		CreatedAt:                    time.Time{},
		ModifiedAt:                   nil,
	}
}
