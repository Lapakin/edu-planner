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

func TestCreateTeachers(t *testing.T) {
	repo := postgres.NewTeacherRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.Teachers
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.TeachersArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateTeachers(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAttachTeachersToAcademicYear(t *testing.T) {
	repo := postgres.NewTeacherRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         map[uint64]domain.Teachers
		expectedError error
	}{
		{
			name:          "OK",
			input:         map[uint64]domain.Teachers{ta.AcademicYear1.ID: ta.TeachersArray},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.AttachTeachersToAcademicYear(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetTeacherByID(t *testing.T) {
	repo := postgres.NewTeacherRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.Teacher
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.Teacher1.ID,
			expectedOutput: ta.Teacher1,
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
			output, err := repo.GetTeacherByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchTeachers(t *testing.T) {
	repo := postgres.NewTeacherRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.Teachers
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.TeachersArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.Teachers{ta.Teacher1},
			expectedError:  nil,
		},
		{
			name:           "WithAcademicYearFilter",
			filters:        f.Filters{domain.AcademicYearIDParam: "1"},
			expectedOutput: ta.TeachersArray,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchTeachers(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateTeachers(t *testing.T) {
	repo := postgres.NewTeacherRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.Teachers
		dataManipulation func(domain.Teachers) domain.Teachers
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.TeachersArray,
			dataManipulation: func(original domain.Teachers) domain.Teachers {
				var modified domain.Teachers
				utils.Copy(original, &modified)
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.TeachersArray,
			dataManipulation: func(original domain.Teachers) domain.Teachers {
				var modified domain.Teachers
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
			err := repo.UpdateTeachers(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteTeachers(t *testing.T) {
	repo := postgres.NewTeacherRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.Teacher2.ID},
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
			err := repo.DeleteTeachers(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
