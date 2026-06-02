package generator

import (
	"strings"
	"testing"
	"time"

	"github.com/Lapakin/edu-planner/internal/domain"
)

// makeFlowCandidate builds a flow-eligible (united, non-lab) lecturer lesson for a
// single group, the shape mergeFlowLessons expects to merge.
func makeFlowCandidate(groupID, subjectID, teacherID uint64) *lesson {
	return &lesson{
		groupID:      groupID,
		subjectID:    subjectID,
		format:       formatUnited,
		flowEligible: true,
		subLessons: []*internalSubLesson{
			{groupID: groupID, teacherID: teacherID},
		},
	}
}

func flowEnabledSettings() *settings {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(4)
	return newSettings(ts, bells, nil, nil, nil) // default restriction → allowFlowLessons=true
}

// --- mergeFlowLessons ---

func TestMergeFlowLessons_MergesAcrossGroups(t *testing.T) {
	cfg := flowEnabledSettings()
	lessons := []*lesson{
		makeFlowCandidate(1, 10, 42),
		makeFlowCandidate(2, 10, 42),
	}

	result := mergeFlowLessons(lessons, cfg)

	if len(result) != 1 {
		t.Fatalf("expected 1 merged flow lesson, got %d", len(result))
	}
	fl := result[0]
	if !fl.isFlow {
		t.Error("merged lesson should be marked as flow")
	}
	if len(fl.subLessons) != 2 {
		t.Errorf("expected 2 sub-lessons, got %d", len(fl.subLessons))
	}
	gids := fl.groupIDs()
	if len(gids) != 2 || !fl.hasGroup(1) || !fl.hasGroup(2) {
		t.Errorf("expected flow to cover groups 1 and 2, got %v", gids)
	}
	if fl.flowID != "flow-10-42" {
		t.Errorf("unexpected flowID %q", fl.flowID)
	}
	if len(fl.flowOrigin) != 2 {
		t.Errorf("expected 2 constituents retained for fallback, got %d", len(fl.flowOrigin))
	}
}

func TestMergeFlowLessons_DisabledReturnsUnchanged(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(4)
	r := &domain.ScheduleRestriction{AllowFlowLessons: false}
	cfg := newSettings(ts, bells, r, nil, nil)

	lessons := []*lesson{
		makeFlowCandidate(1, 10, 42),
		makeFlowCandidate(2, 10, 42),
	}
	result := mergeFlowLessons(lessons, cfg)

	if len(result) != 2 {
		t.Fatalf("expected lessons unchanged when flow disabled, got %d", len(result))
	}
	for _, l := range result {
		if len(l.subLessons) != 1 {
			t.Error("lessons should not be merged when flow disabled")
		}
	}
}

func TestMergeFlowLessons_SingleGroupStaysRegular(t *testing.T) {
	cfg := flowEnabledSettings()
	lessons := []*lesson{makeFlowCandidate(1, 10, 42)}

	result := mergeFlowLessons(lessons, cfg)

	if len(result) != 1 {
		t.Fatalf("expected 1 lesson, got %d", len(result))
	}
	if result[0].isFlow {
		t.Error("a lonely flow candidate must fall back to a regular lesson")
	}
}

func TestMergeFlowLessons_LabNotMerged(t *testing.T) {
	cfg := flowEnabledSettings()
	l1 := makeFlowCandidate(1, 10, 42)
	l2 := makeFlowCandidate(2, 10, 42)
	l1.isLab = true
	l2.isLab = true

	result := mergeFlowLessons([]*lesson{l1, l2}, cfg)

	if len(result) != 2 {
		t.Fatalf("expected lab lessons untouched, got %d results", len(result))
	}
}

func TestMergeFlowLessons_DifferentTeacherNotMerged(t *testing.T) {
	cfg := flowEnabledSettings()
	result := mergeFlowLessons([]*lesson{
		makeFlowCandidate(1, 10, 42),
		makeFlowCandidate(2, 10, 99),
	}, cfg)

	for _, l := range result {
		if l.isFlow {
			t.Error("lessons taught by different teachers must not be merged into a flow")
		}
	}
}

func TestMergeFlowLessons_UnequalCountsLeavesSurplusSeparate(t *testing.T) {
	cfg := flowEnabledSettings()
	lessons := []*lesson{
		makeFlowCandidate(1, 10, 42),
		makeFlowCandidate(1, 10, 42), // group 1 has two
		makeFlowCandidate(2, 10, 42), // group 2 has one
	}

	result := mergeFlowLessons(lessons, cfg)

	flows, regular := partitionFlowLessons(result)
	if len(flows) != 1 {
		t.Fatalf("expected exactly 1 flow lesson, got %d", len(flows))
	}
	if len(regular) != 1 {
		t.Fatalf("expected 1 surplus regular lesson, got %d", len(regular))
	}
	if regular[0].groupID != 1 {
		t.Errorf("surplus lesson should belong to group 1, got %d", regular[0].groupID)
	}
}

// --- lesson helpers ---

func TestLesson_GroupIDsAndHasGroup(t *testing.T) {
	regular := &lesson{groupID: 5, subLessons: []*internalSubLesson{{groupID: 5, teacherID: 1}}}
	if got := regular.groupIDs(); len(got) != 1 || got[0] != 5 {
		t.Errorf("regular lesson groupIDs = %v, want [5]", got)
	}
	if !regular.hasGroup(5) || regular.hasGroup(6) {
		t.Error("hasGroup mismatch for regular lesson")
	}

	flow := &lesson{
		groupID: 3,
		isFlow:  true,
		subLessons: []*internalSubLesson{
			{groupID: 3, teacherID: 1},
			{groupID: 7, teacherID: 1},
		},
	}
	if !flow.hasGroup(3) || !flow.hasGroup(7) {
		t.Error("flow lesson should report both groups")
	}
	if !flow.hasAnyGroup([]uint64{9, 7}) {
		t.Error("hasAnyGroup should detect overlap")
	}
	if flow.hasAnyGroup([]uint64{9, 8}) {
		t.Error("hasAnyGroup should not report a false overlap")
	}
}

func TestLesson_Unmerge(t *testing.T) {
	cfg := flowEnabledSettings()
	merged := mergeFlowLessons([]*lesson{
		makeFlowCandidate(1, 10, 42),
		makeFlowCandidate(2, 10, 42),
	}, cfg)
	fl := merged[0]
	fl.isScheduled = true

	parts := fl.unmerge()
	if len(parts) != 2 {
		t.Fatalf("expected 2 constituents from unmerge, got %d", len(parts))
	}
	for _, p := range parts {
		if p.isFlow {
			t.Error("unmerged constituent must no longer be a flow lesson")
		}
		if p.isScheduled {
			t.Error("unmerged constituent must be reset to unscheduled")
		}
	}
}

// --- flowRoomFeasible ---

func TestFlowRoomFeasible(t *testing.T) {
	base := flowEnabledSettings()
	flow := &lesson{
		isFlow:     true,
		subLessons: []*internalSubLesson{{groupID: 1}, {groupID: 2}},
	}

	// No capacity data → permissive.
	if !flowRoomFeasible(base, flow) {
		t.Error("expected feasible when no room capacities are known")
	}

	cfg := flowEnabledSettings()
	cfg.groupSizes = map[uint64]int{1: 30, 2: 30}
	cfg.roomCapacities = map[uint64]int{101: 40}
	if flowRoomFeasible(cfg, flow) {
		t.Error("combined size 60 must not fit a 40-seat room")
	}

	cfg.roomCapacities = map[uint64]int{101: 40, 102: 80}
	if !flowRoomFeasible(cfg, flow) {
		t.Error("an 80-seat room should hold the combined size of 60")
	}

	cfg.roomCapacities = map[uint64]int{101: 0} // unlimited
	if !flowRoomFeasible(cfg, flow) {
		t.Error("an unlimited room should always be feasible")
	}
}

// --- flowDiagnostics ---

func TestFlowDiagnostics_Metadata(t *testing.T) {
	if md := newFlowDiagnostics().metadata(); md != "" {
		t.Errorf("expected empty metadata when no flows formed, got %q", md)
	}

	d := newFlowDiagnostics()
	d.countFlows([]*lesson{{isFlow: true}, {isFlow: false}}, []*lesson{{isFlow: true}})
	d.recordFallback(&lesson{flowID: "flow-10-42"})
	d.recordFallback(&lesson{flowID: "flow-10-42"}) // duplicate ignored

	md := d.metadata()
	if !strings.Contains(md, `"flow_lessons":2`) {
		t.Errorf("expected 2 flow lessons in metadata, got %q", md)
	}
	if !strings.Contains(md, `"flow_fallbacks":1`) {
		t.Errorf("expected 1 fallback in metadata, got %q", md)
	}
	if !strings.Contains(md, "flow-10-42") {
		t.Errorf("expected fallback flow id in metadata, got %q", md)
	}
}

// --- integration via generation.exec ---

// flowGenerationScenario wires two groups sharing one flow discipline (subject 10,
// teacher 42), each carrying 4 hours → one numerator + one denominator lesson per
// group that should merge into a shared flow lecture.
func flowGenerationScenario(roomCapacity, groupSize int) (*settings, []*Workload) {
	cfg := minimalScenarioSettings()
	cfg.allowFlowLessons = true
	cfg.roomCapacities = map[uint64]int{101: roomCapacity}
	cfg.groupSizes = map[uint64]int{1: groupSize, 2: groupSize}

	workloads := []*Workload{
		{workloadDistributionID: 1, groupID: 1, subjectID: 10, teacherID: 42, totalHours: 4, isFlow: true},
		{workloadDistributionID: 2, groupID: 2, subjectID: 10, teacherID: 42, totalHours: 4, isFlow: true},
	}
	return cfg, workloads
}

func TestGeneration_Exec_FlowMergedAcrossGroups(t *testing.T) {
	cfg, workloads := flowGenerationScenario(100, 30) // combined 60 fits the 100-seat room

	start := newDate(time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC))
	end := newDate(time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC))

	gen := newGeneration(cfg, workloads, []uint64{101}, []uint64{1, 2}, start, end, 42)
	data, err := gen.exec()
	if err != nil {
		t.Fatalf("exec() returned error: %v", err)
	}

	flowFound := false
	for _, daySchedule := range data.Numerator {
		if daySchedule == nil {
			continue
		}
		for _, lessons := range *daySchedule {
			for _, l := range lessons {
				if l.Type != domain.LessonTypeFlow {
					continue
				}
				flowFound = true
				if !l.IsFlow || l.FlowID == "" {
					t.Error("flow lesson should carry IsFlow and a FlowID")
				}
				if len(l.SubLessons) != 2 {
					t.Fatalf("flow lesson should have 2 sub-lessons, got %d", len(l.SubLessons))
				}
				room := l.SubLessons[0].RoomID
				if room == 0 {
					t.Error("flow lesson should have a room assigned")
				}
				for _, sl := range l.SubLessons {
					if sl.RoomID != room {
						t.Error("all groups of a flow lesson must share one room")
					}
				}
			}
		}
	}
	if !flowFound {
		t.Error("expected a flow lesson in the numerator")
	}
	if !strings.Contains(data.Metadata, `"flow_lessons"`) {
		t.Errorf("expected flow metadata, got %q", data.Metadata)
	}
}

func TestGeneration_Exec_FlowFallbackWhenRoomTooSmall(t *testing.T) {
	cfg, workloads := flowGenerationScenario(40, 30) // combined 60 cannot fit any 40-seat room

	start := newDate(time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC))
	end := newDate(time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC))

	gen := newGeneration(cfg, workloads, []uint64{101}, []uint64{1, 2}, start, end, 7)
	data, err := gen.exec()
	if err != nil {
		t.Fatalf("exec() returned error: %v", err)
	}

	for _, daySchedule := range data.Numerator {
		if daySchedule == nil {
			continue
		}
		for _, lessons := range *daySchedule {
			for _, l := range lessons {
				if l.Type == domain.LessonTypeFlow {
					t.Error("no flow lesson should survive when no room fits the combined size")
				}
			}
		}
	}
	if !strings.Contains(data.Metadata, `"flow_fallbacks"`) {
		t.Errorf("expected fallback to be reported in metadata, got %q", data.Metadata)
	}
}
