package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/syllabus/generator"
	"github.com/Lapakin/edu-planner/internal/syllabus/repository"
	"github.com/google/uuid"

	pg "github.com/Lapakin/edu-planner/internal/adapter/db/postgres"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type ScheduleTemplateService struct {
	db  *pg.DB
	rm  repository.RepoManager
	l   *logging.Logger
	cfg domain.GenerationConfig
}

func NewScheduleTemplateService(db *pg.DB, rm repository.RepoManager, l *logging.Logger, cfg domain.GenerationConfig) *ScheduleTemplateService {
	return &ScheduleTemplateService{
		db:  db,
		rm:  rm,
		l:   l,
		cfg: cfg,
	}
}

func (s *ScheduleTemplateService) GenerateScheduleTemplate(ctx context.Context, req *domain.GenerateScheduleRequest) (*domain.ScheduleData, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during schedule template generation. err: %v", err)
			s.l.Debugf("request: %v", req)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return nil, ErrInternal
	}
	defer tx.Rollback()

	// 1. Fetch schedule template settings (use the latest)
	allSettings, err := s.rm.NewScheduleTemplateSettingRepo(tx).FetchScheduleTemplateSettings(ctx, make(f.Filters))
	if err != nil || len(allSettings) == 0 {
		return nil, fmt.Errorf("%w: %w", generator.ErrNoSettings, err)
	}
	templateSetting := allSettings[len(allSettings)-1]

	// 1a. Fetch schedule restriction (use the latest; fall back to defaults if none configured)
	var restriction *domain.ScheduleRestriction
	allRestrictions, err := s.rm.NewScheduleRestrictionRepo(tx).FetchScheduleRestrictions(ctx, make(f.Filters))
	if err == nil && len(allRestrictions) > 0 {
		restriction = allRestrictions[len(allRestrictions)-1]
	}

	// 2. Fetch bell schedules
	bellSchedules, err := s.rm.NewBellScheduleRepo(tx).FetchBellSchedules(ctx, make(f.Filters))
	if err != nil || len(bellSchedules) == 0 {
		return nil, fmt.Errorf("%w: %w", generator.ErrNoBellSchedule, err)
	}

	// 3. Fetch group semesters for this semester
	gsFilters := f.Filters{domain.SemesterIDParam: strconv.FormatUint(req.SemesterID, 10)}
	groupSemesters, err := s.rm.NewGroupRepo(tx).FetchGroupSemesters(ctx, gsFilters)
	if err != nil || len(groupSemesters) == 0 {
		return nil, fmt.Errorf("%w: %w", generator.ErrNoGroupSemesters, err)
	}

	// 4. Extract group IDs and semester date range
	groupIDs := make([]uint64, 0, len(groupSemesters))
	var semesterStart, semesterEnd time.Time
	for i, gs := range groupSemesters {
		groupIDs = append(groupIDs, gs.GroupID)
		if i == 0 || gs.StartDate.Before(semesterStart) {
			semesterStart = gs.StartDate
		}
		if i == 0 || gs.EndDate.After(semesterEnd) {
			semesterEnd = gs.EndDate
		}
	}

	if len(req.GroupIDs) > 0 {
		requestedSet := make(map[uint64]bool)
		for _, id := range req.GroupIDs {
			requestedSet[id] = true
		}
		filtered := make([]uint64, 0)
		for _, id := range groupIDs {
			if requestedSet[id] {
				filtered = append(filtered, id)
			}
		}
		groupIDs = filtered
	}

	if len(groupIDs) == 0 {
		return nil, generator.ErrNoGroupSemesters
	}

	// 5. Fetch groups
	groups, err := s.rm.NewGroupRepo(tx).FetchGroups(ctx, make(f.Filters))
	if err != nil {
		return nil, ErrInternal
	}

	// 6. Fetch workload distributions and filter to relevant groups
	distributions, err := s.rm.NewWorkloadDistributionRepo(tx).FetchWorkloadDistributions(ctx, make(f.Filters))
	if err != nil {
		return nil, ErrInternal
	}
	groupIDSet := make(map[uint64]bool)
	for _, id := range groupIDs {
		groupIDSet[id] = true
	}
	var relevantDistributions domain.WorkloadDistributions
	for _, d := range distributions {
		if groupIDSet[d.GroupID] {
			relevantDistributions = append(relevantDistributions, d)
		}
	}

	// 7. Fetch workload assignments filtered to relevant distributions
	assignments, err := s.rm.NewWorkloadAssignmentRepo(tx).FetchWorkloadAssignments(ctx, make(f.Filters))
	if err != nil {
		return nil, ErrInternal
	}
	distIDSet := make(map[uint64]bool)
	for _, d := range relevantDistributions {
		distIDSet[d.ID] = true
	}
	var relevantAssignments domain.WorkloadAssignments
	for _, a := range assignments {
		if distIDSet[a.WorkloadDistributionID] {
			relevantAssignments = append(relevantAssignments, a)
		}
	}

	// 8. Fetch study plans
	studyPlans, err := s.rm.NewStudyPlanRepo(tx).FetchStudyPlans(ctx, make(f.Filters))
	if err != nil {
		return nil, ErrInternal
	}

	// 9. Fetch rooms
	rooms, err := s.rm.NewRoomRepo(tx).FetchRooms(ctx, make(f.Filters))
	if err != nil {
		return nil, ErrInternal
	}
	roomIDs := make([]uint64, 0, len(rooms))
	for _, r := range rooms {
		roomIDs = append(roomIDs, r.ID)
	}

	// 10. Build workloads
	workloads := generator.BuildWorkloads(relevantDistributions, relevantAssignments, studyPlans, groups)
	if len(workloads) == 0 {
		return nil, generator.ErrNoWorkloads
	}

	// 11. Run generator (read-only — nothing is written to DB)
	gen := generator.New(s.cfg, templateSetting, bellSchedules, restriction, workloads, roomIDs, groupIDs, semesterStart, semesterEnd)
	scheduleData, err := gen.Generate(ctx)
	if err != nil {
		return nil, err
	}

	return scheduleData, nil
}

func (s *ScheduleTemplateService) SaveScheduleTemplate(ctx context.Context, claims *jwt.Claims, req *domain.SaveScheduleTemplateRequest) (*domain.ScheduleTemplate, error) {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during schedule template save. err: %v", err)
			s.l.Debugf("claims: %v, semesterID: %v", claims, req.SemesterID)
		}
	}()

	if req.Data == nil {
		return nil, generator.ErrNoWorkloads
	}

	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		return nil, ErrInternal
	}

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return nil, ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()

	name := "Generated " + currentTime.Format("2006-01-02 15:04:05")
	if req.Name != nil && *req.Name != "" {
		name = *req.Name
	}

	template := &domain.ScheduleTemplate{
		ID:         0,
		SemesterID: &req.SemesterID,
		Name:       &name,
		Data:       dataBytes,
		IsActive:   new(false),
		IsDeleted:  new(false),
		CreatedAt:  currentTime,
		ModifiedAt: nil,
	}

	if err = s.rm.NewScheduleTemplateRepo(tx).CreateScheduleTemplate(ctx, template); err != nil {
		return nil, handleDBError(err)
	}

	var objBytes json.RawMessage
	objBytes, err = json.Marshal(template)
	if err != nil {
		return nil, ErrInternal
	}

	massage := &domain.Massage{
		EventID:    uuid.New().String(),
		ActionType: domain.ActionTypeCreate,
		ObjectType: domain.ObjectTypeScheduleTemplate,
		Object:     objBytes,
		Claims:     claims,
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, domain.Massages{massage}); err != nil {
		return nil, handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, ErrInternal
	}

	return template, nil
}

func (s *ScheduleTemplateService) GetScheduleTemplateByID(ctx context.Context, id uint64) (*domain.ScheduleTemplate, error) {
	template, err := s.rm.NewScheduleTemplateRepo(s.db).GetScheduleTemplateByID(ctx, id)
	if err != nil {
		s.l.Errorf("error during getting schedule template by id. err: %v", err)
		s.l.Debugf("scheduleTemplateID: %v", id)
		return nil, handleDBError(err)
	}
	return template, nil
}

func (s *ScheduleTemplateService) FetchScheduleTemplates(ctx context.Context, filters f.Filters) (domain.ScheduleTemplates, error) {
	templates, err := s.rm.NewScheduleTemplateRepo(s.db).FetchScheduleTemplates(ctx, filters)
	if err != nil {
		s.l.Errorf("error during fetching schedule templates. err: %v", err)
		s.l.Debugf("filters: %v", filters)
		return nil, ErrInternal
	}
	return templates, nil
}

func (s *ScheduleTemplateService) DeleteScheduleTemplates(ctx context.Context, claims *jwt.Claims, ids []uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during deleting schedule templates. err: %v", err)
			s.l.Debugf("claims: %v, scheduleTemplateIDs: %v", claims, ids)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	if err = s.rm.NewScheduleTemplateRepo(tx).DeleteScheduleTemplates(ctx, ids); err != nil {
		return handleDBError(err)
	}

	massages := make(domain.Massages, len(ids))
	for i, id := range ids {
		var objBytes json.RawMessage
		objBytes, err = json.Marshal(id)
		if err != nil {
			return ErrInternal
		}

		massages[i] = &domain.Massage{
			EventID:    uuid.New().String(),
			ActionType: domain.ActionTypeDelete,
			ObjectType: domain.ObjectTypeScheduleTemplate,
			Object:     objBytes,
			Claims:     claims,
		}
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, massages); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *ScheduleTemplateService) ActivateScheduleTemplate(ctx context.Context, claims *jwt.Claims, id uint64) error {
	var err error
	defer func() {
		if err != nil {
			s.l.Errorf("error during activating schedule template. err: %v", err)
			s.l.Debugf("claims: %v, scheduleTemplateID: %v", claims, id)
		}
	}()

	tx, err := s.db.Begintx(ctx)
	if err != nil {
		return ErrInternal
	}
	defer tx.Rollback()

	currentTime := time.Now()
	if err = s.rm.NewScheduleTemplateRepo(tx).ActivateScheduleTemplate(ctx, id, currentTime); err != nil {
		return handleDBError(err)
	}

	var objBytes json.RawMessage
	objBytes, err = json.Marshal(id)
	if err != nil {
		return ErrInternal
	}

	massage := &domain.Massage{
		EventID:    uuid.New().String(),
		ActionType: domain.ActionTypeUpdate,
		ObjectType: domain.ObjectTypeScheduleTemplate,
		Object:     objBytes,
		Claims:     claims,
	}

	if err = s.rm.NewMessageRepo(tx).Write(ctx, domain.Massages{massage}); err != nil {
		return handleDBError(err)
	}

	if err = tx.Commit(); err != nil {
		return ErrInternal
	}

	return nil
}
