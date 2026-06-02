package generator

import (
	"testing"
)

func makeFinderSettings(numBells int) *settings {
	ts := makeTemplateSetting(2, numBells*2, 40)
	bells := makeBellSchedules(numBells)
	return newSettings(ts, bells, nil, nil, nil)
}

func TestFinder_FindUnscheduledNumbers_EmptyBusy(t *testing.T) {
	cfg := makeFinderSettings(4)
	f := newFinder(cfg)

	// When no busy slots, should return the first lesson number
	result := f.findUnscheduledNumbers([]int{})
	if len(result) != 1 {
		t.Fatalf("expected 1 result for empty busy slots, got %v", result)
	}
	if result[0] != cfg.firstLessonNumber {
		t.Errorf("expected first lesson number %d, got %d", cfg.firstLessonNumber, result[0])
	}
}

func TestFinder_FindUnscheduledNumbers_AdjacentToExistingSlot(t *testing.T) {
	cfg := makeFinderSettings(4) // lessons 1..4
	f := newFinder(cfg)

	// Busy at slot 2 → adjacent free slots are 1 and 3
	result := f.findUnscheduledNumbers([]int{2})
	if len(result) == 0 {
		t.Fatal("expected adjacent free slots, got none")
	}
	found1 := false
	found3 := false
	for _, r := range result {
		if r == 1 {
			found1 = true
		}
		if r == 3 {
			found3 = true
		}
	}
	if !found1 || !found3 {
		t.Errorf("expected slots 1 and 3 adjacent to slot 2, got %v", result)
	}
}

func TestFinder_FindUnscheduledNumbers_AtBoundary(t *testing.T) {
	cfg := makeFinderSettings(4) // lessons 1..4
	f := newFinder(cfg)

	// Busy at slot 1 → adjacent is only slot 2 (slot 0 doesn't exist)
	result := f.findUnscheduledNumbers([]int{1})
	for _, r := range result {
		if r < cfg.firstLessonNumber || r > cfg.lastLessonNumber {
			t.Errorf("slot %d is out of range [%d, %d]", r, cfg.firstLessonNumber, cfg.lastLessonNumber)
		}
	}
}

func TestFinder_FindUnscheduledNumbers_Deduplication(t *testing.T) {
	cfg := makeFinderSettings(5) // lessons 1..5
	f := newFinder(cfg)

	// Busy at slots 2 and 4 → adjacent are: 1, 3 (from 2) and 3, 5 (from 4)
	// Slot 3 appears twice but should be deduped
	result := f.findUnscheduledNumbers([]int{2, 4})
	seen := make(map[int]bool)
	for _, r := range result {
		if seen[r] {
			t.Errorf("duplicate slot %d in result %v", r, result)
		}
		seen[r] = true
	}
}

func TestFinder_FindUnscheduledNumbers_CachesResult(t *testing.T) {
	cfg := makeFinderSettings(4)
	f := newFinder(cfg)

	result1 := f.findUnscheduledNumbers([]int{2})
	result2 := f.findUnscheduledNumbers([]int{2})
	if len(result1) != len(result2) {
		t.Errorf("cached result differs: %v vs %v", result1, result2)
	}
}

func TestFinder_FindUnscheduledNumbers_AllBusy_FallbackToAllFree(t *testing.T) {
	cfg := makeFinderSettings(3) // lessons 1..3
	f := newFinder(cfg)

	// Busy at 1, 2, 3 — all busy, no adjacent free slots
	result := f.findUnscheduledNumbers([]int{1, 2, 3})
	// All busy → deduped adjacent slots is empty → fallback returns remaining free slots
	// All are busy, so result should be empty
	if len(result) != 0 {
		t.Errorf("expected empty when all slots busy, got %v", result)
	}
}

func TestFinder_FindAllFreeSlots(t *testing.T) {
	cfg := makeFinderSettings(4) // lessons 1..4
	f := newFinder(cfg)
	sched := newSchedule()
	d := date{year: 2024, month: 9, day: 2}
	groupID := uint64(1)

	// Add a lesson at slot 2
	l := &lesson{groupID: groupID, subjectID: 10}
	sched.add(d, 2, l)

	free := f.findAllFreeSlots(sched, d, []uint64{groupID})
	// Slots 1, 3, 4 should be free
	if len(free) != 3 {
		t.Errorf("expected 3 free slots, got %v", free)
	}
	for _, s := range free {
		if s == 2 {
			t.Error("slot 2 should not be in free list")
		}
	}
}

func TestFinder_FindAllFreeSlots_AllFreeWhenEmpty(t *testing.T) {
	cfg := makeFinderSettings(4)
	f := newFinder(cfg)
	sched := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	free := f.findAllFreeSlots(sched, d, []uint64{1})
	if len(free) != 4 {
		t.Errorf("expected 4 free slots for empty schedule, got %v", free)
	}
}
