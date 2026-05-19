package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Lapakin/edu-planner/internal/adapter/http/rest"
	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/user-management/service"
	"github.com/Lapakin/edu-planner/internal/utils"
)

func NewCreateUsersHandler(svc service.UserSvc) gin.HandlerFunc {
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

		var users domain.Users
		if err = json.Unmarshal(body, &users); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(users) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.CreateUsers(c.Request.Context(), claims, users); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, users)
	}
}

func NewGetUserByIDHandler(svc service.UserSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.StringToUint64(c.Param("userId"))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		user, err := svc.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, user)
	}
}

func NewFetchUsersHandler(svc service.UserSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters, err := rest.CreateFiltersFromQueries(c, rest.Queries{
			{Param: domain.RoleParam, ValidateFunc: nil},
		})
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrParseQuery)
			return
		}

		users, err := svc.FetchUsers(c.Request.Context(), filters)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, users)
	}
}

func NewUpdateUsersHandler(svc service.UserSvc) gin.HandlerFunc {
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

		var users domain.Users
		if err = json.Unmarshal(body, &users); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(users) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.UpdateUsers(c.Request.Context(), claims, users); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, users)
	}
}

func NewDeleteUsersHandler(svc service.UserSvc) gin.HandlerFunc {
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

		var userIDs []uint64
		if err = json.Unmarshal(body, &userIDs); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(userIDs) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.DeleteUsers(c.Request.Context(), claims, userIDs); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, nil)
	}
}

func NewActivateUserHandler(svc service.UserSvc) gin.HandlerFunc {
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

		if err = svc.ActivateUser(c.Request.Context(), claims, userID); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, nil)
	}
}

func NewDeactivateUserHandler(svc service.UserSvc) gin.HandlerFunc {
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

		if err = svc.DeactivateUser(c.Request.Context(), claims, userID); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, nil)
	}
}
