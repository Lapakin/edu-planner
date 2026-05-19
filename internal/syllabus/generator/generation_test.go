package generator

import (
	"testing"
	"time"

	"github.com/Lapakin/edu-planner/internal/domain"
)

// minimalScenarioSettings builds the settings for the minimal test scenario.
func minimalScenarioSettings() *settings {
	ts := makeTemplateSetting(2.0, 8, 40)
	bells := makeBellSchedules(4)
	r := &domain.ScheduleRestriction{
		MinGroupLessonsPerDay:   0,
		MaxGroupLessonsPerDay:   4,
		MaxTeacherLessonsPerDay: 5,
		NoGapsInGroupSchedule:   false,
	}
	return newSettings(ts, bells, r)
}

// minimalScenarioWorkload builds the workload for the minimal scenario.
// totalHours=4, hoursPerLesson=2 → 2 total lessons (1 num + 1 denom).
func minimalScenarioWorkload() *Workload {
	return &Workload{
		workloadDistributionID: 1,
		groupID:                1,
		subjectID:              10,
		teacherID:              42,
		totalHours:             4,
		isSplitting:            false,
	}
}

func TestGeneration_Exec_MinimalScenario(t *testing.T) {
	cfg := minimalScenarioSettings()
	w := minimalScenarioWorkload()

	// 2024-09-02 (Monday) to 2024-09-13 (Friday) — 2 weeks Mon-Fri
	startDate := newDate(time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC))
	endDate := newDate(time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC))

	gen := newGeneration(
		cfg,
		[]*Workload{w},
		[]uint64{101},
		[]uint64{1},
		startDate,
		endDate,
		42,
	)

	data, err := gen.exec()
	if err != nil {
		t.Fatalf("exec() returned error: %v", err)
	}
	if data == nil {
		t.Fatal("exec() returned nil data")
	}

	// Should have numerator schedule with at least one lesson
	if data.Numerator == nil {
		t.Error("expected non-nil numerator")
	}

	// Numerator should contain at least one day
	totalLessons := 0
	for _, daySchedule := range data.Numerator {
		if daySchedule != nil {
			for _, lessons := range *daySchedule {
				totalLessons += len(lessons)
			}
		}
	}
	if totalLessons == 0 {
		t.Error("expected at least one lesson in numerator")
	}
}

func TestGeneration_Exec_ReturnsScheduleData(t *testing.T) {
	cfg := minimalScenarioSettings()
	w := minimalScenarioWorkload()

	startDate := newDate(time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC))
	endDate := newDate(time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC))

	gen := newGeneration(cfg, []*Workload{w}, []uint64{101}, []uint64{1}, startDate, endDate, 99)
	data, err := gen.exec()
	if err != nil {
		t.Fatalf("exec() returned error: %v", err)
	}

	// Denominator should be present for 2-week semester
	if data.Denominator == nil {
		t.Error("expected non-nil denominator for 2-week semester")
	}
}

func TestGeneration_Exec_LessonsHaveRoomAssigned(t *testing.T) {
	cfg := minimalScenarioSettings()
	w := minimalScenarioWorkload()

	startDate := newDate(time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC))
	endDate := newDate(time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC))

	gen := newGeneration(cfg, []*Workload{w}, []uint64{101}, []uint64{1}, startDate, endDate, 7)
	data, err := gen.exec()
	if err != nil {
		t.Fatalf("exec() returned error: %v", err)
	}

	// All sub-lessons should have a non-zero room ID
	for _, daySchedule := range data.Numerator {
		if daySchedule == nil {
			continue
		}
		for _, lessons := range *daySchedule {
			for _, l := range lessons {
				for _, sl := range l.SubLessons {
					if sl.RoomID == 0 {
						t.Error("expected non-zero room ID after assignRooms")
					}
				}
			}
		}
	}
}

func TestGeneration_Exec_SubjectIDPreserved(t *testing.T) {
	cfg := minimalScenarioSettings()
	w := minimalScenarioWorkload()

	startDate := newDate(time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC))
	endDate := newDate(time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC))

	gen := newGeneration(cfg, []*Workload{w}, []uint64{101}, []uint64{1}, startDate, endDate, 13)
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
				if l.SubjectID != 10 {
					t.Errorf("expected subjectID=10, got %d", l.SubjectID)
				}
			}
		}
	}
}
