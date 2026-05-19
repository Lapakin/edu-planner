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

func NewCreateDisciplinesHandler(svc service.DisciplineSvc) gin.HandlerFunc {
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

		var disciplines domain.Disciplines
		if err = json.Unmarshal(body, &disciplines); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(disciplines) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.CreateDisciplines(c.Request.Context(), claims, disciplines); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, disciplines)
	}
}

func NewGetDisciplineByIDHandler(svc service.DisciplineSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		disciplineID, err := utils.StringToUint64(c.Param("disciplineId"))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		discipline, err := svc.GetDisciplineByID(c.Request.Context(), disciplineID)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, discipline)
	}
}

func NewFetchDisciplinesHandler(svc service.DisciplineSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters, err := rest.CreateFiltersFromQueries(c, rest.Queries{
			{Param: domain.IDsParam, ValidateFunc: utils.ValidateSliceUInt64},
			{Param: domain.AcademicYearIDParam, ValidateFunc: utils.ValidateUInt64},
			{Param: domain.CycleCommitteeIDParam, ValidateFunc: utils.ValidateUInt64},
		})
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		disciplines, err := svc.FetchDisciplines(c.Request.Context(), filters)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, disciplines)
	}
}

func NewUpdateDisciplinesHandler(svc service.DisciplineSvc) gin.HandlerFunc {
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

		var disciplines domain.Disciplines
		if err = json.Unmarshal(body, &disciplines); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(disciplines) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.UpdateDisciplines(c.Request.Context(), claims, disciplines); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, disciplines)
	}
}

func NewDeleteDisciplineByIDsHandler(svc service.DisciplineSvc) gin.HandlerFunc {
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

		if err = svc.DeleteDisciplines(c.Request.Context(), claims, ids); err != nil {
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
