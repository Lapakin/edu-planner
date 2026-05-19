package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type WorkloadDistributionRepository struct {
	db sqlx.ExtContext
}

func NewWorkloadDistributionRepository(db sqlx.ExtContext) *WorkloadDistributionRepository {
	return &WorkloadDistributionRepository{
		db: db,
	}
}

func (r WorkloadDistributionRepository) CreateWorkloadDistributions(ctx context.Context, distributions domain.WorkloadDistributions) error {
	builder := q.Insert(ctx).
		Into("workload_distribution").
		Columns(`
			study_plan_id,
			group_id,
			coursework_type,
			coursework_hours,
			classroom_coursework,
			classroom_work,
			laboratory,
			practical,
			exam,
			consultations,
			credit,
			practice,
			okr,
			individual_lessons,
			deductions,
			created_at
		`)

	for _, wd := range distributions {
		builder = builder.Values(
			wd.StudyPlanID,
			wd.GroupID,
			wd.CourseworkType,
			wd.CourseworkHours,
			wd.ClassroomCoursework,
			wd.ClassroomWork,
			wd.Laboratory,
			wd.Practical,
			wd.Exam,
			wd.Consultations,
			wd.Credit,
			wd.Practice,
			wd.OKR,
			wd.IndividualLessons,
			wd.Deductions,
			wd.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&distributions[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r WorkloadDistributionRepository) GetWorkloadDistributionByID(ctx context.Context, id uint64) (*domain.WorkloadDistribution, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			study_plan_id,
			group_id,
			coursework_type,
			coursework_hours,
			classroom_coursework,
			classroom_work,
			laboratory,
			practical,
			exam,
			consultations,
			credit,
			practice,
			okr,
			individual_lessons,
			deductions,
			created_at,
			modified_at
		`).
		From("workload_distribution").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	wd := &domain.WorkloadDistribution{}
	if err = sqlx.GetContext(ctx, r.db, wd, query, args...); err != nil {
		return nil, err
	}

	return wd, nil
}

func (r WorkloadDistributionRepository) FetchWorkloadDistributions(ctx context.Context, filters f.Filters) (domain.WorkloadDistributions, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			study_plan_id,
			group_id,
			coursework_type,
			coursework_hours,
			classroom_coursework,
			classroom_work,
			laboratory,
			practical,
			exam,
			consultations,
			credit,
			practice,
			okr,
			individual_lessons,
			deductions,
			created_at,
			modified_at
		`).
		From("workload_distribution").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
					{Name: "study_plan_id", Column: "study_plan_id", Operator: q.Equals},
					{Name: "group_id", Column: "group_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	wds := make(domain.WorkloadDistributions, 0)
	if err = sqlx.SelectContext(ctx, r.db, &wds, query, args...); err != nil {
		return nil, err
	}

	return wds, nil
}

func (r WorkloadDistributionRepository) UpdateWorkloadDistributions(ctx context.Context, distributions domain.WorkloadDistributions) error {
	builder := q.BatchUpdate(ctx).
		Table("workload_distribution").
		NumberOfRows(len(distributions)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddBigintColumn("study_plan_id").
		AddBigintColumn("group_id").
		AddVarCharColumn("coursework_type").
		AddIntColumn("coursework_hours").
		AddIntColumn("classroom_coursework").
		AddIntColumn("classroom_work").
		AddIntColumn("laboratory").
		AddIntColumn("practical").
		AddIntColumn("exam").
		AddIntColumn("consultations").
		AddIntColumn("credit").
		AddIntColumn("practice").
		AddIntColumn("okr").
		AddIntColumn("individual_lessons").
		AddIntColumn("deductions").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 17*len(distributions))
	for _, wd := range distributions {
		args = append(args,
			wd.ID,
			wd.StudyPlanID,
			wd.GroupID,
			wd.CourseworkType,
			wd.CourseworkHours,
			wd.ClassroomCoursework,
			wd.ClassroomWork,
			wd.Laboratory,
			wd.Practical,
			wd.Exam,
			wd.Consultations,
			wd.Credit,
			wd.Practice,
			wd.OKR,
			wd.IndividualLessons,
			wd.Deductions,
			wd.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&distributions[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r WorkloadDistributionRepository) DeleteWorkloadDistributions(ctx context.Context, ids []uint64) error {
	builder := q.SoftDelete(ctx).
		Table("workload_distribution").
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
