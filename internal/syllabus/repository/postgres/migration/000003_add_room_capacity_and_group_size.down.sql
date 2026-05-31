ALTER TABLE education_group
    DROP COLUMN IF EXISTS student_count;

ALTER TABLE room
    DROP COLUMN IF EXISTS capacity;
