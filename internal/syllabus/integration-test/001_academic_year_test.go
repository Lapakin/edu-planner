package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateAcademicYears(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/academic-years", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.AcademicYearsArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetAcademicYearByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/academic-years/{academicYearId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"academicYearId": ta.AcademicYear1.ID}, ta.AcademicYear1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"academicYearId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchAcademicYears(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/academic-years", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.AcademicYearsArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.AcademicYear1.ID}}, nil, domain.AcademicYears{ta.AcademicYear1}).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateAcademicYears(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/academic-years", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.AcademicYearsArray,
			func(input any) any {
				var modified domain.AcademicYears
				utils.Copy(input, &modified)
				modified[0].Name = "updated-year"
				return modified
			}, http.StatusOK, ta.AcademicYearsArray,
		).
		NewRequestWithBody("NotFound", domain.AcademicYears{&domain.AcademicYear{ID: 0, Name: "ghost"}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestActivateAcademicYear(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/academic-years/{academicYearId}/activate", http.MethodPost, adminToken).
		NewRequestWithoutBody("OK", nil, map[string]any{"academicYearId": ta.AcademicYear1.ID}, http.StatusOK, nil).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"academicYearId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeactivateAcademicYear(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/academic-years/{academicYearId}/deactivate", http.MethodPost, adminToken).
		NewRequestWithoutBody("OK", nil, map[string]any{"academicYearId": ta.AcademicYear1.ID}, http.StatusOK, nil).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"academicYearId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteAcademicYears(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/academic-years/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.AcademicYear2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
