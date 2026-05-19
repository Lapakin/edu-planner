package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type WorkloadAssignmentRepository struct {
	db sqlx.ExtContext
}

func NewWorkloadAssignmentRepository(db sqlx.ExtContext) *WorkloadAssignmentRepository {
	return &WorkloadAssignmentRepository{
		db: db,
	}
}

func (r WorkloadAssignmentRepository) CreateWorkloadAssignments(ctx context.Context, assignments domain.WorkloadAssignments) error {
	builder := q.Insert(ctx).
		Into("workload_assignment").
		Columns(`
			workload_distribution_id,
			teacher_id,
			role_type,
			assigned_hours,
			created_at
		`)

	for _, wa := range assignments {
		builder = builder.Values(
			wa.WorkloadDistributionID,
			wa.TeacherID,
			wa.RoleType,
			wa.AssignedHours,
			wa.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&assignments[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r WorkloadAssignmentRepository) GetWorkloadAssignmentByID(ctx context.Context, id uint64) (*domain.WorkloadAssignment, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			workload_distribution_id,
			teacher_id,
			role_type,
			assigned_hours,
			created_at
		`).
		From("workload_assignment").
		WhereID(id).
		ToSQL()
	if err != nil {
		return nil, err
	}

	wa := &domain.WorkloadAssignment{}
	if err = sqlx.GetContext(ctx, r.db, wa, query, args...); err != nil {
		return nil, err
	}

	return wa, nil
}

func (r WorkloadAssignmentRepository) FetchWorkloadAssignments(ctx context.Context, filters f.Filters) (domain.WorkloadAssignments, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			workload_distribution_id,
			teacher_id,
			role_type,
			assigned_hours,
			created_at
		`).
		From("workload_assignment").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
					{Name: "workload_distribution_id", Column: "workload_distribution_id", Operator: q.Equals},
					{Name: "teacher_id", Column: "teacher_id", Operator: q.Equals},
				},
			},
		}).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	was := make(domain.WorkloadAssignments, 0)
	if err = sqlx.SelectContext(ctx, r.db, &was, query, args...); err != nil {
		return nil, err
	}

	return was, nil
}

func (r WorkloadAssignmentRepository) UpdateWorkloadAssignments(ctx context.Context, assignments domain.WorkloadAssignments) error {
	builder := q.BatchUpdate(ctx).
		Table("workload_assignment").
		NumberOfRows(len(assignments)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddBigintColumn("workload_distribution_id").
		AddBigintColumn("teacher_id").
		AddVarCharColumn("role_type").
		AddIntColumn("assigned_hours")

	args := make([]any, 0, 5*len(assignments))
	for _, wa := range assignments {
		args = append(args,
			wa.ID,
			wa.WorkloadDistributionID,
			wa.TeacherID,
			wa.RoleType,
			wa.AssignedHours,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&assignments[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r WorkloadAssignmentRepository) DeleteWorkloadAssignments(ctx context.Context, ids []uint64) error {
	builder := q.Delete(ctx).
		From("workload_assignment").
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
