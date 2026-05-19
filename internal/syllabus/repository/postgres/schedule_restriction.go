package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type ScheduleRestrictionRepository struct {
	db sqlx.ExtContext
}

func NewScheduleRestrictionRepository(db sqlx.ExtContext) *ScheduleRestrictionRepository {
	return &ScheduleRestrictionRepository{db: db}
}

func (r ScheduleRestrictionRepository) CreateScheduleRestrictions(ctx context.Context, restrictions domain.ScheduleRestrictions) error {
	builder := q.Insert(ctx).
		Into("schedule_restriction").
		Columns(`
			min_group_lessons_per_day,
			max_group_lessons_per_day,
			max_teacher_lessons_per_day,
			no_gaps_in_group_schedule,
			created_at
		`)

	for _, res := range restrictions {
		builder = builder.Values(
			res.MinGroupLessonsPerDay,
			res.MaxGroupLessonsPerDay,
			res.MaxTeacherLessonsPerDay,
			res.NoGapsInGroupSchedule,
			res.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&restrictions[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r ScheduleRestrictionRepository) GetScheduleRestrictionByID(ctx context.Context, id uint64) (*domain.ScheduleRestriction, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			min_group_lessons_per_day,
			max_group_lessons_per_day,
			max_teacher_lessons_per_day,
			no_gaps_in_group_schedule,
			created_at,
			modified_at
		`).
		From("schedule_restriction").
		WhereID(id).
		ToSQL()
	if err != nil {
		return nil, err
	}

	res := &domain.ScheduleRestriction{}
	if err = sqlx.GetContext(ctx, r.db, res, query, args...); err != nil {
		return nil, err
	}

	return res, nil
}

func (r ScheduleRestrictionRepository) FetchScheduleRestrictions(ctx context.Context, filters f.Filters) (domain.ScheduleRestrictions, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			min_group_lessons_per_day,
			max_group_lessons_per_day,
			max_teacher_lessons_per_day,
			no_gaps_in_group_schedule,
			created_at,
			modified_at
		`).
		From("schedule_restriction").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
				},
			},
		}).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	restrictions := make(domain.ScheduleRestrictions, 0)
	if err = sqlx.SelectContext(ctx, r.db, &restrictions, query, args...); err != nil {
		return nil, err
	}

	return restrictions, nil
}

func (r ScheduleRestrictionRepository) UpdateScheduleRestrictions(ctx context.Context, restrictions domain.ScheduleRestrictions) error {
	builder := q.BatchUpdate(ctx).
		Table("schedule_restriction").
		NumberOfRows(len(restrictions)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddIntColumn("min_group_lessons_per_day").
		AddIntColumn("max_group_lessons_per_day").
		AddIntColumn("max_teacher_lessons_per_day").
		AddBoolColumn("no_gaps_in_group_schedule").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 6*len(restrictions))
	for _, res := range restrictions {
		args = append(args,
			res.ID,
			res.MinGroupLessonsPerDay,
			res.MaxGroupLessonsPerDay,
			res.MaxTeacherLessonsPerDay,
			res.NoGapsInGroupSchedule,
			res.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&restrictions[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r ScheduleRestrictionRepository) DeleteScheduleRestrictions(ctx context.Context, ids []uint64) error {
	builder := q.Delete(ctx).
		From("schedule_restriction").
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
