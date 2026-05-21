package service

import (
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
)

func NewServices(db *pg.DB, rm repository.RepoManager, l *logging.Logger, cfg domain.GenerationConfig) *Services {
	return &Services{
		AcademicYearSvc:            NewAcademicYearService(db, rm, l.WithField("service", "AcademicYear")),
		SemesterSvc:                NewSemesterService(db, rm, l.WithField("service", "Semester")),
		DepartmentSvc:              NewDepartmentService(db, rm, l.WithField("service", "Department")),
		RoomSvc:                    NewRoomService(db, rm, l.WithField("service", "Room")),
		SpecialtySvc:               NewSpecialtyService(db, rm, l.WithField("service", "Speciality")),
		CycleCommitteeSvc:          NewCycleCommitteeService(db, rm, l.WithField("service", "CycleCommittee")),
		TeacherSvc:                 NewTeacherService(db, rm, l.WithField("service", "Teacher")),
		DisciplineSvc:              NewDisciplineService(db, rm, l.WithField("service", "Discipline")),
		GroupSvc:                   NewGroupService(db, rm, l.WithField("service", "Group")),
		StudyPlanSvc:               NewStudyPlanService(db, rm, l.WithField("service", "StudyPlan")),
		WorkloadDistributionSvc:    NewWorkloadDistributionService(db, rm, l.WithField("service", "WorkloadDistribution")),
		WorkloadAssignmentSvc:      NewWorkloadAssignmentService(db, rm, l.WithField("service", "WorkloadAssignment")),
		ScheduleTemplateSettingSvc: NewScheduleTemplateSettingService(db, rm, l.WithField("service", "ScheduleTemplateSetting")),
		BellScheduleSvc:            NewBellScheduleService(db, rm, l.WithField("service", "BellSchedule")),
		ScheduleTemplateSvc:        NewScheduleTemplateService(db, rm, l.WithField("service", "ScheduleTemplate"), cfg),
		ScheduleRestrictionSvc:     NewScheduleRestrictionService(db, rm, l.WithField("service", "ScheduleRestriction")),
		TeacherSlotPreferenceSvc:   NewTeacherSlotPreferenceService(db, rm, l.WithField("service", "TeacherSlotPreference")),
		CycleCommitteeLabRoomSvc:   NewCycleCommitteeLabRoomService(db, rm, l.WithField("service", "CycleCommitteeLabRoom")),
		UserSvc:                    NewUserService(db, rm, l.WithField("service", "UserInfo")),
	}
}
