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

func TestCreateGroups(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.Groups
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.GroupsArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateGroups(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAttachGroupsToAcademicYear(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         map[uint64]domain.Groups
		expectedError error
	}{
		{
			name:          "OK",
			input:         map[uint64]domain.Groups{ta.AcademicYear1.ID: ta.GroupsArray},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.AttachGroupsToAcademicYear(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetGroupByID(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.Group
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.Group1.ID,
			expectedOutput: ta.Group1,
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
			output, err := repo.GetGroupByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchGroups(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.Groups
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.GroupsArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.Groups{ta.Group1},
			expectedError:  nil,
		},
		{
			name:           "WithAcademicYearFilter",
			filters:        f.Filters{domain.AcademicYearIDParam: "1"},
			expectedOutput: ta.GroupsArray,
			expectedError:  nil,
		},
		{
			name:           "WithSpecialtyFilter",
			filters:        f.Filters{domain.SpecialtyIDParam: "1"},
			expectedOutput: ta.GroupsArray,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchGroups(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateGroups(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.Groups
		dataManipulation func(domain.Groups) domain.Groups
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.GroupsArray,
			dataManipulation: func(original domain.Groups) domain.Groups {
				var modified domain.Groups
				utils.Copy(original, &modified)
				modified[1].Name = ta.Test
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.GroupsArray,
			dataManipulation: func(original domain.Groups) domain.Groups {
				var modified domain.Groups
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
			err := repo.UpdateGroups(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteGroups(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.Group2.ID},
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
			err := repo.DeleteGroups(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCreateGroupSemesters(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.GroupSemesters
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.GroupSemestersArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateGroupSemesters(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetGroupSemesterByID(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.GroupSemester
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.GroupSemester1.ID,
			expectedOutput: ta.GroupSemester1,
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
			output, err := repo.GetGroupSemesterByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchGroupSemesters(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.GroupSemesters
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.GroupSemestersArray,
			expectedError:  nil,
		},
		{
			name:           "WithGroupIDFilter",
			filters:        f.Filters{domain.GroupIDParam: "1"},
			expectedOutput: domain.GroupSemesters{ta.GroupSemester1},
			expectedError:  nil,
		},
		{
			name:           "WithSemesterIDFilter",
			filters:        f.Filters{domain.SemesterIDParam: "1"},
			expectedOutput: ta.GroupSemestersArray,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchGroupSemesters(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateGroupSemesters(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.GroupSemesters
		dataManipulation func(domain.GroupSemesters) domain.GroupSemesters
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.GroupSemestersArray,
			dataManipulation: func(original domain.GroupSemesters) domain.GroupSemesters {
				var modified domain.GroupSemesters
				utils.Copy(original, &modified)
				modified[1].EndDate = time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.GroupSemestersArray,
			dataManipulation: func(original domain.GroupSemesters) domain.GroupSemesters {
				var modified domain.GroupSemesters
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
			err := repo.UpdateGroupSemesters(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteGroupSemesters(t *testing.T) {
	repo := postgres.NewGroupRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.GroupSemester2.ID},
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
			err := repo.DeleteGroupSemesters(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
