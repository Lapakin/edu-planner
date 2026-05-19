package integration_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/user-management/repository/postgres"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestUpsertAcademicYearInfo(t *testing.T) {
	repo := postgres.NewAcademicYearRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *domain.AcademicYear
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.AcademicYear1,
			expectedError: nil,
		},
		{
			name:          "OK_Second",
			input:         ta.AcademicYear2,
			expectedError: nil,
		},
		{
			name:          "OK_Conflict",
			input:         ta.AcademicYear1,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.UpsertAcademicYear(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteAcademicYearInfo(t *testing.T) {
	repo := postgres.NewAcademicYearRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.AcademicYear2.ID,
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
			err := repo.DeleteAcademicYear(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
