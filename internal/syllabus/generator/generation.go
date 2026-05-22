package generator

import (
	"math/rand"
	"sort"

	"github.com/Lapakin/edu-planner/internal/domain"
)

// generation represents a single independent generation attempt.
type generation struct {
	cfg       *settings
	workloads []*Workload
	roomIDs   []uint64
	groupIDs  []uint64
	startDate date
	endDate   date
	rng       *rand.Rand
}

// newGeneration creates a new generation attempt.
func newGeneration(
	cfg *settings,
	workloads []*Workload,
	roomIDs []uint64,
	groupIDs []uint64,
	startDate, endDate date,
	seed int64,
) *generation {
	return &generation{
		cfg:       cfg,
		workloads: workloads,
		roomIDs:   roomIDs,
		groupIDs:  groupIDs,
		startDate: startDate,
		endDate:   endDate,
		rng:       rand.New(rand.NewSource(seed)),
	}
}

// exec runs the full generation pipeline for one attempt.
func (g *generation) exec() (*domain.ScheduleData, error) {
	// 1. Build coordinator
	coord := newCoordinator(g.cfg, g.roomIDs)

	// 2. Set pattern day limits
	for _, w := range g.workloads {
		p := newPattern(w.groupID, w.teacherID, w.subjectID)
		coord.patternCtrl.setDayLimit(p, g.cfg.maxIdenticalPerDay)
	}

	// 3. Collect semester dates
	numeratorDates, denominatorDates := collectSemesterDates(
		g.startDate.toTime(),
		g.endDate.toTime(),
		g.cfg.educationWeek,
	)

	if len(numeratorDates) == 0 {
		return nil, ErrValidationFailed
	}

	// 4. Distribute workloads
	numWeeks := countWeeks(g.startDate.toTime(), g.endDate.toTime())
	distributor := newWorkloadDistributor(g.cfg)
	numLessons, denomLessons, err := distributor.distributeWorkloads(g.workloads, numWeeks)
	if err != nil {
		return nil, ErrWorkloadDistribution
	}

	// 5. Extract first-week dates — the template represents one repeating week,
	// so all scheduling must happen within these dates only.
	numFirstWeek := firstWeekDates(numeratorDates, g.cfg.educationWeek)

	// 5a. Pre-mark teacher forbidden slots across all dates
	g.applyTeacherForbiddenSlots(coord.availability, append(numeratorDates, denominatorDates...))

	// 5b. Validate generation possibility
	v := newValidator(g.cfg)
	if err := v.validate(numLessons, numFirstWeek, g.groupIDs); err != nil {
		return nil, err
	}
	if len(denominatorDates) > 0 && len(denomLessons) > 0 {
		denomFirstWeekForValidation := firstWeekDates(denominatorDates, g.cfg.educationWeek)
		if err := v.validate(denomLessons, denomFirstWeekForValidation, g.groupIDs); err != nil {
			return nil, err
		}
	}

	// 6. Random scheduling for numerator (first week only)
	randomSched := newRandomScheduler(g.cfg, coord, g.rng)
	unscheduled := randomSched.scheduleRandom(numLessons, numFirstWeek, numFirstWeek)

	// 7. Fix week schedule (gap minimization)
	coord.fixer.fixWeekSchedule(numFirstWeek, g.groupIDs)

	// 8. Recursive scheduling for remaining numerator lessons
	if len(unscheduled) > 0 {
		recursiveSched := newRecursiveScheduler(g.cfg, coord)
		if err := recursiveSched.schedule(unscheduled, numFirstWeek, numFirstWeek); err != nil {
			return nil, err
		}
	}

	// 8a. Post-placement validation (no-gaps and min lessons per day)
	if err := validatePlacement(coord.schedule, numFirstWeek, g.groupIDs, g.cfg); err != nil {
		return nil, err
	}

	// 9. Denominator reproduction
	if len(denominatorDates) > 0 && len(denomLessons) > 0 {
		denomFirstWeek := firstWeekDates(denominatorDates, g.cfg.educationWeek)
		denomReproducer := newDenominatorReproducer(g.cfg, coord, g.rng)
		if err := denomReproducer.reproduce(numFirstWeek, denomFirstWeek, denomLessons, g.groupIDs); err != nil {
			return nil, err
		}

		// Fix denominator gaps
		coord.fixer.fixWeekSchedule(denomFirstWeek, g.groupIDs)

		// Post-placement validation for denominator
		if err := validatePlacement(coord.schedule, denomFirstWeek, g.groupIDs, g.cfg); err != nil {
			return nil, err
		}
	}

	// 10. Room assignment
	coord.roomSelector.assignRooms(coord.schedule)

	// 11. Convert to output format
	return g.convertToScheduleData(coord.schedule, numeratorDates, denominatorDates), nil
}

// applyTeacherForbiddenSlots pre-marks forbidden teacher slots as unavailable.
func (g *generation) applyTeacherForbiddenSlots(avail *availability, allDates []date) {
	for _, d := range allDates {
		wd := d.weekday()
		for teacherID, dayMap := range g.cfg.teacherForbiddenSlots {
			for ln := range dayMap[wd] {
				avail.markUnavailable(d, teacherID, ln)
			}
		}
	}
}

// convertToScheduleData converts internal schedule to the domain ScheduleData format.
func (g *generation) convertToScheduleData(
	sched *schedule,
	numeratorDates []date,
	denominatorDates []date,
) *domain.ScheduleData {
	data := &domain.ScheduleData{
		Numerator:   make(domain.WeekSchedule),
		Denominator: make(domain.WeekSchedule),
		Metadata:    "",
	}

	// Build numerator week from the first week's dates
	g.buildWeekSchedule(data.Numerator, sched, numeratorDates)

	// Build denominator week
	if len(denominatorDates) > 0 {
		g.buildWeekSchedule(data.Denominator, sched, denominatorDates)
	}

	return data
}

// buildWeekSchedule populates a WeekSchedule from scheduled dates.
func (g *generation) buildWeekSchedule(
	weekSchedule domain.WeekSchedule,
	sched *schedule,
	dates []date,
) {
	// Group dates by weekday, use the first occurrence of each weekday
	byWeekday := groupDatesByWeekday(dates)

	for _, wd := range g.cfg.educationWeek {
		datesForDay := byWeekday[wd]
		if len(datesForDay) == 0 {
			continue
		}

		// Use the first date for this weekday
		d := datesForDay[0]
		dayLessons := sched.getDaySchedule(d)
		if dayLessons == nil {
			continue
		}

		daySchedule := make(domain.DaySchedule)

		// Sort lesson numbers
		lessonNumbers := make([]int, 0, len(dayLessons))
		for ln := range dayLessons {
			lessonNumbers = append(lessonNumbers, ln)
		}
		sort.Ints(lessonNumbers)

		for _, ln := range lessonNumbers {
			lessons := dayLessons[ln]
			scheduleLessons := make([]*domain.ScheduleLesson, 0, len(lessons))

			for _, l := range lessons {
				sl := g.convertLesson(l)
				scheduleLessons = append(scheduleLessons, sl)
			}

			daySchedule[ln] = scheduleLessons
		}

		weekSchedule[wd] = &daySchedule
	}
}

// convertLesson converts an internal lesson to a domain ScheduleLesson.
func (g *generation) convertLesson(l *lesson) *domain.ScheduleLesson {
	lessonType := domain.LessonTypeUnited
	switch l.format {
	case formatUnited:
		lessonType = domain.LessonTypeUnited
	case formatSplit:
		lessonType = domain.LessonTypeSplit
	}

	subLessons := make([]*domain.ScheduleSubLesson, 0, len(l.subLessons))
	for _, sl := range l.subLessons {
		subLessons = append(subLessons, &domain.ScheduleSubLesson{
			SubGroupNumber: sl.subGroupNumber,
			GroupID:        sl.groupID,
			RoomID:         sl.roomID,
			TeacherID:      sl.teacherID,
		})
	}

	// Sort sub-lessons by group ID
	sort.Slice(subLessons, func(i, j int) bool {
		return subLessons[i].GroupID < subLessons[j].GroupID
	})

	return &domain.ScheduleLesson{
		Type:       lessonType,
		SubjectID:  l.subjectID,
		IsLab:      l.isLab,
		SubLessons: subLessons,
	}
}
