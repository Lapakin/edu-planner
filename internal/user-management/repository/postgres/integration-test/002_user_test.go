package integration_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/user-management/repository/postgres"
	"github.com/Lapakin/edu-planner/internal/utils"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateUsers(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.Users
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.UsersArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateUsers(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCreateUserProfiles(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.Users
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.UsersArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateUserProfiles(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.User
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.User1.ID,
			expectedOutput: ta.User1,
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
			output, err := repo.GetUserByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchUsers(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		output, err := repo.FetchUsers(ctx, f.Filters{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(output), 2)
	})

	t.Run("WithRoleFilter", func(t *testing.T) {
		output, err := repo.FetchUsers(ctx, f.Filters{domain.RoleParam: ta.User1.Role})
		require.NoError(t, err)
		assert.Equal(t, ta.UsersArray, output)
	})
}

func TestUpdateUsers(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.Users
		dataManipulation func(domain.Users) domain.Users
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.UsersArray,
			dataManipulation: func(original domain.Users) domain.Users {
				var modified domain.Users
				utils.Copy(original, &modified)
				modified[1].Email = "updated@test.com"
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.UsersArray,
			dataManipulation: func(original domain.Users) domain.Users {
				var modified domain.Users
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
			err := repo.UpdateUsers(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateUserProfiles(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.Users
		dataManipulation func(domain.Users) domain.Users
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.UsersArray,
			dataManipulation: func(original domain.Users) domain.Users {
				var modified domain.Users
				utils.Copy(original, &modified)
				modified[1].FirstName = ta.Test
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.UsersArray,
			dataManipulation: func(original domain.Users) domain.Users {
				var modified domain.Users
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
			err := repo.UpdateUserProfiles(ctx, input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestActivateUser(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.User1.ID,
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
			err := repo.ActivateUser(ctx, tc.input, ta.CurrentTime)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeactivateUser(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.User1.ID,
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
			err := repo.DeactivateUser(ctx, tc.input, ta.CurrentTime)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteUsers(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.User1.ID},
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
			err := repo.DeleteUsers(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
