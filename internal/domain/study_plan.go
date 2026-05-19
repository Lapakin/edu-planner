package domain

import (
	"time"
)

type StudyPlan struct {
	ID                   uint64     `json:"id" db:"id"`
	AcademicYearID       *uint64    `json:"academic_year_id" db:"academic_year_id"`
	SpecialtyID          *uint64    `json:"specialty_id" db:"specialty_id"`
	DisciplineID         *uint64    `json:"discipline_id" db:"discipline_id"`
	SemesterNumber       *int       `json:"semester_number" db:"semester_number"`
	EducationalComponent *string    `json:"educational_component" db:"educational_component"`
	OKR                  *int       `json:"okr" db:"okr"`
	CourseworkType       *string    `json:"coursework_type" db:"coursework_type"`
	ClassroomWork        *int       `json:"classroom_work" db:"classroom_work"`
	Lectures             *int       `json:"lectures" db:"lectures"`
	Laboratory           *int       `json:"laboratory" db:"laboratory"`
	Practical            *int       `json:"practical" db:"practical"`
	IndependentWork      *int       `json:"independent_work" db:"independent_work"`
	IndividualLessons    *int       `json:"individual_lessons" db:"individual_lessons"`
	Consultations        *int       `json:"consultations" db:"consultations"`
	Exam                 *int       `json:"exam" db:"exam"`
	Credit               *int       `json:"credit" db:"credit"`
	Practice             *int       `json:"practice" db:"practice"`
	CourseworkHours      *int       `json:"coursework_hours" db:"coursework_hours"`
	ClassroomCoursework  *int       `json:"classroom_coursework" db:"classroom_coursework"`
	AssistingHours       *int       `json:"assisting_hours" db:"assisting_hours"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt           *time.Time `json:"modified_at" db:"modified_at"`
	IsDeleted            bool       `json:"is_deleted" db:"is_deleted"`
}

type StudyPlans []*StudyPlan

func (sp StudyPlans) SetCreatedAt(t time.Time) {
	for _, p := range sp {
		p.CreatedAt = t
	}
}

func (sp StudyPlans) SetModifiedAt(t time.Time) {
	for _, p := range sp {
		p.ModifiedAt = &t
	}
}
