package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type SemesterRepository struct {
	db sqlx.ExtContext
}

func NewSemesterRepository(db sqlx.ExtContext) *SemesterRepository {
	return &SemesterRepository{
		db: db,
	}
}

func (r *SemesterRepository) CreateSemesters(ctx context.Context, semesters domain.Semesters) error {
	builder := q.Insert(ctx).
		Into("semester").
		Columns(`
			academic_year_id,
			period_start,
			period_end,
			created_at
		`)

	for _, se := range semesters {
		builder = builder.Values(
			se.AcademicYearID,
			se.PeriodStart,
			se.PeriodEnd,
			se.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&semesters[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *SemesterRepository) GetSemesterByID(ctx context.Context, id uint64) (*domain.Semester, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			s.id,
			s.academic_year_id,
			s.period_start,
			s.period_end,
			s.created_at,
			s.modified_at
		`).
		From("semester s").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	semester := &domain.Semester{}
	if err = sqlx.GetContext(ctx, r.db, semester, query, args...); err != nil {
		return nil, err
	}

	return semester, nil
}

func (r *SemesterRepository) FetchSemesters(ctx context.Context, filters f.Filters) (domain.Semesters, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			s.id,
			s.academic_year_id,
			s.period_start,
			s.period_end,
			s.created_at,
			s.modified_at
		`).
		From("semester s").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "s.id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "s.academic_year_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("s.id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	semesters := make(domain.Semesters, 0)
	if err = sqlx.SelectContext(ctx, r.db, &semesters, query, args...); err != nil {
		return nil, err
	}

	return semesters, nil
}

func (r *SemesterRepository) UpdateSemesters(ctx context.Context, semesters domain.Semesters) error {
	builder := q.BatchUpdate(ctx).
		Table("semester").
		NumberOfRows(len(semesters)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddBigintColumn("academic_year_id").
		AddTimeStampColumn("period_start").
		AddTimeStampColumn("period_end").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 5*len(semesters))
	for _, se := range semesters {
		args = append(args,
			se.ID,
			se.AcademicYearID,
			se.PeriodStart,
			se.PeriodEnd,
			se.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&semesters[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *SemesterRepository) DeleteSemesters(ctx context.Context, ids []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("semester").
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
