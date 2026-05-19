package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateDepartments(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/departments", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.DepartmentsArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetDepartmentByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/departments/{departmentId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"departmentId": ta.Department1.ID}, ta.Department1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"departmentId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchDepartments(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/departments", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.DepartmentsArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.Department1.ID}}, nil, domain.Departments{ta.Department1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.DepartmentsArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateDepartments(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/departments", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.DepartmentsArray,
			func(input any) any {
				var modified domain.Departments
				utils.Copy(input, &modified)
				modified[0].Name = "updated-department"
				return modified
			}, http.StatusOK, ta.DepartmentsArray,
		).
		NewRequestWithBody("NotFound", domain.Departments{&domain.Department{ID: 0, AcademicYearID: ta.AcademicYear1.ID, Name: "ghost"}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteDepartments(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/departments/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.Department2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
