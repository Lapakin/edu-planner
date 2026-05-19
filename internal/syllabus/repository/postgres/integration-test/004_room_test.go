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

func TestCreateRooms(t *testing.T) {
	repo := postgres.NewRoomRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.Rooms
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.RoomsArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateRooms(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAttachRoomsToAcademicYear(t *testing.T) {
	repo := postgres.NewRoomRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         map[uint64]domain.Rooms
		expectedError error
	}{
		{
			name:          "OK",
			input:         map[uint64]domain.Rooms{ta.AcademicYear1.ID: ta.RoomsArray},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.AttachRoomsToAcademicYear(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetRoomByID(t *testing.T) {
	repo := postgres.NewRoomRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.Room
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.Room1.ID,
			expectedOutput: ta.Room1,
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
			output, err := repo.GetRoomByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchRooms(t *testing.T) {
	repo := postgres.NewRoomRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.Rooms
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.RoomsArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.Rooms{ta.Room1},
			expectedError:  nil,
		},
		{
			name:           "WithAcademicYearFilter",
			filters:        f.Filters{domain.AcademicYearIDParam: "1"},
			expectedOutput: ta.RoomsArray,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchRooms(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateRooms(t *testing.T) {
	repo := postgres.NewRoomRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.Rooms
		dataManipulation func(domain.Rooms) domain.Rooms
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.RoomsArray,
			dataManipulation: func(original domain.Rooms) domain.Rooms {
				var modified domain.Rooms
				utils.Copy(original, &modified)
				modified[1].Name = ta.Test
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.RoomsArray,
			dataManipulation: func(original domain.Rooms) domain.Rooms {
				var modified domain.Rooms
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
			err := repo.UpdateRooms(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteRooms(t *testing.T) {
	repo := postgres.NewRoomRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.Room2.ID},
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
			err := repo.DeleteRooms(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
