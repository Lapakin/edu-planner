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

func TestCreateAcademicYear(t *testing.T) {
	repo := postgres.NewAcademicYearRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.AcademicYears
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.AcademicYearsArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateAcademicYear(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetAcademicYearByID(t *testing.T) {
	repo := postgres.NewAcademicYearRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.AcademicYear
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.AcademicYear1.ID,
			expectedOutput: ta.AcademicYear1,
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
			output, err := repo.GetAcademicYearByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchAcademicYears(t *testing.T) {
	repo := postgres.NewAcademicYearRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.AcademicYears
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.AcademicYearsArray,
			expectedError:  nil,
		},
		{
			name:           "WithFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.AcademicYears{ta.AcademicYear1},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchAcademicYears(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestActivateAcademicYear(t *testing.T) {
	repo := postgres.NewAcademicYearRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.AcademicYear1.ID,
			expectedError: nil,
		},
		{
			name:          "NotFound",
			input:         0,
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.ActivateAcademicYear(ctx, tc.input, ta.CurrentTime)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateAcademicYears(t *testing.T) {
	repo := postgres.NewAcademicYearRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.AcademicYears
		dataManipulation func(original domain.AcademicYears) domain.AcademicYears
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.AcademicYearsArray,
			dataManipulation: func(original domain.AcademicYears) domain.AcademicYears {
				var modified domain.AcademicYears
				utils.Copy(original, &modified)
				modified[1].Name = ta.Test
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.AcademicYearsArray,
			dataManipulation: func(original domain.AcademicYears) domain.AcademicYears {
				var modified domain.AcademicYears
				utils.Copy(original, &modified)
				modified[0].ID = 0
				return modified
			},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.input
			if tc.dataManipulation != nil {
				input = tc.dataManipulation(tc.input)
			}
			err := repo.UpdateAcademicYears(ctx, input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeactivateAcademicYear(t *testing.T) {
	repo := postgres.NewAcademicYearRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.AcademicYear1.ID,
			expectedError: nil,
		},
		{
			name:          "NotFound",
			input:         0,
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.DeactivateAcademicYear(ctx, tc.input, ta.CurrentTime)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteAcademicYears(t *testing.T) {
	repo := postgres.NewAcademicYearRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.AcademicYear2.ID},
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
			err := repo.DeleteAcademicYears(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
