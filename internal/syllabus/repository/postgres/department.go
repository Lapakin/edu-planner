package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type DepartmentRepository struct {
	db sqlx.ExtContext
}

func NewDepartmentRepository(db sqlx.ExtContext) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) CreateDepartments(ctx context.Context, departments domain.Departments) error {
	builder := q.Insert(ctx).
		Into("department").
		Columns(`
			name,
			created_at
		`)

	for _, d := range departments {
		builder = builder.Values(
			d.Name,
			d.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&departments[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *DepartmentRepository) GetDepartmentByID(ctx context.Context, id uint64) (*domain.Department, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			d.id,
			da.academic_year_id,
			d.name,
			d.created_at,
			d.modified_at
		`).
		From("department d").
		InnerJoin("academic_year_to_department AS da ON da.department_id = d.id").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	department := &domain.Department{}
	if err = sqlx.GetContext(ctx, r.db, department, query, args...); err != nil {
		return nil, err
	}

	return department, nil
}

func (r *DepartmentRepository) FetchDepartments(ctx context.Context, filters f.Filters) (domain.Departments, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			d.id,
			da.academic_year_id,
			d.name,
			d.created_at,
			d.modified_at
		`).
		From("department d").
		InnerJoin("academic_year_to_department AS da ON da.department_id = d.id").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "d.id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "da.academic_year_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("d.id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	departments := make(domain.Departments, 0)
	if err = sqlx.SelectContext(ctx, r.db, &departments, query, args...); err != nil {
		return nil, err
	}

	return departments, nil
}

func (r *DepartmentRepository) UpdateDepartments(ctx context.Context, departments domain.Departments) error {
	builder := q.BatchUpdate(ctx).
		Table("department").
		NumberOfRows(len(departments)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddVarCharColumn("name").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 3*len(departments))
	for _, d := range departments {
		args = append(args,
			d.ID,
			d.Name,
			d.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&departments[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *DepartmentRepository) DeleteDepartments(ctx context.Context, ids []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("department").
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

func (r *DepartmentRepository) AttachDepartmentsToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Departments) error {
	builder := q.Insert(ctx).
		Into("academic_year_to_department").
		Columns(`
			academic_year_id,
			department_id
		`)

	for academicYearID, departments := range toAttach {
		for _, d := range departments {
			builder = builder.Values(
				academicYearID,
				d.ID,
			)
		}
	}

	query, args, err := builder.ToSQL()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
