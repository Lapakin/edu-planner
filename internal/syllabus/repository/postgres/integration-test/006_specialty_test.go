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

func TestCreateSpecialties(t *testing.T) {
	repo := postgres.NewSpecialtyRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.Specialties
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.SpecialtiesArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateSpecialties(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAttachSpecialtiesToAcademicYear(t *testing.T) {
	repo := postgres.NewSpecialtyRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         map[uint64]domain.Specialties
		expectedError error
	}{
		{
			name:          "OK",
			input:         map[uint64]domain.Specialties{ta.AcademicYear1.ID: ta.SpecialtiesArray},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.AttachSpecialtiesToAcademicYear(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetSpecialtyByID(t *testing.T) {
	repo := postgres.NewSpecialtyRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.Specialty
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.Specialty1.ID,
			expectedOutput: ta.Specialty1,
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
			output, err := repo.GetSpecialtyByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchSpecialties(t *testing.T) {
	repo := postgres.NewSpecialtyRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.Specialties
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.SpecialtiesArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.Specialties{ta.Specialty1},
			expectedError:  nil,
		},
		{
			name:           "WithAcademicYearFilter",
			filters:        f.Filters{domain.AcademicYearIDParam: "1"},
			expectedOutput: ta.SpecialtiesArray,
			expectedError:  nil,
		},
		{
			name:           "WithDepartmentFilter",
			filters:        f.Filters{domain.DepartmentIDParam: "1"},
			expectedOutput: ta.SpecialtiesArray,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchSpecialties(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateSpecialties(t *testing.T) {
	repo := postgres.NewSpecialtyRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.Specialties
		dataManipulation func(domain.Specialties) domain.Specialties
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.SpecialtiesArray,
			dataManipulation: func(original domain.Specialties) domain.Specialties {
				var modified domain.Specialties
				utils.Copy(original, &modified)
				modified[1].Name = ta.Test
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.SpecialtiesArray,
			dataManipulation: func(original domain.Specialties) domain.Specialties {
				var modified domain.Specialties
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
			err := repo.UpdateSpecialties(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteSpecialties(t *testing.T) {
	repo := postgres.NewSpecialtyRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.Specialty2.ID},
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
			err := repo.DeleteSpecialties(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
