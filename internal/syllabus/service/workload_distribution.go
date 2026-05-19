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

type WorkloadDistributionService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewWorkloadDistributionService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *WorkloadDistributionService {
	return &WorkloadDistributionService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *WorkloadDistributionService) CreateWorkloadDistributions(ctx context.Context, claims *jwt.Claims, distributions domain.WorkloadDistributions) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating workload distributions. err: %v", err)
			s.l.Debugf("claims: %v, distributions: %v", claims, distributions)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	distributions.SetCreatedAt(currentTime)

	if err = s.rm.NewWorkloadDistributionRepo(tx).CreateWorkloadDistributions(ctx, distributions); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(distributions))
	for i, wd := range distributions {
		var objBytes []byte
		objBytes, err = json.Marshal(wd)
		if err != nil {
			return ErrInternal
		}
		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeWorkloadDistribution,
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

func (s *WorkloadDistributionService) GetWorkloadDistributionByID(ctx context.Context, id uint64) (*domain.WorkloadDistribution, error) {
	workload, err := s.rm.NewWorkloadDistributionRepo(s.db).GetWorkloadDistributionByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting workload distribution by id. err: %v", err)
		s.l.Debugf("workloadDistributionID: %v", id)
		return nil, handleDBError(err)
	}
	return workload, nil
}

func (s *WorkloadDistributionService) FetchWorkloadDistributions(ctx context.Context, filters f.Filters) (domain.WorkloadDistributions, error) {
	distributions, err := s.rm.NewWorkloadDistributionRepo(s.db).FetchWorkloadDistributions(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching workload distributions. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return distributions, nil
}

func (s *WorkloadDistributionService) UpdateWorkloadDistributions(ctx context.Context, claims *jwt.Claims, distributions domain.WorkloadDistributions) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating workload distributions. err: %v", err)
			s.l.Debugf("claims: %v, distributions: %v", claims, distributions)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	distributions.SetModifiedAt(currentTime)

	if err = s.rm.NewWorkloadDistributionRepo(tx).UpdateWorkloadDistributions(ctx, distributions); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(distributions))
	for i, wd := range distributions {
		var objBytes []byte
		objBytes, err = json.Marshal(wd)
		if err != nil {
			return ErrInternal
		}
		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeWorkloadDistribution,
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

func (s *WorkloadDistributionService) DeleteWorkloadDistributions(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting workload distributions. err: %v", err)
			s.l.Debugf("claims: %v, workloadDistributionIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewWorkloadDistributionRepo(tx).DeleteWorkloadDistributions(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeWorkloadDistribution,
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
