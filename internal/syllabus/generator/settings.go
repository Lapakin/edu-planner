package generator

import (
	"github.com/Lapakin/edu-planner/internal/domain"
)

// settings wraps schedule template settings combined with bell schedule data and
// an optional restriction that applies hard daily/gap constraints.
type settings struct {
	lessonsPerClass        int
	maxIdenticalPerDay     int
	numberOfLessonsInDay   int
	firstLessonNumber      int
	lastLessonNumber       int
	educationWeek          []domain.Weekday
	maxStudyHoursPerDay    int
	maxTeacherHoursPerWeek int
	maxGroupHoursPerWeek   *int
	// computed hard limits (applied from templateSetting + restriction)
	maxGroupLessonsPerDay        int
	minGroupLessonsPerDay        int
	maxTeacherLessonsPerDay      int
	maxGroupLessonsPerWeekV      int
	maxTeacherLessonsPerWeekV    int
	noGapsRequired               bool
	maxConsecutiveTeacherLessons int
	timePriority                 domain.TimePriority
	allowFlowLessons             bool
	// teacher slot preferences
	teacherForbiddenSlots map[uint64]map[domain.Weekday]map[int]bool
	teacherPreferredSlots map[uint64]map[domain.Weekday]map[int]bool
	// cycle committee lab rooms: cycleCommitteeID → []roomID
	cycleCommitteeLabRooms map[uint64][]uint64
	// room capacities: roomID → number of seats (0 or absent = unlimited/unset)
	roomCapacities map[uint64]int
	// group sizes: groupID → number of students (0 or absent = unknown)
	groupSizes map[uint64]int
}

// newSettings creates a settings instance from domain models.
func newSettings(
	templateSetting *domain.ScheduleTemplateSetting,
	bellSchedules domain.BellSchedules,
	restriction *domain.ScheduleRestriction,
	preferences domain.TeacherSlotPreferences,
	labRooms domain.CycleCommitteeLabRooms,
) *settings {
	educationWeek := templateSetting.StudyDaysMask.ToWeekdays()

	s := &settings{
		lessonsPerClass:              templateSetting.LessonsPerClass,
		maxIdenticalPerDay:           templateSetting.MaxIdenticalLessonsPerDay,
		numberOfLessonsInDay:         0,
		firstLessonNumber:            0,
		lastLessonNumber:             0,
		educationWeek:                educationWeek,
		maxStudyHoursPerDay:          templateSetting.MaxStudyHoursPerDay,
		maxTeacherHoursPerWeek:       templateSetting.MaxTeacherHoursPerWeek,
		maxGroupHoursPerWeek:         templateSetting.MaxGroupLessonHoursPerWeek,
		maxGroupLessonsPerDay:        0,
		minGroupLessonsPerDay:        0,
		maxTeacherLessonsPerDay:      0,
		maxGroupLessonsPerWeekV:      0,
		maxTeacherLessonsPerWeekV:    0,
		noGapsRequired:               false,
		maxConsecutiveTeacherLessons: 0,
		timePriority:                 domain.TimePriorityNone,
		allowFlowLessons:             true,
		teacherForbiddenSlots:        make(map[uint64]map[domain.Weekday]map[int]bool),
		teacherPreferredSlots:        make(map[uint64]map[domain.Weekday]map[int]bool),
		cycleCommitteeLabRooms:       make(map[uint64][]uint64),
		roomCapacities:               make(map[uint64]int),
		groupSizes:                   make(map[uint64]int),
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

	// Compute lessonsPerClass-based group daily max
	computedGroupMax := s.numberOfLessonsInDay
	if s.lessonsPerClass > 0 {
		computed := s.maxStudyHoursPerDay / s.lessonsPerClass
		if computed < computedGroupMax {
			computedGroupMax = computed
		}
	}

	// Compute teacher weekly max
	computedTeacherWeekMax := s.numberOfLessonsInDay * len(s.educationWeek)
	if s.lessonsPerClass > 0 {
		computedTeacherWeekMax = s.maxTeacherHoursPerWeek / s.lessonsPerClass
	}

	// Compute group weekly max
	computedGroupWeekMax := s.numberOfLessonsInDay * len(s.educationWeek)
	if s.maxGroupHoursPerWeek != nil && s.lessonsPerClass > 0 {
		computedGroupWeekMax = *s.maxGroupHoursPerWeek / s.lessonsPerClass
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
	s.maxConsecutiveTeacherLessons = r.MaxConsecutiveTeacherLessons
	s.timePriority = r.TimePriority
	s.allowFlowLessons = r.AllowFlowLessons

	if r.MaxGroupLessonsPerDay > 0 && r.MaxGroupLessonsPerDay < s.maxGroupLessonsPerDay {
		s.maxGroupLessonsPerDay = r.MaxGroupLessonsPerDay
	}
	if r.MaxTeacherLessonsPerDay > 0 {
		s.maxTeacherLessonsPerDay = r.MaxTeacherLessonsPerDay
	}

	// Build teacher slot preference maps
	for _, pref := range preferences {
		switch pref.SlotType {
		case domain.SlotTypeForbidden:
			if s.teacherForbiddenSlots[pref.TeacherID] == nil {
				s.teacherForbiddenSlots[pref.TeacherID] = make(map[domain.Weekday]map[int]bool)
			}
			if s.teacherForbiddenSlots[pref.TeacherID][pref.Weekday] == nil {
				s.teacherForbiddenSlots[pref.TeacherID][pref.Weekday] = make(map[int]bool)
			}
			s.teacherForbiddenSlots[pref.TeacherID][pref.Weekday][pref.LessonNumber] = true
		case domain.SlotTypePreferred:
			if s.teacherPreferredSlots[pref.TeacherID] == nil {
				s.teacherPreferredSlots[pref.TeacherID] = make(map[domain.Weekday]map[int]bool)
			}
			if s.teacherPreferredSlots[pref.TeacherID][pref.Weekday] == nil {
				s.teacherPreferredSlots[pref.TeacherID][pref.Weekday] = make(map[int]bool)
			}
			s.teacherPreferredSlots[pref.TeacherID][pref.Weekday][pref.LessonNumber] = true
		}
	}

	// Build cycle committee lab room map
	for _, lr := range labRooms {
		s.cycleCommitteeLabRooms[lr.CycleCommitteeID] = append(s.cycleCommitteeLabRooms[lr.CycleCommitteeID], lr.RoomID)
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
