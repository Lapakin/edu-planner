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

type SpecialtyService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewSpecialtyService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *SpecialtyService {
	return &SpecialtyService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *SpecialtyService) CreateSpecialties(ctx context.Context, claims *jwt.Claims, specialities domain.Specialties) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating specialties. err: %v", err)
			s.l.Debugf("claims: %v, specialities: %v", claims, specialities)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	specialities.SetCreatedAt(time.Now())

	r := s.rm.NewSpecialtyRepo(tx)
	if err = r.CreateSpecialties(ctx, specialities); err != nil {
		return handleDBError(err)
	}

	toAttach := specialities.GroupByAcademicYear()
	if err = r.AttachSpecialtiesToAcademicYear(ctx, toAttach); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(specialities))
	for i, sp := range specialities {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(sp)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeSpecialty,
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

func (s *SpecialtyService) GetSpecialtyByID(ctx context.Context, specialtyID uint64) (*domain.Specialty, error) {
	specialty, err := s.rm.NewSpecialtyRepo(s.db).GetSpecialtyByID(ctx, specialtyID)
	if err != nil {
		s.l.Errorf("error during getting specialty by id. err: %v", err)
		s.l.Debugf("specialtyID: %v", specialtyID)
		return nil, handleDBError(err)
	}
	return specialty, nil
}

func (s *SpecialtyService) FetchSpecialties(ctx context.Context, filters f.Filters) (domain.Specialties, error) {
	specialties, err := s.rm.NewSpecialtyRepo(s.db).FetchSpecialties(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching specialties. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return specialties, nil
}

func (s *SpecialtyService) UpdateSpecialties(ctx context.Context, claims *jwt.Claims, specialities domain.Specialties) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating specialties. err: %v", err)
			s.l.Debugf("claims: %v, specialities: %v", claims, specialities)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	specialities.SetModifiedAt(time.Now())

	if err = s.rm.NewSpecialtyRepo(tx).UpdateSpecialties(ctx, specialities); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(specialities))
	for i, sp := range specialities {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(sp)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeSpecialty,
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

func (s *SpecialtyService) DeleteSpecialties(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting specialties. err: %v", err)
			s.l.Debugf("claims: %v, ids: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewSpecialtyRepo(tx).DeleteSpecialties(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeSpecialty,
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
