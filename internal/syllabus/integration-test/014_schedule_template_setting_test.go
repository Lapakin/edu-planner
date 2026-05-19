package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateScheduleTemplateSettings(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-template-settings", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.ScheduleTemplateSettingsArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetScheduleTemplateSettingByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-template-settings/{scheduleTemplateSettingId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"scheduleTemplateSettingId": ta.ScheduleTemplateSetting1.ID}, ta.ScheduleTemplateSetting1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"scheduleTemplateSettingId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchScheduleTemplateSettings(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-template-settings", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.ScheduleTemplateSettingsArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.ScheduleTemplateSetting1.ID}}, nil, domain.ScheduleTemplateSettings{ta.ScheduleTemplateSetting1}).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateScheduleTemplateSettings(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-template-settings", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.ScheduleTemplateSettingsArray,
			func(input any) any {
				var modified domain.ScheduleTemplateSettings
				utils.Copy(input, &modified)
				modified[0].MaxStudyHoursPerDay = 9
				return modified
			}, http.StatusOK, ta.ScheduleTemplateSettingsArray,
		).
		NewRequestWithBody("NotFound", domain.ScheduleTemplateSettings{&domain.ScheduleTemplateSetting{ID: 0, HoursPerLesson: 1.5, MaxIdenticalLessonsPerDay: 2, MaxStudyHoursPerDay: 8, MaxTeacherHoursPerWeek: 36}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteScheduleTemplateSettings(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-template-settings/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.ScheduleTemplateSetting2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
