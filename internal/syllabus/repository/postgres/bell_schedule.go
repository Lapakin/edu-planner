package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type BellScheduleRepository struct {
	db sqlx.ExtContext
}

func NewBellScheduleRepository(db sqlx.ExtContext) *BellScheduleRepository {
	return &BellScheduleRepository{
		db: db,
	}
}

func (r BellScheduleRepository) CreateBellSchedules(ctx context.Context, schedules domain.BellSchedules) error {
	builder := q.Insert(ctx).
		Into("bell_schedule").
		Columns(`
			academic_year_id,
			lesson_number,
			start_time,
			end_time
		`)

	for _, s := range schedules {
		builder = builder.Values(
			s.AcademicYearID,
			s.LessonNumber,
			s.StartTime,
			s.EndTime,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&schedules[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r BellScheduleRepository) GetBellScheduleByID(ctx context.Context, id uint64) (*domain.BellSchedule, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			lesson_number,
			start_time,
			end_time
		`).
		From("bell_schedule").
		WhereID(id).
		ToSQL()
	if err != nil {
		return nil, err
	}

	bell := &domain.BellSchedule{}
	if err = sqlx.GetContext(ctx, r.db, bell, query, args...); err != nil {
		return nil, err
	}

	return bell, nil
}

func (r BellScheduleRepository) FetchBellSchedules(ctx context.Context, filters f.Filters) (domain.BellSchedules, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			lesson_number,
			start_time,
			end_time
		`).
		From("bell_schedule").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "academic_year_id", Operator: q.Equals},
				},
			},
		}).
		OrderBy("lesson_number").
		ToSQL()
	if err != nil {
		return nil, err
	}

	bells := make(domain.BellSchedules, 0)
	if err = sqlx.SelectContext(ctx, r.db, &bells, query, args...); err != nil {
		return nil, err
	}

	return bells, nil
}

func (r BellScheduleRepository) UpdateBellSchedules(ctx context.Context, schedules domain.BellSchedules) error {
	builder := q.BatchUpdate(ctx).
		Table("bell_schedule").
		NumberOfRows(len(schedules)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddBigintColumn("academic_year_id").
		AddIntColumn("lesson_number").
		AddTimeColumn("start_time").
		AddTimeColumn("end_time")

	args := make([]any, 0, 5*len(schedules))
	for _, s := range schedules {
		args = append(args,
			s.ID,
			s.AcademicYearID,
			s.LessonNumber,
			s.StartTime,
			s.EndTime,
		)
	}

	result, err := builder.ExecContext(r.db, args)
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

func (r BellScheduleRepository) DeleteBellSchedules(ctx context.Context, ids []uint64) error {
	builder := q.Delete(ctx).
		From("bell_schedule").
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
