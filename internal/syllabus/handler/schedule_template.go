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
	"github.com/Lapakin/edu-planner/internal/syllabus/generator"
	"github.com/Lapakin/edu-planner/internal/syllabus/service"
	"github.com/Lapakin/edu-planner/internal/utils"
)

func NewGenerateScheduleTemplateHandler(svc service.ScheduleTemplateSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrRequestBodyReading)
			return
		}
		defer c.Request.Body.Close()

		var req domain.GenerateScheduleRequest
		if err = json.Unmarshal(body, &req); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if req.SemesterID == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyID)
			return
		}

		data, err := svc.GenerateScheduleTemplate(c.Request.Context(), &req)
		if err != nil {
			switch {
			case errors.Is(err, generator.ErrTimeout):
				rest.RespondError(c, http.StatusGatewayTimeout, err)
			case errors.Is(err, generator.ErrSemesterNotFound):
				rest.RespondError(c, http.StatusNotFound, err)
			case errors.Is(err, generator.ErrNoSettings),
				errors.Is(err, generator.ErrNoBellSchedule),
				errors.Is(err, generator.ErrNoGroupSemesters),
				errors.Is(err, generator.ErrNoWorkloads):
				rest.RespondError(c, http.StatusBadRequest, err)
			case errors.Is(err, generator.ErrValidationFailed):
				rest.RespondError(c, http.StatusUnprocessableEntity, err)
			case errors.Is(err, generator.ErrUnableToSchedule),
				errors.Is(err, generator.ErrWorkloadDistribution):
				rest.RespondError(c, http.StatusConflict, err)
			default:
				rest.RespondError(c, http.StatusInternalServerError, err)
			}
			return
		}

		rest.RespondJSON(c, http.StatusOK, data)
	}
}

func NewSaveScheduleTemplateHandler(svc service.ScheduleTemplateSvc) gin.HandlerFunc {
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

		var req domain.SaveScheduleTemplateRequest
		if err = json.Unmarshal(body, &req); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if req.SemesterID == 0 || req.Data == nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		template, err := svc.SaveScheduleTemplate(c.Request.Context(), claims, &req)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusCreated, template)
	}
}

func NewGetScheduleTemplateByIDHandler(svc service.ScheduleTemplateSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.StringToUint64(c.Param("scheduleTemplateId"))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		template, err := svc.GetScheduleTemplateByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, template)
	}
}

func NewFetchScheduleTemplatesHandler(svc service.ScheduleTemplateSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters, err := rest.CreateFiltersFromQueries(c, rest.Queries{
			{Param: domain.IDsParam, ValidateFunc: utils.ValidateSliceUInt64},
			{Param: domain.SemesterIDParam, ValidateFunc: utils.ValidateUInt64},
		})
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrParseQuery)
			return
		}

		templates, err := svc.FetchScheduleTemplates(c.Request.Context(), filters)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, templates)
	}
}

func NewDeleteScheduleTemplateByIDsHandler(svc service.ScheduleTemplateSvc) gin.HandlerFunc {
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

		if err = svc.DeleteScheduleTemplates(c.Request.Context(), claims, ids); err != nil {
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

func NewActivateScheduleTemplateHandler(svc service.ScheduleTemplateSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.ExtractClaims(c.GetHeader(jwt.AuthHeader))
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, jwt.ErrInvalidToken)
			return
		}

		id, err := utils.StringToUint64(c.Param("scheduleTemplateId"))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		if err = svc.ActivateScheduleTemplate(c.Request.Context(), claims, id); err != nil {
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
