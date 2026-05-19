package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateScheduleRestrictions(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-restrictions", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.ScheduleRestrictionsArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetScheduleRestrictionByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-restrictions/{scheduleRestrictionId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"scheduleRestrictionId": ta.ScheduleRestriction1.ID}, ta.ScheduleRestriction1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"scheduleRestrictionId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchScheduleRestrictions(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-restrictions", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.ScheduleRestrictionsArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.ScheduleRestriction1.ID}}, nil, domain.ScheduleRestrictions{ta.ScheduleRestriction1}).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateScheduleRestrictions(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-restrictions", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.ScheduleRestrictionsArray,
			func(input any) any {
				var modified domain.ScheduleRestrictions
				utils.Copy(input, &modified)
				modified[0].MaxGroupLessonsPerDay = 5
				return modified
			}, http.StatusOK, ta.ScheduleRestrictionsArray,
		).
		NewRequestWithBody("NotFound", domain.ScheduleRestrictions{&domain.ScheduleRestriction{ID: 0, MinGroupLessonsPerDay: 2, MaxGroupLessonsPerDay: 4, MaxTeacherLessonsPerDay: 5}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteScheduleRestrictions(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-restrictions/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.ScheduleRestriction2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
