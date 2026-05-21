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

type TeacherSlotPreferenceService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewTeacherSlotPreferenceService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *TeacherSlotPreferenceService {
	return &TeacherSlotPreferenceService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *TeacherSlotPreferenceService) CreateTeacherSlotPreferences(ctx context.Context, claims *jwt.Claims, preferences domain.TeacherSlotPreferences) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating teacher slot preferences. err: %v", err)
			s.l.Debugf("claims: %v, preferences: %v", claims, preferences)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	preferences.SetCreatedAt(time.Now())

	if err = s.rm.NewTeacherSlotPreferenceRepo(tx).CreateTeacherSlotPreferences(ctx, preferences); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(preferences))
	for i, p := range preferences {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(p); err != nil {
			return ErrInternal
		}
		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeTeacherSlotPreference,
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

func (s *TeacherSlotPreferenceService) GetTeacherSlotPreferenceByID(ctx context.Context, id uint64) (*domain.TeacherSlotPreference, error) {
	preference, err := s.rm.NewTeacherSlotPreferenceRepo(s.db).GetTeacherSlotPreferenceByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting teacher slot preference by id. err: %v", err)
		return nil, handleDBError(err)
	}
	return preference, nil
}

func (s *TeacherSlotPreferenceService) FetchTeacherSlotPreferences(ctx context.Context, filters f.Filters) (domain.TeacherSlotPreferences, error) {
	preferences, err := s.rm.NewTeacherSlotPreferenceRepo(s.db).FetchTeacherSlotPreferences(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching teacher slot preferences. err: %v", err)
		return nil, ErrInternal
	}
	return preferences, nil
}

func (s *TeacherSlotPreferenceService) UpdateTeacherSlotPreferences(ctx context.Context, claims *jwt.Claims, preferences domain.TeacherSlotPreferences) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating teacher slot preferences. err: %v", err)
			s.l.Debugf("claims: %v, preferences: %v", claims, preferences)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	preferences.SetModifiedAt(time.Now())

	if err = s.rm.NewTeacherSlotPreferenceRepo(tx).UpdateTeacherSlotPreferences(ctx, preferences); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(preferences))
	for i, p := range preferences {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(p); err != nil {
			return ErrInternal
		}
		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeTeacherSlotPreference,
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

func (s *TeacherSlotPreferenceService) DeleteTeacherSlotPreferences(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting teacher slot preferences. err: %v", err)
			s.l.Debugf("claims: %v, preferenceIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewTeacherSlotPreferenceRepo(tx).DeleteTeacherSlotPreferences(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeTeacherSlotPreference,
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
