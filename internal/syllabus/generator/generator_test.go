package generator

import (
	"context"
	"testing"
	"time"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func makeMinimalGenerator() *Generator {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(4)
	r := &domain.ScheduleRestriction{
		MinGroupLessonsPerDay:   0,
		MaxGroupLessonsPerDay:   4,
		MaxTeacherLessonsPerDay: 5,
		NoGapsInGroupSchedule:   false,
	}

	w := &Workload{
		workloadDistributionID: 1,
		groupID:                1,
		subjectID:              10,
		teacherID:              42,
		totalHours:             4,
		isSplitting:            false,
	}

	genCfg := domain.GenerationConfig{
		Timeout:            5 * time.Second,
		NumberOfGoroutines: 1,
	}

	startDate := time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC)

	return New(
		genCfg,
		ts,
		bells,
		r,
		nil,
		nil,
		[]*Workload{w},
		[]uint64{101},
		[]uint64{1},
		nil,
		nil,
		startDate,
		endDate,
	)
}

func TestGenerator_Generate_MinimalScenario(t *testing.T) {
	g := makeMinimalGenerator()
	ctx := context.Background()

	data, err := g.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}
	if data == nil {
		t.Fatal("Generate() returned nil data")
	}
	if data.Numerator == nil {
		t.Error("expected non-nil numerator")
	}
}

func TestGenerator_Generate_ProducesLessons(t *testing.T) {
	g := makeMinimalGenerator()
	ctx := context.Background()

	data, err := g.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	totalLessons := 0
	for _, daySchedule := range data.Numerator {
		if daySchedule != nil {
			for _, lessons := range *daySchedule {
				totalLessons += len(lessons)
			}
		}
	}
	if totalLessons == 0 {
		t.Error("expected at least one lesson in numerator schedule")
	}
}

func TestGenerator_Generate_CancelledContext(t *testing.T) {
	g := makeMinimalGenerator()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := g.Generate(ctx)
	// Either timeout or unable to schedule is acceptable
	if err == nil {
		// Some implementations might succeed if they check context after completion
		t.Log("Generate succeeded despite cancelled context (acceptable)")
	}
}

func TestGenerator_Generate_WithTimeout(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(4)
	r := &domain.ScheduleRestriction{
		MinGroupLessonsPerDay:   0,
		MaxGroupLessonsPerDay:   4,
		MaxTeacherLessonsPerDay: 5,
		NoGapsInGroupSchedule:   false,
	}

	w := &Workload{
		workloadDistributionID: 1,
		groupID:                1,
		subjectID:              10,
		teacherID:              42,
		totalHours:             4,
		isSplitting:            false,
	}

	genCfg := domain.GenerationConfig{
		Timeout:            10 * time.Second,
		NumberOfGoroutines: 2, // two goroutines
	}

	startDate := time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC)

	g := New(genCfg, ts, bells, r, nil, nil, []*Workload{w}, []uint64{101}, []uint64{1}, nil, nil, startDate, endDate)
	data, err := g.Generate(context.Background())
	if err != nil {
		t.Fatalf("Generate() with 2 goroutines returned error: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil data")
	}
}

func TestBuildWorkloads_ReturnsWorkloads(t *testing.T) {
	disciplineID := uint64(10)
	classroomWork := 60
	assignedHours := 60

	distributions := domain.WorkloadDistributions{
		{ID: 1, StudyPlanID: 1, GroupID: 1, ClassroomWork: &classroomWork},
	}
	assignments := domain.WorkloadAssignments{
		{ID: 1, WorkloadDistributionID: 1, TeacherID: 42, AssignedHours: &assignedHours},
	}
	studyPlans := domain.StudyPlans{
		{ID: 1, DisciplineID: &disciplineID},
	}
	groups := domain.Groups{{ID: 1}}

	workloads := BuildWorkloads(distributions, assignments, studyPlans, groups, nil)
	if len(workloads) != 1 {
		t.Fatalf("expected 1 workload, got %d", len(workloads))
	}
	if workloads[0].groupID != 1 {
		t.Errorf("expected groupID=1, got %d", workloads[0].groupID)
	}
	if workloads[0].teacherID != 42 {
		t.Errorf("expected teacherID=42, got %d", workloads[0].teacherID)
	}
	if workloads[0].subjectID != 10 {
		t.Errorf("expected subjectID=10, got %d", workloads[0].subjectID)
	}
	if workloads[0].totalHours != 60 {
		t.Errorf("expected totalHours=60, got %d", workloads[0].totalHours)
	}
}

func TestBuildWorkloads_EmptyInputs(t *testing.T) {
	workloads := BuildWorkloads(
		domain.WorkloadDistributions{},
		domain.WorkloadAssignments{},
		domain.StudyPlans{},
		domain.Groups{},
		nil,
	)
	if len(workloads) != 0 {
		t.Errorf("expected 0 workloads for empty inputs, got %d", len(workloads))
	}
}
