package query

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type DeleteBuilder struct {
	ctx     context.Context
	builder sq.DeleteBuilder
}

func HardDelete(ctx context.Context) *DeleteBuilder {
	return &DeleteBuilder{
		ctx:     ctx,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Delete(""),
	}
}

func SoftDelete(ctx context.Context) *UpdateBuilder {
	u := Update(ctx)
	u.Set(IsDeleted, true)
	return u
}

// Delete is an alias for HardDelete for backward compatibility.
func Delete(ctx context.Context) *DeleteBuilder {
	return HardDelete(ctx)
}

// From sets the table to delete from.
func (d *DeleteBuilder) From(table string) *DeleteBuilder {
	d.builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Delete(table)
	return d
}

// Where adds a WHERE clause.
func (d *DeleteBuilder) Where(pred any, args ...any) *DeleteBuilder {
	d.builder = d.builder.Where(pred, args...)
	return d
}

// WhereID adds a WHERE clause for id = value.
func (d *DeleteBuilder) WhereID(id any) *DeleteBuilder {
	d.builder = d.builder.Where(sq.Eq{ID: id})
	return d
}

// IsActive adds a WHERE clause for is_active = true.
func (d *DeleteBuilder) IsActive() *DeleteBuilder {
	d.builder = d.builder.Where(sq.Eq{IsActive: true})
	return d
}

// IsDeleted adds a WHERE clause for is_deleted filtering.
func (d *DeleteBuilder) IsDeleted(deleted bool) *DeleteBuilder {
	d.builder = d.builder.Where(sq.Eq{IsDeleted: deleted})
	return d
}

// WhereMap adds WHERE clauses from a map.
func (d *DeleteBuilder) WhereMap(clauses map[string]any) *DeleteBuilder {
	for key, value := range clauses {
		d.builder = d.builder.Where(sq.Eq{key: value})
	}
	return d
}

// Returning adds RETURNING clause.
func (d *DeleteBuilder) Returning(fields ...string) *DeleteBuilder {
	if len(fields) > 0 {
		d.builder = d.builder.Suffix("RETURNING " + fields[0])
		for idx := 1; idx < len(fields); idx++ {
			d.builder = d.builder.Suffix(", " + fields[idx])
		}
	}
	return d
}

// ToSQL builds and returns the SQL query string and arguments.
func (d *DeleteBuilder) ToSQL() (query string, args []any, err error) {
	query, args, err = d.builder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build DELETE query: %w", err)
	}
	return query, args, nil
}
