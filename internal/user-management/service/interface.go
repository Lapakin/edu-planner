package service

import (
	"context"

	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type Services struct {
	UserSvc         UserSvc
	AuthSvc         AuthSvc
	AcademicYearSvc AcademicYearSvc
}

type UserSvc interface {
	CreateUsers(ctx context.Context, claims *jwt.Claims, users domain.Users) error
	UpdateUsers(ctx context.Context, claims *jwt.Claims, users domain.Users) error
	DeleteUsers(ctx context.Context, claims *jwt.Claims, userIDs []uint64) error
	ActivateUser(ctx context.Context, claims *jwt.Claims, userID uint64) error
	DeactivateUser(ctx context.Context, claims *jwt.Claims, userID uint64) error
	GetUserByID(ctx context.Context, userID uint64) (*domain.User, error)
	FetchUsers(ctx context.Context, filters f.Filters) (domain.Users, error)
}

type AuthSvc interface {
	CreateInvite(ctx context.Context, claims *jwt.Claims, req domain.InviteReq) (*domain.InviteResp, error)
	SetPassword(ctx context.Context, token string, password string) (string, error)
	ResetInvite(ctx context.Context, claims *jwt.Claims, userID uint64) (*domain.InviteResp, error)
	Login(ctx context.Context, email string, password string) (string, error)
}

type AcademicYearSvc interface {
	MassageConsumer(ctx context.Context, massage *domain.Massage) error
}
