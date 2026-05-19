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

type WorkloadAssignmentService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewWorkloadAssignmentService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *WorkloadAssignmentService {
	return &WorkloadAssignmentService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *WorkloadAssignmentService) CreateWorkloadAssignments(ctx context.Context, claims *jwt.Claims, assignments domain.WorkloadAssignments) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating workload assignments. err: %v", err)
			s.l.Debugf("claims: %v, assignments: %v", claims, assignments)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	assignments.SetCreatedAt(currentTime)

	if err = s.rm.NewWorkloadAssignmentRepo(tx).CreateWorkloadAssignments(ctx, assignments); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(assignments))
	for i, wa := range assignments {
		var objBytes []byte
		objBytes, err = json.Marshal(wa)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeWorkloadAssignment,
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

func (s *WorkloadAssignmentService) GetWorkloadAssignmentByID(ctx context.Context, id uint64) (*domain.WorkloadAssignment, error) {
	workload, err := s.rm.NewWorkloadAssignmentRepo(s.db).GetWorkloadAssignmentByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting workload assignment by id. err: %v", err)
		s.l.Debugf("workloadAssignmentID: %v", id)
		return nil, handleDBError(err)
	}
	return workload, nil
}

func (s *WorkloadAssignmentService) FetchWorkloadAssignments(ctx context.Context, filters f.Filters) (domain.WorkloadAssignments, error) {
	workloads, err := s.rm.NewWorkloadAssignmentRepo(s.db).FetchWorkloadAssignments(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching workload assignments. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return workloads, nil
}

func (s *WorkloadAssignmentService) UpdateWorkloadAssignments(ctx context.Context, claims *jwt.Claims, assignments domain.WorkloadAssignments) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating workload assignments. err: %v", err)
			s.l.Debugf("claims: %v, assignments: %v", claims, assignments)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	assignments.SetCreatedAt(currentTime)

	if err = s.rm.NewWorkloadAssignmentRepo(tx).UpdateWorkloadAssignments(ctx, assignments); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(assignments))
	for i, wa := range assignments {
		var objBytes []byte
		objBytes, err = json.Marshal(wa)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeWorkloadAssignment,
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

func (s *WorkloadAssignmentService) DeleteWorkloadAssignments(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting workload assignments. err: %v", err)
			s.l.Debugf("claims: %v, workloadAssignmentIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewWorkloadAssignmentRepo(tx).DeleteWorkloadAssignments(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeWorkloadAssignment,
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
