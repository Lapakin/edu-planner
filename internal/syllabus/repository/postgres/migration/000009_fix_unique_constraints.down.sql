ALTER TABLE teacher_slot_preference
    DROP CONSTRAINT IF EXISTS teacher_slot_preference_academic_year_id_teacher_id_weekday_les,
    ALTER COLUMN academic_year_id DROP NOT NULL,
    ADD CONSTRAINT teacher_slot_preference_teacher_id_weekday_lesson_number_slot_t
        UNIQUE (teacher_id, weekday, lesson_number, slot_type);

ALTER TABLE cycle_committee_lab_room
    DROP CONSTRAINT IF EXISTS cycle_committee_lab_room_academic_year_id_cycle_committee_id_ro,
    ALTER COLUMN academic_year_id DROP NOT NULL,
    ADD CONSTRAINT cycle_committee_lab_room_cycle_committee_id_room_id_key
        UNIQUE (cycle_committee_id, room_id);
