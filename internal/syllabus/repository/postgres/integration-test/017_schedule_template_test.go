package integration_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository/postgres"
	"github.com/Lapakin/edu-planner/internal/utils"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateScheduleTemplate(t *testing.T) {
	repo := postgres.NewScheduleTemplateRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *domain.ScheduleTemplate
		expectedError error
	}{
		{
			name:          "First",
			input:         ta.ScheduleTemplate1,
			expectedError: nil,
		},
		{
			name:          "Second",
			input:         ta.ScheduleTemplate2,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateScheduleTemplate(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetScheduleTemplateByID(t *testing.T) {
	repo := postgres.NewScheduleTemplateRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.ScheduleTemplate
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.ScheduleTemplate1.ID,
			expectedOutput: ta.ScheduleTemplate1,
			expectedError:  nil,
		},
		{
			name:           "NotFound",
			input:          0,
			expectedOutput: nil,
			expectedError:  sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.GetScheduleTemplateByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchScheduleTemplates(t *testing.T) {
	repo := postgres.NewScheduleTemplateRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.ScheduleTemplates
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.ScheduleTemplatesArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.ScheduleTemplates{ta.ScheduleTemplate1},
			expectedError:  nil,
		},
		{
			name:           "WithSemesterIDFilter",
			filters:        f.Filters{domain.SemesterIDParam: "1"},
			expectedOutput: ta.ScheduleTemplatesArray,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchScheduleTemplates(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateScheduleTemplate(t *testing.T) {
	repo := postgres.NewScheduleTemplateRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            *domain.ScheduleTemplate
		dataManipulation func(*domain.ScheduleTemplate) *domain.ScheduleTemplate
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.ScheduleTemplate1,
			dataManipulation: func(original *domain.ScheduleTemplate) *domain.ScheduleTemplate {
				var modified domain.ScheduleTemplate
				utils.Copy(original, &modified)
				modified.Name = new("Updated Template")
				return &modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.ScheduleTemplate1,
			dataManipulation: func(original *domain.ScheduleTemplate) *domain.ScheduleTemplate {
				var modified domain.ScheduleTemplate
				utils.Copy(original, &modified)
				modified.ID = 0
				return &modified
			},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.dataManipulation != nil {
				tc.input = tc.dataManipulation(tc.input)
			}
			err := repo.UpdateScheduleTemplate(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestActivateScheduleTemplate(t *testing.T) {
	repo := postgres.NewScheduleTemplateRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.ScheduleTemplate1.ID,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.ActivateScheduleTemplate(ctx, tc.input, time.Now())
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteScheduleTemplates(t *testing.T) {
	repo := postgres.NewScheduleTemplateRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.ScheduleTemplate2.ID},
			expectedError: nil,
		},
		{
			name:          "NotFound",
			input:         []uint64{0},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.DeleteScheduleTemplates(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
