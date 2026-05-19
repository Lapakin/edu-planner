package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type CycleCommitteeRepository struct {
	db sqlx.ExtContext
}

func NewCycleCommitteeRepository(db sqlx.ExtContext) *CycleCommitteeRepository {
	return &CycleCommitteeRepository{db: db}
}

func (r *CycleCommitteeRepository) CreateCycleCommittees(ctx context.Context, cycleCommittees domain.CycleCommittees) error {
	builder := q.Insert(ctx).
		Into("cycle_committee").
		Columns(`
			user_id,
			name,
			created_at
		`)

	for _, cc := range cycleCommittees {
		builder = builder.Values(
			cc.UserID,
			cc.Name,
			cc.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&cycleCommittees[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *CycleCommitteeRepository) GetCycleCommitteeByID(ctx context.Context, id uint64) (*domain.CycleCommittee, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			cc.id,
			ca.academic_year_id,
			cc.user_id,
			cc.name,
			cc.created_at,
			cc.modified_at
		`).
		From("cycle_committee cc").
		InnerJoin("academic_year_to_cycle_committee AS ca ON ca.cycle_committee_id = cc.id").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	cycleCommittee := &domain.CycleCommittee{}
	if err = sqlx.GetContext(ctx, r.db, cycleCommittee, query, args...); err != nil {
		return nil, err
	}

	return cycleCommittee, nil
}

func (r *CycleCommitteeRepository) FetchCycleCommittees(ctx context.Context, filters f.Filters) (domain.CycleCommittees, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			cc.id,
			ca.academic_year_id,
			cc.user_id,
			cc.name,
			cc.created_at,
			cc.modified_at
		`).
		From("cycle_committee cc").
		InnerJoin("academic_year_to_cycle_committee AS ca ON ca.cycle_committee_id = cc.id").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "cc.id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "ca.academic_year_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("cc.id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	cycleCommittees := make(domain.CycleCommittees, 0)
	if err = sqlx.SelectContext(ctx, r.db, &cycleCommittees, query, args...); err != nil {
		return nil, err
	}

	return cycleCommittees, nil
}

func (r *CycleCommitteeRepository) UpdateCycleCommittees(ctx context.Context, cycleCommittees domain.CycleCommittees) error {
	builder := q.BatchUpdate(ctx).
		Table("cycle_committee").
		PrimaryKey("id").
		NumberOfRows(len(cycleCommittees)).
		AddBigintColumn("id").
		AddBigintColumn("user_id").
		AddVarCharColumn("name").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 4*len(cycleCommittees))
	for _, cc := range cycleCommittees {
		args = append(args,
			cc.ID,
			cc.UserID,
			cc.Name,
			cc.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&cycleCommittees[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *CycleCommitteeRepository) DeleteCycleCommittees(ctx context.Context, ids []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("cycle_committee").
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

func (r *CycleCommitteeRepository) AttachCycleCommitteesToAcademicYear(ctx context.Context, toAttach map[uint64]domain.CycleCommittees) error {
	builder := q.Insert(ctx).
		Into("academic_year_to_cycle_committee").
		Columns(`
			academic_year_id,
			cycle_committee_id
		`)

	for academicYearID, cycleCommittees := range toAttach {
		for _, cc := range cycleCommittees {
			builder = builder.Values(
				academicYearID,
				cc.ID,
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
