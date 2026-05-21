package generator

import "sort"

// roomSelector assigns rooms to scheduled lessons.
type roomSelector struct {
	availableRoomIDs       []uint64
	roomUsageCount         map[uint64]int      // tracks how often each room is used for balancing
	cycleCommitteeLabRooms map[uint64][]uint64 // cycleCommitteeID → preferred room IDs
}

// newRoomSelector creates a new room selector.
func newRoomSelector(roomIDs []uint64, cycleCommitteeLabRooms map[uint64][]uint64) *roomSelector {
	if cycleCommitteeLabRooms == nil {
		cycleCommitteeLabRooms = make(map[uint64][]uint64)
	}
	return &roomSelector{
		availableRoomIDs:       roomIDs,
		roomUsageCount:         make(map[uint64]int),
		cycleCommitteeLabRooms: cycleCommitteeLabRooms,
	}
}

// assignRooms assigns rooms to all lessons in the schedule.
func (rs *roomSelector) assignRooms(sched *schedule) {
	dates := sched.allDates()
	for _, d := range dates {
		daySchedule := sched.getDaySchedule(d)
		if daySchedule == nil {
			continue
		}

		// Collect lesson numbers sorted
		lessonNumbers := make([]int, 0, len(daySchedule))
		for ln := range daySchedule {
			lessonNumbers = append(lessonNumbers, ln)
		}
		sort.Ints(lessonNumbers)

		// Track rooms used this time slot
		for _, ln := range lessonNumbers {
			usedRooms := make(map[uint64]bool)

			// First pass: record already-assigned rooms
			for _, l := range daySchedule[ln] {
				for _, sl := range l.subLessons {
					if sl.roomID != 0 {
						usedRooms[sl.roomID] = true
					}
				}
			}

			// Second pass: assign rooms to sub-lessons without rooms
			for _, l := range daySchedule[ln] {
				for _, sl := range l.subLessons {
					if sl.roomID == 0 {
						roomID := rs.selectRoomForLesson(usedRooms, l.cycleCommitteeID, l.isLab)
						sl.roomID = roomID
						usedRooms[roomID] = true
					}
				}
			}
		}
	}
}

// selectRoomForLesson selects the best room for a lesson, preferring cycle-committee-specific rooms for labs.
func (rs *roomSelector) selectRoomForLesson(usedInSlot map[uint64]bool, cycleCommitteeID uint64, isLab bool) uint64 {
	// Only apply lab room preference for actual lab lessons
	if isLab && cycleCommitteeID != 0 {
		if labRoomIDs, ok := rs.cycleCommitteeLabRooms[cycleCommitteeID]; ok && len(labRoomIDs) > 0 {
			room := rs.selectFromList(usedInSlot, labRoomIDs)
			if room != 0 {
				return room
			}
		}
	}
	return rs.selectRoom(usedInSlot)
}

// selectRoom selects the best room from all available rooms considering usage balance.
func (rs *roomSelector) selectRoom(usedInSlot map[uint64]bool) uint64 {
	return rs.selectFromList(usedInSlot, rs.availableRoomIDs)
}

// selectFromList selects the best room from a given list considering usage balance.
func (rs *roomSelector) selectFromList(usedInSlot map[uint64]bool, roomIDs []uint64) uint64 {
	if len(roomIDs) == 0 {
		return 0
	}

	// Find room with lowest usage that's not already used in this slot
	bestRoom := uint64(0)
	bestCount := int(^uint(0) >> 1) // max int

	for _, roomID := range roomIDs {
		if usedInSlot[roomID] {
			continue
		}
		count := rs.roomUsageCount[roomID]
		if count < bestCount {
			bestCount = count
			bestRoom = roomID
		}
	}

	// If all rooms are used in this slot, just use the least used
	if bestRoom == 0 {
		bestCount = int(^uint(0) >> 1)
		for _, roomID := range roomIDs {
			count := rs.roomUsageCount[roomID]
			if count < bestCount {
				bestCount = count
				bestRoom = roomID
			}
		}
	}

	if bestRoom != 0 {
		rs.roomUsageCount[bestRoom]++
	}
	return bestRoom
}
