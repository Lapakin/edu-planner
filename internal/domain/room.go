package domain

import (
	"time"
)

type RoomType string

const (
	RoomTypeAuditorium RoomType = "auditorium"
	RoomTypeLaboratory RoomType = "laboratory"
)

type Room struct {
	ID             uint64     `json:"id" db:"id"`
	AcademicYearID uint64     `json:"academic_year_id" db:"academic_year_id"`
	Name           string     `json:"name" db:"name"`
	RoomType       RoomType   `json:"room_type" db:"room_type"`
	Capacity       int        `json:"capacity" db:"capacity"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt     *time.Time `json:"modified_at" db:"modified_at"`
}

type Rooms []*Room

func (rs Rooms) SetCreatedAt(t time.Time) {
	for _, r := range rs {
		r.CreatedAt = t
	}
}

func (rs Rooms) SetModifiedAt(t time.Time) {
	for _, r := range rs {
		r.ModifiedAt = &t
	}
}

func (rs Rooms) GroupByAcademicYear() map[uint64]Rooms {
	m := make(map[uint64]Rooms)
	for _, r := range rs {
		m[r.AcademicYearID] = append(m[r.AcademicYearID], r)
	}
	return m
}
