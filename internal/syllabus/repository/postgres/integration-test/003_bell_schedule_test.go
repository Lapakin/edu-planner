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

func TestCreateBellSchedules(t *testing.T) {
	repo := postgres.NewBellScheduleRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.BellSchedules
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.BellSchedulesArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateBellSchedules(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetBellScheduleByID(t *testing.T) {
	repo := postgres.NewBellScheduleRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.BellSchedule
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.BellSchedule1.ID,
			expectedOutput: ta.BellSchedule1,
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
			output, err := repo.GetBellScheduleByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchBellSchedules(t *testing.T) {
	repo := postgres.NewBellScheduleRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.BellSchedules
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.BellSchedulesArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.BellSchedules{ta.BellSchedule1},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchBellSchedules(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateBellSchedules(t *testing.T) {
	repo := postgres.NewBellScheduleRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.BellSchedules
		dataManipulation func(domain.BellSchedules) domain.BellSchedules
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.BellSchedulesArray,
			dataManipulation: func(original domain.BellSchedules) domain.BellSchedules {
				var modified domain.BellSchedules
				utils.Copy(original, &modified)
				modified[0].StartTime = "08:30:00"
				return modified
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.dataManipulation != nil {
				tc.input = tc.dataManipulation(tc.input)
			}
			err := repo.UpdateBellSchedules(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteBellSchedules(t *testing.T) {
	repo := postgres.NewBellScheduleRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.BellSchedule1.ID},
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
			err := repo.DeleteBellSchedules(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
