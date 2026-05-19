package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateCycleCommittees(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committees", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.CycleCommitteesArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetCycleCommitteeByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committees/{cycleCommitteeId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"cycleCommitteeId": ta.CycleCommittee1.ID}, ta.CycleCommittee1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"cycleCommitteeId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchCycleCommittees(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committees", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.CycleCommitteesArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.CycleCommittee1.ID}}, nil, domain.CycleCommittees{ta.CycleCommittee1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.CycleCommitteesArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateCycleCommittees(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committees", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.CycleCommitteesArray,
			func(input any) any {
				var modified domain.CycleCommittees
				utils.Copy(input, &modified)
				modified[0].Name = "updated-cycle-committee"
				return modified
			}, http.StatusOK, ta.CycleCommitteesArray,
		).
		NewRequestWithBody("NotFound", domain.CycleCommittees{&domain.CycleCommittee{ID: 0, AcademicYearID: ta.AcademicYear1.ID, UserID: ta.User1.ID, Name: "ghost"}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteCycleCommittees(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committees/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.CycleCommittee2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
