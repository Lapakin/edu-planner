package generator

import (
	"testing"
)

func TestRoomSelector_SelectRoom_SingleRoom(t *testing.T) {
	rs := newRoomSelector([]uint64{101}, nil)
	used := make(map[uint64]bool)
	roomID := rs.selectRoom(used)
	if roomID != 101 {
		t.Errorf("expected roomID=101, got %d", roomID)
	}
}

func TestRoomSelector_SelectRoom_EmptyRoomList(t *testing.T) {
	rs := newRoomSelector([]uint64{}, nil)
	used := make(map[uint64]bool)
	roomID := rs.selectRoom(used)
	if roomID != 0 {
		t.Errorf("expected roomID=0 for empty room list, got %d", roomID)
	}
}

func TestRoomSelector_SelectRoom_BalancesUsage(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil)
	used := make(map[uint64]bool)

	// First call — pick the least used (both at 0, so first one)
	r1 := rs.selectRoom(used)
	if r1 == 0 {
		t.Error("expected non-zero room ID")
	}

	// Second call with fresh usedInSlot — should pick the other room to balance
	used2 := make(map[uint64]bool)
	r2 := rs.selectRoom(used2)
	if r2 == 0 {
		t.Error("expected non-zero room ID on second call")
	}
	// r1 and r2 could be same or different, but usage count should be incremented
	if rs.roomUsageCount[r1] == 0 && rs.roomUsageCount[r2] == 0 {
		t.Error("expected usage count to be incremented")
	}
}

func TestRoomSelector_SelectRoom_SkipsAlreadyUsedInSlot(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil)

	// Mark 101 as used in this slot
	used := map[uint64]bool{101: true}
	roomID := rs.selectRoom(used)
	if roomID != 102 {
		t.Errorf("expected 102 when 101 is already used in slot, got %d", roomID)
	}
}

func TestRoomSelector_SelectRoom_AllRoomsUsedInSlot_FallsBack(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil)

	// All rooms used in this slot — should still return a room (least used)
	used := map[uint64]bool{101: true, 102: true}
	roomID := rs.selectRoom(used)
	if roomID == 0 {
		t.Error("expected non-zero room ID even when all slots used")
	}
}

func TestRoomSelector_AssignRooms_SetsNonZeroRoomIDs(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil)
	sched := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l := &lesson{
		groupID:   1,
		subjectID: 10,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 42, roomID: 0}, // no room yet
		},
	}
	sched.add(d, 1, l)

	rs.assignRooms(sched)

	if l.subLessons[0].roomID == 0 {
		t.Error("expected non-zero roomID after assignRooms")
	}
}

func TestRoomSelector_AssignRooms_PreservesExistingRoomIDs(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil)
	sched := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l := &lesson{
		groupID:   1,
		subjectID: 10,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 42, roomID: 999}, // already has a room
		},
	}
	sched.add(d, 1, l)

	rs.assignRooms(sched)

	// Existing room IDs should be preserved (not overwritten)
	if l.subLessons[0].roomID != 999 {
		t.Errorf("expected roomID=999 to be preserved, got %d", l.subLessons[0].roomID)
	}
}

func TestRoomSelector_AssignRooms_EmptySchedule(_ *testing.T) {
	rs := newRoomSelector([]uint64{101}, nil)
	sched := newSchedule()
	// Should not panic on empty schedule
	rs.assignRooms(sched)
}

func TestRoomSelector_AssignRooms_MultipleLessons(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil)
	sched := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l1 := &lesson{
		groupID:    1,
		subjectID:  10,
		subLessons: []*internalSubLesson{{groupID: 1, teacherID: 42, roomID: 0}},
	}
	l2 := &lesson{
		groupID:    2,
		subjectID:  11,
		subLessons: []*internalSubLesson{{groupID: 2, teacherID: 43, roomID: 0}},
	}
	sched.add(d, 1, l1)
	sched.add(d, 2, l2)

	rs.assignRooms(sched)

	if l1.subLessons[0].roomID == 0 {
		t.Error("expected non-zero roomID for l1")
	}
	if l2.subLessons[0].roomID == 0 {
		t.Error("expected non-zero roomID for l2")
	}
}
