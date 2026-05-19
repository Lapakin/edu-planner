package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateStudyPlans(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/study-plans", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.StudyPlansArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetStudyPlanByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/study-plans/{studyPlanId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"studyPlanId": ta.StudyPlan1.ID}, ta.StudyPlan1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"studyPlanId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchStudyPlans(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/study-plans", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.StudyPlansArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.StudyPlan1.ID}}, nil, domain.StudyPlans{ta.StudyPlan1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.StudyPlansArray).
		NewOKRequestWithoutBody(map[string]any{domain.SpecialtyIDParam: ta.Specialty1.ID}, nil, domain.StudyPlans{ta.StudyPlan1}).
		NewOKRequestWithoutBody(map[string]any{domain.DisciplineIDParam: ta.Discipline1.ID}, nil, domain.StudyPlans{ta.StudyPlan1}).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateStudyPlans(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/study-plans", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.StudyPlansArray,
			func(input any) any {
				var modified domain.StudyPlans
				utils.Copy(input, &modified)
				modified[0].Lectures = new(25)
				return modified
			}, http.StatusOK, ta.StudyPlansArray,
		).
		NewRequestWithBody("NotFound", domain.StudyPlans{&domain.StudyPlan{ID: 0, AcademicYearID: new(ta.AcademicYear1.ID)}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteStudyPlans(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/study-plans/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.StudyPlan2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
