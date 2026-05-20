package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

// UpdateBuilder wraps squirrel.UpdateBuilder with additional functionality.
type UpdateBuilder struct {
	ctx     context.Context
	builder sq.UpdateBuilder
}

func Update(ctx context.Context) *UpdateBuilder {
	return &UpdateBuilder{
		ctx:     ctx,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Update(""),
	}
}

// Table sets the table to update.
func (u *UpdateBuilder) Table(table string) *UpdateBuilder {
	u.builder = u.builder.Table(table)
	return u
}

// Set sets a column value.
func (u *UpdateBuilder) Set(column string, value any) *UpdateBuilder {
	u.builder = u.builder.Set(column, value)
	return u
}

// SetMap sets multiple column values from a map.
func (u *UpdateBuilder) SetMap(clauses map[string]any) *UpdateBuilder {
	u.builder = u.builder.SetMap(clauses)
	return u
}

// Where adds a WHERE clause.
func (u *UpdateBuilder) Where(pred any, args ...any) *UpdateBuilder {
	u.builder = u.builder.Where(pred, args...)
	return u
}

// WhereID adds a WHERE clause for id = value.
func (u *UpdateBuilder) WhereID(id any) *UpdateBuilder {
	u.builder = u.builder.Where(sq.Eq{ID: id})
	return u
}

// IsActive adds a WHERE clause for is_active = true.
func (u *UpdateBuilder) IsActive() *UpdateBuilder {
	u.builder = u.builder.Where(sq.Eq{IsActive: true})
	return u
}

// IsDeleted adds a WHERE clause for is_deleted = false (not deleted records by default).
func (u *UpdateBuilder) IsDeleted(deleted bool) *UpdateBuilder {
	u.builder = u.builder.Where(sq.Eq{IsDeleted: deleted})
	return u
}

// WhereMap adds WHERE clauses from a map.
func (u *UpdateBuilder) WhereMap(clauses map[string]any) *UpdateBuilder {
	for key, value := range clauses {
		u.builder = u.builder.Where(sq.Eq{key: value})
	}
	return u
}

// Returning adds RETURNING clause.
func (u *UpdateBuilder) Returning(fields ...string) *UpdateBuilder {
	if len(fields) > 0 {
		u.builder = u.builder.Suffix("RETURNING " + fields[0])
		for i := 1; i < len(fields); i++ {
			u.builder = u.builder.Suffix(", " + fields[i])
		}
	}
	return u
}

// ToSQL builds and returns the SQL query string and arguments.
func (u *UpdateBuilder) ToSQL() (query string, args []any, err error) {
	query, args, err = u.builder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build UPDATE query: %w", err)
	}
	return query, args, nil
}

type BatchUpdateBuilder struct {
	ctx          context.Context
	table        string
	pks          []string
	cols         []string
	types        []string
	returning    []string
	expectedRows int
}

func BatchUpdate(ctx context.Context) *BatchUpdateBuilder {
	return &BatchUpdateBuilder{
		ctx:          ctx,
		table:        "",
		pks:          nil,
		cols:         nil,
		types:        nil,
		returning:    nil,
		expectedRows: 0,
	}
}

func (b *BatchUpdateBuilder) Table(table string) *BatchUpdateBuilder {
	b.table = table
	return b
}

func (b *BatchUpdateBuilder) NumberOfRows(n int) *BatchUpdateBuilder {
	b.expectedRows = n
	return b
}

func (b *BatchUpdateBuilder) PrimaryKey(pks ...string) *BatchUpdateBuilder {
	b.pks = pks
	return b
}

func (b *BatchUpdateBuilder) AddBigintColumn(name string) *BatchUpdateBuilder {
	return b.addColumn(name, "::bigint")
}

func (b *BatchUpdateBuilder) AddVarCharColumn(name string) *BatchUpdateBuilder {
	return b.addColumn(name, "::varchar")
}

func (b *BatchUpdateBuilder) AddBoolColumn(name string) *BatchUpdateBuilder {
	return b.addColumn(name, "::bool")
}

func (b *BatchUpdateBuilder) AddTextArrayColumn(name string) *BatchUpdateBuilder {
	return b.addColumn(name, "::text[]")
}

func (b *BatchUpdateBuilder) AddTimeStampColumn(name string) *BatchUpdateBuilder {
	return b.addColumn(name, "::timestamptz")
}

func (b *BatchUpdateBuilder) AddIntColumn(name string) *BatchUpdateBuilder {
	return b.addColumn(name, "::int")
}

func (b *BatchUpdateBuilder) AddFloatColumn(name string) *BatchUpdateBuilder {
	return b.addColumn(name, "::float8")
}

func (b *BatchUpdateBuilder) AddTimeColumn(name string) *BatchUpdateBuilder {
	return b.addColumn(name, "::time")
}

// AddGenericColumn allows custom types if needed.
func (b *BatchUpdateBuilder) addColumn(name, typeCast string) *BatchUpdateBuilder {
	b.cols = append(b.cols, name)
	b.types = append(b.types, typeCast)
	return b
}

func (b *BatchUpdateBuilder) Returning(cols ...string) *BatchUpdateBuilder {
	b.returning = cols
	return b
}

func (b *BatchUpdateBuilder) ReturningID() *BatchUpdateBuilder {
	b.returning = append(b.returning, ID)
	return b
}

func (b *BatchUpdateBuilder) ReturningCreatedAt() *BatchUpdateBuilder {
	b.returning = append(b.returning, CreatedAt)
	return b
}

func (b *BatchUpdateBuilder) QueryxContext(db sqlx.ExtContext, args []any, scanFunc func(*sqlx.Rows, *int) error) error {
	if b.table == "" {
		return errors.New("batch update: table not set")
	}
	if len(b.pks) == 0 {
		return errors.New("batch update: primary keys not set")
	}
	if len(b.cols) == 0 {
		return errors.New("batch update: no columns specified")
	}

	colsPerRow := len(b.cols)
	if len(args)%colsPerRow != 0 {
		return fmt.Errorf("batch update: args count (%d) not divisible by column count (%d)", len(args), colsPerRow)
	}

	totalRows := len(args) / colsPerRow
	if totalRows == 0 {
		return nil // Nothing to update
	}

	// Postgres Parameter Limit is ~65535. Safe limit ~60000.
	// Calculate max rows per chunk.
	maxRowsPerChunk := max(60000/colsPerRow, 1)

	globalIndex := 0 // Tracks the index in the original slice (for the scanner)

	for startRow := 0; startRow < totalRows; startRow += maxRowsPerChunk {
		endRow := min(startRow+maxRowsPerChunk, totalRows)

		// Slice the flat args for this chunk
		startArg := startRow * colsPerRow
		endArg := endRow * colsPerRow
		chunkArgs := args[startArg:endArg]

		query := b.buildSQL(endRow - startRow)

		// Execute Chunk
		rows, err := db.QueryxContext(b.ctx, query, chunkArgs...)
		if err != nil {
			return fmt.Errorf("batch update exec failed: %w", err)
		}

		// Handle Scanning
		// We pass the globalIndex pointer so the user can map results back to their original slice
		if scanFunc != nil {
			err = func() error {
				defer rows.Close()
				return scanFunc(rows, &globalIndex)
			}()
			if err != nil {
				return err
			}
		} else {
			rows.Close()
		}
	}

	if len(b.returning) > 0 && globalIndex != totalRows {
		return sql.ErrNoRows
	}

	return nil
}

// buildSQL constructs the raw query string for N rows.
func (b *BatchUpdateBuilder) buildSQL(numRows int) string {
	var sb strings.Builder

	// 1. UPDATE Clause
	sb.WriteString("UPDATE ")
	sb.WriteString(b.table)
	sb.WriteString(" AS t SET ")

	// 2. SET Clause (exclude PKs from being SET)
	// "col1 = v.col1, col2 = v.col2"
	pkMap := make(map[string]bool)
	for _, pk := range b.pks {
		pkMap[pk] = true
	}

	setParts := make([]string, 0, len(b.cols))
	for _, col := range b.cols {
		if !pkMap[col] {
			setParts = append(setParts, fmt.Sprintf("%q = v.%q", col, col))
		}
	}
	sb.WriteString(strings.Join(setParts, ", "))

	// 3. FROM (VALUES ...)
	sb.WriteString(" FROM (VALUES ")

	// Generate placeholders: ($1, $2::type, ...), ($X, ...)
	placeholders := make([]string, numRows)
	colsCount := len(b.cols)

	for i := range numRows {
		rowPh := make([]string, colsCount)
		for j := range colsCount {
			// Absolute parameter index: $1, $2, etc.
			idx := (i * colsCount) + j + 1
			// Append cast to the placeholder: "$1::bigint"
			rowPh[j] = fmt.Sprintf("$%d%s", idx, b.types[j])
		}
		placeholders[i] = "(" + strings.Join(rowPh, ", ") + ")"
	}
	sb.WriteString(strings.Join(placeholders, ", "))

	// 4. AS Clause
	// "AS v(col1, col2, ...)"
	sb.WriteString(") AS v(")
	quotedCols := make([]string, len(b.cols))
	for i, c := range b.cols {
		quotedCols[i] = fmt.Sprintf("%q", c)
	}
	sb.WriteString(strings.Join(quotedCols, ", "))
	sb.WriteString(")")

	// 5. WHERE Clause (Join on PKs)
	// "WHERE t.id = v.id"
	sb.WriteString(" WHERE ")
	whereParts := make([]string, len(b.pks))
	for i, pk := range b.pks {
		whereParts[i] = fmt.Sprintf("t.%q = v.%q", pk, pk)
	}
	sb.WriteString(strings.Join(whereParts, " AND "))

	// 6. RETURNING Clause (qualify with table alias to avoid ambiguity with VALUES alias)
	if len(b.returning) > 0 {
		sb.WriteString(" RETURNING ")
		qualified := make([]string, len(b.returning))
		for i, col := range b.returning {
			qualified[i] = fmt.Sprintf("t.%q", col)
		}
		sb.WriteString(strings.Join(qualified, ", "))
	}

	return sb.String()
}

func (b *BatchUpdateBuilder) ExecContext(db sqlx.ExtContext, args []any) (sql.Result, error) {
	if b.table == "" {
		return nil, errors.New("batch update: table not set")
	}
	if len(b.pks) == 0 {
		return nil, errors.New("batch update: primary keys not set")
	}
	if len(b.cols) == 0 {
		return nil, errors.New("batch update: no columns specified")
	}

	colsPerRow := len(b.cols)
	if len(args)%colsPerRow != 0 {
		return nil, fmt.Errorf("batch update: args count (%d) not divisible by column count (%d)", len(args), colsPerRow)
	}

	totalRows := len(args) / colsPerRow
	if totalRows == 0 {
		return nil, nil
	}

	query := b.buildSQL(totalRows)

	return db.ExecContext(b.ctx, query, args...)
}
