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

func TestCreateCycleCommittees(t *testing.T) {
	repo := postgres.NewCycleCommitteeRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.CycleCommittees
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.CycleCommitteesArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateCycleCommittees(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAttachCycleCommitteesToAcademicYear(t *testing.T) {
	repo := postgres.NewCycleCommitteeRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         map[uint64]domain.CycleCommittees
		expectedError error
	}{
		{
			name:          "OK",
			input:         map[uint64]domain.CycleCommittees{ta.AcademicYear1.ID: ta.CycleCommitteesArray},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.AttachCycleCommitteesToAcademicYear(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetCycleCommitteeByID(t *testing.T) {
	repo := postgres.NewCycleCommitteeRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.CycleCommittee
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.CycleCommittee1.ID,
			expectedOutput: ta.CycleCommittee1,
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
			output, err := repo.GetCycleCommitteeByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchCycleCommittees(t *testing.T) {
	repo := postgres.NewCycleCommitteeRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.CycleCommittees
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.CycleCommitteesArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.CycleCommittees{ta.CycleCommittee1},
			expectedError:  nil,
		},
		{
			name:           "WithAcademicYearFilter",
			filters:        f.Filters{domain.AcademicYearIDParam: "1"},
			expectedOutput: ta.CycleCommitteesArray,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchCycleCommittees(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateCycleCommittees(t *testing.T) {
	repo := postgres.NewCycleCommitteeRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.CycleCommittees
		dataManipulation func(domain.CycleCommittees) domain.CycleCommittees
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.CycleCommitteesArray,
			dataManipulation: func(original domain.CycleCommittees) domain.CycleCommittees {
				var modified domain.CycleCommittees
				utils.Copy(original, &modified)
				modified[1].Name = ta.Test
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.CycleCommitteesArray,
			dataManipulation: func(original domain.CycleCommittees) domain.CycleCommittees {
				var modified domain.CycleCommittees
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
			err := repo.UpdateCycleCommittees(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteCycleCommittees(t *testing.T) {
	repo := postgres.NewCycleCommitteeRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.CycleCommittee2.ID},
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
			err := repo.DeleteCycleCommittees(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
