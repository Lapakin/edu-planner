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

type DepartmentService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewDepartmentService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *DepartmentService {
	return &DepartmentService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *DepartmentService) CreateDepartments(ctx context.Context, claims *jwt.Claims, departments domain.Departments) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating departments. err: %v", err)
			s.l.Debugf("claims: %v, departments: %v", claims, departments)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	departments.SetCreatedAt(time.Now())

	r := s.rm.NewDepartmentRepo(tx)
	if err = r.CreateDepartments(ctx, departments); err != nil {
		return handleDBError(err)
	}

	toAttach := departments.GroupByAcademicYear()
	if err = r.AttachDepartmentsToAcademicYear(ctx, toAttach); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(departments))
	for i, d := range departments {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(d)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeDepartment,
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

func (s *DepartmentService) GetDepartmentByID(ctx context.Context, departmentID uint64) (*domain.Department, error) {
	department, err := s.rm.NewDepartmentRepo(s.db).GetDepartmentByID(ctx, departmentID)
	if err != nil {
		s.l.Errorf("error during getting department by id. err: %v", err)
		s.l.Debugf("departmentID: %v", departmentID)
		return nil, handleDBError(err)
	}
	return department, nil
}

func (s *DepartmentService) FetchDepartments(ctx context.Context, filters f.Filters) (domain.Departments, error) {
	departments, err := s.rm.NewDepartmentRepo(s.db).FetchDepartments(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching departments. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return departments, nil
}

func (s *DepartmentService) UpdateDepartments(ctx context.Context, claims *jwt.Claims, departments domain.Departments) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating departments. err: %v", err)
			s.l.Debugf("claims: %v, departments: %v", claims, departments)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	departments.SetModifiedAt(time.Now())

	if err = s.rm.NewDepartmentRepo(tx).UpdateDepartments(ctx, departments); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(departments))
	for i, d := range departments {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(d)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeDepartment,
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

func (s *DepartmentService) DeleteDepartments(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting departments. err: %v", err)
			s.l.Debugf("claims: %v, departmentIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewDepartmentRepo(tx).DeleteDepartments(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeDepartment,
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
