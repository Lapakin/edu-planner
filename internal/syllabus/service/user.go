package service

import (
	"context"
	"fmt"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
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

func (s *UserService) MassageConsumer(ctx context.Context, massage *domain.Massage) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during processing massage. err: %v", err)
			s.l.Debugf("massage: %v", massage)
		}
	}()

	switch massage.ActionType {
	case domain.ActionTypeCreate, domain.ActionTypeUpdate:
		var user *domain.User
		if err = json.Unmarshal(massage.Object, &user); err != nil {
			return ErrUnmarshal
		}
		if err = s.upsertUser(ctx, user); err != nil {
			return ErrInternal
		}

	case domain.ActionTypeDelete:
		var id uint64
		if err = json.Unmarshal(massage.Object, &id); err != nil {
			return ErrUnmarshal
		}
		if err = s.deleteUser(ctx, id); err != nil {
			return ErrInternal
		}
	}

	return nil
}

func (s *UserService) upsertUser(ctx context.Context, user *domain.User) error {
	if err := s.rm.NewUserRepo(s.db).UpsertUser(ctx, user); err != nil {
		return fmt.Errorf("error during upserting user info. err: %w", err)
	}
	return nil
}

func (s *UserService) deleteUser(ctx context.Context, id uint64) error {
	if err := s.rm.NewUserRepo(s.db).DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("error during deleting user info. err: %w", err)
	}
	return nil
}
