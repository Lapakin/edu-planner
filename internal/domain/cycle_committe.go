package domain

import (
	"time"
)

type CycleCommittee struct {
	ID             uint64     `json:"id" db:"id"`
	AcademicYearID uint64     `json:"academic_year_id" db:"academic_year_id"`
	UserID         uint64     `json:"user_id" db:"user_id"`
	Name           string     `json:"name" db:"name"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt     *time.Time `json:"modified_at" db:"modified_at"`
}

type CycleCommittees []*CycleCommittee

func (ccs CycleCommittees) SetCreatedAt(t time.Time) {
	for _, cc := range ccs {
		cc.CreatedAt = t
	}
}

func (ccs CycleCommittees) SetModifiedAt(t time.Time) {
	for _, cc := range ccs {
		cc.ModifiedAt = &t
	}
}

func (ccs CycleCommittees) GroupByAcademicYear() map[uint64]CycleCommittees {
	m := make(map[uint64]CycleCommittees)
	for _, cc := range ccs {
		m[cc.AcademicYearID] = append(m[cc.AcademicYearID], cc)
	}
	return m
}
