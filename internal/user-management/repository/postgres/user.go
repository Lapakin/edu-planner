package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type UserRepository struct {
	db sqlx.ExtContext
}

func NewUserRepository(db sqlx.ExtContext) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUsers(ctx context.Context, users domain.Users) error {
	builder := q.Insert(ctx).
		Into(`"user"`).
		Columns(`
			email,
			role,
			created_at
		`)

	for _, u := range users {
		builder = builder.Values(
			u.Email,
			u.Role,
			u.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&users[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *UserRepository) CreateUserProfiles(ctx context.Context, users domain.Users) error {
	builder := q.Insert(ctx).
		Into("user_profile").
		Columns(`
			user_id,
			first_name,
			last_name,
			patronymic,
			created_at
		`)

	for _, u := range users {
		builder = builder.Values(
			u.ID,
			u.FirstName,
			u.LastName,
			u.Patronymic,
			u.CreatedAt,
		)
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

func (r *UserRepository) GetUserByID(ctx context.Context, userID uint64) (*domain.User, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			u.id,
			u.email,
			u.role,
			up.first_name,
			up.last_name,
			up.patronymic,
			u.is_active,
			u.is_deleted,
			u.created_at,
			u.modified_at
		`).
		From(`"user" u`).
		LeftJoin("user_profile up ON u.id = up.user_id").
		Where("u.id = ?", userID).
		Where("u.is_deleted = ?", false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	user := &domain.User{}
	if err = sqlx.GetContext(ctx, r.db, user, query, args...); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FetchUsers(ctx context.Context, filters f.Filters) (domain.Users, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			u.id,
			u.email,
			u.role,
			up.first_name,
			up.last_name,
			up.patronymic,
			u.is_active,
			u.is_deleted,
			u.created_at,
			u.modified_at
		`).
		From(`"user" u`).
		LeftJoin("user_profile up ON u.id = up.user_id").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.RoleParam, Column: "u.role", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("u.id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	users := make(domain.Users, 0)
	if err = sqlx.SelectContext(ctx, r.db, &users, query, args...); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUsers(ctx context.Context, users domain.Users) error {
	builder := q.BatchUpdate(ctx).
		Table(`"user"`).
		NumberOfRows(len(users)).
		PrimaryKey("id").
		AddBigintColumn("id").
		AddVarCharColumn("email").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 3*len(users))
	for _, u := range users {
		args = append(args,
			u.ID,
			u.Email,
			u.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&users[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *UserRepository) UpdateUserProfiles(ctx context.Context, users domain.Users) error {
	builder := q.BatchUpdate(ctx).
		Table("user_profile").
		NumberOfRows(len(users)).
		PrimaryKey("user_id").
		AddBigintColumn("user_id").
		AddVarCharColumn("first_name").
		AddVarCharColumn("last_name").
		AddVarCharColumn("patronymic").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 5*len(users))
	for _, u := range users {
		args = append(args,
			u.ID,
			u.FirstName,
			u.LastName,
			u.Patronymic,
			u.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&users[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *UserRepository) DeleteUsers(ctx context.Context, userIDs []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table(`"user"`).
		WhereID(userIDs).
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

func (r *UserRepository) ActivateUser(ctx context.Context, userID uint64, currentTime time.Time) error {
	query, args, err := q.Update(ctx).
		Table(`"user"`).
		Set("is_active", true).
		Set("modified_at", currentTime).
		WhereID(userID).
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

func (r *UserRepository) DeactivateUser(ctx context.Context, userID uint64, currentTime time.Time) error {
	query, args, err := q.Update(ctx).
		Table(`"user"`).
		Set("is_active", false).
		Set("modified_at", currentTime).
		WhereID(userID).
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
