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

type RoomService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewRoomService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *RoomService {
	return &RoomService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *RoomService) CreateRooms(ctx context.Context, claims *jwt.Claims, rooms domain.Rooms) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating rooms. err: %v", err)
			s.l.Debugf("claims: %v, rooms: %v", claims, rooms)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	rooms.SetCreatedAt(time.Now())

	r := s.rm.NewRoomRepo(tx)
	if err = r.CreateRooms(ctx, rooms); err != nil {
		return handleDBError(err)
	}

	toAttach := rooms.GroupByAcademicYear()
	if err = r.AttachRoomsToAcademicYear(ctx, toAttach); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(rooms))
	for i, r := range rooms {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(r); err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeRoom,
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

func (s *RoomService) GetRoomByID(ctx context.Context, roomID uint64) (*domain.Room, error) {
	room, err := s.rm.NewRoomRepo(s.db).GetRoomByID(ctx, roomID)
	if err != nil {
		s.l.Errorf("error getting room by id. err: %v", err)
		s.l.Debugf("roomID: %v", roomID)
		return nil, handleDBError(err)
	}
	return room, nil
}

func (s *RoomService) FetchRooms(ctx context.Context, filters f.Filters) (domain.Rooms, error) {
	rooms, err := s.rm.NewRoomRepo(s.db).FetchRooms(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching rooms. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return rooms, nil
}

func (s *RoomService) UpdateRooms(ctx context.Context, claims *jwt.Claims, rooms domain.Rooms) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating rooms. err: %v", err)
			s.l.Debugf("claims: %v, rooms: %v", claims, rooms)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	rooms.SetModifiedAt(time.Now())

	if err = s.rm.NewRoomRepo(tx).UpdateRooms(ctx, rooms); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(rooms))
	for i, r := range rooms {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(r); err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeRoom,
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

func (s *RoomService) DeleteRooms(ctx context.Context, claims *jwt.Claims, roomIDs []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting rooms. err: %v", err)
			s.l.Debugf("claims: %v, roomIDs: %v", claims, roomIDs)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewRoomRepo(tx).DeleteRooms(ctx, roomIDs); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(roomIDs))
	for i, id := range roomIDs {
		var objBytes json.RawMessage
		if objBytes, err = json.Marshal(id); err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeDelete,
			ObjectType: domain.ObjectTypeRoom,
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
