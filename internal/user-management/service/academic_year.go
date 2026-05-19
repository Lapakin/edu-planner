package service

import (
	"context"
	"fmt"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/user-management/repository"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
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

func (s *AcademicYearService) MassageConsumer(ctx context.Context, massage *domain.Massage) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during processing massage. err: %v", err)
			s.l.Debugf("massage: %v", massage)
		}
	}()

	switch massage.ActionType {
	case domain.ActionTypeCreate, domain.ActionTypeUpdate:
		var academicYear *domain.AcademicYear
		if err = json.Unmarshal(massage.Object, &academicYear); err != nil {
			return ErrUnmarshal
		}
		if err = s.upsertAcademicYear(ctx, academicYear); err != nil {
			return ErrInternal
		}

	case domain.ActionTypeDelete:
		var id uint64
		if err = json.Unmarshal(massage.Object, &id); err != nil {
			return ErrUnmarshal
		}
		if err = s.DeleteAcademicYear(ctx, id); err != nil {
			return ErrInternal
		}
	}

	return nil
}

func (s *AcademicYearService) upsertAcademicYear(ctx context.Context, academicYear *domain.AcademicYear) error {
	if err := s.rm.NewAcademicYearRepo(s.db).UpsertAcademicYear(ctx, academicYear); err != nil {
		return fmt.Errorf("error during upserting academic year by id. err: %w", err)
	}
	return nil
}

func (s *AcademicYearService) DeleteAcademicYear(ctx context.Context, id uint64) error {
	if err := s.rm.NewAcademicYearRepo(s.db).DeleteAcademicYear(ctx, id); err != nil {
		return fmt.Errorf("error during deleting academic year by id. err: %w", err)
	}
	return nil
}
