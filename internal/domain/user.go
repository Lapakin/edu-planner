package domain

import (
	"time"
)

const (
	RoleAdmin   = "admin"
	RoleDean    = "dean"
	RoleTeacher = "teacher"
)

// IsPrivileged returns true for roles that have administrative access.
func IsPrivileged(role string) bool {
	return role == RoleAdmin || role == RoleDean
}

type User struct {
	ID         uint64     `json:"id" db:"id"`
	Email      string     `json:"email" db:"email"`
	Role       string     `json:"role" db:"role"`
	FirstName  string     `json:"first_name" db:"first_name"`
	LastName   string     `json:"last_name" db:"last_name"`
	Patronymic *string    `json:"patronymic" db:"patronymic"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	IsDeleted  bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
}

func (u *User) Activate(t time.Time) {
	u.IsActive = true
	u.ModifiedAt = &t
}

func (u *User) Deactivate(t time.Time) {
	u.IsActive = false
	u.ModifiedAt = &t
}

type Users []*User

func (us Users) SetCreatedAt(t time.Time) {
	for _, u := range us {
		u.CreatedAt = t
	}
}

func (us Users) SetModifiedAt(t time.Time) {
	for _, u := range us {
		u.ModifiedAt = &t
	}
}
