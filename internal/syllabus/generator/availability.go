package generator

// availabilityStatus represents whether a participant is available at a given slot.
type availabilityStatus int

const (
	statusFree availabilityStatus = iota
	statusUnavailable
	statusHaveLesson
	statusMethodicalDay
)

// availability tracks the free/busy status of all participants (teachers, groups)
// for each date and lesson number.
type availability struct {
	// dateToAvailability maps: date → participantID → lessonNumber → status
	data map[date]map[uint64]map[int]availabilityStatus
}

// newAvailability creates a new availability tracker.
func newAvailability() *availability {
	return &availability{
		data: make(map[date]map[uint64]map[int]availabilityStatus),
	}
}

// getStatus returns the availability status for a participant at a specific date/lesson.
func (a *availability) getStatus(d date, participantID uint64, lessonNumber int) availabilityStatus {
	if a.data[d] == nil {
		return statusFree
	}
	if a.data[d][participantID] == nil {
		return statusFree
	}
	status, ok := a.data[d][participantID][lessonNumber]
	if !ok {
		return statusFree
	}
	return status
}

// setStatus sets the availability status for a participant at a specific date/lesson.
func (a *availability) setStatus(d date, participantID uint64, lessonNumber int, status availabilityStatus) {
	if a.data[d] == nil {
		a.data[d] = make(map[uint64]map[int]availabilityStatus)
	}
	if a.data[d][participantID] == nil {
		a.data[d][participantID] = make(map[int]availabilityStatus)
	}
	a.data[d][participantID][lessonNumber] = status
}

// markBusy marks a participant as having a lesson at the given slot.
func (a *availability) markBusy(d date, participantID uint64, lessonNumber int) {
	a.setStatus(d, participantID, lessonNumber, statusHaveLesson)
}

// markFree clears the busy status for a participant at the given slot.
func (a *availability) markFree(d date, participantID uint64, lessonNumber int) {
	if a.data[d] != nil && a.data[d][participantID] != nil {
		delete(a.data[d][participantID], lessonNumber)
	}
}

// markUnavailable marks a participant as unavailable at the given slot.
func (a *availability) markUnavailable(d date, participantID uint64, lessonNumber int) {
	a.setStatus(d, participantID, lessonNumber, statusUnavailable)
}

// markMethodicalDay marks an entire day as a methodical day for a participant.
func (a *availability) markMethodicalDay(d date, participantID uint64, lessonNumbers []int) {
	for _, ln := range lessonNumbers {
		a.setStatus(d, participantID, ln, statusMethodicalDay)
	}
}

// isFree returns true if the participant is free at the given date/lesson.
func (a *availability) isFree(d date, participantID uint64, lessonNumber int) bool {
	return a.getStatus(d, participantID, lessonNumber) == statusFree
}

// isAvailable returns true if the participant is free or has no specific constraint.
func (a *availability) isAvailable(d date, participantID uint64, lessonNumber int) bool {
	status := a.getStatus(d, participantID, lessonNumber)
	return status == statusFree
}

// areFree returns true if all the given participants are free at the given date/lesson.
func (a *availability) areFree(d date, participantIDs []uint64, lessonNumber int) bool {
	for _, id := range participantIDs {
		if !a.isFree(d, id, lessonNumber) {
			return false
		}
	}
	return true
}

// getBusyLessonNumbers returns which lesson numbers a participant is busy on a given date.
func (a *availability) getBusyLessonNumbers(d date, participantID uint64) []int {
	if a.data[d] == nil || a.data[d][participantID] == nil {
		return nil
	}
	busy := make([]int, 0)
	for ln, status := range a.data[d][participantID] {
		if status != statusFree {
			busy = append(busy, ln)
		}
	}
	return busy
}

// wouldExceedConsecutive returns true if placing a lesson at lessonNumber for
// teacherID on date d would create a consecutive run longer than maxConsecutive.
// A maxConsecutive value of 0 means no limit (always returns false).
func (a *availability) wouldExceedConsecutive(d date, teacherID uint64, lessonNumber, maxConsecutive int) bool {
	if maxConsecutive <= 0 {
		return false
	}
	busy := make(map[int]bool)
	if a.data[d] != nil && a.data[d][teacherID] != nil {
		for ln, status := range a.data[d][teacherID] {
			if status == statusHaveLesson {
				busy[ln] = true
			}
		}
	}
	busy[lessonNumber] = true

	run := 1
	for ln := lessonNumber - 1; busy[ln]; ln-- {
		run++
	}
	for ln := lessonNumber + 1; busy[ln]; ln++ {
		run++
	}
	return run > maxConsecutive
}
