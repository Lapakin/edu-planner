package service

import (
	"context"

	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type Services struct {
	AcademicYearSvc            AcademicYearSvc
	SemesterSvc                SemesterSvc
	DepartmentSvc              DepartmentSvc
	RoomSvc                    RoomSvc
	SpecialtySvc               SpecialtySvc
	CycleCommitteeSvc          CycleCommitteeSvc
	TeacherSvc                 TeacherSvc
	DisciplineSvc              DisciplineSvc
	GroupSvc                   GroupSvc
	StudyPlanSvc               StudyPlanSvc
	WorkloadDistributionSvc    WorkloadDistributionSvc
	WorkloadAssignmentSvc      WorkloadAssignmentSvc
	ScheduleTemplateSettingSvc ScheduleTemplateSettingSvc
	BellScheduleSvc            BellScheduleSvc
	ScheduleTemplateSvc        ScheduleTemplateSvc
	ScheduleRestrictionSvc     ScheduleRestrictionSvc
	TeacherSlotPreferenceSvc   TeacherSlotPreferenceSvc
	CycleCommitteeLabRoomSvc   CycleCommitteeLabRoomSvc
	UserSvc                    UserSvc
}

type AcademicYearSvc interface {
	CreateAcademicYears(ctx context.Context, claims *jwt.Claims, academicYears domain.AcademicYears) error
	GetAcademicYearByID(ctx context.Context, id uint64) (*domain.AcademicYear, error)
	FetchAcademicYears(ctx context.Context, filters f.Filters) (domain.AcademicYears, error)
	UpdateAcademicYears(ctx context.Context, claims *jwt.Claims, academicYears domain.AcademicYears) error
	DeleteAcademicYears(ctx context.Context, claims *jwt.Claims, ids []uint64) error
	ActivateAcademicYear(ctx context.Context, claims *jwt.Claims, id uint64) error
	DeactivateAcademicYear(ctx context.Context, claims *jwt.Claims, id uint64) error
}

type SemesterSvc interface {
	CreateSemesters(ctx context.Context, claims *jwt.Claims, semesters domain.Semesters) error
	GetSemesterByID(ctx context.Context, id uint64) (*domain.Semester, error)
	FetchSemesters(ctx context.Context, filters f.Filters) (domain.Semesters, error)
	UpdateSemesters(ctx context.Context, claims *jwt.Claims, semesters domain.Semesters) error
	DeleteSemesters(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type DepartmentSvc interface {
	CreateDepartments(ctx context.Context, claims *jwt.Claims, departments domain.Departments) error
	GetDepartmentByID(ctx context.Context, id uint64) (*domain.Department, error)
	FetchDepartments(ctx context.Context, filters f.Filters) (domain.Departments, error)
	UpdateDepartments(ctx context.Context, claims *jwt.Claims, departments domain.Departments) error
	DeleteDepartments(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type RoomSvc interface {
	CreateRooms(ctx context.Context, claims *jwt.Claims, rooms domain.Rooms) error
	GetRoomByID(ctx context.Context, id uint64) (*domain.Room, error)
	FetchRooms(ctx context.Context, filters f.Filters) (domain.Rooms, error)
	UpdateRooms(ctx context.Context, claims *jwt.Claims, rooms domain.Rooms) error
	DeleteRooms(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type SpecialtySvc interface {
	CreateSpecialties(ctx context.Context, claims *jwt.Claims, specialities domain.Specialties) error
	GetSpecialtyByID(ctx context.Context, id uint64) (*domain.Specialty, error)
	FetchSpecialties(ctx context.Context, filters f.Filters) (domain.Specialties, error)
	UpdateSpecialties(ctx context.Context, claims *jwt.Claims, specialities domain.Specialties) error
	DeleteSpecialties(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type CycleCommitteeSvc interface {
	CreateCycleCommittees(ctx context.Context, claims *jwt.Claims, cycleCommittees domain.CycleCommittees) error
	GetCycleCommitteeByID(ctx context.Context, id uint64) (*domain.CycleCommittee, error)
	FetchCycleCommittees(ctx context.Context, filters f.Filters) (domain.CycleCommittees, error)
	UpdateCycleCommittees(ctx context.Context, claims *jwt.Claims, cycleCommittees domain.CycleCommittees) error
	DeleteCycleCommittees(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type TeacherSvc interface {
	CreateTeachers(ctx context.Context, claims *jwt.Claims, teachers domain.Teachers) error
	GetTeacherByID(ctx context.Context, id uint64) (*domain.Teacher, error)
	FetchTeachers(ctx context.Context, filters f.Filters) (domain.Teachers, error)
	UpdateTeachers(ctx context.Context, claims *jwt.Claims, teachers domain.Teachers) error
	DeleteTeachers(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type DisciplineSvc interface {
	CreateDisciplines(ctx context.Context, claims *jwt.Claims, disciplines domain.Disciplines) error
	GetDisciplineByID(ctx context.Context, id uint64) (*domain.Discipline, error)
	FetchDisciplines(ctx context.Context, filters f.Filters) (domain.Disciplines, error)
	UpdateDisciplines(ctx context.Context, claims *jwt.Claims, disciplines domain.Disciplines) error
	DeleteDisciplines(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type GroupSvc interface {
	CreateGroups(ctx context.Context, claims *jwt.Claims, groups domain.Groups) error
	GetGroupByID(ctx context.Context, id uint64) (*domain.Group, error)
	FetchGroups(ctx context.Context, filters f.Filters) (domain.Groups, error)
	UpdateGroups(ctx context.Context, claims *jwt.Claims, groups domain.Groups) error
	DeleteGroups(ctx context.Context, claims *jwt.Claims, ids []uint64) error

	CreateGroupSemesters(ctx context.Context, claims *jwt.Claims, groupSemesters domain.GroupSemesters) error
	GetGroupSemesterByID(ctx context.Context, id uint64) (*domain.GroupSemester, error)
	FetchGroupSemesters(ctx context.Context, filters f.Filters) (domain.GroupSemesters, error)
	UpdateGroupSemesters(ctx context.Context, claims *jwt.Claims, groupSemesters domain.GroupSemesters) error
	DeleteGroupSemesters(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type StudyPlanSvc interface {
	CreateStudyPlans(ctx context.Context, claims *jwt.Claims, studyPlans domain.StudyPlans) error
	GetStudyPlanByID(ctx context.Context, id uint64) (*domain.StudyPlan, error)
	FetchStudyPlans(ctx context.Context, filters f.Filters) (domain.StudyPlans, error)
	UpdateStudyPlans(ctx context.Context, claims *jwt.Claims, studyPlans domain.StudyPlans) error
	DeleteStudyPlans(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type WorkloadDistributionSvc interface {
	CreateWorkloadDistributions(ctx context.Context, claims *jwt.Claims, distributions domain.WorkloadDistributions) error
	GetWorkloadDistributionByID(ctx context.Context, id uint64) (*domain.WorkloadDistribution, error)
	FetchWorkloadDistributions(ctx context.Context, filters f.Filters) (domain.WorkloadDistributions, error)
	UpdateWorkloadDistributions(ctx context.Context, claims *jwt.Claims, distributions domain.WorkloadDistributions) error
	DeleteWorkloadDistributions(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type WorkloadAssignmentSvc interface {
	CreateWorkloadAssignments(ctx context.Context, claims *jwt.Claims, assignments domain.WorkloadAssignments) error
	GetWorkloadAssignmentByID(ctx context.Context, id uint64) (*domain.WorkloadAssignment, error)
	FetchWorkloadAssignments(ctx context.Context, filters f.Filters) (domain.WorkloadAssignments, error)
	UpdateWorkloadAssignments(ctx context.Context, claims *jwt.Claims, assignments domain.WorkloadAssignments) error
	DeleteWorkloadAssignments(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type ScheduleTemplateSettingSvc interface {
	CreateScheduleTemplateSettings(ctx context.Context, claims *jwt.Claims, settings domain.ScheduleTemplateSettings) error
	GetScheduleTemplateSettingByID(ctx context.Context, id uint64) (*domain.ScheduleTemplateSetting, error)
	FetchScheduleTemplateSettings(ctx context.Context, filters f.Filters) (domain.ScheduleTemplateSettings, error)
	UpdateScheduleTemplateSettings(ctx context.Context, claims *jwt.Claims, settings domain.ScheduleTemplateSettings) error
	DeleteScheduleTemplateSettings(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type BellScheduleSvc interface {
	CreateBellSchedules(ctx context.Context, claims *jwt.Claims, schedules domain.BellSchedules) error
	GetBellScheduleByID(ctx context.Context, id uint64) (*domain.BellSchedule, error)
	FetchBellSchedules(ctx context.Context, filters f.Filters) (domain.BellSchedules, error)
	UpdateBellSchedules(ctx context.Context, claims *jwt.Claims, schedules domain.BellSchedules) error
	DeleteBellSchedules(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type ScheduleTemplateSvc interface {
	GenerateScheduleTemplate(ctx context.Context, req *domain.GenerateScheduleRequest) (*domain.ScheduleData, error)
	SaveScheduleTemplate(ctx context.Context, claims *jwt.Claims, req *domain.SaveScheduleTemplateRequest) (*domain.ScheduleTemplate, error)
	GetScheduleTemplateByID(ctx context.Context, id uint64) (*domain.ScheduleTemplate, error)
	FetchScheduleTemplates(ctx context.Context, filters f.Filters) (domain.ScheduleTemplates, error)
	DeleteScheduleTemplates(ctx context.Context, claims *jwt.Claims, ids []uint64) error
	ActivateScheduleTemplate(ctx context.Context, claims *jwt.Claims, id uint64) error
}

type ScheduleRestrictionSvc interface {
	CreateScheduleRestrictions(ctx context.Context, claims *jwt.Claims, restrictions domain.ScheduleRestrictions) error
	GetScheduleRestrictionByID(ctx context.Context, id uint64) (*domain.ScheduleRestriction, error)
	FetchScheduleRestrictions(ctx context.Context, filters f.Filters) (domain.ScheduleRestrictions, error)
	UpdateScheduleRestrictions(ctx context.Context, claims *jwt.Claims, restrictions domain.ScheduleRestrictions) error
	DeleteScheduleRestrictions(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type TeacherSlotPreferenceSvc interface {
	CreateTeacherSlotPreferences(ctx context.Context, claims *jwt.Claims, preferences domain.TeacherSlotPreferences) error
	GetTeacherSlotPreferenceByID(ctx context.Context, id uint64) (*domain.TeacherSlotPreference, error)
	FetchTeacherSlotPreferences(ctx context.Context, filters f.Filters) (domain.TeacherSlotPreferences, error)
	UpdateTeacherSlotPreferences(ctx context.Context, claims *jwt.Claims, preferences domain.TeacherSlotPreferences) error
	DeleteTeacherSlotPreferences(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type CycleCommitteeLabRoomSvc interface {
	CreateCycleCommitteeLabRooms(ctx context.Context, claims *jwt.Claims, labRooms domain.CycleCommitteeLabRooms) error
	GetCycleCommitteeLabRoomByID(ctx context.Context, id uint64) (*domain.CycleCommitteeLabRoom, error)
	FetchCycleCommitteeLabRooms(ctx context.Context, filters f.Filters) (domain.CycleCommitteeLabRooms, error)
	UpdateCycleCommitteeLabRooms(ctx context.Context, claims *jwt.Claims, labRooms domain.CycleCommitteeLabRooms) error
	DeleteCycleCommitteeLabRooms(ctx context.Context, claims *jwt.Claims, ids []uint64) error
}

type UserSvc interface {
	MassageConsumer(ctx context.Context, massage *domain.Massage) error
}
