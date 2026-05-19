package rest

import (
	"github.com/gin-gonic/gin"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type Query struct {
	Param        string
	ValidateFunc func(string) error
	Required     bool
}

type Queries []*Query

func CreateFiltersFromQueries(c *gin.Context, queries Queries) (f.Filters, error) {
	filters := make(f.Filters)
	for _, q := range queries {
		value := c.Query(q.Param)
		if len(value) == 0 {
			if q.Required {
				return nil, ErrEmptyValue
			}
			continue
		}
		if q.ValidateFunc != nil {
			if err := q.ValidateFunc(value); err != nil {
				return nil, err
			}
		}
		filters[q.Param] = value
	}
	return filters, nil
}
