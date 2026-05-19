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

type ScheduleTemplateRepository struct {
	db sqlx.ExtContext
}

func NewScheduleTemplateRepository(db sqlx.ExtContext) *ScheduleTemplateRepository {
	return &ScheduleTemplateRepository{
		db: db,
	}
}

func (r ScheduleTemplateRepository) CreateScheduleTemplate(ctx context.Context, template *domain.ScheduleTemplate) error {
	builder := q.Insert(ctx).
		Into("schedule_template").
		Columns(`
			semester_id,
			name,
			data,
			is_active,
			is_deleted,
			created_at
		`)

	builder = builder.Values(
		template.SemesterID,
		template.Name,
		template.Data,
		template.IsActive,
		template.IsDeleted,
		template.CreatedAt,
	)

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&template.ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r ScheduleTemplateRepository) GetScheduleTemplateByID(ctx context.Context, id uint64) (*domain.ScheduleTemplate, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			semester_id,
			name,
			data,
			is_active,
			is_deleted,
			created_at,
			modified_at
		`).
		From("schedule_template").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	st := &domain.ScheduleTemplate{}
	if err = sqlx.GetContext(ctx, r.db, st, query, args...); err != nil {
		return nil, err
	}

	return st, nil
}

func (r ScheduleTemplateRepository) FetchScheduleTemplates(ctx context.Context, filters f.Filters) (domain.ScheduleTemplates, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			semester_id,
			name,
			data,
			is_active,
			is_deleted,
			created_at,
			modified_at
		`).
		From("schedule_template").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
					{Name: domain.SemesterIDParam, Column: "semester_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	templates := make(domain.ScheduleTemplates, 0)
	if err = sqlx.SelectContext(ctx, r.db, &templates, query, args...); err != nil {
		return nil, err
	}

	return templates, nil
}

func (r ScheduleTemplateRepository) UpdateScheduleTemplate(ctx context.Context, template *domain.ScheduleTemplate) error {
	query, args, err := q.Update(ctx).
		Table("schedule_template").
		Set("semester_id", template.SemesterID).
		Set("name", template.Name).
		Set("data", template.Data).
		Set("is_active", template.IsActive).
		Set("modified_at", template.ModifiedAt).
		WhereID(template.ID).
		Returning("created_at").
		ToSQL()
	if err != nil {
		return err
	}

	return sqlx.GetContext(ctx, r.db, &template.CreatedAt, query, args...)
}

func (r ScheduleTemplateRepository) DeleteScheduleTemplates(ctx context.Context, ids []uint64) error {
	builder := q.SoftDelete(ctx).
		Table("schedule_template").
		WhereID(ids)

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

func (r ScheduleTemplateRepository) ActivateScheduleTemplate(ctx context.Context, id uint64, t time.Time) error {
	// Deactivate all templates for the same semester first, then activate this one
	// First get the template to find its semester
	template, err := r.GetScheduleTemplateByID(ctx, id)
	if err != nil {
		return err
	}

	// Deactivate all templates for this semester
	if template.SemesterID != nil {
		deactivateQuery := "UPDATE schedule_template SET is_active = false, modified_at = $1 WHERE semester_id = $2 AND is_deleted = false"
		if _, err = r.db.ExecContext(ctx, deactivateQuery, t, *template.SemesterID); err != nil {
			return err
		}
	}

	// Activate this template
	activateQuery := "UPDATE schedule_template SET is_active = true, modified_at = $1 WHERE id = $2"
	if _, err = r.db.ExecContext(ctx, activateQuery, t, id); err != nil {
		return err
	}

	return nil
}
