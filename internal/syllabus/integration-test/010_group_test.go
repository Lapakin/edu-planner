package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateGroups(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.GroupsArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetGroupByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups/{groupId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"groupId": ta.Group1.ID}, ta.Group1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"groupId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchGroups(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.GroupsArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.Group1.ID}}, nil, domain.Groups{ta.Group1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.GroupsArray).
		NewOKRequestWithoutBody(map[string]any{domain.SpecialtyIDParam: ta.Specialty1.ID}, nil, ta.GroupsArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateGroups(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.GroupsArray,
			func(input any) any {
				var modified domain.Groups
				utils.Copy(input, &modified)
				modified[0].Name = "CS-1-updated"
				return modified
			}, http.StatusOK, ta.GroupsArray,
		).
		NewRequestWithBody("NotFound", domain.Groups{&domain.Group{ID: 0, AcademicYearID: ta.AcademicYear1.ID, SpecialtyID: ta.Specialty1.ID, Name: "ghost", ShortName: "gh"}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteGroups(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.Group2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestCreateGroupSemesters(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups/semesters", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.GroupSemestersArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetGroupSemesterByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups/semesters/{semesterId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"semesterId": ta.GroupSemester1.ID}, ta.GroupSemester1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"semesterId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchGroupSemesters(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups/semesters", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.GroupSemestersArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.GroupSemester1.ID}}, nil, domain.GroupSemesters{ta.GroupSemester1}).
		NewOKRequestWithoutBody(map[string]any{domain.GroupIDParam: ta.Group1.ID}, nil, domain.GroupSemesters{ta.GroupSemester1}).
		NewOKRequestWithoutBody(map[string]any{domain.SemesterIDParam: ta.Semester1.ID}, nil, ta.GroupSemestersArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateGroupSemesters(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups/semesters", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.GroupSemestersArray,
			func(input any) any {
				var modified domain.GroupSemesters
				utils.Copy(input, &modified)
				modified[0].SemesterID = ta.Semester1.ID
				return modified
			}, http.StatusOK, ta.GroupSemestersArray,
		).
		NewRequestWithBody("NotFound", domain.GroupSemesters{&domain.GroupSemester{ID: 0, GroupID: ta.Group1.ID, SemesterID: ta.Semester1.ID}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteGroupSemesters(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/groups/semesters/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.GroupSemester2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
