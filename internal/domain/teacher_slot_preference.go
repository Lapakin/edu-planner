package domain

import "time"

type SlotType string

const (
	SlotTypePreferred SlotType = "preferred"
	SlotTypeForbidden SlotType = "forbidden"
)

type TeacherSlotPreference struct {
	ID             uint64     `json:"id" db:"id"`
	AcademicYearID uint64     `json:"academic_year_id" db:"academic_year_id"`
	TeacherID      uint64     `json:"teacher_id" db:"teacher_id"`
	Weekday        Weekday    `json:"weekday" db:"weekday"`
	LessonNumber   int        `json:"lesson_number" db:"lesson_number"`
	SlotType       SlotType   `json:"slot_type" db:"slot_type"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt     *time.Time `json:"modified_at" db:"modified_at"`
}

type TeacherSlotPreferences []*TeacherSlotPreference

func (t TeacherSlotPreferences) SetCreatedAt(ts time.Time) {
	for _, p := range t {
		p.CreatedAt = ts
	}
}

func (t TeacherSlotPreferences) SetModifiedAt(ts time.Time) {
	for _, p := range t {
		p.ModifiedAt = &ts
	}
}
