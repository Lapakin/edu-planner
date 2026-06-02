package generator

const (
	maxRecursionDepth = 100
	maxRetryPerLesson = 20
)

// recursiveScheduler uses recursive backtracking to schedule lessons
// that couldn't be placed by the random scheduler.
type recursiveScheduler struct {
	cfg   *settings
	coord *coordinator
	depth int
}

// newRecursiveScheduler creates a new recursive scheduler.
func newRecursiveScheduler(cfg *settings, coord *coordinator) *recursiveScheduler {
	return &recursiveScheduler{
		cfg:   cfg,
		coord: coord,
		depth: 0,
	}
}

// schedule attempts to place all unscheduled lessons using recursive backtracking.
func (rs *recursiveScheduler) schedule(unscheduled []*lesson, dates []date, weekDates []date) error {
	for _, l := range unscheduled {
		if l.isScheduled {
			continue
		}

		placed := false
		for retry := 0; retry < maxRetryPerLesson && !placed; retry++ {
			rs.depth = 0
			if rs.tryScheduleLesson(l, dates, weekDates) {
				placed = true
			}
		}

		if !placed {
			return ErrUnableToSchedule
		}
	}
	return nil
}

// tryScheduleLesson tries to place a lesson, potentially displacing others.
func (rs *recursiveScheduler) tryScheduleLesson(l *lesson, dates []date, weekDates []date) bool {
	if rs.depth >= maxRecursionDepth {
		return false
	}
	rs.depth++

	// First: try to find a free slot
	for _, d := range dates {
		if rs.coord.patternCtrl.shouldSkipSubLessons(d, l.subLessons, l.subjectID, l.groupID) {
			continue
		}

		if !rs.coord.limiter.canPlaceLesson(rs.coord.schedule, d, weekDates, l) {
			continue
		}

		freeSlots := rs.coord.finder.findAllFreeSlots(rs.coord.schedule, d, l.groupIDs())

		for _, ln := range freeSlots {
			if rs.isForbiddenForTeachers(d, l.teacherIDs(), ln) {
				continue
			}
			if !rs.coord.schedule.hasConflict(d, ln, l.groupIDs(), l.teacherIDs()) {
				allParticipants := append(l.groupIDs(), l.teacherIDs()...)
				if rs.coord.availability.areFree(d, allParticipants, ln) &&
					!rs.violatesConsecutive(d, l.teacherIDs(), ln) {
					// Place directly
					rs.placeLesson(d, ln, l)
					return true
				}
			}
		}
	}

	// Second: try to displace a conflicting lesson
	for _, d := range dates {
		if rs.coord.patternCtrl.shouldSkipSubLessons(d, l.subLessons, l.subjectID, l.groupID) {
			continue
		}

		freeSlots := rs.coord.finder.findAllFreeSlots(rs.coord.schedule, d, l.groupIDs())

		for _, ln := range freeSlots {
			if rs.isForbiddenForTeachers(d, l.teacherIDs(), ln) {
				continue
			}
			conflicting := rs.coord.schedule.findConflictingLesson(d, ln, l.groupIDs(), l.teacherIDs())
			if conflicting == nil {
				continue
			}

			// Try to displace the conflicting lesson
			rs.removeLesson(d, ln, conflicting)
			rs.placeLesson(d, ln, l)

			// Recursively schedule the displaced lesson
			if rs.tryScheduleLesson(conflicting, dates, weekDates) {
				return true
			}

			// Backtrack: undo the displacement
			rs.removeLesson(d, ln, l)
			rs.placeLesson(d, ln, conflicting)
		}
	}

	return false
}

// isForbiddenForTeachers returns true if the given lesson number is forbidden for any teacher on date d.
func (rs *recursiveScheduler) isForbiddenForTeachers(d date, teacherIDs []uint64, ln int) bool {
	wd := d.weekday()
	for _, tid := range teacherIDs {
		if rs.cfg.teacherForbiddenSlots[tid] != nil && rs.cfg.teacherForbiddenSlots[tid][wd][ln] {
			return true
		}
	}
	return false
}

// violatesConsecutive returns true if placing for any teacher at ln on d would
// exceed the configured maxConsecutiveTeacherLessons limit.
func (rs *recursiveScheduler) violatesConsecutive(d date, teacherIDs []uint64, ln int) bool {
	maxC := rs.cfg.maxConsecutiveTeacherLessons
	if maxC <= 0 {
		return false
	}
	for _, tid := range teacherIDs {
		if rs.coord.availability.wouldExceedConsecutive(d, tid, ln, maxC) {
			return true
		}
	}
	return false
}

// placeLesson places a lesson and updates all tracking.
func (rs *recursiveScheduler) placeLesson(d date, ln int, l *lesson) {
	rs.coord.schedule.add(d, ln, l)
	allParticipants := append(l.groupIDs(), l.teacherIDs()...)
	for _, pid := range allParticipants {
		rs.coord.availability.markBusy(d, pid, ln)
	}
	rs.coord.patternCtrl.increment(d, l.pattern())
}

// removeLesson removes a lesson and updates all tracking.
func (rs *recursiveScheduler) removeLesson(d date, ln int, l *lesson) {
	rs.coord.schedule.remove(d, ln, l)
	allParticipants := append(l.groupIDs(), l.teacherIDs()...)
	for _, pid := range allParticipants {
		rs.coord.availability.markFree(d, pid, ln)
	}
	rs.coord.patternCtrl.decrement(d, l.pattern())
}
