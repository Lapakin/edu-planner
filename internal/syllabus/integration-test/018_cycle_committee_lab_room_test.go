package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateCycleCommitteeLabRooms(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committee-lab-rooms", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.CycleCommitteeLabRoomsArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetCycleCommitteeLabRoomByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committee-lab-rooms/{cycleCommitteeLabRoomId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"cycleCommitteeLabRoomId": ta.CycleCommitteeLabRoom1.ID}, ta.CycleCommitteeLabRoom1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"cycleCommitteeLabRoomId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchCycleCommitteeLabRooms(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committee-lab-rooms", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.CycleCommitteeLabRoomsArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.CycleCommitteeLabRoom1.ID}}, nil, domain.CycleCommitteeLabRooms{ta.CycleCommitteeLabRoom1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.CycleCommitteeLabRoomsArray).
		NewOKRequestWithoutBody(map[string]any{domain.CycleCommitteeIDParam: ta.CycleCommittee1.ID}, nil, ta.CycleCommitteeLabRoomsArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateCycleCommitteeLabRooms(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committee-lab-rooms", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.CycleCommitteeLabRoomsArray,
			func(input any) any {
				var modified domain.CycleCommitteeLabRooms
				utils.Copy(input, &modified)
				modified[0].CycleCommitteeID = ta.CycleCommittee2.ID
				return modified
			}, http.StatusOK, ta.CycleCommitteeLabRoomsArray,
		).
		NewRequestWithBody("NotFound", domain.CycleCommitteeLabRooms{&domain.CycleCommitteeLabRoom{
			ID: 0, CycleCommitteeID: ta.CycleCommittee1.ID, RoomID: ta.Room1.ID,
		}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteCycleCommitteeLabRooms(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/cycle-committee-lab-rooms/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.CycleCommitteeLabRoom2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
