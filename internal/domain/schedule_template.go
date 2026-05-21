package domain

import (
	"time"
)

type Weekday string

const (
	WeekdayMonday    Weekday = "monday"
	WeekdayTuesday   Weekday = "tuesday"
	WeekdayWednesday Weekday = "wednesday"
	WeekdayThursday  Weekday = "thursday"
	WeekdayFriday    Weekday = "friday"
	WeekdaySaturday  Weekday = "saturday"
	WeekdaySunday    Weekday = "sunday"
)

var DefaultEducationWeek = []Weekday{
	WeekdayMonday,
	WeekdayTuesday,
	WeekdayWednesday,
	WeekdayThursday,
	WeekdayFriday,
}

type LessonType string

const (
	LessonTypeUnited LessonType = "united"
	LessonTypeSplit  LessonType = "split"
	LessonTypeFlow   LessonType = "flow"
)

type ScheduleTemplate struct {
	ID         uint64     `json:"id" db:"id"`
	SemesterID *uint64    `json:"semester_id" db:"semester_id"`
	Name       *string    `json:"name" db:"name"`
	Data       []byte     `json:"data" db:"data"`
	IsActive   *bool      `json:"is_active" db:"is_active"`
	IsDeleted  *bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
}

type ScheduleTemplates []*ScheduleTemplate

func (st ScheduleTemplates) SetCreatedAt(t time.Time) {
	for _, s := range st {
		s.CreatedAt = t
	}
}

func (st ScheduleTemplates) SetModifiedAt(t time.Time) {
	for _, s := range st {
		s.ModifiedAt = &t
	}
}

type ScheduleData struct {
	Numerator   WeekSchedule `json:"numerator"`
	Denominator WeekSchedule `json:"denominator"`
	Metadata    string       `json:"metadata,omitempty"`
}

type WeekSchedule map[Weekday]*DaySchedule

type DaySchedule map[int][]*ScheduleLesson

type ScheduleLesson struct {
	Type       LessonType           `json:"type"`
	SubjectID  uint64               `json:"subject_id"`
	SubLessons []*ScheduleSubLesson `json:"sub_lessons"`
}

type ScheduleSubLesson struct {
	SubGroupNumber *uint64 `json:"sub_group_number,omitempty"`
	GroupID        uint64  `json:"group_id"`
	RoomID         uint64  `json:"room_id"`
	TeacherID      uint64  `json:"teacher_id"`
}
