'use strict';

// Each entry: { apiPath, autoYearFilter?, columns[], filters[], extraActions[] }
// Column types: text, number, date, bool, select, badge
// Column.ref: key into App.cache to populate a select
// autoYearFilter: true → GET requests append ?academic_year_id=<header year>
//                       and POST/PUT payloads get academic_year_id auto-injected
const ENTITY_CONFIGS = {

  'academic-years': {
    apiPath: '/api/syllabus/academic-years',
    i18nKey: 'entities.academicYears',
    columns: [
      { key: 'name',       type: 'text',  labelKey: 'fields.name',     required: true },
      { key: 'start_date', type: 'date',  labelKey: 'fields.startDate', required: true },
      { key: 'end_date',   type: 'date',  labelKey: 'fields.endDate',   required: true },
      { key: 'is_active',  type: 'badge', labelKey: 'fields.isActive',  readonly: true,
        values: { true: 'active', false: 'inactive' } },
    ],
    filters: [],
    extraActions: [
      {
        labelFn: () => t('actions.activate'),
        classFn: (row) => row.is_active ? 'hidden' : 'btn btn-sm btn-primary',
        handler: (row, mgr) => App.toggleAcademicYear(row, mgr),
      },
    ],
  },

  'semesters': {
    apiPath: '/api/syllabus/semesters',
    i18nKey: 'entities.semesters',
    autoYearFilter: true,
    columns: [
      { key: 'period_start', type: 'date', labelKey: 'fields.periodStart', required: true },
      { key: 'period_end',   type: 'date', labelKey: 'fields.periodEnd',   required: true },
    ],
    filters: [],
  },

  'departments': {
    apiPath: '/api/syllabus/departments',
    i18nKey: 'entities.departments',
    autoYearFilter: true,
    columns: [
      { key: 'name', type: 'text', labelKey: 'fields.name', required: true },
    ],
    filters: [],
  },

  'rooms': {
    apiPath: '/api/syllabus/rooms',
    i18nKey: 'entities.rooms',
    autoYearFilter: true,
    columns: [
      { key: 'name',      type: 'text',   labelKey: 'fields.name',     required: true },
      { key: 'room_type', type: 'select', labelKey: 'fields.roomType', required: true,
        options: [
          { value: 'auditorium', labelKey: 'roomTypes.auditorium' },
          { value: 'laboratory', labelKey: 'roomTypes.laboratory' },
        ] },
    ],
    filters: [
      { key: 'room_type', type: 'select', labelKey: 'fields.roomType',
        options: [
          { value: 'auditorium', labelKey: 'roomTypes.auditorium' },
          { value: 'laboratory', labelKey: 'roomTypes.laboratory' },
        ] },
    ],
  },

  'specialties': {
    apiPath: '/api/syllabus/specialties',
    i18nKey: 'entities.specialties',
    autoYearFilter: true,
    columns: [
      { key: 'name',            type: 'text',   labelKey: 'fields.name',          required: true },
      { key: 'short_name',      type: 'text',   labelKey: 'fields.shortName' },
      { key: 'specialty_code',  type: 'text',   labelKey: 'fields.specialtyCode', required: true },
      { key: 'department_id',   type: 'select', labelKey: 'fields.department',
        ref: 'departments', refLabel: 'name', required: true },
    ],
    filters: [
      { key: 'department_id', type: 'select', labelKey: 'fields.department',
        ref: 'departments', refLabel: 'name' },
    ],
  },

  'cycle-committees': {
    apiPath: '/api/syllabus/cycle-committees',
    i18nKey: 'entities.cycleCommittees',
    autoYearFilter: true,
    columns: [
      { key: 'name',    type: 'text',   labelKey: 'fields.name', required: true },
      { key: 'user_id', type: 'select', labelKey: 'fields.fullName',
        ref: 'users', refLabel: 'full_name' },
    ],
    filters: [],
  },

  'teachers': {
    apiPath: '/api/syllabus/teachers',
    i18nKey: 'entities.teachers',
    autoYearFilter: true,
    columns: [
      { key: 'user_id', type: 'select', labelKey: 'fields.fullName',
        ref: 'teacher-users', refLabel: 'full_name', required: true },
    ],
    filters: [],
  },

  'disciplines': {
    apiPath: '/api/syllabus/disciplines',
    i18nKey: 'entities.disciplines',
    autoYearFilter: true,
    columns: [
      { key: 'name',               type: 'text',   labelKey: 'fields.name',          required: true },
      { key: 'short_name',         type: 'text',   labelKey: 'fields.shortName' },
      { key: 'cycle_committee_id', type: 'select', labelKey: 'fields.cycleCommittee',
        ref: 'cycle-committees', refLabel: 'name', required: true },
      { key: 'is_splitting',  type: 'bool', labelKey: 'fields.isSplitting' },
      { key: 'is_subvention', type: 'bool', labelKey: 'fields.isSubvention' },
    ],
    filters: [
      { key: 'cycle_committee_id', type: 'select', labelKey: 'fields.cycleCommittee',
        ref: 'cycle-committees', refLabel: 'name' },
    ],
  },

  'groups': {
    apiPath: '/api/syllabus/groups',
    i18nKey: 'entities.groups',
    autoYearFilter: true,
    columns: [
      { key: 'name',       type: 'text', labelKey: 'fields.name',      required: true },
      { key: 'short_name', type: 'text', labelKey: 'fields.shortName' },
      { key: 'specialty_id', type: 'select', labelKey: 'fields.specialty',
        ref: 'specialties', refLabel: 'name', required: true },
      { key: 'is_contract',  type: 'bool', labelKey: 'fields.isContract' },
      { key: 'is_splitting', type: 'bool', labelKey: 'fields.isSplitting' },
      { key: 'education_start', type: 'date', labelKey: 'fields.educationStart', required: true },
      { key: 'education_end',   type: 'date', labelKey: 'fields.educationEnd',   required: true },
    ],
    filters: [
      { key: 'specialty_id', type: 'select', labelKey: 'fields.specialty',
        ref: 'specialties', refLabel: 'name' },
    ],
  },

  'group-semesters': {
    apiPath: '/api/syllabus/groups/semesters',
    i18nKey: 'entities.groupSemesters',
    columns: [
      { key: 'group_id',    type: 'select', labelKey: 'fields.group',
        ref: 'groups', refLabel: 'name', required: true },
      { key: 'semester_id', type: 'select', labelKey: 'fields.semester',
        ref: 'semesters', refLabel: 'id', required: true },
      { key: 'start_date',  type: 'date',   labelKey: 'fields.startDate',  required: true },
      { key: 'end_date',    type: 'date',   labelKey: 'fields.endDate',    required: true },
    ],
    filters: [
      { key: 'semester_id', type: 'select', labelKey: 'fields.semester',
        ref: 'semesters', refLabel: 'id' },
      { key: 'group_id', type: 'select', labelKey: 'fields.group',
        ref: 'groups', refLabel: 'name' },
    ],
  },

  'bell-schedules': {
    apiPath: '/api/syllabus/bell-schedules',
    i18nKey: 'entities.bellSchedules',
    autoYearFilter: true,
    columns: [
      { key: 'lesson_number', type: 'number', labelKey: 'fields.lessonNumber', required: true },
      { key: 'start_time',    type: 'text',   labelKey: 'fields.startTime',    required: true },
      { key: 'end_time',      type: 'text',   labelKey: 'fields.endTime',      required: true },
    ],
    filters: [],
  },

  // ---- Processing ----

  'study-plans': {
    apiPath: '/api/syllabus/study-plans',
    i18nKey: 'entities.studyPlans',
    autoYearFilter: true,
    columns: [
      { key: 'specialty_id',  type: 'select', labelKey: 'fields.specialty',
        ref: 'specialties', refLabel: 'name', required: true },
      { key: 'discipline_id', type: 'select', labelKey: 'fields.discipline',
        ref: 'disciplines', refLabel: 'name', required: true },
      { key: 'semester_number',  type: 'number', labelKey: 'fields.semesterNumber',  required: true },
      { key: 'lectures',         type: 'number', labelKey: 'fields.lecturesHours' },
      { key: 'laboratory',       type: 'number', labelKey: 'fields.laboratoryHours' },
      { key: 'practical',        type: 'number', labelKey: 'fields.practicalHours' },
      { key: 'independent_work', type: 'number', labelKey: 'fields.independentWork' },
      { key: 'exam',             type: 'number', labelKey: 'fields.examHours' },
      { key: 'credit',           type: 'number', labelKey: 'fields.creditHours' },
    ],
    filters: [
      { key: 'specialty_id', type: 'select', labelKey: 'fields.specialty',
        ref: 'specialties', refLabel: 'name' },
    ],
    extraActions: [
      {
        labelFn: () => '→ ' + t('workload.transfer'),
        handler: (row) => App._openTransferToWorkloadModal(row),
      },
    ],
  },

  'workload-distributions': {
    apiPath: '/api/syllabus/workload-distributions',
    i18nKey: 'entities.workloadDistributions',
    columns: [
      { key: 'study_plan_id', type: 'select', labelKey: 'fields.studyPlan',
        ref: 'study-plans', refLabel: 'id', required: true },
      { key: 'group_id', type: 'select', labelKey: 'fields.group',
        ref: 'groups', refLabel: 'name', required: true },
      { key: 'classroom_work', type: 'number', labelKey: 'fields.classroomWork' },
      { key: 'laboratory',     type: 'number', labelKey: 'fields.laboratoryHours' },
      { key: 'practical',      type: 'number', labelKey: 'fields.practicalHours' },
      { key: 'exam',           type: 'number', labelKey: 'fields.examHours' },
    ],
    filters: [
      { key: 'study_plan_id', type: 'select', labelKey: 'fields.studyPlan',
        ref: 'study-plans', refLabel: 'id' },
    ],
  },

  'workload-assignments': {
    apiPath: '/api/syllabus/workload-assignments',
    i18nKey: 'entities.workloadAssignments',
    columns: [
      { key: 'workload_distribution_id', type: 'select', labelKey: 'fields.workloadDistribution',
        ref: 'workload-distributions', refLabel: 'id', required: true },
      { key: 'teacher_id', type: 'select', labelKey: 'fields.teacher',
        ref: 'teachers', refLabel: 'id', required: true },
      { key: 'role_type', type: 'select', labelKey: 'fields.roleType',
        options: [
          { value: 'assistant', labelKey: 'roleTypes.assistant' },
          { value: 'lector',    labelKey: 'roleTypes.lector' },
        ] },
      { key: 'assigned_hours', type: 'number', labelKey: 'fields.assignedHours' },
    ],
    filters: [
      { key: 'workload_distribution_id', type: 'select', labelKey: 'fields.workloadDistribution',
        ref: 'workload-distributions', refLabel: 'id' },
    ],
  },

  // ---- Schedule ----

  'schedule-template-settings': {
    apiPath: '/api/syllabus/schedule-template-settings',
    i18nKey: 'entities.scheduleTemplateSettings',
    columns: [
      { key: 'is_active',                   type: 'badge',  labelKey: 'fields.isActive',                readonly: true,
        values: { true: 'active', false: 'inactive' } },
      { key: 'hours_per_lesson',               type: 'number', labelKey: 'fields.hoursPerLesson',             required: true },
      { key: 'max_identical_lessons_per_day',  type: 'number', labelKey: 'fields.maxIdenticalLessonsPerDay',  required: true },
      { key: 'max_study_hours_per_day',        type: 'number', labelKey: 'fields.maxStudyHoursPerDay',        required: true },
      { key: 'max_teacher_hours_per_week',     type: 'number', labelKey: 'fields.maxTeacherHoursPerWeek',     required: true },
      { key: 'max_group_lesson_hours_per_week', type: 'number', labelKey: 'fields.maxGroupLessonHoursPerWeek' },
    ],
    filters: [],
    extraActions: [
      {
        labelFn: () => t('actions.activate'),
        classFn: (row) => row.is_active ? 'btn btn-sm btn-success' : 'btn btn-sm btn-primary',
        handler: (row, mgr) => App.activateEntity('schedule-template-settings', row, mgr),
        disabledFn: (row) => row.is_active,
      },
    ],
  },

  'schedule-restrictions': {
    apiPath: '/api/syllabus/schedule-restrictions',
    i18nKey: 'entities.scheduleRestrictions',
    columns: [
      { key: 'is_active',                   type: 'badge',  labelKey: 'fields.isActive',                readonly: true,
        values: { true: 'active', false: 'inactive' } },
      { key: 'min_group_lessons_per_day',   type: 'number', labelKey: 'fields.minGroupLessonsPerDay',   required: true },
      { key: 'max_group_lessons_per_day',   type: 'number', labelKey: 'fields.maxGroupLessonsPerDay',   required: true },
      { key: 'max_teacher_lessons_per_day', type: 'number', labelKey: 'fields.maxTeacherLessonsPerDay', required: true },
      { key: 'no_gaps_in_group_schedule',   type: 'bool',   labelKey: 'fields.noGapsInGroupSchedule',  required: true },
    ],
    filters: [],
    extraActions: [
      {
        labelFn: () => t('actions.activate'),
        classFn: (row) => row.is_active ? 'btn btn-sm btn-success' : 'btn btn-sm btn-primary',
        handler: (row, mgr) => App.activateEntity('schedule-restrictions', row, mgr),
        disabledFn: (row) => row.is_active,
      },
    ],
  },

  'users': {
    apiPath: '/api/auth/users',
    i18nKey: 'entities.users',
    columns: [
      { key: 'first_name',  type: 'text',  labelKey: 'fields.firstName',  readonly: true },
      { key: 'last_name',   type: 'text',  labelKey: 'fields.lastName',   readonly: true },
      { key: 'patronymic',  type: 'text',  labelKey: 'fields.patronymic', readonly: true },
      { key: 'email',       type: 'text',  labelKey: 'fields.email',      readonly: true },
      { key: 'role',        type: 'text',  labelKey: 'fields.role',       readonly: true },
      { key: 'is_active',   type: 'badge', labelKey: 'fields.isActive',   readonly: true,
        values: { true: 'active', false: 'inactive' } },
    ],
    filters: [],
    noCreate: true,
    noEdit: true,
    extraActions: [
      {
        labelFn: (row) => row.is_active ? t('users.resetPassword') : t('users.resetInvite'),
        classFn: () => App.state.user?.role !== 'admin' ? 'hidden' : 'btn btn-sm btn-secondary',
        handler: (row) => row.is_active ? App.resetUserPassword(row) : App.resetUserInvite(row),
      },
    ],
  },

  'schedule-templates': {
    apiPath: '/api/syllabus/schedule-templates',
    i18nKey: 'entities.scheduleTemplates',
    columns: [
      { key: 'name',        type: 'text',   labelKey: 'fields.name' },
      { key: 'semester_id', type: 'select', labelKey: 'fields.semester',
        ref: 'semesters', refLabel: 'id' },
      { key: 'is_active',   type: 'badge',  labelKey: 'fields.isActive', readonly: true,
        values: { true: 'active', false: 'inactive' } },
    ],
    filters: [
      { key: 'semester_id', type: 'select', labelKey: 'fields.semester',
        ref: 'semesters', refLabel: 'id' },
    ],
    noEdit: true,
    noCreate: true,
    extraActions: [
      {
        labelFn: () => t('actions.view'),
        classFn: () => 'btn btn-sm btn-secondary',
        handler: (row) => App.viewScheduleTemplate(row),
      },
      {
        labelFn: (row) => row.is_active ? t('actions.deactivate') : t('actions.activate'),
        classFn: (row) => row.is_active ? 'btn btn-sm btn-secondary' : 'btn btn-sm btn-primary',
        handler: (row, mgr) => App.toggleScheduleTemplate(row, mgr),
      },
    ],
  },

};
