package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
)

type UserRepository struct {
	db sqlx.ExtContext
}

func NewUserRepository(db sqlx.ExtContext) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UpsertUser(ctx context.Context, user *domain.User) error {
	query, args, err := q.Insert(ctx).
		Into("user_info").
		Columns("id").
		Values(user.ID).
		OnConflictDoNothing().
		ToSQL()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("user_info").
		WhereID(id).
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
