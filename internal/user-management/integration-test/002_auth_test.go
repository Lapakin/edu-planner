package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Lapakin/edu-planner/internal/domain"
	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

// capturedInviteToken is set by TestCreateInvite and consumed by TestSetPassword / TestLogin.
var capturedInviteToken string

func TestCreateInvite(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/auth/invite", http.MethodPost, adminToken).
		NewRequestWithBody("OK", ta.InviteReq1, http.StatusOK, ta.ExpectedInviteResp1).
		Run(t, ts.URL, []string{"created_at", "modified_at", "id", "invite_link", "is_active", "is_deleted", "patronymic", "first_name", "last_name"})

	ta.NewHTTPCasesRunner("/api/v1/auth/invite", http.MethodPost, deanToken).
		NewRequestWithBody("DeanAccess", ta.DeanInviteReq, http.StatusOK, ta.ExpectedDeanInviteResp).
		Run(t, ts.URL, []string{"created_at", "modified_at", "id", "invite_link", "is_active", "is_deleted", "patronymic", "first_name", "last_name"})

	ta.NewHTTPCasesRunner("/api/v1/auth/invite", http.MethodPost, teacherToken).
		NewRequestWithBody("NoAccess", ta.InviteReq1, http.StatusForbidden, ta.ExpectedResponseHaveNoAccess).
		Run(t, ts.URL, []string{})

	ta.NewHTTPCasesRunner("/api/v1/auth/invite", http.MethodPost, adminToken).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		Run(t, ts.URL, []string{})

	// Create a fresh invite so TestSetPassword has a valid token to consume.
	capturedInviteToken = createInviteAndExtractToken(t, domain.InviteReq{
		Email:     "setpass@test.com",
		FirstName: "SetPass",
		LastName:  "User",
		Role:      "user",
	})
}

func TestSetPassword(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/set-password", http.MethodPost, ta.EmptyToken).
		NewRequestWithBody("InvalidToken",
			&domain.SetPasswordReq{Token: "nonexistent-token", Password: "pass123"},
			http.StatusGone,
			map[string]any{"error": "invite link is invalid or has already been used"},
		).
		NewUnmarshalErrorRequest().
		Run(t, ts.URL, []string{})

	require.NotEmpty(t, capturedInviteToken, "capturedInviteToken must be set by TestCreateInvite")

	ta.NewHTTPCasesRunner("/api/v1/set-password", http.MethodPost, ta.EmptyToken).
		NewRequestWithBody("OK",
			&domain.SetPasswordReq{Token: capturedInviteToken, Password: "securepass123"},
			http.StatusOK,
			ta.ExpectedTokenResp,
		).
		Run(t, ts.URL, []string{"token"})

	ta.NewHTTPCasesRunner("/api/v1/set-password", http.MethodPost, ta.EmptyToken).
		NewRequestWithBody("UsedToken",
			&domain.SetPasswordReq{Token: capturedInviteToken, Password: "securepass123"},
			http.StatusGone,
			map[string]any{"error": "invite link is invalid or has already been used"},
		).
		Run(t, ts.URL, []string{})
}

func TestLogin(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/login", http.MethodPost, ta.EmptyToken).
		NewRequestWithBody("OK",
			&domain.LoginReq{Email: "setpass@test.com", Password: "securepass123"},
			http.StatusOK,
			ta.ExpectedTokenResp,
		).
		NewRequestWithBody("InvalidCredentials",
			&domain.LoginReq{Email: "setpass@test.com", Password: "wrongpassword"},
			http.StatusUnauthorized,
			ta.ExpectedResponseInvalidCreds,
		).
		NewUnmarshalErrorRequest().
		Run(t, ts.URL, []string{"token"})
}

func TestResetInvite(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/users/1/reset-invite", http.MethodPost, adminToken).
		NewRequestWithBody("OK", nil, http.StatusOK, ta.ExpectedResetInviteResp1).
		Run(t, ts.URL, []string{"created_at", "modified_at", "id", "invite_link", "is_active", "is_deleted", "patronymic", "first_name", "last_name", "email"})

	ta.NewHTTPCasesRunner("/api/v1/users/1/reset-invite", http.MethodPost, deanToken).
		NewRequestWithBody("DeanAccess", nil, http.StatusOK, ta.ExpectedResetInviteResp1).
		Run(t, ts.URL, []string{"created_at", "modified_at", "id", "invite_link", "is_active", "is_deleted", "patronymic", "first_name", "last_name", "email"})

	ta.NewHTTPCasesRunner("/api/v1/users/1/reset-invite", http.MethodPost, teacherToken).
		NewRequestWithBody("NoAccess", nil, http.StatusForbidden, ta.ExpectedResponseHaveNoAccess).
		Run(t, ts.URL, []string{})

	ta.NewHTTPCasesRunner("/api/v1/users/0/reset-invite", http.MethodPost, adminToken).
		NewRequestWithBody("NotFound", nil, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, []string{})

	ta.NewHTTPCasesRunner("/api/v1/users/1/reset-invite", http.MethodPost, adminToken).
		NewBadJWTRequest().
		Run(t, ts.URL, []string{})
}

func TestGenerateResetPasswordLink(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/users/1/generate-reset-password-link", http.MethodPost, adminToken).
		NewRequestWithBody("OK", nil, http.StatusOK, ta.ExpectedResetPasswordLinkResp1).
		Run(t, ts.URL, []string{"created_at", "modified_at", "id", "invite_link", "is_active", "is_deleted", "patronymic", "first_name", "last_name", "email"})

	ta.NewHTTPCasesRunner("/api/v1/users/1/generate-reset-password-link", http.MethodPost, deanToken).
		NewRequestWithBody("DeanAccess", nil, http.StatusOK, ta.ExpectedResetPasswordLinkResp1).
		Run(t, ts.URL, []string{"created_at", "modified_at", "id", "invite_link", "is_active", "is_deleted", "patronymic", "first_name", "last_name", "email"})

	ta.NewHTTPCasesRunner("/api/v1/users/1/generate-reset-password-link", http.MethodPost, teacherToken).
		NewRequestWithBody("NoAccess", nil, http.StatusForbidden, ta.ExpectedResponseHaveNoAccess).
		Run(t, ts.URL, []string{})

	ta.NewHTTPCasesRunner("/api/v1/users/0/generate-reset-password-link", http.MethodPost, adminToken).
		NewRequestWithBody("NotFound", nil, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, []string{})

	ta.NewHTTPCasesRunner("/api/v1/users/1/generate-reset-password-link", http.MethodPost, adminToken).
		NewBadJWTRequest().
		Run(t, ts.URL, []string{})
}

var capturedResetPasswordToken string

func TestResetPassword(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/reset-password", http.MethodPost, ta.EmptyToken).
		NewRequestWithBody("InvalidToken",
			&domain.ResetPasswordReq{Token: "nonexistent-token", Password: "pass123"},
			http.StatusGone,
			map[string]any{"error": "invite link is invalid or has already been used"},
		).
		NewUnmarshalErrorRequest().
		Run(t, ts.URL, []string{})

	capturedResetPasswordToken = generateResetPasswordLinkAndExtractToken(t, 1)
	require.NotEmpty(t, capturedResetPasswordToken, "capturedResetPasswordToken must be set")

	ta.NewHTTPCasesRunner("/api/v1/reset-password", http.MethodPost, ta.EmptyToken).
		NewRequestWithBody("OK",
			&domain.ResetPasswordReq{Token: capturedResetPasswordToken, Password: "newpass456"},
			http.StatusOK,
			ta.ExpectedTokenResp,
		).
		Run(t, ts.URL, []string{"token"})

	ta.NewHTTPCasesRunner("/api/v1/reset-password", http.MethodPost, ta.EmptyToken).
		NewRequestWithBody("UsedToken",
			&domain.ResetPasswordReq{Token: capturedResetPasswordToken, Password: "newpass456"},
			http.StatusGone,
			map[string]any{"error": "invite link is invalid or has already been used"},
		).
		Run(t, ts.URL, []string{})
}

// generateResetPasswordLinkAndExtractToken calls the admin endpoint and returns the raw token.
func generateResetPasswordLinkAndExtractToken(t *testing.T, userID int) string {
	t.Helper()

	httpReq, err := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/users/"+strconv.Itoa(userID)+"/generate-reset-password-link", nil)
	require.NoError(t, err)
	httpReq.Header.Set("Authorization", adminToken)

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBytes, _ := io.ReadAll(resp.Body)
	var inviteResp domain.InviteResp
	require.NoError(t, json.Unmarshal(respBytes, &inviteResp))

	parts := strings.SplitN(inviteResp.InviteLink, "token=", 2)
	require.Len(t, parts, 2, "invite_link must contain 'token='")
	return parts[1]
}

// createInviteAndExtractToken is a helper that performs the full invite HTTP call and returns the token.
func createInviteAndExtractToken(t *testing.T, req domain.InviteReq) string {
	t.Helper()

	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth/invite", bytes.NewBuffer(body))
	require.NoError(t, err)
	httpReq.Header.Set("Authorization", adminToken)

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBytes, _ := io.ReadAll(resp.Body)
	var inviteResp domain.InviteResp
	require.NoError(t, json.Unmarshal(respBytes, &inviteResp))

	parts := strings.SplitN(inviteResp.InviteLink, "token=", 2)
	require.Len(t, parts, 2, "invite_link must contain 'token='")
	return parts[1]
}
