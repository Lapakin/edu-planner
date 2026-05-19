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

type GroupService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewGroupService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *GroupService {
	return &GroupService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *GroupService) CreateGroups(ctx context.Context, claims *jwt.Claims, groups domain.Groups) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating groups. err: %v", err)
			s.l.Debugf("claims: %v, groups: %v", claims, groups)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	groups.SetCreatedAt(time.Now())

	r := s.rm.NewGroupRepo(tx)
	if err = r.CreateGroups(ctx, groups); err != nil {
		return handleDBError(err)
	}

	toAttach := groups.GroupByAcademicYearID()
	if err = r.AttachGroupsToAcademicYear(ctx, toAttach); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(groups))
	for i, g := range groups {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(g)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeGroup,
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

func (s *GroupService) GetGroupByID(ctx context.Context, id uint64) (*domain.Group, error) {
	group, err := s.rm.NewGroupRepo(s.db).GetGroupByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting group by id. err: %v", err)
		s.l.Debugf("groupID: %v", id)
		return nil, handleDBError(err)
	}
	return group, nil
}

func (s *GroupService) FetchGroups(ctx context.Context, filters f.Filters) (domain.Groups, error) {
	groups, err := s.rm.NewGroupRepo(s.db).FetchGroups(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching groups. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return groups, nil
}

func (s *GroupService) UpdateGroups(ctx context.Context, claims *jwt.Claims, groups domain.Groups) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating groups. err: %v", err)
			s.l.Debugf("claims: %v, groups: %v", claims, groups)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	groups.SetModifiedAt(time.Now())

	if err = s.rm.NewGroupRepo(tx).UpdateGroups(ctx, groups); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(groups))
	for i, g := range groups {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(g)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeGroup,
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

func (s *GroupService) DeleteGroups(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting groups. err: %v", err)
			s.l.Debugf("claims: %v, groupIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewGroupRepo(tx).DeleteGroups(ctx, ids); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(ids))
	for i, id := range ids {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(map[string]uint64{"id": id})
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeDelete,
			ObjectType: domain.ObjectTypeGroup,
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

func (s *GroupService) CreateGroupSemesters(ctx context.Context, claims *jwt.Claims, groupSemesters domain.GroupSemesters) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating group semesters. err: %v", err)
			s.l.Debugf("claims: %v, groupSemesters: %v", claims, groupSemesters)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	groupSemesters.SetCreatedAt(time.Now())

	if err = s.rm.NewGroupRepo(tx).CreateGroupSemesters(ctx, groupSemesters); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(groupSemesters))
	for i, gs := range groupSemesters {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(gs)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeGroupSemester,
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

func (s *GroupService) GetGroupSemesterByID(ctx context.Context, id uint64) (*domain.GroupSemester, error) {
	groupSemester, err := s.rm.NewGroupRepo(s.db).GetGroupSemesterByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting group semester by id. err: %v", err)
		s.l.Debugf("groupSemesterID: %v", id)
		return nil, handleDBError(err)
	}
	return groupSemester, nil
}

func (s *GroupService) FetchGroupSemesters(ctx context.Context, filters f.Filters) (domain.GroupSemesters, error) {
	groupSemesters, err := s.rm.NewGroupRepo(s.db).FetchGroupSemesters(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching group semesters. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return groupSemesters, nil
}

func (s *GroupService) UpdateGroupSemesters(ctx context.Context, claims *jwt.Claims, groupSemesters domain.GroupSemesters) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating group semesters. err: %v", err)
			s.l.Debugf("claims: %v, groupSemesters: %v", claims, groupSemesters)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	groupSemesters.SetModifiedAt(time.Now())

	if err = s.rm.NewGroupRepo(tx).UpdateGroupSemesters(ctx, groupSemesters); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(groupSemesters))
	for i, gs := range groupSemesters {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(gs)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeGroupSemester,
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

func (s *GroupService) DeleteGroupSemesters(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting group semesters. err: %v", err)
			s.l.Debugf("claims: %v, groupSemesterIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewGroupRepo(tx).DeleteGroupSemesters(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeGroupSemester,
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
