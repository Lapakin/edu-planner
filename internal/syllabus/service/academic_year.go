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

type AcademicYearService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewAcademicYearService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *AcademicYearService {
	return &AcademicYearService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *AcademicYearService) CreateAcademicYears(ctx context.Context, claims *jwt.Claims, academicYears domain.AcademicYears) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating academic years. err: %v", err)
			s.l.Debugf("claims: %v, academicYears: %v", claims, academicYears)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	academicYears.SetCreatedAt(time.Now())

	if err = s.rm.NewAcademicYearRepo(tx).CreateAcademicYear(ctx, academicYears); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(academicYears))
	for i, a := range academicYears {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(a)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeAcademicYear,
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

func (s *AcademicYearService) GetAcademicYearByID(ctx context.Context, id uint64) (*domain.AcademicYear, error) {
	academicYear, err := s.rm.NewAcademicYearRepo(s.db).GetAcademicYearByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting academic year. err %v", err)
		s.l.Debugf("academicYearID: %v", id)
		return nil, handleDBError(err)
	}
	return academicYear, nil
}

func (s *AcademicYearService) FetchAcademicYears(ctx context.Context, filters f.Filters) (domain.AcademicYears, error) {
	academicYears, err := s.rm.NewAcademicYearRepo(s.db).FetchAcademicYears(ctx, filters)
	if err != nil {
		s.l.Errorf("Error during fetching academic years. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return academicYears, nil
}

func (s *AcademicYearService) UpdateAcademicYears(ctx context.Context, claims *jwt.Claims, academicYears domain.AcademicYears) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("Error during updating academic years. err: %v", err)
			s.l.Debugf("claims: %v, academicYears: %v", claims, academicYears)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	academicYears.SetModifiedAt(time.Now())

	if err = s.rm.NewAcademicYearRepo(tx).UpdateAcademicYears(ctx, academicYears); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(academicYears))
	for i, a := range academicYears {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(a)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeAcademicYear,
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

func (s *AcademicYearService) DeleteAcademicYears(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting academic years. err: %v", err)
			s.l.Debugf("claims: %v, academicYearIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewAcademicYearRepo(tx).DeleteAcademicYears(ctx, ids); err != nil {
		return handleDBError(err)
	}

	for _, id := range ids {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(id)
		if err != nil {
			return ErrInternal
		}

		massage := &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeDelete,
			ObjectType: domain.ObjectTypeAcademicYear,
			Object:     objBytes,
			Claims:     claims,
		}

		if err = s.rm.NewMessageRepo(tx).Write(ctx, domain.Massages{massage}); err != nil {
			return handleDBError(err)
		}
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *AcademicYearService) ActivateAcademicYear(ctx context.Context, claims *jwt.Claims, academicYearID uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during activating academic year. err: %v", err)
			s.l.Debugf("claims: %v, academicYearID: %v", claims, academicYearID)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	r := s.rm.NewAcademicYearRepo(tx)
	if err = r.DeactivateAllAcademicYears(ctx, currentTime); err != nil {
		return handleDBError(err)
	}

	if err = s.rm.NewAcademicYearRepo(tx).ActivateAcademicYear(ctx, academicYearID, currentTime); err != nil {
		return handleDBError(err)
	}

	var objBytes json.RawMessage
	objBytes, err = json.Marshal(academicYearID)
	if err != nil {
		return ErrInternal
	}

	massage := &domain.Massage{
		EventID:    uuid.New().String(),
		ActionType: domain.ActionTypeActivate,
		ObjectType: domain.ObjectTypeAcademicYear,
		Object:     objBytes,
		Claims:     claims,
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, domain.Massages{massage}); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}
