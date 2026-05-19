package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateTeachers(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teachers", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.TeachersArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetTeacherByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teachers/{teacherId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"teacherId": ta.Teacher1.ID}, ta.Teacher1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"teacherId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchTeachers(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teachers", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.TeachersArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.Teacher1.ID}}, nil, domain.Teachers{ta.Teacher1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.TeachersArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateTeachers(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teachers", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.TeachersArray,
			func(input any) any {
				var modified domain.Teachers
				utils.Copy(input, &modified)
				modified[0].UserID = ta.User1.ID
				return modified
			}, http.StatusOK, ta.TeachersArray,
		).
		NewRequestWithBody("NotFound", domain.Teachers{&domain.Teacher{ID: 0, AcademicYearID: ta.AcademicYear1.ID, UserID: ta.User1.ID}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteTeachers(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teachers/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.Teacher2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
