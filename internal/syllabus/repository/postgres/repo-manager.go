package postgres

import (
	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/syllabus/repository"
)

type RepoManager struct{}

func NewRepoManager() *RepoManager {
	return &RepoManager{}
}

func (rm *RepoManager) NewAcademicYearRepo(db sqlx.ExtContext) repository.AcademicYearRepository {
	return NewAcademicYearRepository(db)
}

func (rm *RepoManager) NewSemesterRepo(db sqlx.ExtContext) repository.SemesterRepository {
	return NewSemesterRepository(db)
}

func (rm *RepoManager) NewDepartmentRepo(db sqlx.ExtContext) repository.DepartmentRepository {
	return NewDepartmentRepository(db)
}

func (rm *RepoManager) NewRoomRepo(db sqlx.ExtContext) repository.RoomRepository {
	return NewRoomRepository(db)
}

func (rm *RepoManager) NewSpecialtyRepo(db sqlx.ExtContext) repository.SpecialtyRepository {
	return NewSpecialtyRepository(db)
}

func (rm *RepoManager) NewCycleCommitteeRepo(db sqlx.ExtContext) repository.CycleCommitteeRepository {
	return NewCycleCommitteeRepository(db)
}

func (rm *RepoManager) NewTeacherRepo(db sqlx.ExtContext) repository.TeacherRepository {
	return NewTeacherRepository(db)
}

func (rm *RepoManager) NewDisciplineRepo(db sqlx.ExtContext) repository.DisciplineRepository {
	return NewDisciplineRepository(db)
}

func (rm *RepoManager) NewGroupRepo(db sqlx.ExtContext) repository.GroupRepository {
	return NewGroupRepository(db)
}

func (rm *RepoManager) NewMessageRepo(db sqlx.ExtContext) repository.MessageRepository {
	return NewMessageWriter(db)
}

func (rm *RepoManager) NewStudyPlanRepo(db sqlx.ExtContext) repository.StudyPlanRepository {
	return NewStudyPlanRepository(db)
}

func (rm *RepoManager) NewWorkloadDistributionRepo(db sqlx.ExtContext) repository.WorkloadDistributionRepository {
	return NewWorkloadDistributionRepository(db)
}

func (rm *RepoManager) NewWorkloadAssignmentRepo(db sqlx.ExtContext) repository.WorkloadAssignmentRepository {
	return NewWorkloadAssignmentRepository(db)
}

func (rm *RepoManager) NewScheduleTemplateSettingRepo(db sqlx.ExtContext) repository.ScheduleTemplateSettingRepository {
	return NewScheduleTemplateSettingRepository(db)
}

func (rm *RepoManager) NewBellScheduleRepo(db sqlx.ExtContext) repository.BellScheduleRepository {
	return NewBellScheduleRepository(db)
}

func (rm *RepoManager) NewScheduleTemplateRepo(db sqlx.ExtContext) repository.ScheduleTemplateRepository {
	return NewScheduleTemplateRepository(db)
}

func (rm *RepoManager) NewScheduleRestrictionRepo(db sqlx.ExtContext) repository.ScheduleRestrictionRepository {
	return NewScheduleRestrictionRepository(db)
}

func (rm *RepoManager) NewTeacherSlotPreferenceRepo(db sqlx.ExtContext) repository.TeacherSlotPreferenceRepository {
	return NewTeacherSlotPreferenceRepository(db)
}

func (rm *RepoManager) NewCycleCommitteeLabRoomRepo(db sqlx.ExtContext) repository.CycleCommitteeLabRoomRepository {
	return NewCycleCommitteeLabRoomRepository(db)
}

func (rm *RepoManager) NewUserRepo(db sqlx.ExtContext) repository.UserRepository {
	return NewUserRepository(db)
}
