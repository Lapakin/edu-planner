package generator

import "fmt"

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
	groupID     uint64
	subjectID   uint64
	subLessons  []*internalSubLesson
	format      lessonFormat
	isScheduled bool
	lType       internalLessonType
	workloadID  uint64
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
		groupID:     l.groupID,
		subjectID:   l.subjectID,
		subLessons:  nil,
		format:      l.format,
		isScheduled: l.isScheduled,
		lType:       l.lType,
		workloadID:  l.workloadID,
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
