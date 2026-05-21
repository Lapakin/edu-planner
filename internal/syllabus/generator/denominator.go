package generator

import (
	"math/rand"

	"github.com/Lapakin/edu-planner/internal/domain"
)

// denominatorReproducer mirrors the numerator schedule to the denominator
// and fills remaining denominator-specific lessons.
type denominatorReproducer struct {
	cfg   *settings
	coord *coordinator
	rng   *rand.Rand
}

// newDenominatorReproducer creates a new denominator reproducer.
func newDenominatorReproducer(cfg *settings, coord *coordinator, rng *rand.Rand) *denominatorReproducer {
	return &denominatorReproducer{
		cfg:   cfg,
		coord: coord,
		rng:   rng,
	}
}

// reproduce mirrors the numerator schedule to the denominator,
// then schedules remaining denominator-only lessons.
func (dr *denominatorReproducer) reproduce(
	numeratorDates []date,
	denominatorDates []date,
	denominatorLessons []*lesson,
	groupIDs []uint64,
) error {
	// Group numerator dates by weekday
	numByWeekday := groupDatesByWeekday(numeratorDates)
	denomByWeekday := groupDatesByWeekday(denominatorDates)

	// Build lookup: lesson key → lesson, for matching denominator lessons to numerator slots
	denomLessonMap := make(map[uint64][]*lesson) // groupID → unscheduled lessons
	for _, l := range denominatorLessons {
		if !l.isScheduled {
			denomLessonMap[l.groupID] = append(denomLessonMap[l.groupID], l)
		}
	}

	// For each weekday, try to copy numerator pattern to denominator
	for _, wd := range dr.cfg.educationWeek {
		numDatesForDay := numByWeekday[wd]
		denomDatesForDay := denomByWeekday[wd]

		if len(numDatesForDay) == 0 || len(denomDatesForDay) == 0 {
			continue
		}

		// Use the first numerator date as the template
		templateDate := numDatesForDay[0]
		targetDate := denomDatesForDay[0]

		dr.mirrorDay(templateDate, targetDate, denomLessonMap, groupIDs)
	}

	// Schedule remaining unscheduled denominator lessons
	var remaining []*lesson
	for _, lessons := range denomLessonMap {
		for _, l := range lessons {
			if !l.isScheduled {
				remaining = append(remaining, l)
			}
		}
	}

	if len(remaining) > 0 {
		// Remaining denominator lessons must also land in the first week only.
		denomFirstWeek := firstWeekDates(denominatorDates, dr.cfg.educationWeek)

		randomSched := newRandomScheduler(dr.cfg, dr.coord, dr.rng)
		remaining = randomSched.scheduleRandom(remaining, denomFirstWeek, denomFirstWeek)

		if len(remaining) > 0 {
			recursiveSched := newRecursiveScheduler(dr.cfg, dr.coord)
			if err := recursiveSched.schedule(remaining, denomFirstWeek, denomFirstWeek); err != nil {
				return err
			}
		}
	}

	return nil
}

// mirrorDay copies the schedule pattern from a numerator date to a denominator date.
func (dr *denominatorReproducer) mirrorDay(
	sourceDate, targetDate date,
	denomLessons map[uint64][]*lesson,
	_ []uint64,
) {
	daySchedule := dr.coord.schedule.getDaySchedule(sourceDate)
	if daySchedule == nil {
		return
	}

	for ln, lessons := range daySchedule {
		for _, sourceLesson := range lessons {
			// Find a matching denominator lesson for this group/subject
			available := denomLessons[sourceLesson.groupID]
			var matched *lesson
			matchIdx := -1

			for i, dl := range available {
				if dl.subjectID == sourceLesson.subjectID && !dl.isScheduled {
					matched = dl
					matchIdx = i
					break
				}
			}

			if matched == nil {
				continue
			}

			// Check if we can place it at the same slot
			allParticipants := append([]uint64{matched.groupID}, matched.teacherIDs()...)
			if dr.coord.schedule.hasConflict(targetDate, ln, matched.groupID, matched.teacherIDs()) {
				continue
			}
			if !dr.coord.availability.areFree(targetDate, allParticipants, ln) {
				continue
			}

			// Place the lesson
			dr.coord.schedule.add(targetDate, ln, matched)
			for _, pid := range allParticipants {
				dr.coord.availability.markBusy(targetDate, pid, ln)
			}
			dr.coord.patternCtrl.increment(targetDate, matched.pattern())

			// Remove from available pool
			denomLessons[sourceLesson.groupID] = append(available[:matchIdx], available[matchIdx+1:]...)
		}
	}
}

// weekdayIndex returns the index of a weekday (0=Monday .. 4=Friday).
func weekdayIndex(wd domain.Weekday) int {
	switch wd {
	case domain.WeekdayMonday:
		return 0
	case domain.WeekdayTuesday:
		return 1
	case domain.WeekdayWednesday:
		return 2
	case domain.WeekdayThursday:
		return 3
	case domain.WeekdayFriday:
		return 4
	case domain.WeekdaySaturday:
		return 5
	case domain.WeekdaySunday:
		return 6
	default:
		return 0
	}
}
