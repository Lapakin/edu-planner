package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type RepoManager interface {
	NewAcademicYearRepo(db sqlx.ExtContext) AcademicYearRepository
	NewSemesterRepo(db sqlx.ExtContext) SemesterRepository
	NewDepartmentRepo(db sqlx.ExtContext) DepartmentRepository
	NewRoomRepo(db sqlx.ExtContext) RoomRepository
	NewSpecialtyRepo(db sqlx.ExtContext) SpecialtyRepository
	NewCycleCommitteeRepo(db sqlx.ExtContext) CycleCommitteeRepository
	NewTeacherRepo(db sqlx.ExtContext) TeacherRepository
	NewDisciplineRepo(db sqlx.ExtContext) DisciplineRepository
	NewGroupRepo(db sqlx.ExtContext) GroupRepository
	NewMessageRepo(db sqlx.ExtContext) MessageRepository
	NewStudyPlanRepo(db sqlx.ExtContext) StudyPlanRepository
	NewWorkloadDistributionRepo(db sqlx.ExtContext) WorkloadDistributionRepository
	NewWorkloadAssignmentRepo(db sqlx.ExtContext) WorkloadAssignmentRepository
	NewScheduleTemplateSettingRepo(db sqlx.ExtContext) ScheduleTemplateSettingRepository
	NewBellScheduleRepo(db sqlx.ExtContext) BellScheduleRepository
	NewScheduleTemplateRepo(db sqlx.ExtContext) ScheduleTemplateRepository
	NewScheduleRestrictionRepo(db sqlx.ExtContext) ScheduleRestrictionRepository
	NewUserRepo(db sqlx.ExtContext) UserRepository
}

type AcademicYearRepository interface {
	CreateAcademicYear(ctx context.Context, academicYears domain.AcademicYears) error
	GetAcademicYearByID(ctx context.Context, id uint64) (*domain.AcademicYear, error)
	FetchAcademicYears(ctx context.Context, filters f.Filters) (domain.AcademicYears, error)
	UpdateAcademicYears(ctx context.Context, academicYears domain.AcademicYears) error
	DeleteAcademicYears(ctx context.Context, ids []uint64) error
	ActivateAcademicYear(ctx context.Context, id uint64, currentTime time.Time) error
	DeactivateAcademicYear(ctx context.Context, id uint64, currentTime time.Time) error
}

type SemesterRepository interface {
	CreateSemesters(ctx context.Context, semesters domain.Semesters) error
	GetSemesterByID(ctx context.Context, id uint64) (*domain.Semester, error)
	FetchSemesters(ctx context.Context, filters f.Filters) (domain.Semesters, error)
	UpdateSemesters(ctx context.Context, semesters domain.Semesters) error
	DeleteSemesters(ctx context.Context, ids []uint64) error
}

type DepartmentRepository interface {
	CreateDepartments(ctx context.Context, departments domain.Departments) error
	GetDepartmentByID(ctx context.Context, id uint64) (*domain.Department, error)
	FetchDepartments(ctx context.Context, filters f.Filters) (domain.Departments, error)
	UpdateDepartments(ctx context.Context, departments domain.Departments) error
	DeleteDepartments(ctx context.Context, ids []uint64) error
	AttachDepartmentsToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Departments) error
}

type RoomRepository interface {
	CreateRooms(ctx context.Context, rooms domain.Rooms) error
	GetRoomByID(ctx context.Context, id uint64) (*domain.Room, error)
	FetchRooms(ctx context.Context, filters f.Filters) (domain.Rooms, error)
	UpdateRooms(ctx context.Context, rooms domain.Rooms) error
	DeleteRooms(ctx context.Context, ids []uint64) error
	AttachRoomsToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Rooms) error
}

type SpecialtyRepository interface {
	CreateSpecialties(ctx context.Context, specialities domain.Specialties) error
	GetSpecialtyByID(ctx context.Context, id uint64) (*domain.Specialty, error)
	FetchSpecialties(ctx context.Context, filters f.Filters) (domain.Specialties, error)
	UpdateSpecialties(ctx context.Context, specialities domain.Specialties) error
	DeleteSpecialties(ctx context.Context, ids []uint64) error
	AttachSpecialtiesToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Specialties) error
}

type CycleCommitteeRepository interface {
	CreateCycleCommittees(ctx context.Context, cycleCommittees domain.CycleCommittees) error
	GetCycleCommitteeByID(ctx context.Context, id uint64) (*domain.CycleCommittee, error)
	FetchCycleCommittees(ctx context.Context, filters f.Filters) (domain.CycleCommittees, error)
	UpdateCycleCommittees(ctx context.Context, cycleCommittees domain.CycleCommittees) error
	DeleteCycleCommittees(ctx context.Context, ids []uint64) error
	AttachCycleCommitteesToAcademicYear(ctx context.Context, toAttach map[uint64]domain.CycleCommittees) error
}

type TeacherRepository interface {
	CreateTeachers(ctx context.Context, teachers domain.Teachers) error
	GetTeacherByID(ctx context.Context, id uint64) (*domain.Teacher, error)
	FetchTeachers(ctx context.Context, filters f.Filters) (domain.Teachers, error)
	UpdateTeachers(ctx context.Context, teachers domain.Teachers) error
	DeleteTeachers(ctx context.Context, ids []uint64) error
	AttachTeachersToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Teachers) error
}

type DisciplineRepository interface {
	CreateDisciplines(ctx context.Context, disciplines domain.Disciplines) error
	GetDisciplineByID(ctx context.Context, id uint64) (*domain.Discipline, error)
	FetchDisciplines(ctx context.Context, filters f.Filters) (domain.Disciplines, error)
	UpdateDisciplines(ctx context.Context, disciplines domain.Disciplines) error
	DeleteDisciplines(ctx context.Context, ids []uint64) error
	AttachDisciplinesToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Disciplines) error
}

type GroupRepository interface {
	CreateGroups(ctx context.Context, groups domain.Groups) error
	GetGroupByID(ctx context.Context, id uint64) (*domain.Group, error)
	FetchGroups(ctx context.Context, filters f.Filters) (domain.Groups, error)
	UpdateGroups(ctx context.Context, groups domain.Groups) error
	DeleteGroups(ctx context.Context, ids []uint64) error
	AttachGroupsToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Groups) error

	CreateGroupSemesters(ctx context.Context, groupSemesters domain.GroupSemesters) error
	GetGroupSemesterByID(ctx context.Context, id uint64) (*domain.GroupSemester, error)
	FetchGroupSemesters(ctx context.Context, filters f.Filters) (domain.GroupSemesters, error)
	UpdateGroupSemesters(ctx context.Context, groupSemesters domain.GroupSemesters) error
	DeleteGroupSemesters(ctx context.Context, ids []uint64) error
}

type MessageRepository interface {
	Write(ctx context.Context, massages domain.Massages) error
}

type StudyPlanRepository interface {
	CreateStudyPlans(ctx context.Context, studyPlans domain.StudyPlans) error
	GetStudyPlanByID(ctx context.Context, id uint64) (*domain.StudyPlan, error)
	FetchStudyPlans(ctx context.Context, filters f.Filters) (domain.StudyPlans, error)
	UpdateStudyPlans(ctx context.Context, studyPlans domain.StudyPlans) error
	DeleteStudyPlans(ctx context.Context, ids []uint64) error
}

type WorkloadDistributionRepository interface {
	CreateWorkloadDistributions(ctx context.Context, distributions domain.WorkloadDistributions) error
	GetWorkloadDistributionByID(ctx context.Context, id uint64) (*domain.WorkloadDistribution, error)
	FetchWorkloadDistributions(ctx context.Context, filters f.Filters) (domain.WorkloadDistributions, error)
	UpdateWorkloadDistributions(ctx context.Context, distributions domain.WorkloadDistributions) error
	DeleteWorkloadDistributions(ctx context.Context, ids []uint64) error
}

type WorkloadAssignmentRepository interface {
	CreateWorkloadAssignments(ctx context.Context, assignments domain.WorkloadAssignments) error
	GetWorkloadAssignmentByID(ctx context.Context, id uint64) (*domain.WorkloadAssignment, error)
	FetchWorkloadAssignments(ctx context.Context, filters f.Filters) (domain.WorkloadAssignments, error)
	UpdateWorkloadAssignments(ctx context.Context, assignments domain.WorkloadAssignments) error
	DeleteWorkloadAssignments(ctx context.Context, ids []uint64) error
}

type ScheduleTemplateSettingRepository interface {
	CreateScheduleTemplateSettings(ctx context.Context, settings domain.ScheduleTemplateSettings) error
	GetScheduleTemplateSettingByID(ctx context.Context, id uint64) (*domain.ScheduleTemplateSetting, error)
	FetchScheduleTemplateSettings(ctx context.Context, filters f.Filters) (domain.ScheduleTemplateSettings, error)
	UpdateScheduleTemplateSettings(ctx context.Context, settings domain.ScheduleTemplateSettings) error
	DeleteScheduleTemplateSettings(ctx context.Context, ids []uint64) error
}

type BellScheduleRepository interface {
	CreateBellSchedules(ctx context.Context, schedules domain.BellSchedules) error
	GetBellScheduleByID(ctx context.Context, id uint64) (*domain.BellSchedule, error)
	FetchBellSchedules(ctx context.Context, filters f.Filters) (domain.BellSchedules, error)
	UpdateBellSchedules(ctx context.Context, schedules domain.BellSchedules) error
	DeleteBellSchedules(ctx context.Context, ids []uint64) error
}

type ScheduleTemplateRepository interface {
	CreateScheduleTemplate(ctx context.Context, template *domain.ScheduleTemplate) error
	GetScheduleTemplateByID(ctx context.Context, id uint64) (*domain.ScheduleTemplate, error)
	FetchScheduleTemplates(ctx context.Context, filters f.Filters) (domain.ScheduleTemplates, error)
	UpdateScheduleTemplate(ctx context.Context, template *domain.ScheduleTemplate) error
	DeleteScheduleTemplates(ctx context.Context, ids []uint64) error
	ActivateScheduleTemplate(ctx context.Context, id uint64, t time.Time) error
}

type ScheduleRestrictionRepository interface {
	CreateScheduleRestrictions(ctx context.Context, restrictions domain.ScheduleRestrictions) error
	GetScheduleRestrictionByID(ctx context.Context, id uint64) (*domain.ScheduleRestriction, error)
	FetchScheduleRestrictions(ctx context.Context, filters f.Filters) (domain.ScheduleRestrictions, error)
	UpdateScheduleRestrictions(ctx context.Context, restrictions domain.ScheduleRestrictions) error
	DeleteScheduleRestrictions(ctx context.Context, ids []uint64) error
}

type UserRepository interface {
	UpsertUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint64) error
}
