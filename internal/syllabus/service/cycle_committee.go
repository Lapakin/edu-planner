package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type CycleCommitteeService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewCycleCommitteeService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *CycleCommitteeService {
	return &CycleCommitteeService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *CycleCommitteeService) CreateCycleCommittees(ctx context.Context, claims *jwt.Claims, cycleCommittees domain.CycleCommittees) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating cycle committees. err: %v", err)
			s.l.Debugf("claims: %v, cycleCommittees: %v", claims, cycleCommittees)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	cycleCommittees.SetCreatedAt(time.Now())

	r := s.rm.NewCycleCommitteeRepo(tx)
	if err = r.CreateCycleCommittees(ctx, cycleCommittees); err != nil {
		return handleDBError(err)
	}

	toAttach := cycleCommittees.GroupByAcademicYear()
	if err = r.AttachCycleCommitteesToAcademicYear(ctx, toAttach); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(cycleCommittees))
	for i, c := range cycleCommittees {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(c)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeCycleCommittee,
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

func (s *CycleCommitteeService) GetCycleCommitteeByID(ctx context.Context, id uint64) (*domain.CycleCommittee, error) {
	cycleCommittee, err := s.rm.NewCycleCommitteeRepo(s.db).GetCycleCommitteeByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting cycle committee. err %v", err)
		s.l.Debugf("cycleCommitteeID: %v", id)
		return nil, handleDBError(err)
	}
	return cycleCommittee, nil
}

func (s *CycleCommitteeService) FetchCycleCommittees(ctx context.Context, filters f.Filters) (domain.CycleCommittees, error) {
	cycleCommittees, err := s.rm.NewCycleCommitteeRepo(s.db).FetchCycleCommittees(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching cycle committees. err %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return cycleCommittees, nil
}

func (s *CycleCommitteeService) UpdateCycleCommittees(ctx context.Context, claims *jwt.Claims, cycleCommittees domain.CycleCommittees) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating cycle committees. err: %v", err)
			s.l.Debugf("claims: %v, cycleCommittees: %v", claims, cycleCommittees)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	cycleCommittees.SetModifiedAt(time.Now())

	if err = s.rm.NewCycleCommitteeRepo(tx).UpdateCycleCommittees(ctx, cycleCommittees); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(cycleCommittees))
	for i, c := range cycleCommittees {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(c)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeCycleCommittee,
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

func (s *CycleCommitteeService) DeleteCycleCommittees(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting cycle committees. err: %v", err)
			s.l.Debugf("claims: %v, cycleCommitteeIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewCycleCommitteeRepo(tx).DeleteCycleCommittees(ctx, ids); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(ids))
	for i, id := range ids {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(id)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeDelete,
			ObjectType: domain.ObjectTypeCycleCommittee,
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
