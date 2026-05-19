package service

import (
	"context"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository"
	"github.com/google/uuid"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type BellScheduleService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewBellScheduleService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *BellScheduleService {
	return &BellScheduleService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *BellScheduleService) CreateBellSchedules(ctx context.Context, claims *jwt.Claims, schedules domain.BellSchedules) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating bell schedules. err: %v", err)
			s.l.Debugf("claims: %v, schedules: %v", claims, schedules)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewBellScheduleRepo(tx).CreateBellSchedules(ctx, schedules); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(schedules))
	for i, sc := range schedules {
		var objBytes []byte
		objBytes, err = json.Marshal(sc)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeBellSchedule,
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

func (s *BellScheduleService) GetBellScheduleByID(ctx context.Context, id uint64) (*domain.BellSchedule, error) {
	schedule, err := s.rm.NewBellScheduleRepo(s.db).GetBellScheduleByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting bell schedule by id. err: %v", err)
		s.l.Debugf("scheduleID: %v", id)
		return nil, handleDBError(err)
	}
	return schedule, nil
}

func (s *BellScheduleService) FetchBellSchedules(ctx context.Context, filters f.Filters) (domain.BellSchedules, error) {
	schedules, err := s.rm.NewBellScheduleRepo(s.db).FetchBellSchedules(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching bell schedules. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return schedules, nil
}

func (s *BellScheduleService) UpdateBellSchedules(ctx context.Context, claims *jwt.Claims, schedules domain.BellSchedules) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating bell schedules. err: %v", err)
			s.l.Debugf("claims: %v, schedules: %v", claims, schedules)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewBellScheduleRepo(tx).UpdateBellSchedules(ctx, schedules); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(schedules))
	for i, sc := range schedules {
		var objBytes []byte
		objBytes, err = json.Marshal(sc)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeBellSchedule,
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

func (s *BellScheduleService) DeleteBellSchedules(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting bell schedules. err: %v", err)
			s.l.Debugf("claims: %v, scheduleIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewBellScheduleRepo(tx).DeleteBellSchedules(ctx, ids); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(ids))
	for i, id := range ids {
		var objBytes []byte
		objBytes, err = json.Marshal(id)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeDelete,
			ObjectType: domain.ObjectTypeBellSchedule,
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
