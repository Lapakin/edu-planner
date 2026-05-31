package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Lapakin/edu-planner/internal/user-management/repository/postgres"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestAttachUsers(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		userIDs       []uint64
		expectedError error
	}{
		{
			name:          "OK",
			userIDs:       []uint64{ta.User1.ID, ta.User2.ID},
			expectedError: nil,
		},
		{
			name:          "OK_Idempotent",
			userIDs:       []uint64{ta.User1.ID, ta.User2.ID},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.AttachUsers(ctx, ta.AcademicYear1.ID, tc.userIDs)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchActiveUserIDs(t *testing.T) {
	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	require.NoError(t, repo.ActivateUser(ctx, ta.User2.ID, ta.CurrentTime))

	ids, err := repo.FetchActiveUserIDs(ctx)
	require.NoError(t, err)
	assert.Contains(t, ids, ta.User2.ID)
}
