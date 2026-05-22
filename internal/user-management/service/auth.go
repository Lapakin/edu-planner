package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/user-management/repository"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
)

const inviteTokenTTL = 72 * time.Hour

var (
	ErrInviteTokenInvalid = errors.New("invite link is invalid or has already been used")
	ErrForbidden          = errors.New("forbidden")
)

type authSvc struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewAuthService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) AuthSvc {
	return &authSvc{
		db: db,
		rm: rm,
		l:  l,
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *authSvc) CreateInvite(ctx context.Context, claims *jwt.Claims, req domain.InviteReq) (*domain.InviteResp, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during CreateInvite. err: %v", err)
		}
	}()

	if !domain.IsPrivileged(claims.Role) {
		return nil, ErrForbidden
	}

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return nil, ErrInternal
	}
	defer tx.Rollback()

	role := req.Role
	if role == "" {
		role = domain.RoleTeacher
	}

	user := &domain.User{
		ID:         0,
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Patronymic: req.Patronymic,
		Role:       role,
		IsActive:   false,
		IsDeleted:  false,
		CreatedAt:  time.Now(),
		ModifiedAt: nil,
	}

	users := domain.Users{user}
	if err = s.rm.NewUserRepo(tx).CreateUsers(ctx, users); err != nil {
		return nil, handleDBError(err)
	}

	if err = s.rm.NewUserRepo(tx).CreateUserProfiles(ctx, users); err != nil {
		return nil, handleDBError(err)
	}

	rawToken, err := generateToken()
	if err != nil {
		return nil, ErrInternal
	}

	now := time.Now()
	inviteToken := &domain.InviteToken{
		ID:        0,
		UserID:    user.ID,
		Token:     rawToken,
		ExpiresAt: now.Add(inviteTokenTTL),
		UsedAt:    nil,
		CreatedAt: now,
	}

	if err = s.rm.NewAuthRepo(tx).CreateInviteToken(ctx, inviteToken); err != nil {
		return nil, handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, ErrInternal
	}

	return &domain.InviteResp{
		User:       user,
		InviteLink: "/set-password?token=" + rawToken,
	}, nil
}

func (s *authSvc) SetPassword(ctx context.Context, token string, password string) (string, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during SetPassword. err: %v", err)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return "", ErrInternal
	}
	defer tx.Rollback()

	inviteToken, err := s.rm.NewAuthRepo(tx).GetInviteToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrInviteTokenInvalid
			return "", err
		}
		return "", ErrInternal
	}

	if inviteToken.UsedAt != nil || time.Now().After(inviteToken.ExpiresAt) {
		err = ErrInviteTokenInvalid
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrInternal
	}

	now := time.Now()
	cred := &domain.UserCredential{
		UserID:       inviteToken.UserID,
		PasswordHash: string(hashedPassword),
		UpdatedAt:    now,
		User:         nil,
	}

	if err = s.rm.NewAuthRepo(tx).UpsertUserCredential(ctx, cred); err != nil {
		return "", handleDBError(err)
	}

	if err = s.rm.NewAuthRepo(tx).MarkInviteTokenUsed(ctx, token, now); err != nil {
		return "", handleDBError(err)
	}

	// Activate user when they set their password via invite link
	if err = s.rm.NewUserRepo(tx).ActivateUser(ctx, inviteToken.UserID, now); err != nil {
		return "", handleDBError(err)
	}

	// Fetch user for JWT claims
	user, err := s.rm.NewUserRepo(tx).GetUserByID(ctx, inviteToken.UserID)
	if err != nil {
		return "", handleDBError(err)
	}

	tokenString, err := jwt.GenerateToken(int64(user.ID), user.Email, user.Role)
	if err != nil {
		return "", ErrInternal
	}

	refreshToken := &domain.RefreshToken{
		ID:        0,
		UserID:    user.ID,
		TokenHash: tokenString,
		ExpiresAt: now.Add(24 * time.Hour),
		IsRevoked: false,
		CreatedAt: now,
	}
	if err = s.rm.NewAuthRepo(tx).CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return "", ErrInternal
	}

	return tokenString, nil
}

func (s *authSvc) ResetInvite(ctx context.Context, claims *jwt.Claims, userID uint64) (*domain.InviteResp, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during ResetInvite. err: %v", err)
		}
	}()

	if !domain.IsPrivileged(claims.Role) {
		return nil, ErrForbidden
	}

	user, err := s.rm.NewUserRepo(s.db).GetUserByID(ctx, userID)
	if err != nil {
		return nil, handleDBError(err)
	}

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return nil, ErrInternal
	}
	defer tx.Rollback()

	rawToken, err := generateToken()
	if err != nil {
		return nil, ErrInternal
	}

	now := time.Now()
	inviteToken := &domain.InviteToken{
		ID:        0,
		UserID:    userID,
		Token:     rawToken,
		ExpiresAt: now.Add(inviteTokenTTL),
		UsedAt:    nil,
		CreatedAt: now,
	}

	if err = s.rm.NewAuthRepo(tx).CreateInviteToken(ctx, inviteToken); err != nil {
		return nil, handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, ErrInternal
	}

	return &domain.InviteResp{
		User:       user,
		InviteLink: "/set-password?token=" + rawToken,
	}, nil
}

func (s *authSvc) Login(ctx context.Context, email string, password string) (string, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during login. err: %v", err)
			s.l.Debugf("email: %v", email)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return "", ErrInternal
	}
	defer tx.Rollback()

	cred, err := s.rm.NewAuthRepo(tx).GetUserCredentialByEmail(ctx, email)
	if err != nil {
		err = errors.New("invalid email or password")
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(password)); err != nil {
		err = errors.New("invalid email or password")
		return "", err
	}

	tokenString, err := jwt.GenerateToken(int64(cred.User.ID), cred.User.Email, cred.User.Role)
	if err != nil {
		return "", err
	}

	refreshToken := &domain.RefreshToken{
		ID:        0,
		UserID:    cred.User.ID,
		TokenHash: tokenString,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsRevoked: false,
		CreatedAt: time.Time{},
	}
	if err = s.rm.NewAuthRepo(tx).CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return "", ErrInternal
	}

	return tokenString, nil
}
