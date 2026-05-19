package domain

import (
	"time"
)

type Specialty struct {
	ID             uint64     `json:"id" db:"id"`
	AcademicYearID uint64     `json:"academic_year_id" db:"academic_year_id"`
	DepartmentID   uint64     `json:"department_id" db:"department_id"`
	Name           string     `json:"name" db:"name"`
	ShortName      string     `json:"short_name" db:"short_name"`
	SpecialtyCode  string     `json:"specialty_code" db:"specialty_code"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt     *time.Time `json:"modified_at" db:"modified_at"`
}

type Specialties []*Specialty

func (ss Specialties) SetCreatedAt(t time.Time) {
	for _, s := range ss {
		s.CreatedAt = t
	}
}

func (ss Specialties) SetModifiedAt(t time.Time) {
	for _, s := range ss {
		s.ModifiedAt = &t
	}
}

func (ss Specialties) GroupByAcademicYear() map[uint64]Specialties {
	m := make(map[uint64]Specialties)
	for _, s := range ss {
		m[s.AcademicYearID] = append(m[s.AcademicYearID], s)
	}
	return m
}
