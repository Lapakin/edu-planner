package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Lapakin/edu-planner/internal/adapter/http/rest"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/user-management/service"
	"github.com/Lapakin/edu-planner/internal/utils"
)

func NewCreateInviteHandler(svc service.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.ExtractClaims(c.GetHeader(jwt.AuthHeader))
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, jwt.ErrInvalidToken)
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		defer c.Request.Body.Close()

		var req domain.InviteReq
		if err = json.Unmarshal(body, &req); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if req.Email == "" || req.FirstName == "" || req.LastName == "" {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		resp, err := svc.CreateInvite(c.Request.Context(), claims, req)
		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				rest.RespondError(c, http.StatusForbidden, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, resp)
	}
}

func NewSetPasswordHandler(svc service.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		defer c.Request.Body.Close()

		var req domain.SetPasswordReq
		if err = json.Unmarshal(body, &req); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if req.Token == "" || req.Password == "" {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		token, err := svc.SetPassword(c.Request.Context(), req.Token, req.Password)
		if err != nil {
			if errors.Is(err, service.ErrInviteTokenInvalid) {
				rest.RespondError(c, http.StatusGone, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, gin.H{"token": token})
	}
}

func NewResetInviteHandler(svc service.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.ExtractClaims(c.GetHeader(jwt.AuthHeader))
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, jwt.ErrInvalidToken)
			return
		}

		userID, err := utils.StringToUint64(c.Param("userId"))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		resp, err := svc.ResetInvite(c.Request.Context(), claims, userID)
		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				rest.RespondError(c, http.StatusForbidden, err)
				return
			}
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, resp)
	}
}

func NewGenerateResetPasswordLinkHandler(svc service.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.ExtractClaims(c.GetHeader(jwt.AuthHeader))
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, jwt.ErrInvalidToken)
			return
		}

		userID, err := utils.StringToUint64(c.Param("userId"))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		resp, err := svc.GenerateResetPasswordLink(c.Request.Context(), claims, userID)
		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				rest.RespondError(c, http.StatusForbidden, err)
				return
			}
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, resp)
	}
}

func NewResetPasswordHandler(svc service.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		defer c.Request.Body.Close()

		var req domain.ResetPasswordReq
		if err = json.Unmarshal(body, &req); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if req.Token == "" || req.Password == "" {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		token, err := svc.ResetPassword(c.Request.Context(), req.Token, req.Password)
		if err != nil {
			if errors.Is(err, service.ErrInviteTokenInvalid) {
				rest.RespondError(c, http.StatusGone, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, gin.H{"token": token})
	}
}

func NewLoginHandler(svc service.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		c.Request.Body.Close()

		var req domain.LoginReq
		if err = json.Unmarshal(body, &req); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		token, err := svc.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, gin.H{"token": token})
	}
}
