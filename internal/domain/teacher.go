package domain

import (
	"time"
)

type Teacher struct {
	ID             uint64     `json:"id" db:"id"`
	AcademicYearID uint64     `json:"academic_year_id" db:"academic_year_id"`
	UserID         uint64     `json:"user_id" db:"user_id"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt     *time.Time `json:"modified_at" db:"modified_at"`
}

type Teachers []*Teacher

func (ts Teachers) SetCreatedAt(t time.Time) {
	for _, te := range ts {
		te.CreatedAt = t
	}
}

func (ts Teachers) SetModifiedAt(t time.Time) {
	for _, te := range ts {
		te.ModifiedAt = &t
	}
}

func (ts Teachers) GroupByAcademicYear() map[uint64]Teachers {
	m := make(map[uint64]Teachers)
	for _, te := range ts {
		m[te.AcademicYearID] = append(m[te.AcademicYearID], te)
	}
	return m
}
