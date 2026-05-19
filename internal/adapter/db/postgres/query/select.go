package query

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

// SelectBuilder wraps squirrel.SelectBuilder with additional functionality.
type SelectBuilder struct {
	ctx     context.Context
	builder sq.SelectBuilder
}

func Select(ctx context.Context) *SelectBuilder {
	return &SelectBuilder{
		ctx:     ctx,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select(),
	}
}

// Columns sets the columns to select.
func (s *SelectBuilder) Columns(columns ...string) *SelectBuilder {
	s.builder = s.builder.Columns(columns...)
	return s
}

// From sets the table to select from.
func (s *SelectBuilder) From(table string) *SelectBuilder {
	s.builder = s.builder.From(table)
	return s
}

// Where adds a WHERE clause.
func (s *SelectBuilder) Where(pred any, args ...any) *SelectBuilder {
	s.builder = s.builder.Where(pred, args...)
	return s
}

// WhereID adds a WHERE clause for id = value.
func (s *SelectBuilder) WhereID(id any) *SelectBuilder {
	s.builder = s.builder.Where(sq.Eq{ID: id})
	return s
}

// IsActive adds a WHERE clause for is_active = true.
func (s *SelectBuilder) IsActive() *SelectBuilder {
	s.builder = s.builder.Where(sq.Eq{IsActive: true})
	return s
}

// IsDeleted adds a WHERE clause for is_deleted = false (not deleted records by default).
func (s *SelectBuilder) IsDeleted(deleted bool) *SelectBuilder {
	s.builder = s.builder.Where(sq.Eq{IsDeleted: deleted})
	return s
}

// OrderBy adds an ORDER BY clause.
func (s *SelectBuilder) OrderBy(orderBys ...string) *SelectBuilder {
	s.builder = s.builder.OrderBy(orderBys...)
	return s
}

// Limit adds a LIMIT clause.
func (s *SelectBuilder) Limit(limit uint64) *SelectBuilder {
	s.builder = s.builder.Limit(limit)
	return s
}

// Offset adds an OFFSET clause.
func (s *SelectBuilder) Offset(offset uint64) *SelectBuilder {
	s.builder = s.builder.Offset(offset)
	return s
}

// Join adds a JOIN clause.
func (s *SelectBuilder) Join(join string, args ...any) *SelectBuilder {
	s.builder = s.builder.Join(join, args...)
	return s
}

// InnerJoin adds an INNER JOIN clause.
func (s *SelectBuilder) InnerJoin(join string, _ ...any) *SelectBuilder {
	s.builder = s.builder.JoinClause("INNER JOIN " + join)
	return s
}

// LeftJoin adds a LEFT JOIN clause.
func (s *SelectBuilder) LeftJoin(join string, args ...any) *SelectBuilder {
	s.builder = s.builder.LeftJoin(join, args...)
	return s
}

// RightJoin adds a RIGHT JOIN clause.
func (s *SelectBuilder) RightJoin(join string, args ...any) *SelectBuilder {
	s.builder = s.builder.RightJoin(join, args...)
	return s
}

// GroupBy adds a GROUP BY clause.
func (s *SelectBuilder) GroupBy(groupBys ...string) *SelectBuilder {
	s.builder = s.builder.GroupBy(groupBys...)
	return s
}

// Having adds a HAVING clause.
func (s *SelectBuilder) Having(pred any, args ...any) *SelectBuilder {
	s.builder = s.builder.Having(pred, args...)
	return s
}

// ApplyQueryFilters applies filters to the query based on the provided conditions.
func (s *SelectBuilder) ApplyQueryFilters(filters f.Filters, whereConditions Wheres) *SelectBuilder {
	for _, where := range whereConditions {
		s.applyWhere(filters, where)
	}
	return s
}

// applyWhere processes a single Where condition.
func (s *SelectBuilder) applyWhere(filters f.Filters, where *Where) {
	if where == nil {
		return
	}

	var predicates []sq.Sqlizer
	for _, condition := range where.Conditions {
		if condition == nil {
			continue
		}

		value, exists := filters[condition.Name]
		if !exists || value == "" {
			continue
		}

		predicate := s.buildPredicate(condition, value)
		if predicate != nil {
			predicates = append(predicates, predicate)
		}
	}

	if len(predicates) == 0 {
		return
	}

	switch where.Operator {
	case And:
		s.builder = s.builder.Where(sq.And(predicates))
	case Or:
		s.builder = s.builder.Where(sq.Or(predicates))
	default:
		s.builder = s.builder.Where(sq.And(predicates))
	}
}

// buildPredicate creates a squirrel predicate based on the condition.
func (s *SelectBuilder) buildPredicate(condition *Condition, value string) sq.Sqlizer {
	switch condition.Operator {
	case Equals:
		return sq.Eq{condition.Column: value}
	case NotEquals:
		return sq.NotEq{condition.Column: value}
	case GreaterThan:
		return sq.Gt{condition.Column: value}
	case GreaterThanOrEqual:
		return sq.GtOrEq{condition.Column: value}
	case LessThan:
		return sq.Lt{condition.Column: value}
	case LessThanOrEqual:
		return sq.LtOrEq{condition.Column: value}
	case Like:
		return sq.Like{condition.Column: "%" + value + "%"}
	case ILike:
		return sq.ILike{condition.Column: "%" + value + "%"}
	case In:
		return sq.Eq{condition.Column: value}
	case NotIn:
		return sq.NotEq{condition.Column: value}
	case IsNull:
		return sq.Eq{condition.Column: nil}
	case IsNotNull:
		return sq.NotEq{condition.Column: nil}
	case Contains:
		return sq.Expr("? = ANY("+condition.Column+")", value)
	default:
		return sq.Eq{condition.Column: value}
	}
}

// ToSQL builds and returns the SQL query string and arguments.
func (s *SelectBuilder) ToSQL() (query string, args []any, err error) {
	query, args, err = s.builder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build SQL query: %w", err)
	}
	return query, args, nil
}
