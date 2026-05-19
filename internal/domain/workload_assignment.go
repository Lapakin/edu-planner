package domain

import (
	"time"
)

type WorkloadAssignment struct {
	ID                     uint64    `json:"id" db:"id"`
	WorkloadDistributionID uint64    `json:"workload_distribution_id" db:"workload_distribution_id"`
	TeacherID              uint64    `json:"teacher_id" db:"teacher_id"`
	RoleType               *string   `json:"role_type" db:"role_type"`
	AssignedHours          *int      `json:"assigned_hours" db:"assigned_hours"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
}

type WorkloadAssignments []*WorkloadAssignment

func (w WorkloadAssignments) SetCreatedAt(t time.Time) {
	for _, wa := range w {
		wa.CreatedAt = t
	}
}
