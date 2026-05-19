package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type StudyPlanRepository struct {
	db sqlx.ExtContext
}

func NewStudyPlanRepository(db sqlx.ExtContext) *StudyPlanRepository {
	return &StudyPlanRepository{
		db: db,
	}
}

func (r StudyPlanRepository) CreateStudyPlans(ctx context.Context, studyPlans domain.StudyPlans) error {
	builder := q.Insert(ctx).
		Into("study_plan").
		Columns(`
			academic_year_id,
			specialty_id,
			discipline_id,
			semester_number,
			educational_component,
			okr,
			coursework_type,
			classroom_work,
			lectures,
			laboratory,
			practical,
			independent_work,
			individual_lessons,
			consultations,
			exam,
			credit,
			practice,
			coursework_hours,
			classroom_coursework,
			assisting_hours,
			created_at
		`)

	for _, sp := range studyPlans {
		builder = builder.Values(
			sp.AcademicYearID,
			sp.SpecialtyID,
			sp.DisciplineID,
			sp.SemesterNumber,
			sp.EducationalComponent,
			sp.OKR,
			sp.CourseworkType,
			sp.ClassroomWork,
			sp.Lectures,
			sp.Laboratory,
			sp.Practical,
			sp.IndependentWork,
			sp.IndividualLessons,
			sp.Consultations,
			sp.Exam,
			sp.Credit,
			sp.Practice,
			sp.CourseworkHours,
			sp.ClassroomCoursework,
			sp.AssistingHours,
			sp.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&studyPlans[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r StudyPlanRepository) GetStudyPlanByID(ctx context.Context, id uint64) (*domain.StudyPlan, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			specialty_id,
			discipline_id,
			semester_number,
			educational_component,
			okr,
			coursework_type,
			classroom_work,
			lectures,
			laboratory,
			practical,
			independent_work,
			individual_lessons,
			consultations,
			exam,
			credit,
			practice,
			coursework_hours,
			classroom_coursework,
			assisting_hours,
			created_at,
			modified_at
		`).
		From("study_plan").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	sp := &domain.StudyPlan{}
	if err = sqlx.GetContext(ctx, r.db, sp, query, args...); err != nil {
		return nil, err
	}

	return sp, nil
}

func (r StudyPlanRepository) FetchStudyPlans(ctx context.Context, filters f.Filters) (domain.StudyPlans, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			specialty_id,
			discipline_id,
			semester_number,
			educational_component,
			okr,
			coursework_type,
			classroom_work,
			lectures,
			laboratory,
			practical,
			independent_work,
			individual_lessons,
			consultations,
			exam,
			credit,
			practice,
			coursework_hours,
			classroom_coursework,
			assisting_hours,
			created_at,
			modified_at
		`).
		From("study_plan").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "academic_year_id", Operator: q.Equals},
					{Name: domain.SpecialtyIDParam, Column: "specialty_id", Operator: q.Equals},
					{Name: domain.DisciplineIDParam, Column: "discipline_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	sps := make(domain.StudyPlans, 0)
	if err = sqlx.SelectContext(ctx, r.db, &sps, query, args...); err != nil {
		return nil, err
	}

	return sps, nil
}

func (r StudyPlanRepository) UpdateStudyPlans(ctx context.Context, studyPlans domain.StudyPlans) error {
	builder := q.BatchUpdate(ctx).
		Table("study_plan").
		NumberOfRows(len(studyPlans)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddBigintColumn("academic_year_id").
		AddBigintColumn("specialty_id").
		AddBigintColumn("discipline_id").
		AddIntColumn("semester_number").
		AddVarCharColumn("educational_component").
		AddIntColumn("okr").
		AddVarCharColumn("coursework_type").
		AddIntColumn("classroom_work").
		AddIntColumn("lectures").
		AddIntColumn("laboratory").
		AddIntColumn("practical").
		AddIntColumn("independent_work").
		AddIntColumn("individual_lessons").
		AddIntColumn("consultations").
		AddIntColumn("exam").
		AddIntColumn("credit").
		AddIntColumn("practice").
		AddIntColumn("coursework_hours").
		AddIntColumn("classroom_coursework").
		AddIntColumn("assisting_hours").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 22*len(studyPlans))
	for _, sp := range studyPlans {
		args = append(args,
			sp.ID,
			sp.AcademicYearID,
			sp.SpecialtyID,
			sp.DisciplineID,
			sp.SemesterNumber,
			sp.EducationalComponent,
			sp.OKR,
			sp.CourseworkType,
			sp.ClassroomWork,
			sp.Lectures,
			sp.Laboratory,
			sp.Practical,
			sp.IndependentWork,
			sp.IndividualLessons,
			sp.Consultations,
			sp.Exam,
			sp.Credit,
			sp.Practice,
			sp.CourseworkHours,
			sp.ClassroomCoursework,
			sp.AssistingHours,
			sp.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&studyPlans[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r StudyPlanRepository) DeleteStudyPlans(ctx context.Context, ids []uint64) error {
	builder := q.SoftDelete(ctx).
		Table("study_plan").
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
