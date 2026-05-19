package domain

import (
	"time"
)

type WorkloadDistribution struct {
	ID                  uint64     `json:"id" db:"id"`
	StudyPlanID         uint64     `json:"study_plan_id" db:"study_plan_id"`
	GroupID             uint64     `json:"group_id" db:"group_id"`
	CourseworkType      *string    `json:"coursework_type" db:"coursework_type"`
	CourseworkHours     *int       `json:"coursework_hours" db:"coursework_hours"`
	ClassroomCoursework *int       `json:"classroom_coursework" db:"classroom_coursework"`
	ClassroomWork       *int       `json:"classroom_work" db:"classroom_work"`
	Laboratory          *int       `json:"laboratory" db:"laboratory"`
	Practical           *int       `json:"practical" db:"practical"`
	Exam                *int       `json:"exam" db:"exam"`
	Consultations       *int       `json:"consultations" db:"consultations"`
	Credit              *int       `json:"credit" db:"credit"`
	Practice            *int       `json:"practice" db:"practice"`
	OKR                 *int       `json:"okr" db:"okr"`
	IndividualLessons   *int       `json:"individual_lessons" db:"individual_lessons"`
	Deductions          *int       `json:"deductions" db:"deductions"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt          *time.Time `json:"modified_at" db:"modified_at"`
	IsDeleted           bool       `json:"is_deleted" db:"is_deleted"`
}

type WorkloadDistributions []*WorkloadDistribution

func (w WorkloadDistributions) SetCreatedAt(t time.Time) {
	for _, wd := range w {
		wd.CreatedAt = t
	}
}

func (w WorkloadDistributions) SetModifiedAt(t time.Time) {
	for _, wd := range w {
		wd.ModifiedAt = &t
	}
}
