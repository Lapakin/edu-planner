ALTER TABLE schedule_template_setting
    ADD COLUMN IF NOT EXISTS academic_year_id BIGINT REFERENCES academic_year(id),
    ADD COLUMN IF NOT EXISTS lessons_per_class INT NOT NULL DEFAULT 2,
    ADD COLUMN IF NOT EXISTS study_days_mask  INT NOT NULL DEFAULT 31;

ALTER TABLE schedule_restriction
    ADD COLUMN IF NOT EXISTS academic_year_id                BIGINT  REFERENCES academic_year(id),
    ADD COLUMN IF NOT EXISTS max_consecutive_teacher_lessons INT     NOT NULL DEFAULT 4,
    ADD COLUMN IF NOT EXISTS time_priority                   VARCHAR NOT NULL DEFAULT 'none',
    ADD COLUMN IF NOT EXISTS allow_flow_lessons              BOOL    NOT NULL DEFAULT TRUE;
