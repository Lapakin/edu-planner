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
	"github.com/Lapakin/edu-planner/internal/syllabus/service"
	"github.com/Lapakin/edu-planner/internal/utils"
)

func NewCreateWorkloadDistributionsHandler(svc service.WorkloadDistributionSvc) gin.HandlerFunc {
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

		var distributions domain.WorkloadDistributions
		if err = json.Unmarshal(body, &distributions); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(distributions) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.CreateWorkloadDistributions(c.Request.Context(), claims, distributions); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, distributions)
	}
}

func NewGetWorkloadDistributionByIDHandler(svc service.WorkloadDistributionSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.StringToUint64(c.Param("workloadDistributionId"))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		distribution, err := svc.GetWorkloadDistributionByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, distribution)
	}
}

func NewFetchWorkloadDistributionsHandler(svc service.WorkloadDistributionSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters, err := rest.CreateFiltersFromQueries(c, rest.Queries{
			{Param: domain.IDsParam, ValidateFunc: utils.ValidateSliceUInt64},
			{Param: "study_plan_id", ValidateFunc: utils.ValidateUInt64},
			{Param: "group_id", ValidateFunc: utils.ValidateUInt64},
		})
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrParseQuery)
			return
		}

		distributions, err := svc.FetchWorkloadDistributions(c.Request.Context(), filters)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, distributions)
	}
}

func NewUpdateWorkloadDistributionsHandler(svc service.WorkloadDistributionSvc) gin.HandlerFunc {
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

		var distributions domain.WorkloadDistributions
		if err = json.Unmarshal(body, &distributions); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(distributions) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.UpdateWorkloadDistributions(c.Request.Context(), claims, distributions); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, distributions)
	}
}

func NewDeleteWorkloadDistributionByIDsHandler(svc service.WorkloadDistributionSvc) gin.HandlerFunc {
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

		var ids []uint64
		if err = json.Unmarshal(body, &ids); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(ids) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.DeleteWorkloadDistributions(c.Request.Context(), claims, ids); err != nil {
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
