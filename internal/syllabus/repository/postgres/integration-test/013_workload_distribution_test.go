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

func TestCreateWorkloadDistributions(t *testing.T) {
	repo := postgres.NewWorkloadDistributionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.WorkloadDistributions
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.WorkloadDistributionsArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateWorkloadDistributions(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetWorkloadDistributionByID(t *testing.T) {
	repo := postgres.NewWorkloadDistributionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.WorkloadDistribution
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.WorkloadDistribution1.ID,
			expectedOutput: ta.WorkloadDistribution1,
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
			output, err := repo.GetWorkloadDistributionByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchWorkloadDistributions(t *testing.T) {
	repo := postgres.NewWorkloadDistributionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.WorkloadDistributions
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.WorkloadDistributionsArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.WorkloadDistributions{ta.WorkloadDistribution1},
			expectedError:  nil,
		},
		{
			name:           "WithStudyPlanIDFilter",
			filters:        f.Filters{domain.StudyPlanIDParam: "1"},
			expectedOutput: ta.WorkloadDistributionsArray,
			expectedError:  nil,
		},
		{
			name:           "WithGroupIDFilter",
			filters:        f.Filters{domain.GroupIDParam: "1"},
			expectedOutput: domain.WorkloadDistributions{ta.WorkloadDistribution1},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchWorkloadDistributions(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateWorkloadDistributions(t *testing.T) {
	repo := postgres.NewWorkloadDistributionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.WorkloadDistributions
		dataManipulation func(domain.WorkloadDistributions) domain.WorkloadDistributions
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.WorkloadDistributionsArray,
			dataManipulation: func(original domain.WorkloadDistributions) domain.WorkloadDistributions {
				var modified domain.WorkloadDistributions
				utils.Copy(original, &modified)
				modified[1].ClassroomWork = new(20)
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.WorkloadDistributionsArray,
			dataManipulation: func(original domain.WorkloadDistributions) domain.WorkloadDistributions {
				var modified domain.WorkloadDistributions
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
			err := repo.UpdateWorkloadDistributions(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteWorkloadDistributions(t *testing.T) {
	repo := postgres.NewWorkloadDistributionRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.WorkloadDistribution2.ID},
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
			err := repo.DeleteWorkloadDistributions(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
