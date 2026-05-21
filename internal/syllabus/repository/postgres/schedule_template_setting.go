package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type ScheduleTemplateSettingRepository struct {
	db sqlx.ExtContext
}

func NewScheduleTemplateSettingRepository(db sqlx.ExtContext) *ScheduleTemplateSettingRepository {
	return &ScheduleTemplateSettingRepository{
		db: db,
	}
}

func (r ScheduleTemplateSettingRepository) CreateScheduleTemplateSettings(ctx context.Context, settings domain.ScheduleTemplateSettings) error {
	builder := q.Insert(ctx).
		Into("schedule_template_setting").
		Columns(`
			academic_year_id,
			lessons_per_class,
			study_days_mask,
			max_identical_lessons_per_day,
			max_study_hours_per_day,
			max_teacher_hours_per_week,
			max_group_lesson_hours_per_week,
			created_at
		`)

	for _, s := range settings {
		builder = builder.Values(
			s.AcademicYearID,
			s.LessonsPerClass,
			s.StudyDaysMask,
			s.MaxIdenticalLessonsPerDay,
			s.MaxStudyHoursPerDay,
			s.MaxTeacherHoursPerWeek,
			s.MaxGroupLessonHoursPerWeek,
			s.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&settings[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r ScheduleTemplateSettingRepository) GetScheduleTemplateSettingByID(ctx context.Context, id uint64) (*domain.ScheduleTemplateSetting, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			lessons_per_class,
			study_days_mask,
			max_identical_lessons_per_day,
			max_study_hours_per_day,
			max_teacher_hours_per_week,
			max_group_lesson_hours_per_week,
			created_at,
			modified_at
		`).
		From("schedule_template_setting").
		WhereID(id).
		ToSQL()
	if err != nil {
		return nil, err
	}

	s := &domain.ScheduleTemplateSetting{}
	if err = sqlx.GetContext(ctx, r.db, s, query, args...); err != nil {
		return nil, err
	}

	return s, nil
}

func (r ScheduleTemplateSettingRepository) FetchScheduleTemplateSettings(ctx context.Context, filters f.Filters) (domain.ScheduleTemplateSettings, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			lessons_per_class,
			study_days_mask,
			max_identical_lessons_per_day,
			max_study_hours_per_day,
			max_teacher_hours_per_week,
			max_group_lesson_hours_per_week,
			created_at,
			modified_at
		`).
		From("schedule_template_setting").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "academic_year_id", Operator: q.Equals},
				},
			},
		}).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	settings := make(domain.ScheduleTemplateSettings, 0)
	if err = sqlx.SelectContext(ctx, r.db, &settings, query, args...); err != nil {
		return nil, err
	}

	return settings, nil
}

func (r ScheduleTemplateSettingRepository) UpdateScheduleTemplateSettings(ctx context.Context, settings domain.ScheduleTemplateSettings) error {
	builder := q.BatchUpdate(ctx).
		Table("schedule_template_setting").
		NumberOfRows(len(settings)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddBigintColumn("academic_year_id").
		AddIntColumn("lessons_per_class").
		AddIntColumn("study_days_mask").
		AddIntColumn("max_identical_lessons_per_day").
		AddIntColumn("max_study_hours_per_day").
		AddIntColumn("max_teacher_hours_per_week").
		AddIntColumn("max_group_lesson_hours_per_week").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 9*len(settings))
	for _, s := range settings {
		args = append(args,
			s.ID,
			s.AcademicYearID,
			s.LessonsPerClass,
			int(s.StudyDaysMask),
			s.MaxIdenticalLessonsPerDay,
			s.MaxStudyHoursPerDay,
			s.MaxTeacherHoursPerWeek,
			s.MaxGroupLessonHoursPerWeek,
			s.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&settings[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r ScheduleTemplateSettingRepository) DeleteScheduleTemplateSettings(ctx context.Context, ids []uint64) error {
	builder := q.Delete(ctx).
		From("schedule_template_setting").
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
