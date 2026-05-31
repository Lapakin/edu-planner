package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/user-management/repository"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type UserService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewUserService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *UserService {
	return &UserService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *UserService) CreateUsers(ctx context.Context, claims *jwt.Claims, users domain.Users) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during creating users. err: %v", err)
			s.l.Debugf("users: %v", users)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	users.SetCreatedAt(time.Now())

	if err = s.rm.NewUserRepo(tx).CreateUsers(ctx, users); err != nil {
		return handleDBError(err)
	}

	if err = s.rm.NewUserRepo(tx).CreateUserProfiles(ctx, users); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(users))
	for i, u := range users {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(u)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeUser,
			Object:     objBytes,
			Claims:     claims,
		}
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, massages); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID uint64) (*domain.User, error) {
	user, err := s.rm.NewUserRepo(s.db).GetUserByID(ctx, userID)
	if err != nil {
		s.l.Errorf("Error during getting user. err %v", err)
		s.l.Debugf("userID: %v", userID)
		return nil, handleDBError(err)
	}
	return user, nil
}

func (s *UserService) FetchUsers(ctx context.Context, filters f.Filters) (domain.Users, error) {
	users, err := s.rm.NewUserRepo(s.db).FetchUsers(ctx, filters)
	if err != nil {
		s.l.Errorf("Error during fetching users. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return users, nil
}

func (s *UserService) AttachUsers(ctx context.Context, academicYearID uint64, userIDs []uint64) error {
	if err := s.rm.NewUserRepo(s.db).AttachUsers(ctx, academicYearID, userIDs); err != nil {
		s.l.Errorf("Error during attaching users to academic year. err: %v", err)
		s.l.Debugf("academicYearID: %v, userIDs: %v", academicYearID, userIDs)
		return handleDBError(err)
	}
	return nil
}

func (s *UserService) UpdateUsers(ctx context.Context, claims *jwt.Claims, users domain.Users) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during updating users. err: %v", err)
			s.l.Debugf("users: %v", users)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	users.SetModifiedAt(time.Now())

	if err = s.rm.NewUserRepo(tx).UpdateUsers(ctx, users); err != nil {
		return handleDBError(err)
	}

	if err = s.rm.NewUserRepo(tx).UpdateUserProfiles(ctx, users); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(users))
	for i, u := range users {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(u)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeUser,
			Object:     objBytes,
			Claims:     claims,
		}
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, massages); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *UserService) DeleteUsers(ctx context.Context, claims *jwt.Claims, userIDs []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during deleting users. err: %v", err)
			s.l.Debugf("userIDs: %v", userIDs)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewUserRepo(tx).DeleteUsers(ctx, userIDs); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(userIDs))
	for i, id := range userIDs {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(id)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeDelete,
			ObjectType: domain.ObjectTypeUser,
			Object:     objBytes,
			Claims:     claims,
		}
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, massages); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *UserService) ActivateUser(ctx context.Context, claims *jwt.Claims, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during activating user. err: %v", err)
			s.l.Debugf("userID: %v", userID)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewUserRepo(tx).ActivateUser(ctx, userID, time.Now()); err != nil {
		return handleDBError(err)
	}

	var objBytes json.RawMessage
	objBytes, err = json.Marshal(userID)
	if err != nil {
		return ErrInternal
	}

	massages := domain.Massages{
		{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeActivate,
			ObjectType: domain.ObjectTypeUser,
			Object:     objBytes,
			Claims:     claims,
		},
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, massages); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *UserService) DeactivateUser(ctx context.Context, claims *jwt.Claims, userID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during deactivating user. err: %v", err)
			s.l.Debugf("userID: %v", userID)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewUserRepo(tx).DeactivateUser(ctx, userID, time.Now()); err != nil {
		return handleDBError(err)
	}

	var objBytes json.RawMessage
	objBytes, err = json.Marshal(userID)
	if err != nil {
		return ErrInternal
	}

	massages := domain.Massages{
		{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeDeactivate,
			ObjectType: domain.ObjectTypeUser,
			Object:     objBytes,
			Claims:     claims,
		},
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, massages); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}
