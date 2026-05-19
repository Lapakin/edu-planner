package domain

import (
	"time"
)

type Group struct {
	ID             uint64     `json:"id" db:"id"`
	AcademicYearID uint64     `json:"academic_year_id" db:"academic_year_id"`
	SpecialtyID    uint64     `json:"specialty_id" db:"specialty_id"`
	Name           string     `json:"name" db:"name"`
	ShortName      string     `json:"short_name" db:"short_name"`
	IsContract     bool       `json:"is_contract" db:"is_contract"`
	IsSplitting    bool       `json:"is_splitting" db:"is_splitting"`
	EducationStart time.Time  `json:"education_start" db:"education_start"`
	EducationEnd   time.Time  `json:"education_end" db:"education_end"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt     *time.Time `json:"modified_at" db:"modified_at"`
}

type Groups []*Group

func (gs Groups) SetCreatedAt(createdAt time.Time) {
	for _, g := range gs {
		g.CreatedAt = createdAt
	}
}

func (gs Groups) SetModifiedAt(modifiedAt time.Time) {
	for _, g := range gs {
		g.ModifiedAt = &modifiedAt
	}
}

func (gs Groups) GroupByAcademicYearID() map[uint64]Groups {
	m := make(map[uint64]Groups)
	for _, g := range gs {
		m[g.AcademicYearID] = append(m[g.AcademicYearID], g)
	}
	return m
}

type GroupSemester struct {
	ID         uint64     `json:"id" db:"id"`
	GroupID    uint64     `json:"group_id" db:"group_id"`
	SemesterID uint64     `json:"semester_id" db:"semester_id"`
	StartDate  time.Time  `json:"start_date" db:"start_date"`
	EndDate    time.Time  `json:"end_date" db:"end_date"`
	IsDeleted  bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
}

type GroupSemesters []*GroupSemester

func (gs GroupSemesters) SetCreatedAt(createdAt time.Time) {
	for _, g := range gs {
		g.CreatedAt = createdAt
	}
}

func (gs GroupSemesters) SetModifiedAt(modifiedAt time.Time) {
	for _, g := range gs {
		g.ModifiedAt = &modifiedAt
	}
}
