package integration_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Lapakin/edu-planner/internal/user-management/repository/postgres"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateUserCredential(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *domain.UserCredential
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.UserCredential1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateUserCredential(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetUserCredentialByEmail(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		output, err := repo.GetUserCredentialByEmail(ctx, ta.User1.Email)
		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, ta.UserCredential1.UserID, output.UserID)
		assert.Equal(t, ta.UserCredential1.PasswordHash, output.PasswordHash)
		require.NotNil(t, output.User)
		assert.Equal(t, ta.User1.ID, output.User.ID)
		assert.Equal(t, ta.User1.Email, output.User.Email)
	})

	t.Run("NotFound", func(t *testing.T) {
		output, err := repo.GetUserCredentialByEmail(ctx, "nonexistent_email")
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, output)
	})
}

func TestUpsertUserCredential(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *domain.UserCredential
		expectedError error
	}{
		{
			name:          "UpdateExisting",
			input:         ta.UserCredential1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.UpsertUserCredential(ctx, tc.input)
			require.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCreateRefreshToken(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *domain.RefreshToken
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.RefreshToken1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateRefreshToken(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetRefreshToken(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		output, err := repo.GetRefreshToken(ctx, ta.RefreshToken1.TokenHash)
		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, ta.RefreshToken1.ID, output.ID)
		assert.Equal(t, ta.RefreshToken1.UserID, output.UserID)
		assert.Equal(t, ta.RefreshToken1.TokenHash, output.TokenHash)
		assert.Equal(t, ta.RefreshToken1.IsRevoked, output.IsRevoked)
	})

	t.Run("NotFound", func(t *testing.T) {
		output, err := repo.GetRefreshToken(ctx, "nonexistent_token_hash")
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, output)
	})
}

func TestRevokeRefreshToken(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		tokenHash     string
		expectedError error
	}{
		{
			name:          "OK",
			tokenHash:     ta.RefreshToken1.TokenHash,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.RevokeRefreshToken(ctx, tc.tokenHash)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCreateInviteToken(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *domain.InviteToken
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.InviteToken1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateInviteToken(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetInviteToken(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		output, err := repo.GetInviteToken(ctx, ta.InviteToken1.Token)
		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, ta.InviteToken1.ID, output.ID)
		assert.Equal(t, ta.InviteToken1.UserID, output.UserID)
		assert.Equal(t, ta.InviteToken1.Token, output.Token)
		assert.Nil(t, output.UsedAt)
	})

	t.Run("NotFound", func(t *testing.T) {
		output, err := repo.GetInviteToken(ctx, "nonexistent_token")
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, output)
	})
}

func TestMarkInviteTokenUsed(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		token         string
		expectedError error
	}{
		{
			name:          "OK",
			token:         ta.InviteToken1.Token,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.MarkInviteTokenUsed(ctx, tc.token, ta.CurrentTime)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCreateResetPasswordToken(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *domain.ResetPasswordToken
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.ResetPasswordToken1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateResetPasswordToken(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetResetPasswordToken(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		output, err := repo.GetResetPasswordToken(ctx, ta.ResetPasswordToken1.Token)
		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, ta.ResetPasswordToken1.UserID, output.UserID)
		assert.Equal(t, ta.ResetPasswordToken1.Token, output.Token)
		assert.Nil(t, output.UsedAt)
	})

	t.Run("NotFound", func(t *testing.T) {
		output, err := repo.GetResetPasswordToken(ctx, "nonexistent_token")
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, output)
	})
}

func TestMarkResetPasswordTokenUsed(t *testing.T) {
	repo := postgres.NewAuthRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		token         string
		expectedError error
	}{
		{
			name:          "OK",
			token:         ta.ResetPasswordToken1.Token,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.MarkResetPasswordTokenUsed(ctx, tc.token, ta.CurrentTime)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
