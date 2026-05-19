package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateUsers(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/users", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.UsersArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetUserByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/users/{userId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"userId": ta.User1.ID}, ta.User1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"userId": 0}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchUsers(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/users", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.UsersArray).
		NewOKRequestWithoutBody(map[string]any{domain.RoleParam: ta.User1.Role}, nil, ta.UsersArray).
		NewOKRequestWithoutBody(map[string]any{domain.RoleParam: "admin"}, nil, domain.Users{}).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateUsers(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/users", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.UsersArray,
			func(input any) any {
				var modified domain.Users
				utils.Copy(input, &modified)
				modified[0].Email = "updated1@test.com"
				return modified
			}, http.StatusOK, ta.UsersArray,
		).
		NewRequestWithBody("NotFound", domain.Users{&domain.User{ID: 0, Email: "ghost@test.com"}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestActivateUser(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/users/{userId}/activate", http.MethodPost, adminToken).
		NewRequestWithoutBody("OK", nil, map[string]any{"userId": ta.User1.ID}, http.StatusOK, nil).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"userId": 0}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeactivateUser(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/users/{userId}/deactivate", http.MethodPost, adminToken).
		NewRequestWithoutBody("OK", nil, map[string]any{"userId": ta.User1.ID}, http.StatusOK, nil).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"userId": 0}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteUsers(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/users/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.User2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{0}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
