package generator

import "sort"

// roomSelector assigns rooms to scheduled lessons.
type roomSelector struct {
	availableRoomIDs []uint64
	roomUsageCount   map[uint64]int // tracks how often each room is used for balancing
}

// newRoomSelector creates a new room selector.
func newRoomSelector(roomIDs []uint64) *roomSelector {
	return &roomSelector{
		availableRoomIDs: roomIDs,
		roomUsageCount:   make(map[uint64]int),
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
						roomID := rs.selectRoom(usedRooms)
						sl.roomID = roomID
						usedRooms[roomID] = true
					}
				}
			}
		}
	}
}

// selectRoom selects the best room considering usage balance.
func (rs *roomSelector) selectRoom(usedInSlot map[uint64]bool) uint64 {
	if len(rs.availableRoomIDs) == 0 {
		return 0
	}

	// Find room with lowest usage that's not already used in this slot
	bestRoom := uint64(0)
	bestCount := int(^uint(0) >> 1) // max int

	for _, roomID := range rs.availableRoomIDs {
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
		for _, roomID := range rs.availableRoomIDs {
			count := rs.roomUsageCount[roomID]
			if count < bestCount {
				bestCount = count
				bestRoom = roomID
			}
		}
	}

	rs.roomUsageCount[bestRoom]++
	return bestRoom
}
