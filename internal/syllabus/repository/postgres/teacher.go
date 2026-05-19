package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type TeacherRepository struct {
	db sqlx.ExtContext
}

func NewTeacherRepository(db sqlx.ExtContext) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) CreateTeachers(ctx context.Context, teachers domain.Teachers) error {
	builder := q.Insert(ctx).
		Into("teacher").
		Columns(`
			id,
			user_id,
			created_at
		`)

	for _, te := range teachers {
		builder = builder.Values(
			te.ID,
			te.UserID,
			te.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&teachers[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *TeacherRepository) GetTeacherByID(ctx context.Context, id uint64) (*domain.Teacher, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			t.id,
			t.user_id,
			ta.academic_year_id,
			t.created_at,
			t.modified_at
		`).
		From("teacher t").
		InnerJoin("academic_year_to_teacher AS ta ON ta.teacher_id = t.id").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	teacher := &domain.Teacher{}
	if err = sqlx.GetContext(ctx, r.db, teacher, query, args...); err != nil {
		return nil, err
	}

	return teacher, nil
}

func (r *TeacherRepository) FetchTeachers(ctx context.Context, filters f.Filters) (domain.Teachers, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			t.id,
			t.user_id,
			ta.academic_year_id,
			t.created_at,
			t.modified_at
		`).
		From("teacher t").
		InnerJoin("academic_year_to_teacher AS ta ON ta.teacher_id = t.id").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "t.id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "ta.academic_year_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("t.id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	teachers := make(domain.Teachers, 0)
	if err = sqlx.SelectContext(ctx, r.db, &teachers, query, args...); err != nil {
		return nil, err
	}

	return teachers, nil
}

func (r *TeacherRepository) UpdateTeachers(ctx context.Context, teachers domain.Teachers) error {
	builder := q.BatchUpdate(ctx).
		Table("teacher").
		PrimaryKey("id").
		NumberOfRows(len(teachers)).
		AddBigintColumn("id").
		AddBigintColumn("user_id").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 3*len(teachers))
	for _, te := range teachers {
		args = append(args,
			te.ID,
			te.UserID,
			te.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&teachers[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *TeacherRepository) DeleteTeachers(ctx context.Context, ids []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("teacher").
		WhereID(ids).
		ToSQL()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TeacherRepository) AttachTeachersToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Teachers) error {
	builder := q.Insert(ctx).
		Into("academic_year_to_teacher").
		Columns(`
			academic_year_id,
			teacher_id
		`)

	for academicYearID, teachers := range toAttach {
		for _, teacher := range teachers {
			builder = builder.Values(
				academicYearID,
				teacher.ID,
			)
		}
	}

	query, args, err := builder.ToSQL()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
