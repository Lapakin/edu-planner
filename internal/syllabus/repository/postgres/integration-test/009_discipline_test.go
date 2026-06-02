package integration_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository/postgres"
	"github.com/Lapakin/edu-planner/internal/utils"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateDisciplines(t *testing.T) {
	repo := postgres.NewDisciplineRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.Disciplines
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.DisciplinesArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateDisciplines(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAttachDisciplinesToAcademicYear(t *testing.T) {
	repo := postgres.NewDisciplineRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         map[uint64]domain.Disciplines
		expectedError error
	}{
		{
			name:          "OK",
			input:         map[uint64]domain.Disciplines{ta.AcademicYear1.ID: ta.DisciplinesArray},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.AttachDisciplinesToAcademicYear(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetDisciplineByID(t *testing.T) {
	repo := postgres.NewDisciplineRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.Discipline
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.Discipline1.ID,
			expectedOutput: ta.Discipline1,
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
			output, err := repo.GetDisciplineByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchDisciplines(t *testing.T) {
	repo := postgres.NewDisciplineRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.Disciplines
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.DisciplinesArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.Disciplines{ta.Discipline1},
			expectedError:  nil,
		},
		{
			name:           "WithAcademicYearFilter",
			filters:        f.Filters{domain.AcademicYearIDParam: "1"},
			expectedOutput: ta.DisciplinesArray,
			expectedError:  nil,
		},
		{
			name:           "WithCycleCommitteeFilter",
			filters:        f.Filters{domain.CycleCommitteeIDParam: "1"},
			expectedOutput: ta.DisciplinesArray,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchDisciplines(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

// TestFetchDisciplines_FlowFlagPersists verifies the is_flow column round-trips
// through create and fetch: Discipline1 is a regular discipline (false) and
// Discipline2 is marked as a flow discipline (true).
func TestFetchDisciplines_FlowFlagPersists(t *testing.T) {
	repo := postgres.NewDisciplineRepository(db)
	ctx := context.Background()

	flowDiscipline, err := repo.GetDisciplineByID(ctx, ta.Discipline2.ID)
	require.NoError(t, err)
	assert.True(t, flowDiscipline.IsFlow, "Discipline2 should be a flow discipline")

	regularDiscipline, err := repo.GetDisciplineByID(ctx, ta.Discipline1.ID)
	require.NoError(t, err)
	assert.False(t, regularDiscipline.IsFlow, "Discipline1 should not be a flow discipline")
}

func TestUpdateDisciplines(t *testing.T) {
	repo := postgres.NewDisciplineRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.Disciplines
		dataManipulation func(domain.Disciplines) domain.Disciplines
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.DisciplinesArray,
			dataManipulation: func(original domain.Disciplines) domain.Disciplines {
				var modified domain.Disciplines
				utils.Copy(original, &modified)
				modified[1].Name = ta.Test
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.DisciplinesArray,
			dataManipulation: func(original domain.Disciplines) domain.Disciplines {
				var modified domain.Disciplines
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
			err := repo.UpdateDisciplines(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteDisciplines(t *testing.T) {
	repo := postgres.NewDisciplineRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.Discipline2.ID},
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
			err := repo.DeleteDisciplines(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
