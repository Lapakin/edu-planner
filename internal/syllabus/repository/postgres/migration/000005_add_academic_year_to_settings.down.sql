
ALTER TABLE schedule_restriction
    DROP COLUMN IF EXISTS allow_flow_lessons,
    DROP COLUMN IF EXISTS time_priority,
    DROP COLUMN IF EXISTS max_consecutive_teacher_lessons,
    DROP COLUMN IF EXISTS academic_year_id;

ALTER TABLE schedule_template_setting
    DROP COLUMN IF EXISTS study_days_mask,
    DROP COLUMN IF EXISTS lessons_per_class,
    DROP COLUMN IF EXISTS academic_year_id;
