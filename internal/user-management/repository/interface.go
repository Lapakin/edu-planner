package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type RepoManager interface {
	NewUserRepo(db sqlx.ExtContext) UserRepository
	NewMessageRepo(db sqlx.ExtContext) MessageRepository
	NewAuthRepo(db sqlx.ExtContext) AuthRepository
	NewAcademicYearRepo(db sqlx.ExtContext) AcademicYearRepository
}

type UserRepository interface {
	CreateUsers(ctx context.Context, users domain.Users) error
	CreateUserProfiles(ctx context.Context, users domain.Users) error
	GetUserByID(ctx context.Context, userID uint64) (*domain.User, error)
	FetchUsers(ctx context.Context, filters f.Filters) (domain.Users, error)
	FetchActiveUserIDs(ctx context.Context) ([]uint64, error)
	AttachUsers(ctx context.Context, academicYearID uint64, userIDs []uint64) error
	UpdateUsers(ctx context.Context, users domain.Users) error
	UpdateUserProfiles(ctx context.Context, users domain.Users) error
	DeleteUsers(ctx context.Context, userIDs []uint64) error
	ActivateUser(ctx context.Context, userID uint64, currentTime time.Time) error
	DeactivateUser(ctx context.Context, userID uint64, currentTime time.Time) error
}

type MessageRepository interface {
	Write(ctx context.Context, messages domain.Massages) error
}

type AuthRepository interface {
	GetUserCredentialByEmail(ctx context.Context, email string) (*domain.UserCredential, error)
	CreateUserCredential(ctx context.Context, cred *domain.UserCredential) error
	UpsertUserCredential(ctx context.Context, cred *domain.UserCredential) error
	CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	CreateInviteToken(ctx context.Context, token *domain.InviteToken) error
	GetInviteToken(ctx context.Context, token string) (*domain.InviteToken, error)
	MarkInviteTokenUsed(ctx context.Context, token string, usedAt time.Time) error
	CreateResetPasswordToken(ctx context.Context, token *domain.ResetPasswordToken) error
	GetResetPasswordToken(ctx context.Context, token string) (*domain.ResetPasswordToken, error)
	MarkResetPasswordTokenUsed(ctx context.Context, token string, usedAt time.Time) error
}

type AcademicYearRepository interface {
	UpsertAcademicYear(ctx context.Context, academicYear *domain.AcademicYear) error
	DeleteAcademicYear(ctx context.Context, id uint64) error
}
