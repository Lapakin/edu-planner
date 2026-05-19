package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateDisciplines(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/disciplines", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.DisciplinesArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetDisciplineByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/disciplines/{disciplineId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"disciplineId": ta.Discipline1.ID}, ta.Discipline1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"disciplineId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchDisciplines(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/disciplines", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.DisciplinesArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.Discipline1.ID}}, nil, domain.Disciplines{ta.Discipline1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.DisciplinesArray).
		NewOKRequestWithoutBody(map[string]any{domain.CycleCommitteeIDParam: ta.CycleCommittee1.ID}, nil, ta.DisciplinesArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateDisciplines(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/disciplines", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.DisciplinesArray,
			func(input any) any {
				var modified domain.Disciplines
				utils.Copy(input, &modified)
				modified[0].Name = "updated-discipline"
				return modified
			}, http.StatusOK, ta.DisciplinesArray,
		).
		NewRequestWithBody("NotFound", domain.Disciplines{&domain.Discipline{ID: 0, AcademicYearID: ta.AcademicYear1.ID, CycleCommitteeID: ta.CycleCommittee1.ID, Name: "ghost", ShortName: "gh"}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteDisciplines(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/disciplines/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.Discipline2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
