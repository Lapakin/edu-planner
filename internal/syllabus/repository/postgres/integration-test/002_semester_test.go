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

func TestCreateSemesters(t *testing.T) {
	repo := postgres.NewSemesterRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.Semesters
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.SemestersArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateSemesters(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetSemesterByID(t *testing.T) {
	repo := postgres.NewSemesterRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.Semester
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.Semester1.ID,
			expectedOutput: ta.Semester1,
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
			output, err := repo.GetSemesterByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchSemesters(t *testing.T) {
	repo := postgres.NewSemesterRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.Semesters
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.SemestersArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.Semesters{ta.Semester1},
			expectedError:  nil,
		},
		{
			name:           "WithAcademicYearFilter",
			filters:        f.Filters{domain.AcademicYearIDParam: "1"},
			expectedOutput: ta.SemestersArray,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchSemesters(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateSemesters(t *testing.T) {
	repo := postgres.NewSemesterRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.Semesters
		dataManipulation func(domain.Semesters) domain.Semesters
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.SemestersArray,
			dataManipulation: func(original domain.Semesters) domain.Semesters {
				var modified domain.Semesters
				utils.Copy(original, &modified)
				modified[1].PeriodEnd = modified[0].PeriodEnd.Add(7 * 24 * time.Hour)
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.SemestersArray,
			dataManipulation: func(original domain.Semesters) domain.Semesters {
				var modified domain.Semesters
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
			err := repo.UpdateSemesters(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteSemesters(t *testing.T) {
	repo := postgres.NewSemesterRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.Semester2.ID},
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
			err := repo.DeleteSemesters(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
