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

func TestCreateStudyPlans(t *testing.T) {
	repo := postgres.NewStudyPlanRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.StudyPlans
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.StudyPlansArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateStudyPlans(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetStudyPlanByID(t *testing.T) {
	repo := postgres.NewStudyPlanRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.StudyPlan
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.StudyPlan1.ID,
			expectedOutput: ta.StudyPlan1,
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
			output, err := repo.GetStudyPlanByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchStudyPlans(t *testing.T) {
	repo := postgres.NewStudyPlanRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.StudyPlans
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.StudyPlansArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.StudyPlans{ta.StudyPlan1},
			expectedError:  nil,
		},
		{
			name:           "WithAcademicYearFilter",
			filters:        f.Filters{domain.AcademicYearIDParam: "1"},
			expectedOutput: ta.StudyPlansArray,
			expectedError:  nil,
		},
		{
			name:           "WithSpecialtyFilter",
			filters:        f.Filters{domain.SpecialtyIDParam: "1"},
			expectedOutput: domain.StudyPlans{ta.StudyPlan1},
			expectedError:  nil,
		},
		{
			name:           "WithDisciplineFilter",
			filters:        f.Filters{domain.DisciplineIDParam: "1"},
			expectedOutput: domain.StudyPlans{ta.StudyPlan1},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchStudyPlans(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateStudyPlans(t *testing.T) {
	repo := postgres.NewStudyPlanRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.StudyPlans
		dataManipulation func(domain.StudyPlans) domain.StudyPlans
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.StudyPlansArray,
			dataManipulation: func(original domain.StudyPlans) domain.StudyPlans {
				var modified domain.StudyPlans
				utils.Copy(original, &modified)
				modified[1].SemesterNumber = new(3)
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.StudyPlansArray,
			dataManipulation: func(original domain.StudyPlans) domain.StudyPlans {
				var modified domain.StudyPlans
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
			err := repo.UpdateStudyPlans(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteStudyPlans(t *testing.T) {
	repo := postgres.NewStudyPlanRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.StudyPlan2.ID},
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
			err := repo.DeleteStudyPlans(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
