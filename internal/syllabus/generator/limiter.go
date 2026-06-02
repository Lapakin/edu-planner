package generator

// limiter enforces daily and weekly lesson limits per group and teacher.
type limiter struct {
	cfg *settings
}

// newLimiter creates a new limiter with the given settings.
func newLimiter(cfg *settings) *limiter {
	return &limiter{cfg: cfg}
}

// canAddLessonForGroup checks if a group can have another lesson on the given date.
func (lim *limiter) canAddLessonForGroup(sched *schedule, d date, groupID uint64) bool {
	current := sched.getGroupDayLessonCount(d, groupID)
	return current < lim.cfg.maxGroupLessonsPerDay
}

// canAddLessonForTeacher checks if a teacher has not exceeded the daily lesson limit.
func (lim *limiter) canAddLessonForTeacher(sched *schedule, d date, teacherID uint64) bool {
	current := len(sched.getTeacherDayLessons(d, teacherID))
	return current < lim.cfg.maxTeacherLessonsPerDay
}

// canAddWeekLessonForGroup checks if a group can have another lesson in the week.
func (lim *limiter) canAddWeekLessonForGroup(sched *schedule, weekDates []date, groupID uint64) bool {
	current := sched.getGroupWeekLessonCount(weekDates, groupID)
	return current < lim.cfg.maxGroupLessonsPerWeek()
}

// canAddWeekLessonForTeacher checks if a teacher can have another lesson in the week.
func (lim *limiter) canAddWeekLessonForTeacher(sched *schedule, weekDates []date, teacherID uint64) bool {
	current := sched.getTeacherWeekLessonCount(weekDates, teacherID)
	return current < lim.cfg.maxTeacherLessonsPerWeek()
}

// canPlaceLesson checks all limits for placing a lesson.
func (lim *limiter) canPlaceLesson(sched *schedule, d date, weekDates []date, l *lesson) bool {
	for _, gid := range l.groupIDs() {
		if !lim.canAddLessonForGroup(sched, d, gid) {
			return false
		}
		if !lim.canAddWeekLessonForGroup(sched, weekDates, gid) {
			return false
		}
	}

	for _, tid := range l.teacherIDs() {
		if !lim.canAddLessonForTeacher(sched, d, tid) {
			return false
		}
		if !lim.canAddWeekLessonForTeacher(sched, weekDates, tid) {
			return false
		}
	}

	return true
}
