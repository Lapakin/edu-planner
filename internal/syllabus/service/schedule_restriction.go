package service

import (
	"context"
	"time"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository"
	"github.com/google/uuid"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type ScheduleRestrictionService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewScheduleRestrictionService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *ScheduleRestrictionService {
	return &ScheduleRestrictionService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *ScheduleRestrictionService) CreateScheduleRestrictions(ctx context.Context, claims *jwt.Claims, restrictions domain.ScheduleRestrictions) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating schedule restrictions. err: %v", err)
			s.l.Debugf("claims: %v, restrictions: %v", claims, restrictions)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	restrictions.SetCreatedAt(time.Now())

	if err = s.rm.NewScheduleRestrictionRepo(tx).CreateScheduleRestrictions(ctx, restrictions); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(restrictions))
	for i, r := range restrictions {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(r); err != nil {
			return ErrInternal
		}
		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeScheduleRestriction,
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

func (s *ScheduleRestrictionService) GetScheduleRestrictionByID(ctx context.Context, id uint64) (*domain.ScheduleRestriction, error) {
	restriction, err := s.rm.NewScheduleRestrictionRepo(s.db).GetScheduleRestrictionByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting schedule restriction by id. err: %v", err)
		return nil, handleDBError(err)
	}
	return restriction, nil
}

func (s *ScheduleRestrictionService) FetchScheduleRestrictions(ctx context.Context, filters f.Filters) (domain.ScheduleRestrictions, error) {
	restrictions, err := s.rm.NewScheduleRestrictionRepo(s.db).FetchScheduleRestrictions(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching schedule restrictions. err: %v", err)
		return nil, ErrInternal
	}
	return restrictions, nil
}

func (s *ScheduleRestrictionService) UpdateScheduleRestrictions(ctx context.Context, claims *jwt.Claims, restrictions domain.ScheduleRestrictions) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating schedule restrictions. err: %v", err)
			s.l.Debugf("claims: %v, restrictions: %v", claims, restrictions)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	restrictions.SetModifiedAt(time.Now())

	if err = s.rm.NewScheduleRestrictionRepo(tx).UpdateScheduleRestrictions(ctx, restrictions); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(restrictions))
	for i, r := range restrictions {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(r); err != nil {
			return ErrInternal
		}
		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeScheduleRestriction,
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

func (s *ScheduleRestrictionService) DeleteScheduleRestrictions(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting schedule restrictions. err: %v", err)
			s.l.Debugf("claims: %v, restrictionIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewScheduleRestrictionRepo(tx).DeleteScheduleRestrictions(ctx, ids); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(ids))
	for i, id := range ids {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(id); err != nil {
			return ErrInternal
		}
		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeDelete,
			ObjectType: domain.ObjectTypeScheduleRestriction,
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
