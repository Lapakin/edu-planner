package domain

import "time"

type CycleCommitteeLabRoom struct {
	ID               uint64     `json:"id" db:"id"`
	AcademicYearID   uint64     `json:"academic_year_id" db:"academic_year_id"`
	CycleCommitteeID uint64     `json:"cycle_committee_id" db:"cycle_committee_id"`
	RoomID           uint64     `json:"room_id" db:"room_id"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt       *time.Time `json:"modified_at" db:"modified_at"`
}

type CycleCommitteeLabRooms []*CycleCommitteeLabRoom

func (c CycleCommitteeLabRooms) SetCreatedAt(t time.Time) {
	for _, r := range c {
		r.CreatedAt = t
	}
}

func (c CycleCommitteeLabRooms) SetModifiedAt(t time.Time) {
	for _, r := range c {
		r.ModifiedAt = &t
	}
}
