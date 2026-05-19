package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateWorkloadAssignments(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-assignments", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.WorkloadAssignmentsArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetWorkloadAssignmentByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-assignments/{workloadAssignmentId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"workloadAssignmentId": ta.WorkloadAssignment1.ID}, ta.WorkloadAssignment1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"workloadAssignmentId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchWorkloadAssignments(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-assignments", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.WorkloadAssignmentsArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.WorkloadAssignment1.ID}}, nil, domain.WorkloadAssignments{ta.WorkloadAssignment1}).
		NewOKRequestWithoutBody(map[string]any{"workload_distribution_id": ta.WorkloadDistribution1.ID}, nil, ta.WorkloadAssignmentsArray).
		NewOKRequestWithoutBody(map[string]any{"teacher_id": ta.Teacher1.ID}, nil, domain.WorkloadAssignments{ta.WorkloadAssignment1}).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateWorkloadAssignments(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-assignments", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.WorkloadAssignmentsArray,
			func(input any) any {
				var modified domain.WorkloadAssignments
				utils.Copy(input, &modified)
				modified[0].AssignedHours = new(55)
				return modified
			}, http.StatusOK, ta.WorkloadAssignmentsArray,
		).
		NewRequestWithBody("NotFound", domain.WorkloadAssignments{&domain.WorkloadAssignment{ID: 0, WorkloadDistributionID: ta.WorkloadDistribution1.ID, TeacherID: ta.Teacher1.ID}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteWorkloadAssignments(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/workload-assignments/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.WorkloadAssignment2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
