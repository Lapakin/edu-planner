package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

type InsertBuilder struct {
	ctx       context.Context
	builder   sq.InsertBuilder
	table     string
	columns   []string
	values    [][]any
	returning []string
}

// Insert creates a new InsertBuilder with Dollar placeholder format.
func Insert(ctx context.Context) *InsertBuilder {
	return &InsertBuilder{
		ctx:       ctx,
		builder:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert(""),
		table:     "",
		columns:   nil,
		values:    nil,
		returning: nil,
	}
}

// Into sets the table to insert into.
func (i *InsertBuilder) Into(table string) *InsertBuilder {
	i.table = table
	i.builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert(table)
	return i
}

// Columns sets the columns to insert.
func (i *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	i.columns = columns
	i.builder = i.builder.Columns(columns...)
	return i
}

// Values adds a row of values to insert.
func (i *InsertBuilder) Values(values ...any) *InsertBuilder {
	i.values = append(i.values, values)
	i.builder = i.builder.Values(values...)
	return i
}

// SetMap sets column values from a map.
func (i *InsertBuilder) SetMap(clauses map[string]any) *InsertBuilder {
	i.builder = i.builder.SetMap(clauses)
	return i
}

// Returning adds RETURNING clause.
func (i *InsertBuilder) Returning(fields ...string) *InsertBuilder {
	i.returning = fields
	if len(fields) > 0 {
		i.builder = i.builder.Suffix("RETURNING " + fields[0])
		for idx := 1; idx < len(fields); idx++ {
			i.builder = i.builder.Suffix(", " + fields[idx])
		}
	}
	return i
}

// ReturningID adds RETURNING id clause.
func (i *InsertBuilder) ReturningID() *InsertBuilder {
	i.builder = i.builder.Suffix("RETURNING " + ID)
	return i
}

// OnConflict adds ON CONFLICT clause for upsert operations.
func (i *InsertBuilder) OnConflict(clause string) *InsertBuilder {
	i.builder = i.builder.Suffix("ON CONFLICT " + clause)
	return i
}

// OnConflictDoNothing adds ON CONFLICT DO NOTHING clause.
func (i *InsertBuilder) OnConflictDoNothing() *InsertBuilder {
	i.builder = i.builder.Suffix("ON CONFLICT DO NOTHING")
	return i
}

// OnConflictDoUpdate adds ON CONFLICT DO UPDATE clause.
func (i *InsertBuilder) OnConflictDoUpdate(conflictColumn string, updateColumns ...string) *InsertBuilder {
	var clause strings.Builder
	fmt.Fprintf(&clause, "ON CONFLICT (%s) DO UPDATE SET ", conflictColumn)
	for idx, col := range updateColumns {
		if idx > 0 {
			clause.WriteString(", ")
		}
		fmt.Fprintf(&clause, "%s = EXCLUDED.%s", col, col)
	}
	i.builder = i.builder.Suffix(clause.String())
	return i
}

// ToSQL builds and returns the SQL query string and arguments.
func (i *InsertBuilder) ToSQL() (query string, args []any, err error) {
	query, args, err = i.builder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build INSERT query: %w", err)
	}
	return query, args, nil
}

// getNumberOfColumns returns the number of columns in the insert.
func (i *InsertBuilder) getNumberOfColumns() int {
	return len(i.columns)
}

// QueryxContext executes the insert query with automatic batching when exceeding MaxParams.
func (i *InsertBuilder) QueryxContext(db sqlx.ExtContext, scanFunc func(*sqlx.Rows, *int) error) error {
	n := i.getNumberOfColumns()

	counter := 0
	executeQuery := func(query string, args []any) error {
		rows, err := db.QueryxContext(i.ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		if err = scanFunc(rows, &counter); err != nil {
			return err
		}

		return rows.Err()
	}

	if len(i.values)*n < MaxParams {
		query, args, err := i.ToSQL()
		if err != nil {
			return err
		}
		return executeQuery(query, args)
	}

	constChunkSize := MaxParams / n
	for len(i.values) > 0 {
		chunkSize := min(len(i.values), constChunkSize)

		chunkValues := i.values[:chunkSize]
		i.values = i.values[chunkSize:]

		chunkInsert := &InsertBuilder{
			ctx:       i.ctx,
			builder:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert(i.table),
			table:     i.table,
			columns:   i.columns,
			values:    nil,
			returning: i.returning,
		}

		chunkInsert.builder = chunkInsert.builder.Columns(i.columns...)

		for _, v := range chunkValues {
			chunkInsert.builder = chunkInsert.builder.Values(v...)
		}

		if len(i.returning) > 0 {
			chunkInsert.builder = chunkInsert.builder.Suffix("RETURNING " + i.returning[0])
			for idx := 1; idx < len(i.returning); idx++ {
				chunkInsert.builder = chunkInsert.builder.Suffix(", " + i.returning[idx])
			}
		}

		query, args, err := chunkInsert.ToSQL()
		if err != nil {
			return err
		}

		if err = executeQuery(query, args); err != nil {
			return err
		}
	}

	return nil
}
