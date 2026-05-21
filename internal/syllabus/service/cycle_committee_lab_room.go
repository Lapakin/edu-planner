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

type CycleCommitteeLabRoomService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewCycleCommitteeLabRoomService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *CycleCommitteeLabRoomService {
	return &CycleCommitteeLabRoomService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *CycleCommitteeLabRoomService) CreateCycleCommitteeLabRooms(ctx context.Context, claims *jwt.Claims, labRooms domain.CycleCommitteeLabRooms) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating cycle committee lab rooms. err: %v", err)
			s.l.Debugf("claims: %v, labRooms: %v", claims, labRooms)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	labRooms.SetCreatedAt(time.Now())

	if err = s.rm.NewCycleCommitteeLabRoomRepo(tx).CreateCycleCommitteeLabRooms(ctx, labRooms); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(labRooms))
	for i, lr := range labRooms {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(lr); err != nil {
			return ErrInternal
		}
		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeCycleCommitteeLabRoom,
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

func (s *CycleCommitteeLabRoomService) GetCycleCommitteeLabRoomByID(ctx context.Context, id uint64) (*domain.CycleCommitteeLabRoom, error) {
	labRoom, err := s.rm.NewCycleCommitteeLabRoomRepo(s.db).GetCycleCommitteeLabRoomByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting cycle committee lab room by id. err: %v", err)
		return nil, handleDBError(err)
	}
	return labRoom, nil
}

func (s *CycleCommitteeLabRoomService) FetchCycleCommitteeLabRooms(ctx context.Context, filters f.Filters) (domain.CycleCommitteeLabRooms, error) {
	labRooms, err := s.rm.NewCycleCommitteeLabRoomRepo(s.db).FetchCycleCommitteeLabRooms(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching cycle committee lab rooms. err: %v", err)
		return nil, ErrInternal
	}
	return labRooms, nil
}

func (s *CycleCommitteeLabRoomService) UpdateCycleCommitteeLabRooms(ctx context.Context, claims *jwt.Claims, labRooms domain.CycleCommitteeLabRooms) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating cycle committee lab rooms. err: %v", err)
			s.l.Debugf("claims: %v, labRooms: %v", claims, labRooms)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	labRooms.SetModifiedAt(time.Now())

	if err = s.rm.NewCycleCommitteeLabRoomRepo(tx).UpdateCycleCommitteeLabRooms(ctx, labRooms); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(labRooms))
	for i, lr := range labRooms {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(lr); err != nil {
			return ErrInternal
		}
		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeCycleCommitteeLabRoom,
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

func (s *CycleCommitteeLabRoomService) DeleteCycleCommitteeLabRooms(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting cycle committee lab rooms. err: %v", err)
			s.l.Debugf("claims: %v, labRoomIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewCycleCommitteeLabRoomRepo(tx).DeleteCycleCommitteeLabRooms(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeCycleCommitteeLabRoom,
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
