package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type CycleCommitteeLabRoomRepository struct {
	db sqlx.ExtContext
}

func NewCycleCommitteeLabRoomRepository(db sqlx.ExtContext) *CycleCommitteeLabRoomRepository {
	return &CycleCommitteeLabRoomRepository{db: db}
}

func (r CycleCommitteeLabRoomRepository) CreateCycleCommitteeLabRooms(ctx context.Context, labRooms domain.CycleCommitteeLabRooms) error {
	builder := q.Insert(ctx).
		Into("cycle_committee_lab_room").
		Columns(`
			academic_year_id,
			cycle_committee_id,
			room_id,
			created_at
		`)

	for _, lr := range labRooms {
		builder = builder.Values(
			lr.AcademicYearID,
			lr.CycleCommitteeID,
			lr.RoomID,
			lr.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&labRooms[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r CycleCommitteeLabRoomRepository) GetCycleCommitteeLabRoomByID(ctx context.Context, id uint64) (*domain.CycleCommitteeLabRoom, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			cycle_committee_id,
			room_id,
			created_at,
			modified_at
		`).
		From("cycle_committee_lab_room").
		WhereID(id).
		ToSQL()
	if err != nil {
		return nil, err
	}

	lr := &domain.CycleCommitteeLabRoom{}
	if err = sqlx.GetContext(ctx, r.db, lr, query, args...); err != nil {
		return nil, err
	}

	return lr, nil
}

func (r CycleCommitteeLabRoomRepository) FetchCycleCommitteeLabRooms(ctx context.Context, filters f.Filters) (domain.CycleCommitteeLabRooms, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			academic_year_id,
			cycle_committee_id,
			room_id,
			created_at,
			modified_at
		`).
		From("cycle_committee_lab_room").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "academic_year_id", Operator: q.Equals},
					{Name: domain.CycleCommitteeIDParam, Column: "cycle_committee_id", Operator: q.Equals},
				},
			},
		}).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	labRooms := make(domain.CycleCommitteeLabRooms, 0)
	if err = sqlx.SelectContext(ctx, r.db, &labRooms, query, args...); err != nil {
		return nil, err
	}

	return labRooms, nil
}

func (r CycleCommitteeLabRoomRepository) UpdateCycleCommitteeLabRooms(ctx context.Context, labRooms domain.CycleCommitteeLabRooms) error {
	builder := q.BatchUpdate(ctx).
		Table("cycle_committee_lab_room").
		NumberOfRows(len(labRooms)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddBigintColumn("academic_year_id").
		AddBigintColumn("cycle_committee_id").
		AddBigintColumn("room_id").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 5*len(labRooms))
	for _, lr := range labRooms {
		args = append(args,
			lr.ID,
			lr.AcademicYearID,
			lr.CycleCommitteeID,
			lr.RoomID,
			lr.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&labRooms[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r CycleCommitteeLabRoomRepository) DeleteCycleCommitteeLabRooms(ctx context.Context, ids []uint64) error {
	builder := q.Delete(ctx).
		From("cycle_committee_lab_room").
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
