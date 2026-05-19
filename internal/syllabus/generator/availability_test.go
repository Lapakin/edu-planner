package generator

import (
	"sort"
	"testing"
)

func TestNewAvailability_StartsEmpty(t *testing.T) {
	a := newAvailability()
	if a == nil {
		t.Fatal("expected non-nil availability")
	}
	if a.data == nil {
		t.Fatal("expected non-nil data map")
	}
}

func TestAvailability_GetStatus_DefaultFree(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	status := a.getStatus(d, 1, 1)
	if status != statusFree {
		t.Errorf("expected statusFree for unknown slot, got %v", status)
	}
}

func TestAvailability_SetAndGetStatus(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.setStatus(d, 1, 1, statusHaveLesson)
	if got := a.getStatus(d, 1, 1); got != statusHaveLesson {
		t.Errorf("expected statusHaveLesson, got %v", got)
	}
}

func TestAvailability_MarkBusy(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.markBusy(d, 42, 3)
	if a.getStatus(d, 42, 3) != statusHaveLesson {
		t.Error("expected statusHaveLesson after markBusy")
	}
}

func TestAvailability_MarkFree(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.markBusy(d, 42, 3)
	a.markFree(d, 42, 3)
	if a.getStatus(d, 42, 3) != statusFree {
		t.Error("expected statusFree after markFree")
	}
}

func TestAvailability_MarkFree_NonExistentIsNoop(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	// Should not panic
	a.markFree(d, 99, 5)
	if a.getStatus(d, 99, 5) != statusFree {
		t.Error("expected statusFree for never-set slot")
	}
}

func TestAvailability_MarkUnavailable(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.markUnavailable(d, 10, 2)
	if a.getStatus(d, 10, 2) != statusUnavailable {
		t.Error("expected statusUnavailable after markUnavailable")
	}
}

func TestAvailability_MarkMethodicalDay(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.markMethodicalDay(d, 7, []int{1, 2, 3})
	for _, ln := range []int{1, 2, 3} {
		if a.getStatus(d, 7, ln) != statusMethodicalDay {
			t.Errorf("expected statusMethodicalDay at lesson %d", ln)
		}
	}
}

func TestAvailability_IsFree_TrueWhenFree(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	if !a.isFree(d, 1, 1) {
		t.Error("expected isFree=true for unset slot")
	}
}

func TestAvailability_IsFree_FalseWhenBusy(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.markBusy(d, 1, 1)
	if a.isFree(d, 1, 1) {
		t.Error("expected isFree=false after markBusy")
	}
}

func TestAvailability_IsAvailable_TrueWhenFree(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	if !a.isAvailable(d, 1, 1) {
		t.Error("expected isAvailable=true for unset slot")
	}
}

func TestAvailability_IsAvailable_FalseWhenBusy(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.markBusy(d, 1, 1)
	if a.isAvailable(d, 1, 1) {
		t.Error("expected isAvailable=false after markBusy")
	}
}

func TestAvailability_AreFree_AllFree(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	if !a.areFree(d, []uint64{1, 2, 3}, 1) {
		t.Error("expected areFree=true for all unset participants")
	}
}

func TestAvailability_AreFree_OneNotFree(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.markBusy(d, 2, 1)
	if a.areFree(d, []uint64{1, 2, 3}, 1) {
		t.Error("expected areFree=false when one participant is busy")
	}
}

func TestAvailability_AreFree_EmptyList(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	if !a.areFree(d, []uint64{}, 1) {
		t.Error("expected areFree=true for empty participant list")
	}
}

func TestAvailability_GetBusyLessonNumbers_EmptyWhenNoneSet(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	busy := a.getBusyLessonNumbers(d, 1)
	if len(busy) != 0 {
		t.Errorf("expected empty busy list, got %v", busy)
	}
}

func TestAvailability_GetBusyLessonNumbers_ReturnsBusySlots(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.markBusy(d, 5, 1)
	a.markBusy(d, 5, 3)
	a.markUnavailable(d, 5, 2)

	busy := a.getBusyLessonNumbers(d, 5)
	sort.Ints(busy)
	if len(busy) != 3 {
		t.Errorf("expected 3 busy slots, got %v", busy)
	}
	if busy[0] != 1 || busy[1] != 2 || busy[2] != 3 {
		t.Errorf("expected [1,2,3], got %v", busy)
	}
}

func TestAvailability_GetBusyLessonNumbers_ExcludesFreeSlots(t *testing.T) {
	a := newAvailability()
	d := date{year: 2024, month: 9, day: 2}
	a.markBusy(d, 5, 1)
	a.markFree(d, 5, 1) // freed it

	busy := a.getBusyLessonNumbers(d, 5)
	if len(busy) != 0 {
		t.Errorf("expected empty after markFree, got %v", busy)
	}
}
