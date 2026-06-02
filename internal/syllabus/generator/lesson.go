package generator

import (
	"fmt"
	"slices"
)

// lessonFormat represents whether a lesson is united or split.
type lessonFormat int

const (
	formatUnited lessonFormat = iota
	formatSplit
)

// internalLessonType distinguishes scheduling behavior.
type internalLessonType int

const (
	lessonTypeRegular internalLessonType = iota
	lessonTypeClassHour
	lessonTypeOneTeacherMultiGroups
	lessonTypeParallel
)

// lesson is the internal mutable representation of a lesson during generation.
type lesson struct {
	groupID          uint64
	subjectID        uint64
	subLessons       []*internalSubLesson
	format           lessonFormat
	isScheduled      bool
	lType            internalLessonType
	workloadID       uint64
	cycleCommitteeID uint64
	isLab            bool
	// flowEligible marks a lesson that may be merged into a flow (a flow
	// discipline's non-lab lecture). Eligibility alone does not make it a flow —
	// only the merged result is a flow lesson.
	flowEligible bool
	// isFlow marks a flow (потокова) lesson: a single lecture shared by several
	// groups at the same slot. A flow lesson carries one sub-lesson per group, all
	// taught by the same teacher.
	isFlow bool
	// flowID groups the lessons that belong to the same flow stream (typically
	// subjectID + teacherID). Empty for non-flow lessons.
	flowID string
	// flowOrigin holds the per-group constituent lessons a flow lesson was merged
	// from. It is used to fall back to separate lessons when no shared slot can be
	// found. Nil for non-flow lessons.
	flowOrigin []*lesson
}

// internalSubLesson represents a teacher assignment within a lesson.
type internalSubLesson struct {
	subGroupNumber *uint64
	groupID        uint64
	teacherID      uint64
	roomID         uint64
}

// pattern returns the unique identifier for this lesson's course pattern.
func (l *lesson) pattern() pattern {
	if len(l.subLessons) == 0 {
		return pattern(fmt.Sprintf("%d:0:%d", l.groupID, l.subjectID))
	}
	return pattern(fmt.Sprintf("%d:%d:%d", l.groupID, l.subLessons[0].teacherID, l.subjectID))
}

// groupIDs returns all distinct group IDs that attend this lesson. For a regular
// lesson this is just its single group; for a flow lesson it is every group in
// the stream.
func (l *lesson) groupIDs() []uint64 {
	if len(l.subLessons) == 0 {
		return []uint64{l.groupID}
	}
	ids := make([]uint64, 0, len(l.subLessons))
	seen := make(map[uint64]bool, len(l.subLessons))
	for _, sl := range l.subLessons {
		if !seen[sl.groupID] {
			seen[sl.groupID] = true
			ids = append(ids, sl.groupID)
		}
	}
	if len(ids) == 0 {
		return []uint64{l.groupID}
	}
	return ids
}

// hasGroup reports whether the given group attends this lesson.
func (l *lesson) hasGroup(groupID uint64) bool {
	return slices.Contains(l.groupIDs(), groupID)
}

// hasAnyGroup reports whether this lesson shares any group with the given set.
func (l *lesson) hasAnyGroup(groupIDs []uint64) bool {
	return slices.ContainsFunc(groupIDs, l.hasGroup)
}

// unmerge restores the per-group constituent lessons of a flow lesson so they can
// be scheduled separately (fallback when no shared slot exists). The constituents
// are reset to an unscheduled, non-flow state. Returns nil for non-flow lessons.
func (l *lesson) unmerge() []*lesson {
	if !l.isFlow {
		return nil
	}
	for _, c := range l.flowOrigin {
		c.isScheduled = false
		c.isFlow = false
		c.flowID = ""
		c.flowOrigin = nil
	}
	return l.flowOrigin
}

// teacherIDs returns all teacher IDs involved in this lesson.
func (l *lesson) teacherIDs() []uint64 {
	ids := make([]uint64, 0, len(l.subLessons))
	for _, sl := range l.subLessons {
		ids = append(ids, sl.teacherID)
	}
	return ids
}

// clone creates a deep copy of the lesson.
func (l *lesson) clone() *lesson {
	c := &lesson{
		groupID:          l.groupID,
		subjectID:        l.subjectID,
		subLessons:       nil,
		format:           l.format,
		isScheduled:      l.isScheduled,
		lType:            l.lType,
		workloadID:       l.workloadID,
		cycleCommitteeID: l.cycleCommitteeID,
		isLab:            l.isLab,
		flowEligible:     l.flowEligible,
		isFlow:           l.isFlow,
		flowID:           l.flowID,
		flowOrigin:       l.flowOrigin,
	}
	c.subLessons = make([]*internalSubLesson, len(l.subLessons))
	for i, sl := range l.subLessons {
		c.subLessons[i] = &internalSubLesson{
			subGroupNumber: sl.subGroupNumber,
			groupID:        sl.groupID,
			teacherID:      sl.teacherID,
			roomID:         sl.roomID,
		}
	}
	return c
}
