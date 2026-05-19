package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type SpecialtyRepository struct {
	db sqlx.ExtContext
}

func NewSpecialtyRepository(db sqlx.ExtContext) *SpecialtyRepository {
	return &SpecialtyRepository{
		db: db,
	}
}

func (r *SpecialtyRepository) CreateSpecialties(ctx context.Context, specialties domain.Specialties) error {
	builder := q.Insert(ctx).
		Into("specialty").
		Columns(`
			department_id,
			name,
			short_name,
			specialty_code,
			created_at
		`)

	for _, sp := range specialties {
		builder = builder.Values(
			sp.DepartmentID,
			sp.Name,
			sp.ShortName,
			sp.SpecialtyCode,
			sp.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&specialties[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *SpecialtyRepository) GetSpecialtyByID(ctx context.Context, id uint64) (*domain.Specialty, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			s.id,
			sa.academic_year_id,
			s.department_id,
			s.name,
			s.short_name,
			s.specialty_code,
			s.created_at,
			s.modified_at
		`).
		From("specialty s").
		InnerJoin("academic_year_to_specialty AS sa ON sa.specialty_id = s.id").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	specialty := &domain.Specialty{}
	if err = sqlx.GetContext(ctx, r.db, specialty, query, args...); err != nil {
		return nil, err
	}

	return specialty, nil
}

func (r *SpecialtyRepository) FetchSpecialties(ctx context.Context, filters f.Filters) (domain.Specialties, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			s.id,
			sa.academic_year_id,
			s.department_id,
			s.name,
			s.short_name,
			s.specialty_code,
			s.created_at,
			s.modified_at
		`).
		From("specialty s").
		InnerJoin("academic_year_to_specialty AS sa ON sa.specialty_id = s.id").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "s.id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "sa.academic_year_id", Operator: q.Equals},
					{Name: domain.DepartmentIDParam, Column: "s.department_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("s.id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	specialties := make(domain.Specialties, 0)
	if err = sqlx.SelectContext(ctx, r.db, &specialties, query, args...); err != nil {
		return nil, err
	}

	return specialties, nil
}

func (r *SpecialtyRepository) UpdateSpecialties(ctx context.Context, specialities domain.Specialties) error {
	builder := q.BatchUpdate(ctx).
		Table("specialty").
		NumberOfRows(len(specialities)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddBigintColumn("department_id").
		AddVarCharColumn("name").
		AddVarCharColumn("short_name").
		AddVarCharColumn("specialty_code").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 6*len(specialities))
	for _, sp := range specialities {
		args = append(args,
			sp.ID,
			sp.DepartmentID,
			sp.Name,
			sp.ShortName,
			sp.SpecialtyCode,
			sp.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&specialities[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *SpecialtyRepository) DeleteSpecialties(ctx context.Context, ids []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("specialty").
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

func (r *SpecialtyRepository) AttachSpecialtiesToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Specialties) error {
	builder := q.Insert(ctx).
		Into("academic_year_to_specialty").
		Columns(`
			academic_year_id,
			specialty_id
		`)

	for academicYearID, specialties := range toAttach {
		for _, sp := range specialties {
			builder = builder.Values(
				academicYearID,
				sp.ID,
			)
		}
	}

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
