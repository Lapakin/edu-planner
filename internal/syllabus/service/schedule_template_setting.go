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

type ScheduleTemplateSettingService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewScheduleTemplateSettingService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *ScheduleTemplateSettingService {
	return &ScheduleTemplateSettingService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *ScheduleTemplateSettingService) CreateScheduleTemplateSettings(ctx context.Context, claims *jwt.Claims, settings domain.ScheduleTemplateSettings) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating schedule template settings. err: %v", err)
			s.l.Debugf("claims: %v, settings: %v", claims, settings)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	settings.SetCreatedAt(currentTime)

	if err = s.rm.NewScheduleTemplateSettingRepo(tx).CreateScheduleTemplateSettings(ctx, settings); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(settings))
	for i, se := range settings {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(se); err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeScheduleTemplateSetting,
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

func (s *ScheduleTemplateSettingService) GetScheduleTemplateSettingByID(ctx context.Context, id uint64) (*domain.ScheduleTemplateSetting, error) {
	setting, err := s.rm.NewScheduleTemplateSettingRepo(s.db).GetScheduleTemplateSettingByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting schedule template setting by id. err: %v", err)
		s.l.Debugf("scheduleTemplateSettingID: %v", id)
		return nil, handleDBError(err)
	}
	return setting, nil
}

func (s *ScheduleTemplateSettingService) FetchScheduleTemplateSettings(ctx context.Context, filters f.Filters) (domain.ScheduleTemplateSettings, error) {
	settings, err := s.rm.NewScheduleTemplateSettingRepo(s.db).FetchScheduleTemplateSettings(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching schedule template settings. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return settings, nil
}

func (s *ScheduleTemplateSettingService) UpdateScheduleTemplateSettings(ctx context.Context, claims *jwt.Claims, settings domain.ScheduleTemplateSettings) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating schedule template settings. err: %v", err)
			s.l.Debugf("claims: %v, settings: %v", claims, settings)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	settings.SetModifiedAt(currentTime)

	if err = s.rm.NewScheduleTemplateSettingRepo(tx).UpdateScheduleTemplateSettings(ctx, settings); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(settings))
	for i, se := range settings {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(se); err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeScheduleTemplateSetting,
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

func (s *ScheduleTemplateSettingService) DeleteScheduleTemplateSettings(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting schedule template settings. err: %v", err)
			s.l.Debugf("claims: %v, scheduleTemplateSettingIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewScheduleTemplateSettingRepo(tx).DeleteScheduleTemplateSettings(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeScheduleTemplateSetting,
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
