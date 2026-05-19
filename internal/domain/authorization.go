package domain

import "time"

type UserCredential struct {
	UserID       uint64    `db:"user_id"`
	PasswordHash string    `db:"password_hash"`
	UpdatedAt    time.Time `db:"updated_at"`
	User         *User     `db:"-"`
}

type RefreshToken struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	IsRevoked bool      `db:"is_revoked"`
	CreatedAt time.Time `db:"created_at"`
}

type InviteToken struct {
	ID        uint64     `db:"id"`
	UserID    uint64     `db:"user_id"`
	Token     string     `db:"token"`
	ExpiresAt time.Time  `db:"expires_at"`
	UsedAt    *time.Time `db:"used_at"`
	CreatedAt time.Time  `db:"created_at"`
}

type InviteReq struct {
	Email      string  `json:"email"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Patronymic *string `json:"patronymic,omitempty"`
	Role       string  `json:"role"`
}

type InviteResp struct {
	User       *User  `json:"user"`
	InviteLink string `json:"invite_link"`
}

type SetPasswordReq struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
