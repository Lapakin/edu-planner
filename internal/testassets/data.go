package testassets

import (
	"time"

	"github.com/Lapakin/edu-planner/internal/domain"
)

//go:fix inline
func Ptr[T any](v T) *T { return new(v) }

const (
	GarbageString = "garbage"
	EmptyString   = ""
	EmptyToken    = ""
)

var (
	ExpectedErrorUnmarshal       = map[string]any{"error": "cannot unmarshal"}
	ExpectedResponseEmptyBody    = map[string]any{"error": "request must contain elements in array"}
	ExpectedResponseBadJWT       = map[string]any{"error": "Authorization header required"}
	ExpectedResponseHaveNoAccess = map[string]any{"error": "forbidden"}
	ExpectedResponseNotFound     = map[string]any{"error": "not found"}
	ExpectedResponseInvalidCreds = map[string]any{"error": "invalid email or password"}
)

var DefaultIgnoredFields = []string{"created_at", "modified_at"}

var (
	Test        = "test"
	CurrentTime = time.Now()

	AcademicYear1 = &domain.AcademicYear{
		ID:        1,
		Name:      "2023-2024",
		StartDate: time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, time.June, 30, 0, 0, 0, 0, time.UTC),
	}
	AcademicYear2 = &domain.AcademicYear{
		ID:        2,
		Name:      "2024-2025",
		StartDate: time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, time.June, 30, 0, 0, 0, 0, time.UTC),
	}
	AcademicYearsArray = domain.AcademicYears{AcademicYear1, AcademicYear2}

	Semester1 = &domain.Semester{
		ID:             1,
		AcademicYearID: AcademicYear1.ID,
		PeriodStart:    time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC),
	}
	Semester2 = &domain.Semester{
		ID:             2,
		AcademicYearID: AcademicYear1.ID,
		PeriodStart:    time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, time.June, 30, 0, 0, 0, 0, time.UTC),
	}
	SemestersArray = domain.Semesters{Semester1, Semester2}

	BellSchedule1 = &domain.BellSchedule{
		ID:           1,
		LessonNumber: 1,
		StartTime:    "08:00:00",
		EndTime:      "09:20:00",
	}
	BellSchedule2 = &domain.BellSchedule{
		ID:           2,
		LessonNumber: 2,
		StartTime:    "09:35:00",
		EndTime:      "10:55:00",
	}
	BellSchedulesArray = domain.BellSchedules{BellSchedule1, BellSchedule2}

	Room1 = &domain.Room{
		ID:             1,
		AcademicYearID: AcademicYear1.ID,
		Name:           "101",
		RoomType:       domain.RoomTypeAuditorium,
	}
	Room2 = &domain.Room{
		ID:             2,
		AcademicYearID: AcademicYear1.ID,
		Name:           "Lab-1",
		RoomType:       domain.RoomTypeLaboratory,
	}
	RoomsArray = domain.Rooms{Room1, Room2}

	Department1 = &domain.Department{
		ID:             1,
		AcademicYearID: AcademicYear1.ID,
		Name:           "Computer Science Department",
	}
	Department2 = &domain.Department{
		ID:             2,
		AcademicYearID: AcademicYear1.ID,
		Name:           "Physics Department",
	}
	DepartmentsArray = domain.Departments{Department1, Department2}

	Specialty1 = &domain.Specialty{
		ID:             1,
		AcademicYearID: AcademicYear1.ID,
		DepartmentID:   Department1.ID,
		Name:           "Computer Science",
		ShortName:      "CS",
		SpecialtyCode:  "122",
	}
	Specialty2 = &domain.Specialty{
		ID:             2,
		AcademicYearID: AcademicYear1.ID,
		DepartmentID:   Department1.ID,
		Name:           "Software Engineering",
		ShortName:      "SE",
		SpecialtyCode:  "121",
	}
	SpecialtiesArray = domain.Specialties{Specialty1, Specialty2}

	User1 = &domain.User{
		ID:        1,
		Email:     "user1@test.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      domain.RoleTeacher,
	}
	User2 = &domain.User{
		ID:        2,
		Email:     "user2@test.com",
		FirstName: "Jane",
		LastName:  "Smith",
		Role:      domain.RoleTeacher,
	}
	UsersArray = domain.Users{User1, User2}

	CycleCommittee1 = &domain.CycleCommittee{
		ID:             1,
		AcademicYearID: AcademicYear1.ID,
		UserID:         User1.ID,
		Name:           "Math Cycle Committee",
	}
	CycleCommittee2 = &domain.CycleCommittee{
		ID:             2,
		AcademicYearID: AcademicYear1.ID,
		UserID:         User1.ID,
		Name:           "Physics Cycle Committee",
	}
	CycleCommitteesArray = domain.CycleCommittees{CycleCommittee1, CycleCommittee2}

	Teacher1 = &domain.Teacher{
		ID:             1,
		AcademicYearID: AcademicYear1.ID,
		UserID:         User1.ID,
	}
	Teacher2 = &domain.Teacher{
		ID:             2,
		AcademicYearID: AcademicYear1.ID,
		UserID:         User1.ID,
	}
	TeachersArray = domain.Teachers{Teacher1, Teacher2}

	Group1 = &domain.Group{
		ID:             1,
		AcademicYearID: AcademicYear1.ID,
		SpecialtyID:    Specialty1.ID,
		Name:           "CS-1",
		ShortName:      "CS1",
		EducationStart: time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC),
		EducationEnd:   time.Date(2027, time.June, 30, 0, 0, 0, 0, time.UTC),
	}
	Group2 = &domain.Group{
		ID:             2,
		AcademicYearID: AcademicYear1.ID,
		SpecialtyID:    Specialty1.ID,
		Name:           "CS-2",
		ShortName:      "CS2",
		EducationStart: time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC),
		EducationEnd:   time.Date(2027, time.June, 30, 0, 0, 0, 0, time.UTC),
	}
	GroupsArray = domain.Groups{Group1, Group2}

	Discipline1 = &domain.Discipline{
		ID:               1,
		AcademicYearID:   AcademicYear1.ID,
		CycleCommitteeID: CycleCommittee1.ID,
		Name:             "Mathematics",
		ShortName:        "Math",
		IsSplitting:      false,
		IsSubvention:     false,
	}
	Discipline2 = &domain.Discipline{
		ID:               2,
		AcademicYearID:   AcademicYear1.ID,
		CycleCommitteeID: CycleCommittee1.ID,
		Name:             "Physics",
		ShortName:        "Phys",
		IsSplitting:      true,
		IsSubvention:     false,
	}
	DisciplinesArray = domain.Disciplines{Discipline1, Discipline2}

	GroupSemester1 = &domain.GroupSemester{
		ID:         1,
		GroupID:    Group1.ID,
		SemesterID: Semester1.ID,
		StartDate:  time.Date(2023, time.September, 1, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC),
	}
	GroupSemester2 = &domain.GroupSemester{
		ID:         2,
		GroupID:    Group2.ID,
		SemesterID: Semester1.ID,
		StartDate:  time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2024, time.June, 30, 0, 0, 0, 0, time.UTC),
	}
	GroupSemestersArray = domain.GroupSemesters{GroupSemester1, GroupSemester2}

	StudyPlan1 = &domain.StudyPlan{
		ID:             1,
		AcademicYearID: new(AcademicYear1.ID),
		SpecialtyID:    new(Specialty1.ID),
		DisciplineID:   new(Discipline1.ID),
		SemesterNumber: new(1),
		Lectures:       new(30),
	}
	StudyPlan2 = &domain.StudyPlan{
		ID:             2,
		AcademicYearID: new(AcademicYear1.ID),
		SpecialtyID:    new(Specialty2.ID),
		DisciplineID:   new(Discipline2.ID),
		SemesterNumber: new(2),
		Lectures:       new(20),
	}
	StudyPlansArray = domain.StudyPlans{StudyPlan1, StudyPlan2}

	WorkloadDistribution1 = &domain.WorkloadDistribution{
		ID:          1,
		StudyPlanID: StudyPlan1.ID,
		GroupID:     Group1.ID,
	}
	WorkloadDistribution2 = &domain.WorkloadDistribution{
		ID:          2,
		StudyPlanID: StudyPlan1.ID,
		GroupID:     Group2.ID,
	}
	WorkloadDistributionsArray = domain.WorkloadDistributions{WorkloadDistribution1, WorkloadDistribution2}

	WorkloadAssignment1 = &domain.WorkloadAssignment{
		ID:                     1,
		WorkloadDistributionID: WorkloadDistribution1.ID,
		TeacherID:              Teacher1.ID,
		RoleType:               new("lecturer"),
		AssignedHours:          new(60),
	}
	WorkloadAssignment2 = &domain.WorkloadAssignment{
		ID:                     2,
		WorkloadDistributionID: WorkloadDistribution1.ID,
		TeacherID:              Teacher2.ID,
		RoleType:               new("assistant"),
		AssignedHours:          new(30),
	}
	WorkloadAssignmentsArray = domain.WorkloadAssignments{WorkloadAssignment1, WorkloadAssignment2}

	ScheduleTemplateSetting1 = &domain.ScheduleTemplateSetting{
		ID:                         1,
		HoursPerLesson:             1.5,
		MaxIdenticalLessonsPerDay:  2,
		MaxStudyHoursPerDay:        8,
		MaxTeacherHoursPerWeek:     36,
		MaxGroupLessonHoursPerWeek: new(20),
	}
	ScheduleTemplateSetting2 = &domain.ScheduleTemplateSetting{
		ID:                        2,
		HoursPerLesson:            2.0,
		MaxIdenticalLessonsPerDay: 3,
		MaxStudyHoursPerDay:       10,
		MaxTeacherHoursPerWeek:    40,
	}
	ScheduleTemplateSettingsArray = domain.ScheduleTemplateSettings{ScheduleTemplateSetting1, ScheduleTemplateSetting2}

	ScheduleRestriction1 = &domain.ScheduleRestriction{
		ID:                      1,
		MinGroupLessonsPerDay:   2,
		MaxGroupLessonsPerDay:   4,
		MaxTeacherLessonsPerDay: 5,
		NoGapsInGroupSchedule:   true,
	}
	ScheduleRestriction2 = &domain.ScheduleRestriction{
		ID:                      2,
		MinGroupLessonsPerDay:   1,
		MaxGroupLessonsPerDay:   3,
		MaxTeacherLessonsPerDay: 4,
		NoGapsInGroupSchedule:   false,
	}
	ScheduleRestrictionsArray = domain.ScheduleRestrictions{ScheduleRestriction1, ScheduleRestriction2}

	ScheduleTemplate1 = &domain.ScheduleTemplate{
		ID:         1,
		SemesterID: new(Semester1.ID),
		Name:       new("Template 1"),
		IsActive:   new(false),
		IsDeleted:  new(false),
	}
	ScheduleTemplate2 = &domain.ScheduleTemplate{
		ID:         2,
		SemesterID: new(Semester1.ID),
		Name:       new("Template 2"),
		IsActive:   new(false),
		IsDeleted:  new(false),
	}
	ScheduleTemplatesArray = domain.ScheduleTemplates{ScheduleTemplate1, ScheduleTemplate2}

	UserCredential1 = &domain.UserCredential{
		UserID:       1,
		PasswordHash: "test_hash_1",
		UpdatedAt:    CurrentTime,
	}

	RefreshToken1 = &domain.RefreshToken{
		UserID:    2,
		TokenHash: "refresh_token_hash_1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsRevoked: false,
		CreatedAt: CurrentTime,
	}

	InviteReq1 = domain.InviteReq{
		Email:     "invited@test.com",
		FirstName: "Invited",
		LastName:  "User",
		Role:      domain.RoleTeacher,
	}
	DeanInviteReq = domain.InviteReq{
		Email:     "dean-invited@test.com",
		FirstName: "DeanInvited",
		LastName:  "User",
		Role:      domain.RoleTeacher,
	}
	InvitedUser1 = &domain.User{
		Email:     "invited@test.com",
		FirstName: "Invited",
		LastName:  "User",
		Role:      domain.RoleTeacher,
	}
	SetPasswordReq1 = &domain.SetPasswordReq{
		Password: "securepass123",
	}
	LoginReq1 = &domain.LoginReq{
		Email:    "invited@test.com",
		Password: "securepass123",
	}
	LoginReqWrongPassword = &domain.LoginReq{
		Email:    "invited@test.com",
		Password: "wrongpassword",
	}

	InviteToken1 = &domain.InviteToken{
		UserID:    1,
		Token:     "test_invite_token_1",
		ExpiresAt: time.Now().Add(72 * time.Hour),
		CreatedAt: CurrentTime,
	}

	ExpectedInviteResp1 = &domain.InviteResp{
		User: &domain.User{
			Email:     InviteReq1.Email,
			Role:      InviteReq1.Role,
			FirstName: "x",
			LastName:  "x",
		},
		InviteLink: "x",
	}
	ExpectedDeanInviteResp = &domain.InviteResp{
		User: &domain.User{
			Email:     DeanInviteReq.Email,
			Role:      DeanInviteReq.Role,
			FirstName: "x",
			LastName:  "x",
		},
		InviteLink: "x",
	}

	ExpectedResetInviteResp1 = &domain.InviteResp{
		User: &domain.User{
			Email:     User1.Email,
			Role:      User1.Role,
			FirstName: "x",
			LastName:  "x",
		},
		InviteLink: "x",
	}

	ExpectedTokenResp = map[string]any{"token": "x"}
)
