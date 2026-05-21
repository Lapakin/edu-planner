package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
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
		if err = s.upsertUser(ctx, massage.Claims, user); err != nil {
			return ErrInternal
		}

	case domain.ActionTypeDelete:
		var id uint64
		if err = json.Unmarshal(massage.Object, &id); err != nil {
			return ErrUnmarshal
		}
		if err = s.deleteUser(ctx, massage.Claims, id); err != nil {
			return ErrInternal
		}
	}

	return nil
}

func (s *UserService) upsertUser(ctx context.Context, claims *jwt.Claims, user *domain.User) error {
	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return fmt.Errorf("error during starting transaction. err: %w", err)
	}
	defer tx.Rollback()

	if err = s.rm.NewUserRepo(tx).UpsertUser(ctx, user); err != nil {
		return fmt.Errorf("error during upserting user info. err: %w", err)
	}

	if err = s.syncRelatedEntities(ctx, tx, claims, user); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error during committing user info transaction. err: %w", err)
	}

	return nil
}

func (s *UserService) deleteUser(ctx context.Context, claims *jwt.Claims, id uint64) error {
	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return fmt.Errorf("error during starting transaction. err: %w", err)
	}
	defer tx.Rollback()

	if err = s.rm.NewUserRepo(tx).DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("error during deleting user info. err: %w", err)
	}

	if err = s.deleteRelatedEntities(ctx, tx, claims, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error during committing transaction. err: %w", err)
	}

	return nil
}

func (s *UserService) syncRelatedEntities(ctx context.Context, tx *sqlx.Tx, claims *jwt.Claims, user *domain.User) error {
	switch user.Role {
	case domain.RoleTeacher:
		_, err := s.rm.NewTeacherRepo(tx).GetTeacherByUserID(ctx, user.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("error during checking teacher existence. err: %w", err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			if err := s.createTeacherEntity(ctx, tx, claims, user); err != nil {
				return fmt.Errorf("error during creating teacher entity. err: %w", err)
			}
		}
	default:
		if err := s.deleteRelatedEntities(ctx, tx, claims, user.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) deleteRelatedEntities(ctx context.Context, tx *sqlx.Tx, claims *jwt.Claims, userID uint64) error {
	teacher, err := s.rm.NewTeacherRepo(tx).GetTeacherByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("error during getting teacher by user id. err: %w", err)
	}

	if err = s.deleteTeacherEntity(ctx, tx, claims, teacher); err != nil {
		return fmt.Errorf("error during deleting related entities. err: %w", err)
	}

	return nil
}

func (s *UserService) createTeacherEntity(ctx context.Context, tx *sqlx.Tx, claims *jwt.Claims, user *domain.User) error {
	academicYear, err := s.rm.NewAcademicYearRepo(tx).GetActiveAcademicYear(ctx)
	if err != nil {
		return fmt.Errorf("error during getting active academic year. err: %w", err)
	}

	teacher := &domain.Teacher{
		ID:             0,
		UserID:         user.ID,
		AcademicYearID: academicYear.ID,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		ModifiedAt:     nil,
	}

	r := s.rm.NewTeacherRepo(tx)
	if err = r.CreateTeachers(ctx, domain.Teachers{teacher}); err != nil {
		return fmt.Errorf("error during creating teacher entity. err: %w", err)
	}

	toAttach := map[uint64]domain.Teachers{academicYear.ID: {teacher}}
	if err = r.AttachTeachersToAcademicYear(ctx, toAttach); err != nil {
		return fmt.Errorf("error during attaching teacher to academic year. err: %w", err)
	}

	objBytes, err := json.Marshal(teacher)
	if err != nil {
		return ErrInternal
	}

	massage := &domain.Massage{
		EventID:    uuid.New().String(),
		ActionType: domain.ActionTypeCreate,
		ObjectType: domain.ObjectTypeTeacher,
		Object:     objBytes,
		Claims:     claims,
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, domain.Massages{massage}); err != nil {
		return fmt.Errorf("error during writing massage. err: %w", err)
	}

	return nil
}

func (s *UserService) deleteTeacherEntity(ctx context.Context, tx *sqlx.Tx, claims *jwt.Claims, teacher *domain.Teacher) error {
	r := s.rm.NewTeacherRepo(tx)
	if err := r.DeleteTeachers(ctx, []uint64{teacher.ID}); err != nil {
		return fmt.Errorf("error during deleting teacher entity: %w", err)
	}

	objBytes, err := json.Marshal(teacher.ID)
	if err != nil {
		return ErrInternal
	}

	massage := &domain.Massage{
		EventID:    uuid.New().String(),
		ActionType: domain.ActionTypeDelete,
		ObjectType: domain.ObjectTypeTeacher,
		Object:     objBytes,
		Claims:     claims,
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, domain.Massages{massage}); err != nil {
		return fmt.Errorf("error during writing delete massage. err: %w", err)
	}

	return nil
}
