package domain

import (
	"time"
)

type Semester struct {
	ID             uint64     `json:"id" db:"id"`
	AcademicYearID uint64     `json:"academic_year_id" db:"academic_year_id"`
	PeriodStart    time.Time  `json:"period_start" db:"period_start"`
	PeriodEnd      time.Time  `json:"period_end" db:"period_end"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt     *time.Time `json:"modified_at" db:"modified_at"`
}

type Semesters []*Semester

func (ss Semesters) SetCreatedAt(t time.Time) {
	for _, s := range ss {
		s.CreatedAt = t
	}
}

func (ss Semesters) SetModifiedAt(t time.Time) {
	for _, s := range ss {
		s.ModifiedAt = &t
	}
}
