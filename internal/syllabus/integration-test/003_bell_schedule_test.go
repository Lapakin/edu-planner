package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateBellSchedules(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/bell-schedules", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.BellSchedulesArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetBellScheduleByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/bell-schedules/{bellScheduleId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"bellScheduleId": ta.BellSchedule1.ID}, ta.BellSchedule1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"bellScheduleId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchBellSchedules(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/bell-schedules", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.BellSchedulesArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.BellSchedule1.ID}}, nil, domain.BellSchedules{ta.BellSchedule1}).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateBellSchedules(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/bell-schedules", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.BellSchedulesArray,
			func(input any) any {
				var modified domain.BellSchedules
				utils.Copy(input, &modified)
				modified[0].StartTime = "08:30:00"
				return modified
			}, http.StatusOK, ta.BellSchedulesArray,
		).
		NewRequestWithBody("NotFound", domain.BellSchedules{&domain.BellSchedule{ID: 0, LessonNumber: 0, StartTime: "00:00:00", EndTime: "00:00:00"}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteBellSchedules(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/bell-schedules/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.BellSchedule2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
