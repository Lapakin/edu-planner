package generator

import (
	"fmt"
	"testing"
)

func makeUnitedLesson(groupID, subjectID uint64, teacherID uint64) *lesson {
	return &lesson{
		groupID:   groupID,
		subjectID: subjectID,
		format:    formatUnited,
		lType:     lessonTypeRegular,
		subLessons: []*internalSubLesson{
			{groupID: groupID, teacherID: teacherID},
		},
	}
}

func TestLesson_Pattern_WithSubLesson(t *testing.T) {
	l := makeUnitedLesson(1, 10, 42)
	p := l.pattern()
	expected := pattern(fmt.Sprintf("%d:%d:%d", 1, 42, 10))
	if p != expected {
		t.Errorf("expected pattern %q, got %q", expected, p)
	}
}

func TestLesson_Pattern_NoSubLesson(t *testing.T) {
	l := &lesson{groupID: 1, subjectID: 10}
	p := l.pattern()
	expected := pattern(fmt.Sprintf("%d:0:%d", 1, 10))
	if p != expected {
		t.Errorf("expected pattern %q, got %q", expected, p)
	}
}

func TestLesson_TeacherIDs_Empty(t *testing.T) {
	l := &lesson{groupID: 1, subjectID: 10}
	ids := l.teacherIDs()
	if len(ids) != 0 {
		t.Errorf("expected empty teacher IDs, got %v", ids)
	}
}

func TestLesson_TeacherIDs_OneTeacher(t *testing.T) {
	l := makeUnitedLesson(1, 10, 42)
	ids := l.teacherIDs()
	if len(ids) != 1 || ids[0] != 42 {
		t.Errorf("expected [42], got %v", ids)
	}
}

func TestLesson_TeacherIDs_MultipleTeachers(t *testing.T) {
	l := &lesson{
		groupID:   1,
		subjectID: 10,
		format:    formatSplit,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 10},
			{groupID: 1, teacherID: 20},
		},
	}
	ids := l.teacherIDs()
	if len(ids) != 2 {
		t.Fatalf("expected 2 teacher IDs, got %v", ids)
	}
	if ids[0] != 10 || ids[1] != 20 {
		t.Errorf("expected [10, 20], got %v", ids)
	}
}

func TestLesson_Clone_DeepCopy(t *testing.T) {
	sg := uint64(1)
	original := &lesson{
		groupID:    1,
		subjectID:  10,
		format:     formatUnited,
		lType:      lessonTypeRegular,
		workloadID: 5,
		subLessons: []*internalSubLesson{
			{groupID: 1, teacherID: 42, roomID: 101, subGroupNumber: &sg},
		},
	}

	cloned := original.clone()

	if cloned == original {
		t.Error("clone should be a different pointer")
	}
	if cloned.groupID != original.groupID {
		t.Errorf("groupID mismatch: %d vs %d", cloned.groupID, original.groupID)
	}
	if cloned.subjectID != original.subjectID {
		t.Errorf("subjectID mismatch")
	}
	if cloned.format != original.format {
		t.Errorf("format mismatch")
	}
	if cloned.lType != original.lType {
		t.Errorf("lType mismatch")
	}
	if cloned.workloadID != original.workloadID {
		t.Errorf("workloadID mismatch")
	}
	if len(cloned.subLessons) != len(original.subLessons) {
		t.Fatalf("subLessons length mismatch")
	}
	if cloned.subLessons[0] == original.subLessons[0] {
		t.Error("subLessons[0] should be different pointers")
	}
	if cloned.subLessons[0].teacherID != 42 {
		t.Errorf("expected cloned teacherID=42, got %d", cloned.subLessons[0].teacherID)
	}
}

func TestLesson_Clone_ModifyCloneDoesNotAffectOriginal(t *testing.T) {
	original := makeUnitedLesson(1, 10, 42)
	cloned := original.clone()

	cloned.groupID = 99
	cloned.subLessons[0].teacherID = 999

	if original.groupID != 1 {
		t.Error("modifying clone affected original groupID")
	}
	if original.subLessons[0].teacherID != 42 {
		t.Error("modifying clone affected original teacherID")
	}
}

func TestLesson_IsScheduled_DefaultFalse(t *testing.T) {
	l := makeUnitedLesson(1, 10, 42)
	if l.isScheduled {
		t.Error("expected isScheduled=false for new lesson")
	}
}
