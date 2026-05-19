CREATE TABLE IF NOT EXISTS academic_year
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR   NOT NULL,
    start_date  TIMESTAMP NOT NULL,
    end_date    TIMESTAMP NOT NULL,
    is_active   BOOLEAN   NOT NULL DEFAULT FALSE,
    is_deleted  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL,
    modified_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS semester
(
    id               BIGSERIAL PRIMARY KEY,
    academic_year_id BIGINT    NOT NULL REFERENCES academic_year (id),
    period_start     TIMESTAMP NOT NULL,
    period_end       TIMESTAMP NOT NULL,
    is_deleted       BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMP NOT NULL,
    modified_at      TIMESTAMP
);

CREATE TABLE IF NOT EXISTS department
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR   NOT NULL,
    is_deleted  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL,
    modified_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS academic_year_to_department
(
    academic_year_id BIGINT NOT NULL REFERENCES academic_year (id),
    department_id    BIGINT NOT NULL REFERENCES department (id)
);

CREATE TABLE IF NOT EXISTS room
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR   NOT NULL,
    room_type   VARCHAR   NOT NULL,
    is_deleted  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL,
    modified_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS academic_year_to_room
(
    academic_year_id BIGINT NOT NULL REFERENCES academic_year (id),
    room_id          BIGINT NOT NULL REFERENCES room (id)
);

CREATE TABLE IF NOT EXISTS specialty
(
    id             BIGSERIAL PRIMARY KEY,
    department_id  BIGINT    NOT NULL REFERENCES department (id),
    name           VARCHAR   NOT NULL,
    short_name     VARCHAR   NOT NULL,
    specialty_code VARCHAR   NOT NULL,
    is_deleted     BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMP NOT NULL,
    modified_at    TIMESTAMP
);

CREATE TABLE IF NOT EXISTS academic_year_to_specialty
(
    academic_year_id BIGINT NOT NULL REFERENCES academic_year (id),
    specialty_id     BIGINT NOT NULL REFERENCES specialty (id)
);

CREATE TABLE IF NOT EXISTS user_info
(
    id         BIGINT  PRIMARY KEY,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS cycle_committee
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT    NOT NULL REFERENCES user_info (id),
    name        VARCHAR   NOT NULL,
    is_deleted  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL,
    modified_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS academic_year_to_cycle_committee
(
    academic_year_id   BIGINT NOT NULL REFERENCES academic_year (id),
    cycle_committee_id BIGINT NOT NULL REFERENCES cycle_committee (id)
);

CREATE TABLE IF NOT EXISTS teacher
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT    NOT NULL REFERENCES user_info (id),
    is_deleted  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL,
    modified_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS academic_year_to_teacher
(
    academic_year_id BIGINT NOT NULL REFERENCES academic_year (id),
    teacher_id       BIGINT NOT NULL REFERENCES teacher (id)
);

CREATE TABLE IF NOT EXISTS discipline
(
    id                 BIGSERIAL PRIMARY KEY,
    cycle_committee_id BIGINT    NOT NULL REFERENCES cycle_committee (id),
    name               VARCHAR   NOT NULL,
    short_name         VARCHAR   NOT NULL,
    is_splitting       BOOLEAN   NOT NULL DEFAULT FALSE,
    is_subvention      BOOLEAN   NOT NULL DEFAULT FALSE,
    is_deleted         BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMP NOT NULL,
    modified_at        TIMESTAMP
);

CREATE TABLE IF NOT EXISTS academic_year_to_discipline
(
    academic_year_id BIGINT NOT NULL REFERENCES academic_year (id),
    discipline_id    BIGINT NOT NULL REFERENCES discipline (id)
);

CREATE TABLE IF NOT EXISTS education_group
(
    id              BIGSERIAL PRIMARY KEY,
    specialty_id    BIGINT    NOT NULL REFERENCES specialty (id),
    name            VARCHAR   NOT NULL,
    short_name      VARCHAR   NOT NULL,
    is_contract     BOOLEAN   NOT NULL DEFAULT FALSE,
    is_splitting    BOOLEAN   NOT NULL DEFAULT FALSE,
    education_start TIMESTAMP NOT NULL,
    education_end   TIMESTAMP NOT NULL,
    is_deleted      BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMP NOT NULL,
    modified_at     TIMESTAMP
);

CREATE TABLE IF NOT EXISTS academic_year_to_education_group
(
    academic_year_id BIGINT NOT NULL REFERENCES academic_year (id),
    group_id         BIGINT NOT NULL REFERENCES education_group (id)
);

CREATE TABLE IF NOT EXISTS education_group_semester
(
    id          BIGSERIAL PRIMARY KEY,
    group_id    BIGINT    NOT NULL REFERENCES education_group (id),
    semester_id BIGINT    NOT NULL REFERENCES semester (id),
    start_date  TIMESTAMP NOT NULL,
    end_date    TIMESTAMP NOT NULL,
    is_deleted  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL,
    modified_at TIMESTAMP
);

CREATE TABLE study_plan
(
    id                    BIGSERIAL PRIMARY KEY,
    academic_year_id      BIGINT,
    specialty_id          BIGINT,
    discipline_id         BIGINT,
    semester_number       INT,
    educational_component VARCHAR,
    okr                   INT,
    coursework_type       VARCHAR,
    classroom_work        INT,
    lectures              INT,
    laboratory            INT,
    practical             INT,
    independent_work      INT,
    individual_lessons    INT,
    consultations         INT,
    exam                  INT,
    credit                INT,
    practice              INT,
    coursework_hours      INT,
    classroom_coursework  INT,
    assisting_hours       INT,
    created_at            TIMESTAMP NOT NULL,
    modified_at           TIMESTAMP,
    is_deleted            BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS schedule_template
(
    id          bigserial PRIMARY KEY,
    semester_id bigint,
    name        varchar,
    data        jsonb,
    is_active   bool,
    created_at  timestamp,
    modified_at timestamp,
    is_deleted  bool
);
CREATE TABLE IF NOT EXISTS event
(
    id                   bigserial PRIMARY KEY,
    chain_id             uuid,
    schedule_template_id bigint,
    discipline_id        bigint,
    teacher_id           bigint,
    group_id             bigint,
    room_id              bigint,
    lesson_number        int,
    subgroup_number      int,
    created_at           timestamp,
    modified_at          timestamp,
    is_deleted           bool
);
CREATE TABLE IF NOT EXISTS event_rrule
(
    id           bigserial PRIMARY KEY,
    event_id     bigint,
    rrule_string varchar,
    valid_from   date,
    valid_until  date,
    created_at   timestamp,
    modified_at  timestamp
);
CREATE TABLE IF NOT EXISTS workload_distribution
(
    id                   bigserial PRIMARY KEY,
    study_plan_id        bigint,
    group_id             bigint,
    coursework_type      varchar,
    coursework_hours     int,
    classroom_coursework int,
    classroom_work       int,
    laboratory           int,
    practical            int,
    exam                 int,
    consultations        int,
    credit               int,
    practice             int,
    okr                  int,
    individual_lessons   int,
    deductions           int,
    created_at           timestamp,
    modified_at          timestamp,
    is_deleted           bool NOT NULL DEFAULT FALSE
);
CREATE TABLE IF NOT EXISTS workload_assignment
(
    id                       bigserial PRIMARY KEY,
    workload_distribution_id bigint NOT NULL,
    teacher_id               bigint NOT NULL,
    role_type                varchar,
    assigned_hours           int,
    created_at               timestamp
);
CREATE TABLE IF NOT EXISTS schedule_template_setting
(
    id                              bigserial PRIMARY KEY,
    hours_per_lesson                numeric(4, 2) NOT NULL,
    max_identical_lessons_per_day   int           NOT NULL,
    max_study_hours_per_day         int           NOT NULL,
    max_teacher_hours_per_week      int           NOT NULL,
    max_group_lesson_hours_per_week int,
    created_at                      timestamp,
    modified_at                     timestamp
);
CREATE TABLE IF NOT EXISTS bell_schedule
(
    id            bigserial PRIMARY KEY,
    lesson_number int UNIQUE NOT NULL,
    start_time    time       NOT NULL,
    end_time      time       NOT NULL
);