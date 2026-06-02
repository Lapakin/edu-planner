package generator

import (
	"slices"
	"sort"
)

// schedule is the internal mutable schedule storage used during generation.
// It maps date → lessonNumber → []*lesson.
type schedule struct {
	dateToLessons map[date]map[int][]*lesson
}

// newSchedule creates an empty schedule.
func newSchedule() *schedule {
	return &schedule{
		dateToLessons: make(map[date]map[int][]*lesson),
	}
}

// add places a lesson at the given date and lesson number.
func (s *schedule) add(d date, lessonNumber int, l *lesson) {
	if s.dateToLessons[d] == nil {
		s.dateToLessons[d] = make(map[int][]*lesson)
	}
	s.dateToLessons[d][lessonNumber] = append(s.dateToLessons[d][lessonNumber], l)
	l.isScheduled = true
}

// remove removes a lesson from the given date and lesson number.
func (s *schedule) remove(d date, lessonNumber int, l *lesson) {
	if s.dateToLessons[d] == nil {
		return
	}
	lessons := s.dateToLessons[d][lessonNumber]
	for i, existing := range lessons {
		if existing == l {
			s.dateToLessons[d][lessonNumber] = append(lessons[:i], lessons[i+1:]...)
			l.isScheduled = false
			return
		}
	}
}

// getLessons returns all lessons at a specific date and lesson number.
func (s *schedule) getLessons(d date, lessonNumber int) []*lesson {
	if s.dateToLessons[d] == nil {
		return nil
	}
	return s.dateToLessons[d][lessonNumber]
}

// getDaySchedule returns all lessons for a specific date.
func (s *schedule) getDaySchedule(d date) map[int][]*lesson {
	return s.dateToLessons[d]
}

// getGroupDayLessons returns all lessons for a specific group on a date.
func (s *schedule) getGroupDayLessons(d date, groupID uint64) []*lesson {
	daySchedule := s.dateToLessons[d]
	if daySchedule == nil {
		return nil
	}
	var result []*lesson
	for _, lessons := range daySchedule {
		for _, l := range lessons {
			if l.hasGroup(groupID) {
				result = append(result, l)
			}
		}
	}
	return result
}

// getGroupDayLessonCount returns the number of lessons for a group on a date.
func (s *schedule) getGroupDayLessonCount(d date, groupID uint64) int {
	return len(s.getGroupDayLessons(d, groupID))
}

// getTeacherDayLessons returns all lessons for a specific teacher on a date.
func (s *schedule) getTeacherDayLessons(d date, teacherID uint64) []*lesson {
	daySchedule := s.dateToLessons[d]
	if daySchedule == nil {
		return nil
	}
	var result []*lesson
	for _, lessons := range daySchedule {
		for _, l := range lessons {
			for _, sl := range l.subLessons {
				if sl.teacherID == teacherID {
					result = append(result, l)
					break
				}
			}
		}
	}
	return result
}

// getGroupWeekLessonCount counts lessons for a group across multiple dates (a week).
func (s *schedule) getGroupWeekLessonCount(dates []date, groupID uint64) int {
	count := 0
	for _, d := range dates {
		count += s.getGroupDayLessonCount(d, groupID)
	}
	return count
}

// getTeacherWeekLessonCount counts lessons for a teacher across multiple dates.
func (s *schedule) getTeacherWeekLessonCount(dates []date, teacherID uint64) int {
	count := 0
	for _, d := range dates {
		count += len(s.getTeacherDayLessons(d, teacherID))
	}
	return count
}

// getScheduledLessonNumbersForGroups returns sorted lesson numbers occupied on a
// date by any of the given groups (the union of their busy slots). It is used to
// place a flow lesson, which must be free for every participating group.
func (s *schedule) getScheduledLessonNumbersForGroups(d date, groupIDs []uint64) []int {
	if len(groupIDs) == 1 {
		return s.getScheduledLessonNumbers(d, groupIDs[0])
	}
	numbersMap := make(map[int]bool)
	for _, gid := range groupIDs {
		for _, ln := range s.getScheduledLessonNumbers(d, gid) {
			numbersMap[ln] = true
		}
	}
	numbers := make([]int, 0, len(numbersMap))
	for ln := range numbersMap {
		numbers = append(numbers, ln)
	}
	sort.Ints(numbers)
	return numbers
}

// getScheduledLessonNumbers returns sorted lesson numbers that have lessons on a date for a group.
func (s *schedule) getScheduledLessonNumbers(d date, groupID uint64) []int {
	daySchedule := s.dateToLessons[d]
	if daySchedule == nil {
		return nil
	}
	numbersMap := make(map[int]bool)
	for ln, lessons := range daySchedule {
		for _, l := range lessons {
			if l.hasGroup(groupID) {
				numbersMap[ln] = true
				break
			}
		}
	}
	numbers := make([]int, 0, len(numbersMap))
	for ln := range numbersMap {
		numbers = append(numbers, ln)
	}
	sort.Ints(numbers)
	return numbers
}

// hasConflict checks if placing a lesson for the given groups and teachers would
// conflict with existing lessons at the slot. groupIDs may hold several groups for
// a flow lesson; a conflict arises if any of them (or any teacher) is already busy.
func (s *schedule) hasConflict(d date, lessonNumber int, groupIDs []uint64, teacherIDs []uint64) bool {
	lessons := s.getLessons(d, lessonNumber)
	for _, existing := range lessons {
		// Group conflict
		if existing.hasAnyGroup(groupIDs) {
			return true
		}
		// Teacher conflict
		for _, existingTID := range existing.teacherIDs() {
			if slices.Contains(teacherIDs, existingTID) {
				return true
			}
		}
	}
	return false
}

// findConflictingLesson returns the lesson that conflicts at a given slot.
func (s *schedule) findConflictingLesson(d date, lessonNumber int, groupIDs []uint64, teacherIDs []uint64) *lesson {
	lessons := s.getLessons(d, lessonNumber)
	for _, existing := range lessons {
		if existing.hasAnyGroup(groupIDs) {
			return existing
		}
		for _, existingTID := range existing.teacherIDs() {
			if slices.Contains(teacherIDs, existingTID) {
				return existing
			}
		}
	}
	return nil
}

// allDates returns all dates that have at least one lesson scheduled.
func (s *schedule) allDates() []date {
	dates := make([]date, 0, len(s.dateToLessons))
	for d := range s.dateToLessons {
		dates = append(dates, d)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].before(dates[j])
	})
	return dates
}
