package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateWorkloadDistributions(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-distributions", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.WorkloadDistributionsArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetWorkloadDistributionByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-distributions/{workloadDistributionId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"workloadDistributionId": ta.WorkloadDistribution1.ID}, ta.WorkloadDistribution1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"workloadDistributionId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchWorkloadDistributions(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-distributions", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.WorkloadDistributionsArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.WorkloadDistribution1.ID}}, nil, domain.WorkloadDistributions{ta.WorkloadDistribution1}).
		NewOKRequestWithoutBody(map[string]any{"study_plan_id": ta.StudyPlan1.ID}, nil, ta.WorkloadDistributionsArray).
		NewOKRequestWithoutBody(map[string]any{"group_id": ta.Group1.ID}, nil, domain.WorkloadDistributions{ta.WorkloadDistribution1}).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateWorkloadDistributions(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-distributions", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.WorkloadDistributionsArray,
			func(input any) any {
				var modified domain.WorkloadDistributions
				utils.Copy(input, &modified)
				modified[0].GroupID = ta.Group1.ID
				return modified
			}, http.StatusOK, ta.WorkloadDistributionsArray,
		).
		NewRequestWithBody("NotFound", domain.WorkloadDistributions{&domain.WorkloadDistribution{ID: 0, StudyPlanID: ta.StudyPlan1.ID, GroupID: ta.Group1.ID}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteWorkloadDistributions(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-distributions/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.WorkloadDistribution2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
