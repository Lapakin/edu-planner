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

type DisciplineService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewDisciplineService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *DisciplineService {
	return &DisciplineService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s DisciplineService) CreateDisciplines(ctx context.Context, claims *jwt.Claims, disciplines domain.Disciplines) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating disciplines. err: %v", err)
			s.l.Debugf("claims: %v, disciplines: %v", claims, disciplines)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	disciplines.SetCreatedAt(time.Now())

	r := s.rm.NewDisciplineRepo(tx)
	if err = r.CreateDisciplines(ctx, disciplines); err != nil {
		return handleDBError(err)
	}

	toAttach := disciplines.GroupByAcademicYear()
	if err = r.AttachDisciplinesToAcademicYear(ctx, toAttach); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(disciplines))
	for i, d := range disciplines {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(d); err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeDiscipline,
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

func (s DisciplineService) GetDisciplineByID(ctx context.Context, id uint64) (*domain.Discipline, error) {
	discipline, err := s.rm.NewDisciplineRepo(s.db).GetDisciplineByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting discipline by id. err: %v", err)
		s.l.Debugf("disciplineID: %v", id)
		return nil, handleDBError(err)
	}
	return discipline, nil
}

func (s DisciplineService) FetchDisciplines(ctx context.Context, filters f.Filters) (domain.Disciplines, error) {
	disciplines, err := s.rm.NewDisciplineRepo(s.db).FetchDisciplines(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching disciplines. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return disciplines, nil
}

func (s DisciplineService) UpdateDisciplines(ctx context.Context, claims *jwt.Claims, disciplines domain.Disciplines) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating disciplines. err: %v", err)
			s.l.Debugf("claims: %v, disciplines: %v", claims, disciplines)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	disciplines.SetModifiedAt(time.Now())

	if err = s.rm.NewDisciplineRepo(tx).UpdateDisciplines(ctx, disciplines); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(disciplines))
	for i, d := range disciplines {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(d); err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeDiscipline,
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

func (s DisciplineService) DeleteDisciplines(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting disciplines. err: %v", err)
			s.l.Debugf("claims: %v, disciplineIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewDisciplineRepo(tx).DeleteDisciplines(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeDiscipline,
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
