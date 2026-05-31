package generator

import "sort"

// roomSelector assigns rooms to scheduled lessons.
type roomSelector struct {
	availableRoomIDs       []uint64
	roomUsageCount         map[uint64]int      // tracks how often each room is used for balancing
	cycleCommitteeLabRooms map[uint64][]uint64 // cycleCommitteeID → preferred room IDs
	roomCapacities         map[uint64]int      // roomID → seats (0 or absent = unlimited/unset)
	groupSizes             map[uint64]int      // groupID → students (0 or absent = unknown)
}

// newRoomSelector creates a new room selector.
func newRoomSelector(
	roomIDs []uint64,
	cycleCommitteeLabRooms map[uint64][]uint64,
	roomCapacities map[uint64]int,
	groupSizes map[uint64]int,
) *roomSelector {
	if cycleCommitteeLabRooms == nil {
		cycleCommitteeLabRooms = make(map[uint64][]uint64)
	}
	if roomCapacities == nil {
		roomCapacities = make(map[uint64]int)
	}
	if groupSizes == nil {
		groupSizes = make(map[uint64]int)
	}
	return &roomSelector{
		availableRoomIDs:       roomIDs,
		roomUsageCount:         make(map[uint64]int),
		cycleCommitteeLabRooms: cycleCommitteeLabRooms,
		roomCapacities:         roomCapacities,
		groupSizes:             groupSizes,
	}
}

// requiredCapacity returns how many seats a sub-lesson needs. A split sub-lesson
// holds a single subgroup (≈ half the group); a united sub-lesson holds the whole
// group. A return of 0 means the size is unknown and imposes no constraint.
func (rs *roomSelector) requiredCapacity(l *lesson, sl *internalSubLesson) int {
	size := rs.groupSizes[sl.groupID]
	if size <= 0 {
		return 0
	}
	if l.format == formatSplit && sl.subGroupNumber != nil {
		return (size + 1) / 2 // ceil(size/2)
	}
	return size
}

// fits reports whether a room can hold the required number of seats. A room with
// capacity 0 is treated as unlimited (capacity unset).
func (rs *roomSelector) fits(roomID uint64, required int) bool {
	if required <= 0 {
		return true
	}
	capacity := rs.roomCapacities[roomID]
	return capacity == 0 || capacity >= required
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
						required := rs.requiredCapacity(l, sl)
						roomID := rs.selectRoomForLesson(usedRooms, l.cycleCommitteeID, l.isLab, required)
						sl.roomID = roomID
						usedRooms[roomID] = true
					}
				}
			}
		}
	}
}

// selectRoomForLesson selects the best room for a lesson, preferring cycle-committee-specific rooms for labs.
// required is the number of seats the lesson needs (0 = no constraint).
func (rs *roomSelector) selectRoomForLesson(usedInSlot map[uint64]bool, cycleCommitteeID uint64, isLab bool, required int) uint64 {
	// Only apply lab room preference for actual lab lessons
	if isLab && cycleCommitteeID != 0 {
		if labRoomIDs, ok := rs.cycleCommitteeLabRooms[cycleCommitteeID]; ok && len(labRoomIDs) > 0 {
			room := rs.selectFromList(usedInSlot, labRoomIDs, required)
			if room != 0 {
				return room
			}
		}
	}
	return rs.selectRoom(usedInSlot, required)
}

// selectRoom selects the best room from all available rooms considering usage balance.
func (rs *roomSelector) selectRoom(usedInSlot map[uint64]bool, required int) uint64 {
	return rs.selectFromList(usedInSlot, rs.availableRoomIDs, required)
}

// selectFromList selects the best room from a given list considering usage balance,
// restricting the choice to rooms that can hold the required number of seats.
func (rs *roomSelector) selectFromList(usedInSlot map[uint64]bool, roomIDs []uint64, required int) uint64 {
	if len(roomIDs) == 0 {
		return 0
	}

	// Find room with lowest usage that's not already used in this slot and fits.
	bestRoom := uint64(0)
	bestCount := int(^uint(0) >> 1) // max int

	for _, roomID := range roomIDs {
		if usedInSlot[roomID] || !rs.fits(roomID, required) {
			continue
		}
		count := rs.roomUsageCount[roomID]
		if count < bestCount {
			bestCount = count
			bestRoom = roomID
		}
	}

	// If all fitting rooms are used in this slot, just use the least used fitting one.
	if bestRoom == 0 {
		bestCount = int(^uint(0) >> 1)
		for _, roomID := range roomIDs {
			if !rs.fits(roomID, required) {
				continue
			}
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
