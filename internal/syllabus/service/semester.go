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

type SemesterService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewSemesterService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *SemesterService {
	return &SemesterService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s SemesterService) CreateSemesters(ctx context.Context, claims *jwt.Claims, semesters domain.Semesters) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating semesters. err: %v", err)
			s.l.Debugf("claims: %v, semesters: %v", claims, semesters)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	semesters.SetCreatedAt(currentTime)

	if err = s.rm.NewSemesterRepo(tx).CreateSemesters(ctx, semesters); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(semesters))
	for i, se := range semesters {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(se); err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeSemester,
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

func (s SemesterService) GetSemesterByID(ctx context.Context, semesterID uint64) (*domain.Semester, error) {
	semester, err := s.rm.NewSemesterRepo(s.db).GetSemesterByID(ctx, semesterID)
	if err != nil {
		s.l.Errorf("error during getting semester by id. err: %v", err)
		s.l.Debugf("semesterID: %v", semesterID)
		return nil, handleDBError(err)
	}
	return semester, nil
}

func (s SemesterService) FetchSemesters(ctx context.Context, filters f.Filters) (domain.Semesters, error) {
	semesters, err := s.rm.NewSemesterRepo(s.db).FetchSemesters(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching semesters. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return semesters, nil
}

func (s SemesterService) UpdateSemesters(ctx context.Context, claims *jwt.Claims, semesters domain.Semesters) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating semesters. err: %v", err)
			s.l.Debugf("claims: %v, semesters: %v", claims, semesters)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	semesters.SetModifiedAt(currentTime)

	if err = s.rm.NewSemesterRepo(tx).UpdateSemesters(ctx, semesters); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(semesters))
	for i, se := range semesters {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(se); err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeSemester,
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

func (s SemesterService) DeleteSemesters(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting semesters. err: %v", err)
			s.l.Debugf("claims: %v, semesterIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewSemesterRepo(tx).DeleteSemesters(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeSemester,
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
