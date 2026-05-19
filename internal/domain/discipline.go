package domain

import (
	"time"
)

type Discipline struct {
	ID               uint64     `json:"id" db:"id"`
	AcademicYearID   uint64     `json:"academic_year_id" db:"academic_year_id"`
	CycleCommitteeID uint64     `json:"cycle_committee_id" db:"cycle_committee_id"`
	Name             string     `json:"name" db:"name"`
	ShortName        string     `json:"short_name" db:"short_name"`
	IsSplitting      bool       `json:"is_splitting" db:"is_splitting"`
	IsSubvention     bool       `json:"is_subvention" db:"is_subvention"`
	IsDeleted        bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt       *time.Time `json:"modified_at" db:"modified_at"`
}

type Disciplines []*Discipline

func (ds Disciplines) SetCreatedAt(t time.Time) {
	for _, d := range ds {
		d.CreatedAt = t
	}
}

func (ds Disciplines) SetModifiedAt(t time.Time) {
	for _, d := range ds {
		d.ModifiedAt = &t
	}
}

func (ds Disciplines) GroupByAcademicYear() map[uint64]Disciplines {
	m := make(map[uint64]Disciplines)
	for _, d := range ds {
		m[d.AcademicYearID] = append(m[d.AcademicYearID], d)
	}
	return m
}
