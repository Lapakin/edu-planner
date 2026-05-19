CREATE TABLE IF NOT EXISTS schedule_restriction
(
    id                          bigserial PRIMARY KEY,
    min_group_lessons_per_day   int       NOT NULL DEFAULT 2,
    max_group_lessons_per_day   int       NOT NULL DEFAULT 4,
    max_teacher_lessons_per_day int       NOT NULL DEFAULT 5,
    no_gaps_in_group_schedule   bool      NOT NULL DEFAULT true,
    created_at                  timestamp NOT NULL,
    modified_at                 timestamp
);
