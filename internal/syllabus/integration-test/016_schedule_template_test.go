package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestGetScheduleTemplateByIDNotFound(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-templates/{scheduleTemplateId}", http.MethodGet, adminToken).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"scheduleTemplateId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchScheduleTemplatesEmpty(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-templates", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, domain.ScheduleTemplates{}).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteScheduleTemplatesNotFound(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/schedule-templates/delete", http.MethodPost, adminToken).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
