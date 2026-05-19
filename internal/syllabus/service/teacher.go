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

type TeacherService struct {
	db *pg.DB
	rm repository.RepoManager
	l  *logging.Logger
}

func NewTeacherService(db *pg.DB, rm repository.RepoManager, l *logging.Logger) *TeacherService {
	return &TeacherService{
		db: db,
		rm: rm,
		l:  l,
	}
}

func (s *TeacherService) CreateTeachers(ctx context.Context, claims *jwt.Claims, teachers domain.Teachers) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during creating teachers. err: %v", err)
			s.l.Debugf("claims: %v, teachers: %v", claims, teachers)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	teachers.SetCreatedAt(time.Now())

	r := s.rm.NewTeacherRepo(tx)
	if err = r.CreateTeachers(ctx, teachers); err != nil {
		return handleDBError(err)
	}

	toAttach := teachers.GroupByAcademicYear()
	if err = r.AttachTeachersToAcademicYear(ctx, toAttach); err != nil {
		return handleDBError(err)
	}

	messages := make(domain.Massages, len(teachers))
	for i, t := range teachers {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(t)
		if err != nil {
			return ErrInternal
		}

		messages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeCreate,
			ObjectType: domain.ObjectTypeTeacher,
			Object:     objBytes,
			Claims:     claims,
		}
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, messages); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *TeacherService) GetTeacherByID(ctx context.Context, id uint64) (*domain.Teacher, error) {
	teacher, err := s.rm.NewTeacherRepo(s.db).GetTeacherByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting teacher by id. err: %v", err)
		s.l.Debugf("teacherID: %v", id)
		return nil, handleDBError(err)
	}
	return teacher, nil
}

func (s *TeacherService) FetchTeachers(ctx context.Context, filters f.Filters) (domain.Teachers, error) {
	teachers, err := s.rm.NewTeacherRepo(s.db).FetchTeachers(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching teachers. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, handleDBError(err)
	}
	return teachers, nil
}

func (s *TeacherService) UpdateTeachers(ctx context.Context, claims *jwt.Claims, teachers domain.Teachers) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during updating teachers. err: %v", err)
			s.l.Debugf("claims: %v, teachers: %v", claims, teachers)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	teachers.SetModifiedAt(time.Now())

	if err = s.rm.NewTeacherRepo(tx).UpdateTeachers(ctx, teachers); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(teachers))
	for i, t := range teachers {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(t)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeUpdate,
			ObjectType: domain.ObjectTypeTeacher,
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

func (s *TeacherService) DeleteTeachers(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting teachers. err: %v", err)
			s.l.Debugf("claims: %v, teacherIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewTeacherRepo(tx).DeleteTeachers(ctx, ids); err != nil {
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
			ObjectType: domain.ObjectTypeTeacher,
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
