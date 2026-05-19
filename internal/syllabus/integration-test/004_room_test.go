package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateRooms(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/rooms", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.RoomsArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetRoomByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/rooms/{roomId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"roomId": ta.Room1.ID}, ta.Room1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"roomId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchRooms(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/rooms", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.RoomsArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.Room1.ID}}, nil, domain.Rooms{ta.Room1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.RoomsArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateRooms(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/rooms", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.RoomsArray,
			func(input any) any {
				var modified domain.Rooms
				utils.Copy(input, &modified)
				modified[0].Name = "updated-room"
				return modified
			}, http.StatusOK, ta.RoomsArray,
		).
		NewRequestWithBody("NotFound", domain.Rooms{&domain.Room{ID: 0, AcademicYearID: ta.AcademicYear1.ID, Name: "ghost", RoomType: domain.RoomTypeAuditorium}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteRooms(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/rooms/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.Room2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
