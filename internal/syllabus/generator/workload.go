package generator

import (
	"math"

	"github.com/Lapakin/edu-planner/internal/domain"
)

// Workload represents the teaching load for a specific teacher-group-subject combination.
type Workload struct {
	workloadDistributionID uint64
	groupID                uint64
	subjectID              uint64
	teacherID              uint64
	totalHours             int
	isSplitting            bool
	subGroupNumber         *uint64
	cycleCommitteeID       uint64
	isLab                  bool
	isFlow                 bool
}

// workloadDistributor distributes Workloads across numerator and denominator weeks.
type workloadDistributor struct {
	cfg *settings
}

// newWorkloadDistributor creates a new workload distributor.
func newWorkloadDistributor(cfg *settings) *workloadDistributor {
	return &workloadDistributor{cfg: cfg}
}

// distributedWorkload holds the distribution of a Workload across weeks.
type distributedWorkload struct {
	workload           *Workload
	numeratorLessons   int
	denominatorLessons int
}

// distributeWorkloads distributes all Workloads across numerator and denominator weeks.
// It returns lessons to be scheduled in each period.
func (wd *workloadDistributor) distributeWorkloads(
	workloads []*Workload,
	numWeeks int,
) ([]*lesson, []*lesson, error) {
	if numWeeks < 1 {
		numWeeks = 1
	}

	// Split into 2 periods (numerator + denominator)
	weeksPerPeriod := max(numWeeks/2, 1)

	var numeratorLessons, denominatorLessons []*lesson

	for _, w := range workloads {
		dist, err := wd.distributeWorkload(w, weeksPerPeriod)
		if err != nil {
			return nil, nil, err
		}

		numLessons := wd.createLessons(w, dist.numeratorLessons)
		denomLessons := wd.createLessons(w, dist.denominatorLessons)

		numeratorLessons = append(numeratorLessons, numLessons...)
		denominatorLessons = append(denominatorLessons, denomLessons...)
	}

	// Merge flow-eligible lessons (same subject + same teacher across several
	// groups) into shared flow lessons, independently per week period.
	numeratorLessons = mergeFlowLessons(numeratorLessons, wd.cfg)
	denominatorLessons = mergeFlowLessons(denominatorLessons, wd.cfg)

	return numeratorLessons, denominatorLessons, nil
}

// distributeWorkload distributes a single Workload between numerator/denominator.
func (wd *workloadDistributor) distributeWorkload(w *Workload, weeksPerPeriod int) (*distributedWorkload, error) {
	if w.totalHours <= 0 {
		return &distributedWorkload{workload: w, numeratorLessons: 0, denominatorLessons: 0}, nil
	}

	hoursPerLesson := float64(wd.cfg.lessonsPerClass)
	if hoursPerLesson <= 0 {
		hoursPerLesson = 2
	}

	totalLessons := int(math.Ceil(float64(w.totalHours) / hoursPerLesson))

	// Distribute evenly between numerator and denominator
	lessonsPerPeriod := totalLessons / 2
	remainder := totalLessons % 2

	numLessons := lessonsPerPeriod
	denomLessons := lessonsPerPeriod

	// Assign remainder to numerator
	if remainder > 0 {
		numLessons++
	}

	// Now distribute within each period across weeks
	numPerWeek := numLessons
	denomPerWeek := denomLessons

	if weeksPerPeriod > 1 {
		numPerWeek = int(math.Ceil(float64(numLessons) / float64(weeksPerPeriod)))
		denomPerWeek = int(math.Ceil(float64(denomLessons) / float64(weeksPerPeriod)))
	}

	return &distributedWorkload{
		workload:           w,
		numeratorLessons:   numPerWeek,
		denominatorLessons: denomPerWeek,
	}, nil
}

// createLessons creates lesson instances from a Workload for a given number of lessons.
func (wd *workloadDistributor) createLessons(w *Workload, count int) []*lesson {
	lessons := make([]*lesson, 0, count)
	for range count {
		l := &lesson{
			groupID:          w.groupID,
			subjectID:        w.subjectID,
			subLessons:       nil,
			format:           formatUnited,
			isScheduled:      false,
			lType:            lessonTypeRegular,
			workloadID:       w.workloadDistributionID,
			cycleCommitteeID: w.cycleCommitteeID,
			isLab:            w.isLab,
			// Only lecturer (non-lab) lessons of a flow discipline are flow-eligible.
			// Eligibility is promoted to an actual flow only when merged.
			flowEligible: w.isFlow && !w.isLab,
			isFlow:       false,
			flowID:       "",
			flowOrigin:   nil,
		}

		if w.isSplitting && w.subGroupNumber != nil {
			l.format = formatSplit
		}

		sl := &internalSubLesson{
			subGroupNumber: w.subGroupNumber,
			groupID:        w.groupID,
			teacherID:      w.teacherID,
			roomID:         0,
		}
		l.subLessons = []*internalSubLesson{sl}
		lessons = append(lessons, l)
	}
	return lessons
}

// buildWorkloadsFromDomain converts domain workload distributions and assignments to internal Workloads.
func buildWorkloadsFromDomain(
	distributions domain.WorkloadDistributions,
	assignments domain.WorkloadAssignments,
	studyPlans domain.StudyPlans,
	groups domain.Groups,
	disciplines domain.Disciplines,
) []*Workload {
	// Build lookup maps
	assignmentsByDistID := make(map[uint64]domain.WorkloadAssignments)
	for _, a := range assignments {
		assignmentsByDistID[a.WorkloadDistributionID] = append(assignmentsByDistID[a.WorkloadDistributionID], a)
	}

	studyPlanByID := make(map[uint64]*domain.StudyPlan)
	for _, sp := range studyPlans {
		studyPlanByID[sp.ID] = sp
	}

	groupByID := make(map[uint64]*domain.Group)
	for _, g := range groups {
		groupByID[g.ID] = g
	}

	disciplineByID := make(map[uint64]*domain.Discipline)
	for _, d := range disciplines {
		disciplineByID[d.ID] = d
	}

	var workloads []*Workload

	for _, dist := range distributions {
		distAssignments := assignmentsByDistID[dist.ID]

		sp := studyPlanByID[dist.StudyPlanID]
		if sp == nil || sp.DisciplineID == nil {
			continue
		}
		subjectID := *sp.DisciplineID

		g := groupByID[dist.GroupID]
		if g == nil {
			continue
		}

		classroomWork := 0
		if dist.ClassroomWork != nil {
			classroomWork = *dist.ClassroomWork
		}
		labHours := 0
		if dist.Laboratory != nil {
			labHours = *dist.Laboratory
		}
		practHours := 0
		if dist.Practical != nil {
			practHours = *dist.Practical
		}

		if classroomWork <= 0 && labHours <= 0 && practHours <= 0 {
			continue
		}

		// Resolve cycle committee ID and flow flag from the discipline
		var cycleCommitteeID uint64
		var isFlowDiscipline bool
		if disc := disciplineByID[subjectID]; disc != nil {
			cycleCommitteeID = disc.CycleCommitteeID
			isFlowDiscipline = disc.IsFlow
		}

		for _, a := range distAssignments {
			isLab := a.RoleType != nil && *a.RoleType != "Лектор"

			hours := 0
			switch {
			case a.AssignedHours != nil:
				hours = *a.AssignedHours
			case isLab:
				// Lab/practical teacher: use laboratory + practical hours from distribution
				hours = labHours + practHours
			default:
				// Lecturer: use classroom hours minus the lab/practical portion
				lecHours := classroomWork - labHours - practHours
				if lecHours > 0 {
					hours = lecHours
				} else {
					hours = classroomWork
				}
			}

			if hours <= 0 {
				continue
			}

			w := &Workload{
				workloadDistributionID: dist.ID,
				groupID:                dist.GroupID,
				subjectID:              subjectID,
				teacherID:              a.TeacherID,
				totalHours:             hours,
				isSplitting:            g.IsSplitting,
				subGroupNumber:         nil,
				cycleCommitteeID:       cycleCommitteeID,
				isLab:                  isLab,
				isFlow:                 isFlowDiscipline,
			}

			workloads = append(workloads, w)
		}
	}

	return workloads
}
