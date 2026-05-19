package integration_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository/postgres"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestUpsertUser(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *domain.User
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.User1,
			expectedError: nil,
		},
		{
			name:          "OK_Second",
			input:         ta.User2,
			expectedError: nil,
		},
		{
			name:          "OK_Conflict",
			input:         ta.User1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.UpsertUser(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.User2.ID,
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
			err := repo.DeleteUser(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
