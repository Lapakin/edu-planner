package domain

import (
	"github.com/Lapakin/edu-planner/internal/adapter/json"
	"github.com/Lapakin/edu-planner/internal/adapter/jwt"
)

const (
	ActionTypeCreate     = "create"
	ActionTypeUpdate     = "update"
	ActionTypeDelete     = "delete"
	ActionTypeActivate   = "activate"
	ActionTypeDeactivate = "deactivate"
	ActionTypeLogin      = "login"
)

const (
	ObjectTypeUser                    = "user"
	ObjectTypeAcademicYear            = "academic_year"
	ObjectTypeSemester                = "semester"
	ObjectTypeBellSchedule            = "bell_schedule"
	ObjectTypeDepartment              = "department"
	ObjectTypeRoom                    = "room"
	ObjectTypeSpecialty               = "specialty"
	ObjectTypeCycleCommittee          = "cycle_committee"
	ObjectTypeTeacher                 = "teacher"
	ObjectTypeDiscipline              = "discipline"
	ObjectTypeGroup                   = "group"
	ObjectTypeGroupSemester           = "group_semester"
	ObjectTypeStudyPlan               = "study_plan"
	ObjectTypeWorkloadAssignment      = "workload_assignment"
	ObjectTypeWorkloadDistribution    = "workload_distribution"
	ObjectTypeScheduleTemplateSetting = "schedule_template_setting"
	ObjectTypeScheduleTemplate        = "schedule_template"
	ObjectTypeScheduleRestriction     = "schedule_restriction"
)

type Massage struct {
	EventID    string          `json:"event_id" db:"event_id"`
	ActionType string          `json:"action_type" db:"action_type"`
	ObjectType string          `json:"object_type" db:"object_type"`
	Object     json.RawMessage `json:"object" db:"object"`
	Claims     *jwt.Claims     `json:"claims"`
}

type Massages []*Massage
