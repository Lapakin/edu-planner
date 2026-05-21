CREATE TABLE IF NOT EXISTS teacher_slot_preference (
    id               bigserial PRIMARY KEY,
    academic_year_id BIGINT REFERENCES academic_year(id),
    teacher_id       BIGINT NOT NULL,
    weekday          VARCHAR NOT NULL,
    lesson_number    INT NOT NULL,
    slot_type        VARCHAR NOT NULL DEFAULT 'preferred',
    created_at       TIMESTAMP NOT NULL,
    modified_at      TIMESTAMP,
    UNIQUE (teacher_id, weekday, lesson_number, slot_type)
);
