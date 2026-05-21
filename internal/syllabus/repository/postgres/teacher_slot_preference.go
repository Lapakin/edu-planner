package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type TeacherSlotPreferenceRepository struct {
	db sqlx.ExtContext
}

func NewTeacherSlotPreferenceRepository(db sqlx.ExtContext) *TeacherSlotPreferenceRepository {
	return &TeacherSlotPreferenceRepository{db: db}
}

func (r TeacherSlotPreferenceRepository) CreateTeacherSlotPreferences(ctx context.Context, preferences domain.TeacherSlotPreferences) error {
	builder := q.Insert(ctx).
		Into("teacher_slot_preference").
		Columns(`
			academic_year_id,
			teacher_id,
			weekday,
			lesson_number,
			slot_type,
			created_at
		`)

	for _, p := range preferences {
		builder = builder.Values(
			p.AcademicYearID,
			p.TeacherID,
			string(p.Weekday),
			p.LessonNumber,
			string(p.SlotType),
			p.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&preferences[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r TeacherSlotPreferenceRepository) GetTeacherSlotPreferenceByID(ctx context.Context, id uint64) (*domain.TeacherSlotPreference, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			teacher_id,
			weekday,
			lesson_number,
			slot_type,
			created_at,
			modified_at
		`).
		From("teacher_slot_preference").
		WhereID(id).
		ToSQL()
	if err != nil {
		return nil, err
	}

	p := &domain.TeacherSlotPreference{}
	if err = sqlx.GetContext(ctx, r.db, p, query, args...); err != nil {
		return nil, err
	}

	return p, nil
}

func (r TeacherSlotPreferenceRepository) FetchTeacherSlotPreferences(ctx context.Context, filters f.Filters) (domain.TeacherSlotPreferences, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			teacher_id,
			weekday,
			lesson_number,
			slot_type,
			created_at,
			modified_at
		`).
		From("teacher_slot_preference").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "academic_year_id", Operator: q.Equals},
					{Name: domain.TeacherIDParam, Column: "teacher_id", Operator: q.Equals},
				},
			},
		}).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	preferences := make(domain.TeacherSlotPreferences, 0)
	if err = sqlx.SelectContext(ctx, r.db, &preferences, query, args...); err != nil {
		return nil, err
	}

	return preferences, nil
}

func (r TeacherSlotPreferenceRepository) UpdateTeacherSlotPreferences(ctx context.Context, preferences domain.TeacherSlotPreferences) error {
	builder := q.BatchUpdate(ctx).
		Table("teacher_slot_preference").
		NumberOfRows(len(preferences)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddBigintColumn("academic_year_id").
		AddBigintColumn("teacher_id").
		AddVarCharColumn("weekday").
		AddIntColumn("lesson_number").
		AddVarCharColumn("slot_type").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 7*len(preferences))
	for _, p := range preferences {
		args = append(args,
			p.ID,
			p.AcademicYearID,
			p.TeacherID,
			string(p.Weekday),
			p.LessonNumber,
			string(p.SlotType),
			p.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&preferences[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r TeacherSlotPreferenceRepository) DeleteTeacherSlotPreferences(ctx context.Context, ids []uint64) error {
	builder := q.Delete(ctx).
		From("teacher_slot_preference").
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
