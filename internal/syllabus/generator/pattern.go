package generator

import "fmt"

// pattern is a key that uniquely identifies a course pattern (group+teacher+subject).
type pattern string

// newPattern creates a pattern from group, teacher, and subject IDs.
func newPattern(groupID, teacherID, subjectID uint64) pattern {
	return pattern(fmt.Sprintf("%d:%d:%d", groupID, teacherID, subjectID))
}

// patternController enforces pattern occurrence limits per day
// and minimum day intervals between occurrences.
type patternController struct {
	occurrences map[date]map[pattern]int
	dayLimits   map[pattern]int
	dayBetween  map[pattern]int
}

// newPatternController creates a new pattern controller.
func newPatternController() *patternController {
	return &patternController{
		occurrences: make(map[date]map[pattern]int),
		dayLimits:   make(map[pattern]int),
		dayBetween:  make(map[pattern]int),
	}
}

// setDayLimit sets the max occurrences per day for a pattern.
func (pc *patternController) setDayLimit(p pattern, limit int) {
	pc.dayLimits[p] = limit
}

// setDayBetween sets the minimum days between occurrences for a pattern.
func (pc *patternController) setDayBetween(p pattern, days int) {
	pc.dayBetween[p] = days
}

// increment records an occurrence of the pattern on the given date.
func (pc *patternController) increment(d date, p pattern) {
	if pc.occurrences[d] == nil {
		pc.occurrences[d] = make(map[pattern]int)
	}
	pc.occurrences[d][p]++
}

// decrement removes an occurrence of the pattern on the given date.
func (pc *patternController) decrement(d date, p pattern) {
	if pc.occurrences[d] == nil {
		return
	}
	pc.occurrences[d][p]--
	if pc.occurrences[d][p] <= 0 {
		delete(pc.occurrences[d], p)
	}
}

// shouldSkip returns true if the pattern should not be placed on the given date
// because it would violate day limits or day-between constraints.
func (pc *patternController) shouldSkip(d date, p pattern) bool {
	// Check day limit
	if limit, ok := pc.dayLimits[p]; ok {
		current := 0
		if pc.occurrences[d] != nil {
			current = pc.occurrences[d][p]
		}
		if current >= limit {
			return true
		}
	}

	// Check day-between constraint
	if between, ok := pc.dayBetween[p]; ok && between > 0 {
		for delta := 1; delta <= between; delta++ {
			checkDate := d.addDays(-delta)
			if pc.occurrences[checkDate] != nil && pc.occurrences[checkDate][p] > 0 {
				return true
			}
			checkDate = d.addDays(delta)
			if pc.occurrences[checkDate] != nil && pc.occurrences[checkDate][p] > 0 {
				return true
			}
		}
	}

	return false
}

// shouldSkipSubLessons checks all sub-lessons' patterns against constraints.
func (pc *patternController) shouldSkipSubLessons(d date, subLessons []*internalSubLesson, subjectID, groupID uint64) bool {
	for _, sl := range subLessons {
		p := newPattern(groupID, sl.teacherID, subjectID)
		if pc.shouldSkip(d, p) {
			return true
		}
	}
	return false
}
