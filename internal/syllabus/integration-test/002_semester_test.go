package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateSemesters(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/semesters", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.SemestersArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetSemesterByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/semesters/{semesterId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"semesterId": ta.Semester1.ID}, ta.Semester1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"semesterId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchSemesters(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/semesters", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.SemestersArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.Semester1.ID}}, nil, domain.Semesters{ta.Semester1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.SemestersArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateSemesters(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/semesters", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.SemestersArray,
			func(input any) any {
				var modified domain.Semesters
				utils.Copy(input, &modified)
				modified[0].AcademicYearID = ta.AcademicYear1.ID
				return modified
			}, http.StatusOK, ta.SemestersArray,
		).
		NewRequestWithBody("NotFound", domain.Semesters{&domain.Semester{ID: 0, AcademicYearID: ta.AcademicYear1.ID}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteSemesters(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/semesters/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.Semester2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
