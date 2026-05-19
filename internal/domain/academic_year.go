package domain

import "time"

type AcademicYear struct {
	ID         uint64     `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	StartDate  time.Time  `json:"start_date" db:"start_date"`
	EndDate    time.Time  `json:"end_date" db:"end_date"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	IsDeleted  bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
}

func (a *AcademicYear) Activate(t time.Time) {
	a.IsActive = true
	a.ModifiedAt = &t
}

func (a *AcademicYear) Deactivate(t time.Time) {
	a.IsActive = false
	a.ModifiedAt = &t
}

type AcademicYears []*AcademicYear

func (ay AcademicYears) SetCreatedAt(t time.Time) {
	for _, a := range ay {
		a.CreatedAt = t
	}
}

func (ay AcademicYears) SetModifiedAt(t time.Time) {
	for _, a := range ay {
		a.ModifiedAt = &t
	}
}
