package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
)

type MessageWriter struct {
	db sqlx.ExtContext
}

func NewMessageWriter(db sqlx.ExtContext) *MessageWriter {
	return &MessageWriter{db: db}
}

func (r *MessageWriter) Write(ctx context.Context, messages domain.Massages) error {
	builder := q.Insert(ctx).
		Into("massage").
		Columns(`
			event_id,
			action_type,
			object_type,
			object
		`)

	for _, m := range messages {
		builder = builder.Values(
			m.EventID,
			m.ActionType,
			m.ObjectType,
			m.Object,
		)
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
