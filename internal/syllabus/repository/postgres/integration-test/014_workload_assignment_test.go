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

func TestCreateWorkloadAssignments(t *testing.T) {
	repo := postgres.NewWorkloadAssignmentRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.WorkloadAssignments
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.WorkloadAssignmentsArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateWorkloadAssignments(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetWorkloadAssignmentByID(t *testing.T) {
	repo := postgres.NewWorkloadAssignmentRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.WorkloadAssignment
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.WorkloadAssignment1.ID,
			expectedOutput: ta.WorkloadAssignment1,
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
			output, err := repo.GetWorkloadAssignmentByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchWorkloadAssignments(t *testing.T) {
	repo := postgres.NewWorkloadAssignmentRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.WorkloadAssignments
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.WorkloadAssignmentsArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.WorkloadAssignments{ta.WorkloadAssignment1},
			expectedError:  nil,
		},
		{
			name:           "WithWorkloadDistributionIDFilter",
			filters:        f.Filters{"workload_distribution_id": "1"},
			expectedOutput: ta.WorkloadAssignmentsArray,
			expectedError:  nil,
		},
		{
			name:           "WithTeacherIDFilter",
			filters:        f.Filters{"teacher_id": "1"},
			expectedOutput: domain.WorkloadAssignments{ta.WorkloadAssignment1},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchWorkloadAssignments(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateWorkloadAssignments(t *testing.T) {
	repo := postgres.NewWorkloadAssignmentRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.WorkloadAssignments
		dataManipulation func(domain.WorkloadAssignments) domain.WorkloadAssignments
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.WorkloadAssignmentsArray,
			dataManipulation: func(original domain.WorkloadAssignments) domain.WorkloadAssignments {
				var modified domain.WorkloadAssignments
				utils.Copy(original, &modified)
				modified[1].AssignedHours = new(90)
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.WorkloadAssignmentsArray,
			dataManipulation: func(original domain.WorkloadAssignments) domain.WorkloadAssignments {
				var modified domain.WorkloadAssignments
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
			err := repo.UpdateWorkloadAssignments(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteWorkloadAssignments(t *testing.T) {
	repo := postgres.NewWorkloadAssignmentRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.WorkloadAssignment2.ID},
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
			err := repo.DeleteWorkloadAssignments(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
