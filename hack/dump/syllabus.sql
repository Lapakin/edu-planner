-- Demo seed data for syllabus database
-- Apply after migrations with: make db-dump-up
--
-- Provides a complete dataset for schedule generation including
-- clear numerator/denominator differences:
--
--   Програмування  38h → 2 lessons/week numerator, 1 lesson/week denominator
--   Фізичне виховання 2h → numerator only (absent in denominator)
--   All other disciplines → same count in both periods
--
-- Academic year 2025-2026, 2 semesters
-- 1 department, 2 specialties (КН, ІСТ)
-- 7 rooms (4 auditorium + 3 laboratory)
-- 2 cycle committees, 10 disciplines
-- 5 teachers (matching user ids 2-6 from user-management)
-- 6 student groups:
--   КН-31, КН-32, ІСТ-31, ІСТ-32  — semesters 1 & 2 (5th semester curriculum)
--   КН-21, ІСТ-21                   — semester 2 only  (6th semester curriculum)
-- Bell schedule: 5 lesson slots (08:00-15:30)
--
-- Semester 1 per-group weekly lesson counts:
--   КН-31 numerator:   Прог(2)+БД(1)+Мережі(1)+ОС(1)+АСД(1)+ФВ(1) = 7/week
--   КН-31 denominator: Прог(1)+БД(1)+Мережі(1)+ОС(1)+АСД(1)       = 5/week
--   КН-32 numerator:   Прог(2)+БД(1)+Мережі(1)+АСД(1)+ФВ(1)        = 6/week
--   КН-32 denominator: Прог(1)+БД(1)+Мережі(1)+АСД(1)              = 4/week
--   ІСТ-31 numerator:  Мат(1)+БД(1)+Мережі(1)+ФВ(1)                = 4/week
--   ІСТ-31 denominator:Мат(1)+БД(1)+Мережі(1)                      = 3/week
--   ІСТ-32 same as ІСТ-31
--
-- Semester 2 additional groups:
--   КН-21: КГ(1)+СА(1)+Веб(1) = 3/week both periods
--   ІСТ-21: КГ(1)+СА(1)+Веб(1) = 3/week both periods
--
-- To generate a schedule: POST /api/syllabus/schedule-templates/generate
-- with body: {"semester_id": 1}  or  {"semester_id": 2}


-- ─────────────────────────────────────────────────────────────
-- Academic year
-- ─────────────────────────────────────────────────────────────

INSERT INTO academic_year (id, name, start_date, end_date, is_active, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES (1, '2025-2026', '2025-09-01 00:00:00', '2026-06-30 00:00:00', TRUE, FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Semesters  (semester 1 = Fall 2025, ~19 weeks → 9 num + 9 denom + 1 extra)
-- ─────────────────────────────────────────────────────────────

INSERT INTO semester (id, academic_year_id, period_start, period_end, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 1, '2025-09-01 00:00:00', '2026-01-17 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (2, 1, '2026-02-02 00:00:00', '2026-06-14 00:00:00', FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Department
-- ─────────────────────────────────────────────────────────────

INSERT INTO department (id, name, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES (1, 'Комп''ютерні науки та інформаційні технології', FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_department (academic_year_id, department_id)
VALUES (1, 1)
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Rooms
-- ─────────────────────────────────────────────────────────────

INSERT INTO room (id, name, room_type, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (1, '101',     'auditorium', FALSE, '2025-09-01 00:00:00'),
    (2, '102',     'auditorium', FALSE, '2025-09-01 00:00:00'),
    (3, '201',     'auditorium', FALSE, '2025-09-01 00:00:00'),
    (4, 'Лаб 301', 'laboratory', FALSE, '2025-09-01 00:00:00'),
    (5, 'Лаб 302', 'laboratory', FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_room (academic_year_id, room_id)
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5)
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Specialties
-- ─────────────────────────────────────────────────────────────

INSERT INTO specialty (id, department_id, name, short_name, specialty_code, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 1, 'Комп''ютерні науки',                'КН',  '122', FALSE, '2025-09-01 00:00:00'),
    (2, 1, 'Інформаційні системи і технології',  'ІСТ', '126', FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_specialty (academic_year_id, specialty_id)
VALUES (1, 1), (1, 2)
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Cycle committees
-- ─────────────────────────────────────────────────────────────

INSERT INTO cycle_committee (id, user_id, name, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 2, 'Циклова комісія математичних та загальнотехнічних дисциплін', FALSE, '2025-09-01 00:00:00'),
    (2, 3, 'Циклова комісія комп''ютерних наук та ІТ',                   FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_cycle_committee (academic_year_id, cycle_committee_id)
VALUES (1, 1), (1, 2)
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Teachers (user ids 2-6 from user-management DB)
-- ─────────────────────────────────────────────────────────────

INSERT INTO teacher (id, user_id, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 2, FALSE, '2025-09-01 00:00:00'),  -- Петренко І.М.  — Програмування, ФВ
    (2, 3, FALSE, '2025-09-01 00:00:00'),  -- Коваленко М.В. — АСД, ОС, ФВ
    (3, 4, FALSE, '2025-09-01 00:00:00'),  -- Шевченко О.І.  — Бази даних
    (4, 5, FALSE, '2025-09-01 00:00:00'),  -- Бондаренко Н.П.— Комп'ютерні мережі
    (5, 6, FALSE, '2025-09-01 00:00:00')   -- Марченко Д.О.  — Вища математика
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_teacher (academic_year_id, teacher_id)
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5)
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Disciplines
-- ─────────────────────────────────────────────────────────────

INSERT INTO discipline (id, cycle_committee_id, name, short_name, is_splitting, is_subvention, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 1, 'Вища математика',               'Мат',  FALSE, FALSE, FALSE, '2025-09-01 00:00:00'),
    (2, 1, 'Алгоритми та структури даних',  'АСД',  FALSE, FALSE, FALSE, '2025-09-01 00:00:00'),
    (3, 2, 'Програмування',                 'Прог', TRUE,  FALSE, FALSE, '2025-09-01 00:00:00'),
    (4, 2, 'Бази даних',                    'БД',   FALSE, FALSE, FALSE, '2025-09-01 00:00:00'),
    (5, 2, 'Комп''ютерні мережі',           'Мер',  FALSE, FALSE, FALSE, '2025-09-01 00:00:00'),
    (6, 2, 'Операційні системи',            'ОС',   FALSE, FALSE, FALSE, '2025-09-01 00:00:00'),
    (7, 1, 'Фізичне виховання',             'ФВ',   FALSE, FALSE, FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_discipline (academic_year_id, discipline_id)
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7)
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Student groups
-- ─────────────────────────────────────────────────────────────

INSERT INTO education_group (id, specialty_id, name, short_name, is_contract, is_splitting,
                              education_start, education_end, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 1, 'КН-31',  'КН-31',  FALSE, FALSE, '2022-09-01 00:00:00', '2026-06-30 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (2, 1, 'КН-32',  'КН-32',  FALSE, TRUE,  '2022-09-01 00:00:00', '2026-06-30 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (3, 2, 'ІСТ-31', 'ІСТ-31', FALSE, FALSE, '2022-09-01 00:00:00', '2026-06-30 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (4, 2, 'ІСТ-32', 'ІСТ-32', FALSE, FALSE, '2022-09-01 00:00:00', '2026-06-30 00:00:00', FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_education_group (academic_year_id, group_id)
VALUES (1, 1), (1, 2), (1, 3), (1, 4)
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Group semesters (both semesters for all groups)
-- ─────────────────────────────────────────────────────────────

INSERT INTO education_group_semester (id, group_id, semester_id, start_date, end_date, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    -- Semester 1 (Fall 2025)
    (1, 1, 1, '2025-09-01 00:00:00', '2026-01-17 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (2, 2, 1, '2025-09-01 00:00:00', '2026-01-17 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (3, 3, 1, '2025-09-01 00:00:00', '2026-01-17 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (4, 4, 1, '2025-09-01 00:00:00', '2026-01-17 00:00:00', FALSE, '2025-09-01 00:00:00'),
    -- Semester 2 (Spring 2026)
    (5, 1, 2, '2026-02-02 00:00:00', '2026-06-14 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (6, 2, 2, '2026-02-02 00:00:00', '2026-06-14 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (7, 3, 2, '2026-02-02 00:00:00', '2026-06-14 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (8, 4, 2, '2026-02-02 00:00:00', '2026-06-14 00:00:00', FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Bell schedule (5 lesson slots)
-- ─────────────────────────────────────────────────────────────

INSERT INTO bell_schedule (id, lesson_number, start_time, end_time)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 1, '08:00', '09:20'),
    (2, 2, '09:30', '10:50'),
    (3, 3, '11:10', '12:30'),
    (4, 4, '12:40', '14:00'),
    (5, 5, '14:10', '15:30')
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Schedule template settings
-- max_group_lesson_hours_per_week=16 → 8 lessons/week max per group
-- ─────────────────────────────────────────────────────────────

INSERT INTO schedule_template_setting (
    id, hours_per_lesson, max_identical_lessons_per_day,
    max_study_hours_per_day, max_teacher_hours_per_week,
    max_group_lesson_hours_per_week, created_at
)
OVERRIDING SYSTEM VALUE
VALUES (1, 2.00, 1, 8, 36, 36, '2025-09-01 00:00:00')
ON CONFLICT (id) DO UPDATE SET
    max_teacher_hours_per_week      = EXCLUDED.max_teacher_hours_per_week,
    max_group_lesson_hours_per_week = EXCLUDED.max_group_lesson_hours_per_week;


-- ─────────────────────────────────────────────────────────────
-- Study plans
-- classroom_work = total contact hours for the semester
-- ─────────────────────────────────────────────────────────────
--
-- КН specialty (specialty_id=1):
--   Програмування  (disc 3): 38h → 2/week numerator, 1/week denominator  ← DIFFERENT
--   Бази даних      (disc 4): 32h → 1/week both
--   Мережі          (disc 5): 16h → 1/week both
--   ОС              (disc 6): 16h → 1/week both
--   АСД             (disc 2): 20h → 1/week both
--   Фіз. виховання  (disc 7):  2h → 1/week numerator, 0/week denominator ← DIFFERENT
--
-- ІСТ specialty (specialty_id=2):
--   Вища математика (disc 1): 32h → 1/week both
--   Бази даних      (disc 4): 32h → 1/week both
--   Мережі          (disc 5): 32h → 1/week both
--   Фіз. виховання  (disc 7):  2h → 1/week numerator, 0/week denominator ← DIFFERENT

INSERT INTO study_plan (
    id, academic_year_id, specialty_id, discipline_id, semester_number,
    classroom_work, lectures, laboratory, practical, independent_work,
    created_at, is_deleted
)
OVERRIDING SYSTEM VALUE
VALUES
    -- КН specialty
    (1, 1, 1, 3, 5,  38, 38,  0,  0, 22, '2025-09-01 00:00:00', FALSE),  -- Програмування (38h)
    (2, 1, 1, 4, 5,  32, 16,  0, 16, 32, '2025-09-01 00:00:00', FALSE),  -- Бази даних
    (3, 1, 1, 5, 5,  16, 16,  0,  0, 16, '2025-09-01 00:00:00', FALSE),  -- Мережі
    (4, 1, 1, 6, 5,  16, 16,  0,  0, 16, '2025-09-01 00:00:00', FALSE),  -- ОС
    (8, 1, 1, 2, 5,  20, 20,  0,  0, 20, '2025-09-01 00:00:00', FALSE),  -- АСД
    (9, 1, 1, 7, 5,   2,  0,  0,  2,  0, '2025-09-01 00:00:00', FALSE),  -- Фіз.виховання (2h)
    -- ІСТ specialty
    (5, 1, 2, 1, 5,  32, 32,  0,  0, 32, '2025-09-01 00:00:00', FALSE),  -- Вища математика
    (6, 1, 2, 4, 5,  32, 16,  0, 16, 32, '2025-09-01 00:00:00', FALSE),  -- Бази даних
    (7, 1, 2, 5, 5,  32, 16,  0, 16, 32, '2025-09-01 00:00:00', FALSE),  -- Мережі
    (10, 1, 2, 7, 5,  2,  0,  0,  2,  0, '2025-09-01 00:00:00', FALSE)   -- Фіз.виховання (2h)
ON CONFLICT (id) DO UPDATE SET
    classroom_work = EXCLUDED.classroom_work,
    lectures       = EXCLUDED.lectures;


-- ─────────────────────────────────────────────────────────────
-- Workload distributions  (group × study_plan)
-- classroom_work drives the lesson count in the generator
-- ─────────────────────────────────────────────────────────────

INSERT INTO workload_distribution (
    id, study_plan_id, group_id, classroom_work, laboratory, practical,
    created_at, is_deleted
)
OVERRIDING SYSTEM VALUE
VALUES
    -- КН-31 (group 1)
    ( 1,  1, 1, 38,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Програмування  38h
    ( 2,  2, 1, 32,  0, 16, '2025-09-01 00:00:00', FALSE),  -- Бази даних
    ( 3,  3, 1, 16,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Мережі
    ( 4,  4, 1, 16,  0,  0, '2025-09-01 00:00:00', FALSE),  -- ОС
    (14,  8, 1, 20,  0,  0, '2025-09-01 00:00:00', FALSE),  -- АСД
    (16,  9, 1,  2,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Фіз.виховання  2h

    -- КН-32 (group 2, is_splitting=TRUE)
    ( 5,  1, 2, 38,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Програмування  38h
    ( 6,  2, 2, 32,  0, 16, '2025-09-01 00:00:00', FALSE),  -- Бази даних
    ( 7,  3, 2, 16,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Мережі
    (15,  8, 2, 20,  0,  0, '2025-09-01 00:00:00', FALSE),  -- АСД
    (17,  9, 2,  2,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Фіз.виховання  2h

    -- ІСТ-31 (group 3)
    ( 8,  5, 3, 32,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Вища математика
    ( 9,  6, 3, 32,  0, 16, '2025-09-01 00:00:00', FALSE),  -- Бази даних
    (10,  7, 3, 32,  0, 16, '2025-09-01 00:00:00', FALSE),  -- Мережі
    (18, 10, 3,  2,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Фіз.виховання  2h

    -- ІСТ-32 (group 4)
    (11,  5, 4, 32,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Вища математика
    (12,  6, 4, 32,  0, 16, '2025-09-01 00:00:00', FALSE),  -- Бази даних
    (13,  7, 4, 32,  0, 16, '2025-09-01 00:00:00', FALSE),  -- Мережі
    (19, 10, 4,  2,  0,  0, '2025-09-01 00:00:00', FALSE)   -- Фіз.виховання  2h
ON CONFLICT (id) DO UPDATE SET
    classroom_work = EXCLUDED.classroom_work;


-- ─────────────────────────────────────────────────────────────
-- Workload assignments  (teacher × distribution)
-- assigned_hours = hours this teacher covers; the generator uses
-- this value directly as totalHours for lesson-count calculation
-- ─────────────────────────────────────────────────────────────
--
-- Teacher 1 (Петренко)   → Програмування (КН-31, КН-32) + ФВ (КН-31, ІСТ-31)
-- Teacher 2 (Коваленко)  → АСД (КН-31, КН-32) + ОС (КН-31) + ФВ (КН-32, ІСТ-32)
-- Teacher 3 (Шевченко)   → Бази даних (all 4 groups)
-- Teacher 4 (Бондаренко) → Комп'ютерні мережі (all 4 groups)
-- Teacher 5 (Марченко)   → Вища математика (ІСТ-31, ІСТ-32)

INSERT INTO workload_assignment (
    id, workload_distribution_id, teacher_id, role_type, assigned_hours, created_at
)
OVERRIDING SYSTEM VALUE
VALUES
    -- КН-31
    ( 1,  1, 1, 'lecturer', 38, '2025-09-01 00:00:00'),  -- Петренко  → Прог   / КН-31
    ( 2,  2, 3, 'lecturer', 32, '2025-09-01 00:00:00'),  -- Шевченко  → БД     / КН-31
    ( 3,  3, 4, 'lecturer', 16, '2025-09-01 00:00:00'),  -- Бондаренко→ Мережі / КН-31
    ( 4,  4, 2, 'lecturer', 16, '2025-09-01 00:00:00'),  -- Коваленко → ОС     / КН-31
    (14, 14, 2, 'lecturer', 20, '2025-09-01 00:00:00'),  -- Коваленко → АСД    / КН-31
    (16, 16, 1, 'lecturer',  2, '2025-09-01 00:00:00'),  -- Петренко  → ФВ     / КН-31

    -- КН-32
    ( 5,  5, 1, 'lecturer', 38, '2025-09-01 00:00:00'),  -- Петренко  → Прог   / КН-32
    ( 6,  6, 3, 'lecturer', 32, '2025-09-01 00:00:00'),  -- Шевченко  → БД     / КН-32
    ( 7,  7, 4, 'lecturer', 16, '2025-09-01 00:00:00'),  -- Бондаренко→ Мережі / КН-32
    (15, 15, 2, 'lecturer', 20, '2025-09-01 00:00:00'),  -- Коваленко → АСД    / КН-32
    (17, 17, 2, 'lecturer',  2, '2025-09-01 00:00:00'),  -- Коваленко → ФВ     / КН-32

    -- ІСТ-31
    ( 8,  8, 5, 'lecturer', 32, '2025-09-01 00:00:00'),  -- Марченко  → Мат    / ІСТ-31
    ( 9,  9, 3, 'lecturer', 32, '2025-09-01 00:00:00'),  -- Шевченко  → БД     / ІСТ-31
    (10, 10, 4, 'lecturer', 32, '2025-09-01 00:00:00'),  -- Бондаренко→ Мережі / ІСТ-31
    (18, 18, 1, 'lecturer',  2, '2025-09-01 00:00:00'),  -- Петренко  → ФВ     / ІСТ-31

    -- ІСТ-32
    (11, 11, 5, 'lecturer', 32, '2025-09-01 00:00:00'),  -- Марченко  → Мат    / ІСТ-32
    (12, 12, 3, 'lecturer', 32, '2025-09-01 00:00:00'),  -- Шевченко  → БД     / ІСТ-32
    (13, 13, 4, 'lecturer', 32, '2025-09-01 00:00:00'),  -- Бондаренко→ Мережі / ІСТ-32
    (19, 19, 2, 'lecturer',  2, '2025-09-01 00:00:00')   -- Коваленко → ФВ     / ІСТ-32
ON CONFLICT (id) DO UPDATE SET
    assigned_hours = EXCLUDED.assigned_hours;


-- ─────────────────────────────────────────────────────────────
-- Additional rooms (103, Лаб 303)
-- ─────────────────────────────────────────────────────────────

INSERT INTO room (id, name, room_type, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (6, '103',     'auditorium', FALSE, '2025-09-01 00:00:00'),
    (7, 'Лаб 303', 'laboratory', FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_room (academic_year_id, room_id)
VALUES (1, 6), (1, 7)
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Additional disciplines  (КГ, СА, Веб, Дискр)
-- ─────────────────────────────────────────────────────────────

INSERT INTO discipline (id, cycle_committee_id, name, short_name, is_splitting, is_subvention, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    ( 8, 2, 'Комп''ютерна графіка', 'КГ',    TRUE,  FALSE, FALSE, '2025-09-01 00:00:00'),
    ( 9, 1, 'Системний аналіз',      'СА',    FALSE, FALSE, FALSE, '2025-09-01 00:00:00'),
    (10, 2, 'Веб-розробка',          'Веб',   FALSE, FALSE, FALSE, '2025-09-01 00:00:00'),
    (11, 1, 'Дискретна математика',  'Дискр', FALSE, FALSE, FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_discipline (academic_year_id, discipline_id)
VALUES (1, 8), (1, 9), (1, 10), (1, 11)
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Additional student groups (КН-21, ІСТ-21) — semester 2 only
-- 2023-entry groups; 3rd year (6th semester) in spring 2026
-- КН-21 is_splitting=TRUE to exercise split lessons (via КГ)
-- ─────────────────────────────────────────────────────────────

INSERT INTO education_group (id, specialty_id, name, short_name, is_contract, is_splitting,
                              education_start, education_end, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (5, 1, 'КН-21',  'КН-21',  FALSE, TRUE,  '2023-09-01 00:00:00', '2027-06-30 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (6, 2, 'ІСТ-21', 'ІСТ-21', FALSE, FALSE, '2023-09-01 00:00:00', '2027-06-30 00:00:00', FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO academic_year_to_education_group (academic_year_id, group_id)
VALUES (1, 5), (1, 6)
ON CONFLICT DO NOTHING;

INSERT INTO education_group_semester (id, group_id, semester_id, start_date, end_date, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    ( 9, 5, 2, '2026-02-02 00:00:00', '2026-06-14 00:00:00', FALSE, '2025-09-01 00:00:00'),
    (10, 6, 2, '2026-02-02 00:00:00', '2026-06-14 00:00:00', FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- Study plans for КН-21 and ІСТ-21 (semester_number=6)
--
-- 64h → 32 total lessons → 16/period → ceil(16/9)=2/week both periods
-- КН-21: КГ(2)+СА(2)+Веб(2)+Дискр(2) = 8/week both periods
-- ІСТ-21: КГ(2)+СА(2)+Веб(2)+Дискр(2) = 8/week both periods
-- ─────────────────────────────────────────────────────────────

INSERT INTO study_plan (
    id, academic_year_id, specialty_id, discipline_id, semester_number,
    classroom_work, lectures, laboratory, practical, independent_work,
    created_at, is_deleted
)
OVERRIDING SYSTEM VALUE
VALUES
    (11, 1, 1,  8, 6,  64, 32, 32,  0, 32, '2025-09-01 00:00:00', FALSE),  -- КГ    for КН
    (12, 1, 1,  9, 6,  64, 64,  0,  0, 32, '2025-09-01 00:00:00', FALSE),  -- СА    for КН
    (13, 1, 1, 10, 6,  64, 32,  0, 32, 32, '2025-09-01 00:00:00', FALSE),  -- Веб   for КН
    (14, 1, 1, 11, 6,  64, 64,  0,  0, 32, '2025-09-01 00:00:00', FALSE),  -- Дискр for КН
    (15, 1, 2,  8, 6,  64, 32, 32,  0, 32, '2025-09-01 00:00:00', FALSE),  -- КГ    for ІСТ
    (16, 1, 2,  9, 6,  64, 64,  0,  0, 32, '2025-09-01 00:00:00', FALSE),  -- СА    for ІСТ
    (17, 1, 2, 10, 6,  64, 32,  0, 32, 32, '2025-09-01 00:00:00', FALSE),  -- Веб   for ІСТ
    (18, 1, 2, 11, 6,  64, 64,  0,  0, 32, '2025-09-01 00:00:00', FALSE)   -- Дискр for ІСТ
ON CONFLICT (id) DO UPDATE SET
    classroom_work = EXCLUDED.classroom_work;


-- ─────────────────────────────────────────────────────────────
-- Workload distributions for КН-21 (group 5) and ІСТ-21 (group 6)
-- ─────────────────────────────────────────────────────────────

INSERT INTO workload_distribution (
    id, study_plan_id, group_id, classroom_work, laboratory, practical,
    created_at, is_deleted
)
OVERRIDING SYSTEM VALUE
VALUES
    -- КН-21 (group 5, is_splitting=TRUE)
    (20, 11, 5, 64, 32,  0, '2025-09-01 00:00:00', FALSE),  -- КГ
    (21, 12, 5, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- СА
    (22, 13, 5, 64,  0, 32, '2025-09-01 00:00:00', FALSE),  -- Веб
    (23, 14, 5, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Дискр
    -- ІСТ-21 (group 6)
    (24, 15, 6, 64, 32,  0, '2025-09-01 00:00:00', FALSE),  -- КГ
    (25, 16, 6, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- СА
    (26, 17, 6, 64,  0, 32, '2025-09-01 00:00:00', FALSE),  -- Веб
    (27, 18, 6, 64,  0,  0, '2025-09-01 00:00:00', FALSE)   -- Дискр
ON CONFLICT (id) DO UPDATE SET
    classroom_work = EXCLUDED.classroom_work;


-- ─────────────────────────────────────────────────────────────
-- Workload assignments for new groups
--
-- Semester 2 teacher loads (lessons/week in numerator, 2h each, max 20h/wk):
--   Teacher 1 (Петренко):   Прог KН-31(2)+Прог КН-32(2)+ФВ×2(1+1)+Веб КН-21(2) → 8 num, 4 denom → 16h/8h
--   Teacher 2 (Коваленко):  АСД×2(1+1)+ОС(1)+ФВ×2(1+1 num)+СА КН-21(2) → 7 num, 5 denom → 14h/10h
--   Teacher 3 (Шевченко):   БД×4(1+1+1+1)+Веб ІСТ-21(2)+Дискр ІСТ-21(2) → 8/wk → 16h
--   Teacher 4 (Бондаренко): Мережі×4(1+1+1+1)+КГ КН-21(2) → 6/wk → 12h
--   Teacher 5 (Марченко):   Мат×2(1+1)+Дискр КН-21(2)+КГ ІСТ-21(2)+СА ІСТ-21(2) → 8/wk → 16h
-- ─────────────────────────────────────────────────────────────

INSERT INTO workload_assignment (
    id, workload_distribution_id, teacher_id, role_type, assigned_hours, created_at
)
OVERRIDING SYSTEM VALUE
VALUES
    -- КН-21
    (20, 20, 4, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Бондаренко → КГ    / КН-21
    (21, 21, 2, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Коваленко  → СА    / КН-21
    (22, 22, 1, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Петренко   → Веб   / КН-21
    (23, 23, 5, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Марченко   → Дискр / КН-21
    -- ІСТ-21
    (24, 24, 5, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Марченко   → КГ    / ІСТ-21
    (25, 25, 5, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Марченко   → СА    / ІСТ-21
    (26, 26, 3, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Шевченко   → Веб   / ІСТ-21
    (27, 27, 3, 'lecturer', 64, '2025-09-01 00:00:00')   -- Шевченко   → Дискр / ІСТ-21
ON CONFLICT (id) DO UPDATE SET
    assigned_hours = EXCLUDED.assigned_hours;


-- ─────────────────────────────────────────────────────────────
-- Extra study plans for existing groups using new disciplines
-- (64h → 2 lessons/week in both periods)
-- ─────────────────────────────────────────────────────────────
-- ОС for КН sem 6 and Мат for ІСТ sem 6 are needed for new
-- groups КН-21 and ІСТ-21 so their weekly load reaches 10/wk.

INSERT INTO study_plan (
    id, academic_year_id, specialty_id, discipline_id, semester_number,
    classroom_work, lectures, laboratory, practical, independent_work,
    created_at, is_deleted
)
OVERRIDING SYSTEM VALUE
VALUES
    (19, 1, 1,  6, 6,  64, 64,  0,  0, 32, '2025-09-01 00:00:00', FALSE),  -- ОС   for КН  sem 6 (КН-21)
    (20, 1, 2,  1, 6,  64, 64,  0,  0, 32, '2025-09-01 00:00:00', FALSE)   -- Мат  for ІСТ sem 6 (ІСТ-21)
ON CONFLICT (id) DO UPDATE SET
    classroom_work = EXCLUDED.classroom_work;


-- ─────────────────────────────────────────────────────────────
-- Extra workload distributions so existing groups reach
-- 10+ lessons/week (≈ 2/day) for realistic schedule testing:
--
--   КН-31: + СА(2) + Веб(2) + Дискр(2) → was 7 num → now 13 num ≈ 2.6/day
--   КН-32: + СА(2) + Веб(2) + Дискр(2) → was 6 num → now 12 num ≈ 2.4/day
--   ІСТ-31: + СА(2) + Веб(2) + Дискр(2) → was 4 num → now 10 num = 2/day
--   ІСТ-32: same → 10 num = 2/day
--   КН-21: + ОС(2) → was 8 num → now 10 num = 2/day
--   ІСТ-21: + Мат(2) → was 8 num → now 10 num = 2/day
-- ─────────────────────────────────────────────────────────────

INSERT INTO workload_distribution (
    id, study_plan_id, group_id, classroom_work, laboratory, practical,
    created_at, is_deleted
)
OVERRIDING SYSTEM VALUE
VALUES
    -- СА (disc 9) — added to all 4 existing groups via study_plans 12 (КН) and 16 (ІСТ)
    (28, 12, 1, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- СА / КН-31
    (29, 12, 2, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- СА / КН-32
    (30, 16, 3, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- СА / ІСТ-31
    (31, 16, 4, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- СА / ІСТ-32
    -- Веб (disc 10) — study_plans 13 (КН) and 17 (ІСТ)
    (32, 13, 1, 64,  0, 32, '2025-09-01 00:00:00', FALSE),  -- Веб / КН-31
    (33, 13, 2, 64,  0, 32, '2025-09-01 00:00:00', FALSE),  -- Веб / КН-32
    (34, 17, 3, 64,  0, 32, '2025-09-01 00:00:00', FALSE),  -- Веб / ІСТ-31
    (35, 17, 4, 64,  0, 32, '2025-09-01 00:00:00', FALSE),  -- Веб / ІСТ-32
    -- Дискр (disc 11) — study_plans 14 (КН) and 18 (ІСТ)
    (36, 14, 1, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Дискр / КН-31
    (37, 14, 2, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Дискр / КН-32
    (38, 18, 3, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Дискр / ІСТ-31
    (39, 18, 4, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- Дискр / ІСТ-32
    -- ОС (disc 6) for КН-21 and Мат (disc 1) for ІСТ-21
    (40, 19, 5, 64,  0,  0, '2025-09-01 00:00:00', FALSE),  -- ОС  / КН-21
    (41, 20, 6, 64,  0,  0, '2025-09-01 00:00:00', FALSE)   -- Мат / ІСТ-21
ON CONFLICT (id) DO UPDATE SET
    classroom_work = EXCLUDED.classroom_work;


-- ─────────────────────────────────────────────────────────────
-- Workload assignments for extra distributions
--
-- Semester 2 peak teacher loads (numerator weeks):
--   Teacher 2 (Коваленко):  existing(7)+СА×4(8)+ОС КН-21(2) = 17 × 2h = 34h ← OK
--   Teacher 3 (Шевченко):   existing(8)+Веб×4(8)            = 16 × 2h = 32h ← OK
--   Teacher 4 (Бондаренко): existing(6)+Дискр×4(8)          = 14 × 2h = 28h ← OK
--   Teacher 5 (Марченко):   existing(8)+Мат ІСТ-21(2)       = 10 × 2h = 20h ← OK
--   (all within max_teacher_hours_per_week = 36h)
-- ─────────────────────────────────────────────────────────────

INSERT INTO workload_assignment (
    id, workload_distribution_id, teacher_id, role_type, assigned_hours, created_at
)
OVERRIDING SYSTEM VALUE
VALUES
    -- СА → Teacher 2 (Коваленко) for all 4 existing groups
    (28, 28, 2, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Коваленко → СА    / КН-31
    (29, 29, 2, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Коваленко → СА    / КН-32
    (30, 30, 2, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Коваленко → СА    / ІСТ-31
    (31, 31, 2, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Коваленко → СА    / ІСТ-32
    -- Веб → Teacher 3 (Шевченко) for all 4 existing groups
    (32, 32, 3, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Шевченко  → Веб   / КН-31
    (33, 33, 3, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Шевченко  → Веб   / КН-32
    (34, 34, 3, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Шевченко  → Веб   / ІСТ-31
    (35, 35, 3, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Шевченко  → Веб   / ІСТ-32
    -- Дискр → Teacher 4 (Бондаренко) for all 4 existing groups
    (36, 36, 4, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Бондаренко→ Дискр / КН-31
    (37, 37, 4, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Бондаренко→ Дискр / КН-32
    (38, 38, 4, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Бондаренко→ Дискр / ІСТ-31
    (39, 39, 4, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Бондаренко→ Дискр / ІСТ-32
    -- Extra disciplines for new groups
    (40, 40, 2, 'lecturer', 64, '2025-09-01 00:00:00'),  -- Коваленко → ОС    / КН-21
    (41, 41, 5, 'lecturer', 64, '2025-09-01 00:00:00')   -- Марченко  → Мат   / ІСТ-21
ON CONFLICT (id) DO UPDATE SET
    assigned_hours = EXCLUDED.assigned_hours;


-- ─────────────────────────────────────────────────────────────
-- Schedule restriction (default constraints for generation)
-- ─────────────────────────────────────────────────────────────

INSERT INTO schedule_restriction (
    id, min_group_lessons_per_day, max_group_lessons_per_day,
    max_teacher_lessons_per_day, no_gaps_in_group_schedule, created_at
)
OVERRIDING SYSTEM VALUE
VALUES (1, 2, 4, 5, TRUE, '2025-09-01 00:00:00')
ON CONFLICT (id) DO UPDATE SET
    min_group_lessons_per_day   = EXCLUDED.min_group_lessons_per_day,
    max_group_lessons_per_day   = EXCLUDED.max_group_lessons_per_day,
    max_teacher_lessons_per_day = EXCLUDED.max_teacher_lessons_per_day,
    no_gaps_in_group_schedule   = EXCLUDED.no_gaps_in_group_schedule;


-- ─────────────────────────────────────────────────────────────
-- Reset all sequences to avoid conflicts with future inserts
-- ─────────────────────────────────────────────────────────────

SELECT setval(pg_get_serial_sequence('academic_year',             'id'), (SELECT MAX(id) FROM academic_year));
SELECT setval(pg_get_serial_sequence('semester',                  'id'), (SELECT MAX(id) FROM semester));
SELECT setval(pg_get_serial_sequence('department',                'id'), (SELECT MAX(id) FROM department));
SELECT setval(pg_get_serial_sequence('room',                      'id'), (SELECT MAX(id) FROM room));
SELECT setval(pg_get_serial_sequence('specialty',                 'id'), (SELECT MAX(id) FROM specialty));
SELECT setval(pg_get_serial_sequence('cycle_committee',           'id'), (SELECT MAX(id) FROM cycle_committee));
SELECT setval(pg_get_serial_sequence('teacher',                   'id'), (SELECT MAX(id) FROM teacher));
SELECT setval(pg_get_serial_sequence('discipline',                'id'), (SELECT MAX(id) FROM discipline));
SELECT setval(pg_get_serial_sequence('education_group',           'id'), (SELECT MAX(id) FROM education_group));
SELECT setval(pg_get_serial_sequence('education_group_semester',  'id'), (SELECT MAX(id) FROM education_group_semester));
SELECT setval(pg_get_serial_sequence('bell_schedule',             'id'), (SELECT MAX(id) FROM bell_schedule));
SELECT setval(pg_get_serial_sequence('schedule_template_setting', 'id'), (SELECT MAX(id) FROM schedule_template_setting));
SELECT setval(pg_get_serial_sequence('study_plan',                'id'), (SELECT MAX(id) FROM study_plan));
SELECT setval(pg_get_serial_sequence('workload_distribution',     'id'), (SELECT MAX(id) FROM workload_distribution));
SELECT setval(pg_get_serial_sequence('workload_assignment',       'id'), (SELECT MAX(id) FROM workload_assignment));
SELECT setval(pg_get_serial_sequence('schedule_restriction',      'id'), (SELECT MAX(id) FROM schedule_restriction));
