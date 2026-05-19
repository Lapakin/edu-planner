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

type StudyPlanService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewStudyPlanService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *StudyPlanService {
	return &StudyPlanService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *StudyPlanService) CreateStudyPlans(ctx context.Context, claims *jwt.Claims, studyPlans domain.StudyPlans) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating study plans. err: %v", err)
			s.l.Debugf("claims: %v, studyPlans: %v", claims, studyPlans)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	studyPlans.SetCreatedAt(currentTime)

	if err = s.rm.NewStudyPlanRepo(tx).CreateStudyPlans(ctx, studyPlans); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(studyPlans))
	for i, sp := range studyPlans {
		var objBytes []byte
		objBytes, err = json.Marshal(sp)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeStudyPlan,
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

func (s *StudyPlanService) GetStudyPlanByID(ctx context.Context, id uint64) (*domain.StudyPlan, error) {
	studyPlan, err := s.rm.NewStudyPlanRepo(s.db).GetStudyPlanByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting study plan by id. err: %v", err)
		s.l.Debugf("studyPlanID: %v", id)
		return nil, handleDBError(err)
	}
	return studyPlan, nil
}

func (s *StudyPlanService) FetchStudyPlans(ctx context.Context, filters f.Filters) (domain.StudyPlans, error) {
	studyPlans, err := s.rm.NewStudyPlanRepo(s.db).FetchStudyPlans(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching study plans. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return studyPlans, nil
}

func (s *StudyPlanService) UpdateStudyPlans(ctx context.Context, claims *jwt.Claims, studyPlans domain.StudyPlans) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating study plans. err: %v", err)
			s.l.Debugf("claims: %v, studyPlans: %v", claims, studyPlans)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	studyPlans.SetModifiedAt(currentTime)

	if err = s.rm.NewStudyPlanRepo(tx).UpdateStudyPlans(ctx, studyPlans); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(studyPlans))
	for i, sp := range studyPlans {
		var objBytes []byte
		objBytes, err = json.Marshal(sp)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeStudyPlan,
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

func (s *StudyPlanService) DeleteStudyPlans(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting study plans. err: %v", err)
			s.l.Debugf("claims: %v, studyPlanIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewStudyPlanRepo(tx).DeleteStudyPlans(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeStudyPlan,
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
