package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateSpecialties(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/specialties", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.SpecialtiesArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetSpecialtyByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/specialties/{specialtyId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"specialtyId": ta.Specialty1.ID}, ta.Specialty1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"specialtyId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchSpecialties(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/specialties", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.SpecialtiesArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.Specialty1.ID}}, nil, domain.Specialties{ta.Specialty1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.SpecialtiesArray).
		NewOKRequestWithoutBody(map[string]any{domain.DepartmentIDParam: ta.Department1.ID}, nil, ta.SpecialtiesArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateSpecialties(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/specialties", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.SpecialtiesArray,
			func(input any) any {
				var modified domain.Specialties
				utils.Copy(input, &modified)
				modified[0].Name = "updated-specialty"
				return modified
			}, http.StatusOK, ta.SpecialtiesArray,
		).
		NewRequestWithBody("NotFound", domain.Specialties{&domain.Specialty{ID: 0, AcademicYearID: ta.AcademicYear1.ID, DepartmentID: ta.Department1.ID, Name: "ghost", ShortName: "gh", SpecialtyCode: "000"}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteSpecialties(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/specialties/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.Specialty2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
