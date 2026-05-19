package domain

import (
	"time"
)

type Department struct {
	ID             uint64     `json:"id" db:"id"`
	AcademicYearID uint64     `json:"academic_year_id" db:"academic_year_id"`
	Name           string     `json:"name" db:"name"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt     *time.Time `json:"modified_at" db:"modified_at"`
}

type Departments []*Department

func (ds Departments) SetCreatedAt(t time.Time) {
	for _, d := range ds {
		d.CreatedAt = t
	}
}

func (ds Departments) SetModifiedAt(t time.Time) {
	for _, d := range ds {
		d.ModifiedAt = &t
	}
}

func (ds Departments) GroupByAcademicYear() map[uint64]Departments {
	m := make(map[uint64]Departments)
	for _, d := range ds {
		m[d.AcademicYearID] = append(m[d.AcademicYearID], d)
	}
	return m
}
