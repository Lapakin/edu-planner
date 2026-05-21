package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Lapakin/edu-planner/internal/middleware"
	"github.com/Lapakin/edu-planner/internal/syllabus/handler"
	"github.com/Lapakin/edu-planner/internal/syllabus/service"
)

func NewRouter(svc *service.Services) *gin.Engine {
	r := gin.New()

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.JWTMiddleware())
	apiV1.Handle(http.MethodGet, "/academic-years/:academicYearId", handler.NewGetAcademicYearByIDHandler(svc.AcademicYearSvc))
	apiV1.Handle(http.MethodPost, "/academic-years", handler.NewCreateAcademicYearsHandler(svc.AcademicYearSvc))
	apiV1.Handle(http.MethodGet, "/academic-years", handler.NewFetchAcademicYearsHandler(svc.AcademicYearSvc))
	apiV1.Handle(http.MethodPut, "/academic-years", handler.NewUpdateAcademicYearsHandler(svc.AcademicYearSvc))
	apiV1.Handle(http.MethodPost, "/academic-years/delete", handler.NewDeleteAcademicYearByIDsHandler(svc.AcademicYearSvc))

	apiV1.Handle(http.MethodPost, "/academic-years/:academicYearId/activate", handler.NewActivateAcademicYearByIDHandler(svc.AcademicYearSvc))
	apiV1.Handle(http.MethodPost, "/academic-years/:academicYearId/deactivate", handler.NewDeactivateAcademicYearByIDHandler(svc.AcademicYearSvc))

	apiV1.Handle(http.MethodGet, "/semesters/:semesterId", handler.NewGetSemesterByIDHandler(svc.SemesterSvc))
	apiV1.Handle(http.MethodPost, "/semesters", handler.NewCreateSemestersHandler(svc.SemesterSvc))
	apiV1.Handle(http.MethodGet, "/semesters", handler.NewFetchSemestersHandler(svc.SemesterSvc))
	apiV1.Handle(http.MethodPut, "/semesters", handler.NewUpdateSemestersHandler(svc.SemesterSvc))
	apiV1.Handle(http.MethodPost, "/semesters/delete", handler.NewDeleteSemesterByIDsHandler(svc.SemesterSvc))

	apiV1.Handle(http.MethodGet, "/bell-schedules/:bellScheduleId", handler.NewGetBellScheduleByIDHandler(svc.BellScheduleSvc))
	apiV1.Handle(http.MethodPost, "/bell-schedules", handler.NewCreateBellSchedulesHandler(svc.BellScheduleSvc))
	apiV1.Handle(http.MethodGet, "/bell-schedules", handler.NewFetchBellSchedulesHandler(svc.BellScheduleSvc))
	apiV1.Handle(http.MethodPut, "/bell-schedules", handler.NewUpdateBellSchedulesHandler(svc.BellScheduleSvc))
	apiV1.Handle(http.MethodPost, "/bell-schedules/delete", handler.NewDeleteBellScheduleByIDsHandler(svc.BellScheduleSvc))

	apiV1.Handle(http.MethodGet, "/departments/:departmentId", handler.NewGetDepartmentByIDHandler(svc.DepartmentSvc))
	apiV1.Handle(http.MethodPost, "/departments", handler.NewCreateDepartmentsHandler(svc.DepartmentSvc))
	apiV1.Handle(http.MethodGet, "/departments", handler.NewFetchDepartmentsHandler(svc.DepartmentSvc))
	apiV1.Handle(http.MethodPut, "/departments", handler.NewUpdateDepartmentsHandler(svc.DepartmentSvc))
	apiV1.Handle(http.MethodPost, "/departments/delete", handler.NewDeleteDepartmentByIDsHandler(svc.DepartmentSvc))

	apiV1.Handle(http.MethodGet, "/rooms/:roomId", handler.NewGetRoomByIDHandler(svc.RoomSvc))
	apiV1.Handle(http.MethodPost, "/rooms", handler.NewCreateRoomsHandler(svc.RoomSvc))
	apiV1.Handle(http.MethodGet, "/rooms", handler.NewFetchRoomsHandler(svc.RoomSvc))
	apiV1.Handle(http.MethodPut, "/rooms", handler.NewUpdateRoomsHandler(svc.RoomSvc))
	apiV1.Handle(http.MethodPost, "/rooms/delete", handler.NewDeleteRoomByIDsHandler(svc.RoomSvc))

	apiV1.Handle(http.MethodGet, "/specialties/:specialtyId", handler.NewGetSpecialtyByIDHandler(svc.SpecialtySvc))
	apiV1.Handle(http.MethodPost, "/specialties", handler.NewCreateSpecialtiesHandler(svc.SpecialtySvc))
	apiV1.Handle(http.MethodGet, "/specialties", handler.NewFetchSpecialtiesHandler(svc.SpecialtySvc))
	apiV1.Handle(http.MethodPut, "/specialties", handler.NewUpdateSpecialtiesHandler(svc.SpecialtySvc))
	apiV1.Handle(http.MethodPost, "/specialties/delete", handler.NewDeleteSpecialtyByIDsHandler(svc.SpecialtySvc))

	apiV1.Handle(http.MethodGet, "/cycle-committees/:cycleCommitteeId", handler.NewGetCycleCommitteeByIDHandler(svc.CycleCommitteeSvc))
	apiV1.Handle(http.MethodPost, "/cycle-committees", handler.NewCreateCycleCommitteesHandler(svc.CycleCommitteeSvc))
	apiV1.Handle(http.MethodGet, "/cycle-committees", handler.NewFetchCycleCommitteesHandler(svc.CycleCommitteeSvc))
	apiV1.Handle(http.MethodPut, "/cycle-committees", handler.NewUpdateCycleCommitteesHandler(svc.CycleCommitteeSvc))
	apiV1.Handle(http.MethodPost, "/cycle-committees/delete", handler.NewDeleteCycleCommitteeByIDsHandler(svc.CycleCommitteeSvc))

	apiV1.Handle(http.MethodGet, "/teachers/:teacherId", handler.NewGetTeacherByIDHandler(svc.TeacherSvc))
	apiV1.Handle(http.MethodPost, "/teachers", handler.NewCreateTeachersHandler(svc.TeacherSvc))
	apiV1.Handle(http.MethodGet, "/teachers", handler.NewFetchTeachersHandler(svc.TeacherSvc))
	apiV1.Handle(http.MethodPut, "/teachers", handler.NewUpdateTeachersHandler(svc.TeacherSvc))
	apiV1.Handle(http.MethodPost, "/teachers/delete", handler.NewDeleteTeacherByIDsHandler(svc.TeacherSvc))

	apiV1.Handle(http.MethodGet, "/disciplines/:disciplineId", handler.NewGetDisciplineByIDHandler(svc.DisciplineSvc))
	apiV1.Handle(http.MethodPost, "/disciplines", handler.NewCreateDisciplinesHandler(svc.DisciplineSvc))
	apiV1.Handle(http.MethodGet, "/disciplines", handler.NewFetchDisciplinesHandler(svc.DisciplineSvc))
	apiV1.Handle(http.MethodPut, "/disciplines", handler.NewUpdateDisciplinesHandler(svc.DisciplineSvc))
	apiV1.Handle(http.MethodPost, "/disciplines/delete", handler.NewDeleteDisciplineByIDsHandler(svc.DisciplineSvc))

	apiV1.Handle(http.MethodGet, "/groups/:groupId", handler.NewGetGroupByIDHandler(svc.GroupSvc))
	apiV1.Handle(http.MethodPost, "/groups", handler.NewCreateGroupsHandler(svc.GroupSvc))
	apiV1.Handle(http.MethodGet, "/groups", handler.NewFetchGroupsHandler(svc.GroupSvc))
	apiV1.Handle(http.MethodPut, "/groups", handler.NewUpdateGroupsHandler(svc.GroupSvc))
	apiV1.Handle(http.MethodPost, "/groups/delete", handler.NewDeleteGroupByIDsHandler(svc.GroupSvc))

	apiV1.Handle(http.MethodGet, "/groups/semesters/:semesterId", handler.NewGetGroupSemesterByIDHandler(svc.GroupSvc))
	apiV1.Handle(http.MethodPost, "/groups/semesters", handler.NewCreateGroupSemestersHandler(svc.GroupSvc))
	apiV1.Handle(http.MethodGet, "/groups/semesters", handler.NewFetchGroupSemestersHandler(svc.GroupSvc))
	apiV1.Handle(http.MethodPut, "/groups/semesters", handler.NewUpdateGroupSemestersHandler(svc.GroupSvc))
	apiV1.Handle(http.MethodPost, "/groups/semesters/delete", handler.NewDeleteGroupSemesterByIDsHandler(svc.GroupSvc))

	apiV1.Handle(http.MethodGet, "/study-plans/:studyPlanId", handler.NewGetStudyPlanByIDHandler(svc.StudyPlanSvc))
	apiV1.Handle(http.MethodPost, "/study-plans", handler.NewCreateStudyPlansHandler(svc.StudyPlanSvc))
	apiV1.Handle(http.MethodGet, "/study-plans", handler.NewFetchStudyPlansHandler(svc.StudyPlanSvc))
	apiV1.Handle(http.MethodPut, "/study-plans", handler.NewUpdateStudyPlansHandler(svc.StudyPlanSvc))
	apiV1.Handle(http.MethodPost, "/study-plans/delete", handler.NewDeleteStudyPlanByIDsHandler(svc.StudyPlanSvc))

	apiV1.Handle(http.MethodGet, "/workload-distributions/:workloadDistributionId", handler.NewGetWorkloadDistributionByIDHandler(svc.WorkloadDistributionSvc))
	apiV1.Handle(http.MethodPost, "/workload-distributions", handler.NewCreateWorkloadDistributionsHandler(svc.WorkloadDistributionSvc))
	apiV1.Handle(http.MethodGet, "/workload-distributions", handler.NewFetchWorkloadDistributionsHandler(svc.WorkloadDistributionSvc))
	apiV1.Handle(http.MethodPut, "/workload-distributions", handler.NewUpdateWorkloadDistributionsHandler(svc.WorkloadDistributionSvc))
	apiV1.Handle(http.MethodPost, "/workload-distributions/delete", handler.NewDeleteWorkloadDistributionByIDsHandler(svc.WorkloadDistributionSvc))

	apiV1.Handle(http.MethodGet, "/workload-assignments/:workloadAssignmentId", handler.NewGetWorkloadAssignmentByIDHandler(svc.WorkloadAssignmentSvc))
	apiV1.Handle(http.MethodPost, "/workload-assignments", handler.NewCreateWorkloadAssignmentsHandler(svc.WorkloadAssignmentSvc))
	apiV1.Handle(http.MethodGet, "/workload-assignments", handler.NewFetchWorkloadAssignmentsHandler(svc.WorkloadAssignmentSvc))
	apiV1.Handle(http.MethodPut, "/workload-assignments", handler.NewUpdateWorkloadAssignmentsHandler(svc.WorkloadAssignmentSvc))
	apiV1.Handle(http.MethodPost, "/workload-assignments/delete", handler.NewDeleteWorkloadAssignmentByIDsHandler(svc.WorkloadAssignmentSvc))

	apiV1.Handle(http.MethodGet, "/schedule-template-settings/:scheduleTemplateSettingId", handler.NewGetScheduleTemplateSettingByIDHandler(svc.ScheduleTemplateSettingSvc))
	apiV1.Handle(http.MethodPost, "/schedule-template-settings", handler.NewCreateScheduleTemplateSettingsHandler(svc.ScheduleTemplateSettingSvc))
	apiV1.Handle(http.MethodGet, "/schedule-template-settings", handler.NewFetchScheduleTemplateSettingsHandler(svc.ScheduleTemplateSettingSvc))
	apiV1.Handle(http.MethodPut, "/schedule-template-settings", handler.NewUpdateScheduleTemplateSettingsHandler(svc.ScheduleTemplateSettingSvc))
	apiV1.Handle(http.MethodPost, "/schedule-template-settings/delete", handler.NewDeleteScheduleTemplateSettingByIDsHandler(svc.ScheduleTemplateSettingSvc))

	apiV1.Handle(http.MethodGet, "/schedule-restrictions/:scheduleRestrictionId", handler.NewGetScheduleRestrictionByIDHandler(svc.ScheduleRestrictionSvc))
	apiV1.Handle(http.MethodPost, "/schedule-restrictions", handler.NewCreateScheduleRestrictionsHandler(svc.ScheduleRestrictionSvc))
	apiV1.Handle(http.MethodGet, "/schedule-restrictions", handler.NewFetchScheduleRestrictionsHandler(svc.ScheduleRestrictionSvc))
	apiV1.Handle(http.MethodPut, "/schedule-restrictions", handler.NewUpdateScheduleRestrictionsHandler(svc.ScheduleRestrictionSvc))
	apiV1.Handle(http.MethodPost, "/schedule-restrictions/delete", handler.NewDeleteScheduleRestrictionByIDsHandler(svc.ScheduleRestrictionSvc))

	apiV1.Handle(http.MethodGet, "/teacher-slot-preferences/:teacherSlotPreferenceId", handler.NewGetTeacherSlotPreferenceByIDHandler(svc.TeacherSlotPreferenceSvc))
	apiV1.Handle(http.MethodPost, "/teacher-slot-preferences", handler.NewCreateTeacherSlotPreferencesHandler(svc.TeacherSlotPreferenceSvc))
	apiV1.Handle(http.MethodGet, "/teacher-slot-preferences", handler.NewFetchTeacherSlotPreferencesHandler(svc.TeacherSlotPreferenceSvc))
	apiV1.Handle(http.MethodPut, "/teacher-slot-preferences", handler.NewUpdateTeacherSlotPreferencesHandler(svc.TeacherSlotPreferenceSvc))
	apiV1.Handle(http.MethodPost, "/teacher-slot-preferences/delete", handler.NewDeleteTeacherSlotPreferenceByIDsHandler(svc.TeacherSlotPreferenceSvc))

	apiV1.Handle(http.MethodGet, "/cycle-committee-lab-rooms/:cycleCommitteeLabRoomId", handler.NewGetCycleCommitteeLabRoomByIDHandler(svc.CycleCommitteeLabRoomSvc))
	apiV1.Handle(http.MethodPost, "/cycle-committee-lab-rooms", handler.NewCreateCycleCommitteeLabRoomsHandler(svc.CycleCommitteeLabRoomSvc))
	apiV1.Handle(http.MethodGet, "/cycle-committee-lab-rooms", handler.NewFetchCycleCommitteeLabRoomsHandler(svc.CycleCommitteeLabRoomSvc))
	apiV1.Handle(http.MethodPut, "/cycle-committee-lab-rooms", handler.NewUpdateCycleCommitteeLabRoomsHandler(svc.CycleCommitteeLabRoomSvc))
	apiV1.Handle(http.MethodPost, "/cycle-committee-lab-rooms/delete", handler.NewDeleteCycleCommitteeLabRoomByIDsHandler(svc.CycleCommitteeLabRoomSvc))

	apiV1.Handle(http.MethodPost, "/schedule-templates/generate", handler.NewGenerateScheduleTemplateHandler(svc.ScheduleTemplateSvc))
	apiV1.Handle(http.MethodPost, "/schedule-templates/save", handler.NewSaveScheduleTemplateHandler(svc.ScheduleTemplateSvc))
	apiV1.Handle(http.MethodGet, "/schedule-templates/:scheduleTemplateId", handler.NewGetScheduleTemplateByIDHandler(svc.ScheduleTemplateSvc))
	apiV1.Handle(http.MethodGet, "/schedule-templates", handler.NewFetchScheduleTemplatesHandler(svc.ScheduleTemplateSvc))
	apiV1.Handle(http.MethodPost, "/schedule-templates/delete", handler.NewDeleteScheduleTemplateByIDsHandler(svc.ScheduleTemplateSvc))
	apiV1.Handle(http.MethodPost, "/schedule-templates/:scheduleTemplateId/activate", handler.NewActivateScheduleTemplateHandler(svc.ScheduleTemplateSvc))

	r.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r
}
