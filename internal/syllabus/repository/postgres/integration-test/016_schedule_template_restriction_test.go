package integration_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository/postgres"
	"github.com/Lapakin/edu-planner/internal/utils"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateScheduleRestrictions(t *testing.T) {
	repo := postgres.NewScheduleRestrictionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.ScheduleRestrictions
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.ScheduleRestrictionsArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateScheduleRestrictions(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetScheduleRestrictionByID(t *testing.T) {
	repo := postgres.NewScheduleRestrictionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.ScheduleRestriction
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.ScheduleRestriction1.ID,
			expectedOutput: ta.ScheduleRestriction1,
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
			output, err := repo.GetScheduleRestrictionByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchScheduleRestrictions(t *testing.T) {
	repo := postgres.NewScheduleRestrictionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.ScheduleRestrictions
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.ScheduleRestrictionsArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.ScheduleRestrictions{ta.ScheduleRestriction1},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchScheduleRestrictions(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateScheduleRestrictions(t *testing.T) {
	repo := postgres.NewScheduleRestrictionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.ScheduleRestrictions
		dataManipulation func(domain.ScheduleRestrictions) domain.ScheduleRestrictions
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.ScheduleRestrictionsArray,
			dataManipulation: func(original domain.ScheduleRestrictions) domain.ScheduleRestrictions {
				var modified domain.ScheduleRestrictions
				utils.Copy(original, &modified)
				modified[1].MaxGroupLessonsPerDay = 5
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.ScheduleRestrictionsArray,
			dataManipulation: func(original domain.ScheduleRestrictions) domain.ScheduleRestrictions {
				var modified domain.ScheduleRestrictions
				utils.Copy(original, &modified)
				modified[0].ID = 0
				return modified
			},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.dataManipulation != nil {
				tc.input = tc.dataManipulation(tc.input)
			}
			err := repo.UpdateScheduleRestrictions(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteScheduleRestrictions(t *testing.T) {
	repo := postgres.NewScheduleRestrictionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.ScheduleRestriction2.ID},
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
			err := repo.DeleteScheduleRestrictions(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
