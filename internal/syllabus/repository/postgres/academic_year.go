package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type AcademicYearRepository struct {
	db sqlx.ExtContext
}

func NewAcademicYearRepository(db sqlx.ExtContext) *AcademicYearRepository {
	return &AcademicYearRepository{
		db: db,
	}
}

func (r *AcademicYearRepository) CreateAcademicYear(ctx context.Context, academicYears domain.AcademicYears) error {
	builder := q.Insert(ctx).
		Into("academic_year").
		Columns(`
			name,
			start_date,
			end_date,
			created_at
		`)

	for _, a := range academicYears {
		builder = builder.Values(
			a.Name,
			a.StartDate,
			a.EndDate,
			a.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&academicYears[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *AcademicYearRepository) GetAcademicYearByID(ctx context.Context, id uint64) (*domain.AcademicYear, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			name,
			start_date,
			end_date,
			is_active,
			is_deleted,
			created_at,
			modified_at
		`).
		From("academic_year").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	academicYear := &domain.AcademicYear{}
	if err = sqlx.GetContext(ctx, r.db, academicYear, query, args...); err != nil {
		return nil, err
	}

	return academicYear, nil
}

func (r *AcademicYearRepository) FetchAcademicYears(ctx context.Context, filters f.Filters) (domain.AcademicYears, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			name,
			start_date,
			end_date,
			is_active,
			is_deleted,
			created_at,
			modified_at
		`).
		From("academic_year").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	academicYears := make(domain.AcademicYears, 0)
	if err = sqlx.SelectContext(ctx, r.db, &academicYears, query, args...); err != nil {
		return nil, err
	}

	return academicYears, nil
}

func (r *AcademicYearRepository) UpdateAcademicYears(ctx context.Context, academicYears domain.AcademicYears) error {
	builder := q.BatchUpdate(ctx).
		Table("academic_year").
		NumberOfRows(len(academicYears)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddVarCharColumn("name").
		AddTimeStampColumn("start_date").
		AddTimeStampColumn("end_date").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 5*len(academicYears))
	for _, a := range academicYears {
		args = append(args,
			a.ID,
			a.Name,
			a.StartDate,
			a.EndDate,
			a.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&academicYears[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *AcademicYearRepository) DeleteAcademicYears(ctx context.Context, ids []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("academic_year").
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

func (r *AcademicYearRepository) ActivateAcademicYear(ctx context.Context, id uint64, currentTime time.Time) error {
	query, args, err := q.Update(ctx).
		Table("academic_year").
		Set("is_active", true).
		Set("modified_at", currentTime).
		WhereID(id).
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

func (r *AcademicYearRepository) DeactivateAcademicYear(ctx context.Context, id uint64, currentTime time.Time) error {
	query, args, err := q.Update(ctx).
		Table("academic_year").
		Set("is_active", false).
		Set("modified_at", currentTime).
		WhereID(id).
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
