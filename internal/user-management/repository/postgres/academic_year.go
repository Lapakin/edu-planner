package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
)

type AcademicYearRepository struct {
	db sqlx.ExtContext
}

func NewAcademicYearRepository(db sqlx.ExtContext) *AcademicYearRepository {
	return &AcademicYearRepository{db: db}
}

func (r *AcademicYearRepository) UpsertAcademicYear(ctx context.Context, academicYear *domain.AcademicYear) error {
	query, args, err := q.Insert(ctx).
		Into("academic_year_info").
		Columns("id").
		Values(academicYear.ID).
		OnConflictDoNothing().
		ToSQL()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *AcademicYearRepository) DeleteAcademicYear(ctx context.Context, id uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("academic_year_info").
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
