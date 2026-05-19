package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type RoomRepository struct {
	db sqlx.ExtContext
}

func NewRoomRepository(db sqlx.ExtContext) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r RoomRepository) CreateRooms(ctx context.Context, rooms domain.Rooms) error {
	builder := q.Insert(ctx).
		Into("room").
		Columns(`
			name,
			room_type,
			created_at
		`)

	for _, rm := range rooms {
		builder = builder.Values(
			rm.Name,
			rm.RoomType,
			rm.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&rooms[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r RoomRepository) GetRoomByID(ctx context.Context, id uint64) (*domain.Room, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			r.id,
			ra.academic_year_id,
			r.name,
			r.room_type,
			r.created_at,
			r.modified_at
		`).
		From("room r").
		InnerJoin("academic_year_to_room AS ra ON ra.room_id = r.id").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	room := &domain.Room{}
	if err = sqlx.GetContext(ctx, r.db, room, query, args...); err != nil {
		return nil, err
	}

	return room, nil
}

func (r RoomRepository) FetchRooms(ctx context.Context, filters f.Filters) (domain.Rooms, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			r.id,
			ra.academic_year_id,
			r.name,
			r.room_type,
			r.created_at,
			r.modified_at
		`).
		From("room r").
		InnerJoin("academic_year_to_room AS ra ON ra.room_id = r.id").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "r.id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "ra.academic_year_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("r.id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	rooms := make(domain.Rooms, 0)
	if err = sqlx.SelectContext(ctx, r.db, &rooms, query, args...); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r RoomRepository) UpdateRooms(ctx context.Context, rooms domain.Rooms) error {
	builder := q.BatchUpdate(ctx).
		Table("room").
		NumberOfRows(len(rooms)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddVarCharColumn("name").
		AddVarCharColumn("room_type").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 4*len(rooms))
	for _, rm := range rooms {
		args = append(args,
			rm.ID,
			rm.Name,
			rm.RoomType,
			rm.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&rooms[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r RoomRepository) DeleteRooms(ctx context.Context, ids []uint64) error {
	builder := q.SoftDelete(ctx).
		Table("room").
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

func (r RoomRepository) AttachRoomsToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Rooms) error {
	builder := q.Insert(ctx).
		Into("academic_year_to_room").
		Columns(`
			academic_year_id,
			room_id
		`)

	for academicYearID, rooms := range toAttach {
		for _, rm := range rooms {
			builder = builder.Values(
				academicYearID,
				rm.ID,
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
