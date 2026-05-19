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

func NewCreateSpecialtiesHandler(svc service.SpecialtySvc) gin.HandlerFunc {
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

		var specialties domain.Specialties
		if err = json.Unmarshal(body, &specialties); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(specialties) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.CreateSpecialties(c.Request.Context(), claims, specialties); err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, specialties)
	}
}

func NewGetSpecialtyByIDHandler(svc service.SpecialtySvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		specialtyID, err := utils.StringToUint64(c.Param("specialtyId"))
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		specialty, err := svc.GetSpecialtyByID(c.Request.Context(), specialtyID)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, specialty)
	}
}

func NewFetchSpecialtiesHandler(svc service.SpecialtySvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters, err := rest.CreateFiltersFromQueries(c, rest.Queries{
			{Param: domain.IDsParam, ValidateFunc: utils.ValidateSliceUInt64},
			{Param: domain.AcademicYearIDParam, ValidateFunc: utils.ValidateUInt64},
			{Param: domain.DepartmentIDParam, ValidateFunc: utils.ValidateSliceUInt64},
		})
		if err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrConvID)
			return
		}

		specialties, err := svc.FetchSpecialties(c.Request.Context(), filters)
		if err != nil {
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, specialties)
	}
}

func NewUpdateSpecialtiesHandler(svc service.SpecialtySvc) gin.HandlerFunc {
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

		var specialties domain.Specialties
		if err = json.Unmarshal(body, &specialties); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}

		if len(specialties) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.UpdateSpecialties(c.Request.Context(), claims, specialties); err != nil {
			if errors.Is(err, service.ErrNotFound) {
				rest.RespondError(c, http.StatusNotFound, err)
				return
			}
			rest.RespondError(c, http.StatusInternalServerError, err)
			return
		}

		rest.RespondJSON(c, http.StatusOK, specialties)
	}
}

func NewDeleteSpecialtyByIDsHandler(svc service.SpecialtySvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.ExtractClaims(c.GetHeader(jwt.AuthHeader))
		if err != nil {
			rest.RespondError(c, http.StatusUnauthorized, jwt.ErrInvalidToken)
			return
		}

		var ids []uint64
		if err = json.NewDecoder(c.Request.Body).Decode(&ids); err != nil {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrUnmarshal)
			return
		}
		defer c.Request.Body.Close()

		if len(ids) == 0 {
			rest.RespondError(c, http.StatusBadRequest, rest.ErrEmptyBody)
			return
		}

		if err = svc.DeleteSpecialties(c.Request.Context(), claims, ids); err != nil {
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
