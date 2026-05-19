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

func NewCreateWorkloadAssignmentsHandler(svc service.WorkloadAssignmentSvc) gin.HandlerFunc {
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

		var assignments domain.WorkloadAssignments
		if err = json.Unmarshal(body, &assignments); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(assignments) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.CreateWorkloadAssignments(c.Request.Context(), claims, assignments); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, assignments)
	}
}

func NewGetWorkloadAssignmentByIDHandler(svc service.WorkloadAssignmentSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.StringToUint64(c.Param("workloadAssignmentId"))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		assignment, err := svc.GetWorkloadAssignmentByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, assignment)
	}
}

func NewFetchWorkloadAssignmentsHandler(svc service.WorkloadAssignmentSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters, err := rest.CreateFiltersFromQueries(c, rest.Queries{
			{Param: domain.IDsParam, ValidateFunc: utils.ValidateSliceUInt64},
			{Param: "workload_distribution_id", ValidateFunc: utils.ValidateUInt64},
			{Param: "teacher_id", ValidateFunc: utils.ValidateUInt64},
		})
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrParseQuery)
			return
		}

		assignments, err := svc.FetchWorkloadAssignments(c.Request.Context(), filters)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, assignments)
	}
}

func NewUpdateWorkloadAssignmentsHandler(svc service.WorkloadAssignmentSvc) gin.HandlerFunc {
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

		var assignments domain.WorkloadAssignments
		if err = json.Unmarshal(body, &assignments); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(assignments) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.UpdateWorkloadAssignments(c.Request.Context(), claims, assignments); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, assignments)
	}
}

func NewDeleteWorkloadAssignmentByIDsHandler(svc service.WorkloadAssignmentSvc) gin.HandlerFunc {
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

		if err = svc.DeleteWorkloadAssignments(c.Request.Context(), claims, ids); err != nil {
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
