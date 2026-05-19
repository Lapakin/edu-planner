package generator

import (
	"github.com/Lapakin/edu-planner/internal/domain"
)

// settings wraps schedule template settings combined with bell schedule data and
// an optional restriction that applies hard daily/gap constraints.
type settings struct {
	hoursPerLesson         float64
	maxIdenticalPerDay     int
	numberOfLessonsInDay   int
	firstLessonNumber      int
	lastLessonNumber       int
	educationWeek          []domain.Weekday
	maxStudyHoursPerDay    int
	maxTeacherHoursPerWeek int
	maxGroupHoursPerWeek   *int
	// computed hard limits (applied from templateSetting + restriction)
	maxGroupLessonsPerDay     int
	minGroupLessonsPerDay     int
	maxTeacherLessonsPerDay   int
	maxGroupLessonsPerWeekV   int
	maxTeacherLessonsPerWeekV int
	noGapsRequired            bool
}

// newSettings creates a settings instance from domain models.
func newSettings(
	templateSetting *domain.ScheduleTemplateSetting,
	bellSchedules domain.BellSchedules,
	restriction *domain.ScheduleRestriction,
) *settings {
	s := &settings{
		hoursPerLesson:            templateSetting.HoursPerLesson,
		maxIdenticalPerDay:        templateSetting.MaxIdenticalLessonsPerDay,
		numberOfLessonsInDay:      0,
		firstLessonNumber:         0,
		lastLessonNumber:          0,
		educationWeek:             domain.DefaultEducationWeek,
		maxStudyHoursPerDay:       templateSetting.MaxStudyHoursPerDay,
		maxTeacherHoursPerWeek:    templateSetting.MaxTeacherHoursPerWeek,
		maxGroupHoursPerWeek:      templateSetting.MaxGroupLessonHoursPerWeek,
		maxGroupLessonsPerDay:     0,
		minGroupLessonsPerDay:     0,
		maxTeacherLessonsPerDay:   0,
		maxGroupLessonsPerWeekV:   0,
		maxTeacherLessonsPerWeekV: 0,
		noGapsRequired:            false,
	}

	if len(bellSchedules) > 0 {
		s.firstLessonNumber = bellSchedules[0].LessonNumber
		s.lastLessonNumber = bellSchedules[0].LessonNumber
		for _, bs := range bellSchedules {
			if bs.LessonNumber < s.firstLessonNumber {
				s.firstLessonNumber = bs.LessonNumber
			}
			if bs.LessonNumber > s.lastLessonNumber {
				s.lastLessonNumber = bs.LessonNumber
			}
		}
		s.numberOfLessonsInDay = s.lastLessonNumber - s.firstLessonNumber + 1
	}

	// Compute hours-based group daily max
	computedGroupMax := s.numberOfLessonsInDay
	if s.hoursPerLesson > 0 {
		computed := int(float64(s.maxStudyHoursPerDay) / s.hoursPerLesson)
		if computed < computedGroupMax {
			computedGroupMax = computed
		}
	}

	// Compute teacher weekly max
	computedTeacherWeekMax := s.numberOfLessonsInDay * len(s.educationWeek)
	if s.hoursPerLesson > 0 {
		computedTeacherWeekMax = int(float64(s.maxTeacherHoursPerWeek) / s.hoursPerLesson)
	}

	// Compute group weekly max
	computedGroupWeekMax := s.numberOfLessonsInDay * len(s.educationWeek)
	if s.maxGroupHoursPerWeek != nil && s.hoursPerLesson > 0 {
		computedGroupWeekMax = int(float64(*s.maxGroupHoursPerWeek) / s.hoursPerLesson)
	}

	s.maxGroupLessonsPerDay = computedGroupMax
	s.maxTeacherLessonsPerDay = computedGroupMax // default: same as group
	s.maxGroupLessonsPerWeekV = computedGroupWeekMax
	s.maxTeacherLessonsPerWeekV = computedTeacherWeekMax

	// Apply restriction (fall back to defaults if nil)
	r := restriction
	if r == nil {
		r = domain.DefaultScheduleRestriction()
	}

	s.minGroupLessonsPerDay = r.MinGroupLessonsPerDay
	s.noGapsRequired = r.NoGapsInGroupSchedule

	if r.MaxGroupLessonsPerDay > 0 && r.MaxGroupLessonsPerDay < s.maxGroupLessonsPerDay {
		s.maxGroupLessonsPerDay = r.MaxGroupLessonsPerDay
	}
	if r.MaxTeacherLessonsPerDay > 0 {
		s.maxTeacherLessonsPerDay = r.MaxTeacherLessonsPerDay
	}

	return s
}

// maxLessonsPerDay returns the max number of lessons a group can have per day.
func (s *settings) maxLessonsPerDay() int {
	return s.maxGroupLessonsPerDay
}

// maxTeacherLessonsPerWeek returns the max lessons a teacher can have per week.
func (s *settings) maxTeacherLessonsPerWeek() int {
	return s.maxTeacherLessonsPerWeekV
}

// maxGroupLessonsPerWeek returns the max lessons a group can have per week.
func (s *settings) maxGroupLessonsPerWeek() int {
	return s.maxGroupLessonsPerWeekV
}

// lessonNumbers returns all valid lesson numbers for a day.
func (s *settings) lessonNumbers() []int {
	numbers := make([]int, 0, s.numberOfLessonsInDay)
	for i := s.firstLessonNumber; i <= s.lastLessonNumber; i++ {
		numbers = append(numbers, i)
	}
	return numbers
}
