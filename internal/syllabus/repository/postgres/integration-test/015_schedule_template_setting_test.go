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

func TestCreateScheduleTemplateSettings(t *testing.T) {
	repo := postgres.NewScheduleTemplateSettingRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         domain.ScheduleTemplateSettings
		expectedError error
	}{
		{
			name:          "OK",
			input:         ta.ScheduleTemplateSettingsArray,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.CreateScheduleTemplateSettings(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetScheduleTemplateSettingByID(t *testing.T) {
	repo := postgres.NewScheduleTemplateSettingRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		input          uint64
		expectedOutput *domain.ScheduleTemplateSetting
		expectedError  error
	}{
		{
			name:           "OK",
			input:          ta.ScheduleTemplateSetting1.ID,
			expectedOutput: ta.ScheduleTemplateSetting1,
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
			output, err := repo.GetScheduleTemplateSettingByID(ctx, tc.input)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestFetchScheduleTemplateSettings(t *testing.T) {
	repo := postgres.NewScheduleTemplateSettingRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name           string
		filters        f.Filters
		expectedOutput domain.ScheduleTemplateSettings
		expectedError  error
	}{
		{
			name:           "OK",
			filters:        f.Filters{},
			expectedOutput: ta.ScheduleTemplateSettingsArray,
			expectedError:  nil,
		},
		{
			name:           "WithIDFilter",
			filters:        f.Filters{domain.IDsParam: "1"},
			expectedOutput: domain.ScheduleTemplateSettings{ta.ScheduleTemplateSetting1},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := repo.FetchScheduleTemplateSettings(ctx, tc.filters)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestUpdateScheduleTemplateSettings(t *testing.T) {
	repo := postgres.NewScheduleTemplateSettingRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name             string
		input            domain.ScheduleTemplateSettings
		dataManipulation func(domain.ScheduleTemplateSettings) domain.ScheduleTemplateSettings
		expectedError    error
	}{
		{
			name:  "OK",
			input: ta.ScheduleTemplateSettingsArray,
			dataManipulation: func(original domain.ScheduleTemplateSettings) domain.ScheduleTemplateSettings {
				var modified domain.ScheduleTemplateSettings
				utils.Copy(original, &modified)
				modified[1].MaxStudyHoursPerDay = 12
				return modified
			},
			expectedError: nil,
		},
		{
			name:  "NotFound",
			input: ta.ScheduleTemplateSettingsArray,
			dataManipulation: func(original domain.ScheduleTemplateSettings) domain.ScheduleTemplateSettings {
				var modified domain.ScheduleTemplateSettings
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
			err := repo.UpdateScheduleTemplateSettings(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteScheduleTemplateSettings(t *testing.T) {
	repo := postgres.NewScheduleTemplateSettingRepository(db)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         []uint64
		expectedError error
	}{
		{
			name:          "OK",
			input:         []uint64{ta.ScheduleTemplateSetting2.ID},
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
			err := repo.DeleteScheduleTemplateSettings(ctx, tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
