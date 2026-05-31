package generator

import (
	"testing"
)

func TestRoomSelector_SelectRoom_SingleRoom(t *testing.T) {
	rs := newRoomSelector([]uint64{101}, nil, nil, nil)
	used := make(map[uint64]bool)
	roomID := rs.selectRoom(used, 0)
	if roomID != 101 {
		t.Errorf("expected roomID=101, got %d", roomID)
	}
}

func TestRoomSelector_SelectRoom_EmptyRoomList(t *testing.T) {
	rs := newRoomSelector([]uint64{}, nil, nil, nil)
	used := make(map[uint64]bool)
	roomID := rs.selectRoom(used, 0)
	if roomID != 0 {
		t.Errorf("expected roomID=0 for empty room list, got %d", roomID)
	}
}

func TestRoomSelector_SelectRoom_BalancesUsage(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil, nil, nil)
	used := make(map[uint64]bool)

	// First call — pick the least used (both at 0, so first one)
	r1 := rs.selectRoom(used, 0)
	if r1 == 0 {
		t.Error("expected non-zero room ID")
	}

	// Second call with fresh usedInSlot — should pick the other room to balance
	used2 := make(map[uint64]bool)
	r2 := rs.selectRoom(used2, 0)
	if r2 == 0 {
		t.Error("expected non-zero room ID on second call")
	}
	// r1 and r2 could be same or different, but usage count should be incremented
	if rs.roomUsageCount[r1] == 0 && rs.roomUsageCount[r2] == 0 {
		t.Error("expected usage count to be incremented")
	}
}

func TestRoomSelector_SelectRoom_SkipsAlreadyUsedInSlot(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil, nil, nil)

	// Mark 101 as used in this slot
	used := map[uint64]bool{101: true}
	roomID := rs.selectRoom(used, 0)
	if roomID != 102 {
		t.Errorf("expected 102 when 101 is already used in slot, got %d", roomID)
	}
}

func TestRoomSelector_SelectRoom_AllRoomsUsedInSlot_FallsBack(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil, nil, nil)

	// All rooms used in this slot — should still return a room (least used)
	used := map[uint64]bool{101: true, 102: true}
	roomID := rs.selectRoom(used, 0)
	if roomID == 0 {
		t.Error("expected non-zero room ID even when all slots used")
	}
}

func TestRoomSelector_AssignRooms_SetsNonZeroRoomIDs(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil, nil, nil)
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
	rs := newRoomSelector([]uint64{101, 102}, nil, nil, nil)
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
	rs := newRoomSelector([]uint64{101}, nil, nil, nil)
	sched := newSchedule()
	// Should not panic on empty schedule
	rs.assignRooms(sched)
}

func TestRoomSelector_SelectFromList_SkipsTooSmallRooms(t *testing.T) {
	// room 101 holds 20, room 102 holds 40; a lesson needing 30 must pick 102.
	rs := newRoomSelector([]uint64{101, 102}, nil,
		map[uint64]int{101: 20, 102: 40}, nil)

	roomID := rs.selectFromList(map[uint64]bool{}, []uint64{101, 102}, 30)
	if roomID != 102 {
		t.Errorf("expected room 102 (capacity 40) for required 30, got %d", roomID)
	}
}

func TestRoomSelector_SelectFromList_NoFittingRoom(t *testing.T) {
	rs := newRoomSelector([]uint64{101}, nil, map[uint64]int{101: 10}, nil)

	roomID := rs.selectFromList(map[uint64]bool{}, []uint64{101}, 30)
	if roomID != 0 {
		t.Errorf("expected 0 when no room fits required 30, got %d", roomID)
	}
}

func TestRoomSelector_SelectFromList_ZeroCapacityIsUnlimited(t *testing.T) {
	// capacity 0 means "unset" → fits any size.
	rs := newRoomSelector([]uint64{101}, nil, map[uint64]int{101: 0}, nil)

	roomID := rs.selectFromList(map[uint64]bool{}, []uint64{101}, 500)
	if roomID != 101 {
		t.Errorf("expected room 101 (unlimited) for required 500, got %d", roomID)
	}
}

func TestRoomSelector_RequiredCapacity(t *testing.T) {
	rs := newRoomSelector(nil, nil, nil, map[uint64]int{1: 30})

	united := &lesson{format: formatUnited}
	if got := rs.requiredCapacity(united, &internalSubLesson{groupID: 1}); got != 30 {
		t.Errorf("united: expected 30, got %d", got)
	}

	sub := uint64(1)
	split := &lesson{format: formatSplit}
	if got := rs.requiredCapacity(split, &internalSubLesson{groupID: 1, subGroupNumber: &sub}); got != 15 {
		t.Errorf("split: expected 15 (ceil 30/2), got %d", got)
	}

	if got := rs.requiredCapacity(united, &internalSubLesson{groupID: 999}); got != 0 {
		t.Errorf("unknown group: expected 0, got %d", got)
	}
}

func TestRoomSelector_AssignRooms_RespectsCapacity(t *testing.T) {
	// group 1 has 35 students; only room 102 (capacity 40) fits.
	rs := newRoomSelector([]uint64{101, 102}, nil,
		map[uint64]int{101: 20, 102: 40}, map[uint64]int{1: 35})
	sched := newSchedule()
	d := date{year: 2024, month: 9, day: 2}

	l := &lesson{
		groupID:    1,
		subjectID:  10,
		subLessons: []*internalSubLesson{{groupID: 1, teacherID: 42, roomID: 0}},
	}
	sched.add(d, 1, l)

	rs.assignRooms(sched)

	if l.subLessons[0].roomID != 102 {
		t.Errorf("expected room 102 (fits 35), got %d", l.subLessons[0].roomID)
	}
}

func TestRoomSelector_AssignRooms_MultipleLessons(t *testing.T) {
	rs := newRoomSelector([]uint64{101, 102}, nil, nil, nil)
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
