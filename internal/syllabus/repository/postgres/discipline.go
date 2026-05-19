package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type DisciplineRepository struct {
	db sqlx.ExtContext
}

func NewDisciplineRepository(db sqlx.ExtContext) *DisciplineRepository {
	return &DisciplineRepository{
		db: db,
	}
}

func (r *DisciplineRepository) CreateDisciplines(ctx context.Context, disciplines domain.Disciplines) error {
	builder := q.Insert(ctx).
		Into("discipline").
		Columns(`
			cycle_committee_id,
			name,
			short_name,
			is_splitting,
			is_subvention,
			created_at
		`)

	for _, d := range disciplines {
		builder = builder.Values(
			d.CycleCommitteeID,
			d.Name,
			d.ShortName,
			d.IsSplitting,
			d.IsSubvention,
			d.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&disciplines[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *DisciplineRepository) GetDisciplineByID(ctx context.Context, id uint64) (*domain.Discipline, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			d.id,
			da.academic_year_id,
			d.cycle_committee_id,
			d.name,
			d.short_name,
			d.is_splitting,	
			d.is_subvention,
			d.created_at,
			d.modified_at
		`).
		From("discipline d").
		InnerJoin("academic_year_to_discipline AS da ON da.discipline_id = d.id").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	discipline := &domain.Discipline{}
	if err = sqlx.GetContext(ctx, r.db, discipline, query, args...); err != nil {
		return nil, err
	}

	return discipline, nil
}

func (r *DisciplineRepository) FetchDisciplines(ctx context.Context, filters f.Filters) (domain.Disciplines, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			d.id,
			da.academic_year_id,
			d.cycle_committee_id,
			d.name,
			d.short_name,
			d.is_splitting,
			d.is_subvention,
			d.created_at,
			d.modified_at
		`).
		From("discipline d").
		InnerJoin("academic_year_to_discipline AS da ON da.discipline_id = d.id").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "d.id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "da.academic_year_id", Operator: q.Equals},
					{Name: domain.CycleCommitteeIDParam, Column: "d.cycle_committee_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("d.id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	disciplines := make(domain.Disciplines, 0)
	if err = sqlx.SelectContext(ctx, r.db, &disciplines, query, args...); err != nil {
		return nil, err
	}

	return disciplines, nil
}

func (r *DisciplineRepository) UpdateDisciplines(ctx context.Context, disciplines domain.Disciplines) error {
	builder := q.BatchUpdate(ctx).
		Table("discipline").
		PrimaryKey("id").
		NumberOfRows(len(disciplines)).
		AddBigintColumn("id").
		AddBigintColumn("cycle_committee_id").
		AddVarCharColumn("name").
		AddVarCharColumn("short_name").
		AddBoolColumn("is_splitting").
		AddBoolColumn("is_subvention").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 7*len(disciplines))
	for _, d := range disciplines {
		args = append(args,
			d.ID,
			d.CycleCommitteeID,
			d.Name,
			d.ShortName,
			d.IsSplitting,
			d.IsSubvention,
			d.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&disciplines[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *DisciplineRepository) DeleteDisciplines(ctx context.Context, ids []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("discipline").
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

func (r *DisciplineRepository) AttachDisciplinesToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Disciplines) error {
	builder := q.Insert(ctx).
		Into("academic_year_to_discipline").
		Columns(`
			academic_year_id,
			discipline_id
		`)

	for academicYearID, disciplines := range toAttach {
		for _, d := range disciplines {
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
