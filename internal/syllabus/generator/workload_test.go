package generator

import (
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
)

func makeWorkload(totalHours int, isSplitting bool) *Workload {
	return &Workload{
		workloadDistributionID: 1,
		groupID:                1,
		subjectID:              10,
		teacherID:              42,
		totalHours:             totalHours,
		isSplitting:            isSplitting,
	}
}

func makeWorkloadDistributor(lessonsPerClass int) *workloadDistributor {
	ts := makeTemplateSetting(lessonsPerClass, 8, 40)
	bells := makeBellSchedules(5)
	cfg := newSettings(ts, bells, nil, nil, nil)
	return newWorkloadDistributor(cfg)
}

// --- distributeWorkload ---

func TestDistributeWorkload_ZeroHoursProducesNoLessons(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	w := makeWorkload(0, false)

	dist, err := wd.distributeWorkload(w, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dist.numeratorLessons != 0 || dist.denominatorLessons != 0 {
		t.Errorf("expected 0/0 lessons for zero hours, got %d/%d", dist.numeratorLessons, dist.denominatorLessons)
	}
}

func TestDistributeWorkload_EvenSplitBetweenNumeratorAndDenominator(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	w := makeWorkload(8, false) // 8h / 2h = 4 lessons total → 2 num + 2 denom

	dist, err := wd.distributeWorkload(w, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dist.numeratorLessons != 2 || dist.denominatorLessons != 2 {
		t.Errorf("expected 2/2, got %d/%d", dist.numeratorLessons, dist.denominatorLessons)
	}
}

func TestDistributeWorkload_OddTotalAssignsRemainderToNumerator(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	w := makeWorkload(6, false) // 6/2 = 3 lessons → 2 num + 1 denom

	dist, err := wd.distributeWorkload(w, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dist.numeratorLessons != 2 || dist.denominatorLessons != 1 {
		t.Errorf("expected 2/1, got %d/%d", dist.numeratorLessons, dist.denominatorLessons)
	}
}

func TestDistributeWorkload_MultipleWeeksDividesPerWeek(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	w := makeWorkload(16, false) // 16/2 = 8 lessons → 4 num + 4 denom; 2 per period per week over 2 weeks

	dist, err := wd.distributeWorkload(w, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dist.numeratorLessons != 2 || dist.denominatorLessons != 2 {
		t.Errorf("expected 2/2 per week, got %d/%d", dist.numeratorLessons, dist.denominatorLessons)
	}
}

// --- createLessons ---

func TestCreateLessons_UnitedFormat(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	w := makeWorkload(4, false)

	lessons := wd.createLessons(w, 3)
	if len(lessons) != 3 {
		t.Fatalf("expected 3 lessons, got %d", len(lessons))
	}
	for i, l := range lessons {
		if l.format != formatUnited {
			t.Errorf("lesson %d: expected formatUnited, got %v", i, l.format)
		}
		if len(l.subLessons) != 1 {
			t.Errorf("lesson %d: expected 1 sub-lesson, got %d", i, len(l.subLessons))
		}
		if l.subLessons[0].teacherID != w.teacherID {
			t.Errorf("lesson %d: expected teacherID %d, got %d", i, w.teacherID, l.subLessons[0].teacherID)
		}
	}
}

func TestCreateLessons_SplitFormatWhenSplittingWithSubGroup(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	subGroup := uint64(1)
	w := makeWorkload(4, true)
	w.subGroupNumber = &subGroup

	lessons := wd.createLessons(w, 2)
	if len(lessons) != 2 {
		t.Fatalf("expected 2 lessons, got %d", len(lessons))
	}
	for i, l := range lessons {
		if l.format != formatSplit {
			t.Errorf("lesson %d: expected formatSplit, got %v", i, l.format)
		}
	}
}

func TestCreateLessons_SplittingWithoutSubGroupIsUnited(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	w := makeWorkload(4, true) // isSplitting=true but no subGroupNumber

	lessons := wd.createLessons(w, 2)
	for i, l := range lessons {
		if l.format != formatUnited {
			t.Errorf("lesson %d: expected formatUnited when subGroupNumber is nil, got %v", i, l.format)
		}
	}
}

func TestCreateLessons_ZeroCountReturnsEmpty(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	w := makeWorkload(4, false)

	lessons := wd.createLessons(w, 0)
	if len(lessons) != 0 {
		t.Errorf("expected empty slice for count=0, got %d", len(lessons))
	}
}

// --- distributeWorkloads (full pipeline) ---

func TestDistributeWorkloads_SingleWorkload(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	w := makeWorkload(4, false) // 4/2 = 2 total → 1 num + 1 denom per week

	numLessons, denomLessons, err := wd.distributeWorkloads([]*Workload{w}, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(numLessons) == 0 {
		t.Error("expected at least one numerator lesson")
	}
	if len(denomLessons) == 0 {
		t.Error("expected at least one denominator lesson")
	}
}

func TestDistributeWorkloads_EmptyWorkloadsReturnsEmpty(t *testing.T) {
	wd := makeWorkloadDistributor(2)

	numLessons, denomLessons, err := wd.distributeWorkloads([]*Workload{}, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(numLessons) != 0 || len(denomLessons) != 0 {
		t.Errorf("expected empty results for empty workloads, got %d/%d", len(numLessons), len(denomLessons))
	}
}

func TestDistributeWorkloads_MultipleWorkloads(t *testing.T) {
	wd := makeWorkloadDistributor(2)
	w1 := makeWorkload(4, false)
	w2 := &Workload{
		workloadDistributionID: 2,
		groupID:                2,
		subjectID:              20,
		teacherID:              43,
		totalHours:             8,
		isSplitting:            false,
	}

	numLessons, _, err := wd.distributeWorkloads([]*Workload{w1, w2}, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(numLessons) < 2 {
		t.Errorf("expected at least 2 numerator lessons for 2 workloads, got %d", len(numLessons))
	}
}

// --- buildWorkloadsFromDomain ---

func TestBuildWorkloadsFromDomain_SkipsDistributionWithNoStudyPlan(t *testing.T) {
	dist := domain.WorkloadDistributions{
		{ID: 1, StudyPlanID: 99, GroupID: 1}, // study plan 99 does not exist
	}
	assignments := domain.WorkloadAssignments{
		{ID: 1, WorkloadDistributionID: 1, TeacherID: 1},
	}
	studyPlans := domain.StudyPlans{}
	groups := domain.Groups{{ID: 1}}

	result := buildWorkloadsFromDomain(dist, assignments, studyPlans, groups, nil)
	if len(result) != 0 {
		t.Errorf("expected 0 workloads when study plan is missing, got %d", len(result))
	}
}

func TestBuildWorkloadsFromDomain_SkipsZeroClassroomHours(t *testing.T) {
	disciplineID := uint64(10)
	classroomWork := 0
	dist := domain.WorkloadDistributions{
		{ID: 1, StudyPlanID: 1, GroupID: 1, ClassroomWork: &classroomWork},
	}
	assignments := domain.WorkloadAssignments{
		{ID: 1, WorkloadDistributionID: 1, TeacherID: 1},
	}
	studyPlans := domain.StudyPlans{
		{ID: 1, DisciplineID: &disciplineID},
	}
	groups := domain.Groups{{ID: 1}}

	result := buildWorkloadsFromDomain(dist, assignments, studyPlans, groups, nil)
	if len(result) != 0 {
		t.Errorf("expected 0 workloads for zero classroom hours, got %d", len(result))
	}
}

func TestBuildWorkloadsFromDomain_BuildsWorkloadCorrectly(t *testing.T) {
	disciplineID := uint64(10)
	classroomWork := 60
	dist := domain.WorkloadDistributions{
		{ID: 1, StudyPlanID: 1, GroupID: 1, ClassroomWork: &classroomWork},
	}
	assignedHours := 60
	assignments := domain.WorkloadAssignments{
		{ID: 1, WorkloadDistributionID: 1, TeacherID: 42, AssignedHours: &assignedHours},
	}
	studyPlans := domain.StudyPlans{
		{ID: 1, DisciplineID: &disciplineID},
	}
	groups := domain.Groups{{ID: 1}}

	result := buildWorkloadsFromDomain(dist, assignments, studyPlans, groups, nil)
	if len(result) != 1 {
		t.Fatalf("expected 1 workload, got %d", len(result))
	}
	w := result[0]
	if w.groupID != 1 {
		t.Errorf("expected groupID=1, got %d", w.groupID)
	}
	if w.subjectID != disciplineID {
		t.Errorf("expected subjectID=%d, got %d", disciplineID, w.subjectID)
	}
	if w.teacherID != 42 {
		t.Errorf("expected teacherID=42, got %d", w.teacherID)
	}
	if w.totalHours != assignedHours {
		t.Errorf("expected totalHours=%d, got %d", assignedHours, w.totalHours)
	}
}
