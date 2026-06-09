-- Seed data for user-management database
-- Based on dataset.sql (12 teachers from Department of Information Systems)
-- Apply after migrations with: make db-dump-up
--
-- User 1 = admin / head of department
-- Users 2–12 = department teachers
-- All users share the same password: admin  (bcrypt hash below)

INSERT INTO "user" (id, email, role, is_active, is_deleted, created_at)
    OVERRIDING SYSTEM VALUE
VALUES (1, 'pachuchyi.d@itis.edu.ua', 'admin', TRUE, FALSE, '2025-09-01 00:00:00'),
       (2, 'kovalenko.o@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (3, 'tkachenko.i@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (4, 'boiko.m@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (5, 'shevchenko.s@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (6, 'melnyk.a@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (7, 'koval.m@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (8, 'lysenko.y@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (9, 'savchenko.a@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (10, 'bondarenko.v@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (11, 'moroz.t@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (12, 'kravchenko.k@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (13, 'oliinyk.v@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
       (14, 'pavlenko.d@itis.edu.ua', 'teacher', TRUE, FALSE, '2025-09-01 00:00:00') ON CONFLICT DO NOTHING;

INSERT INTO user_profile (user_id, first_name, last_name, patronymic, created_at)
VALUES (1, 'Дмитро', 'Пахучий', 'Олександрович', '2025-09-01 00:00:00'),
       (2, 'Олена', 'Коваленко', 'Іванівна', '2025-09-01 00:00:00'),
       (3, 'Іван', 'Ткаченко', 'Петрович', '2025-09-01 00:00:00'),
       (4, 'Марія', 'Бойко', 'Василівна', '2025-09-01 00:00:00'),
       (5, 'Сергій', 'Шевченко', 'Олександрович', '2025-09-01 00:00:00'),
       (6, 'Анна', 'Мельник', 'Володимирівна', '2025-09-01 00:00:00'),
       (7, 'Максим', 'Коваль', 'Андрійович', '2025-09-01 00:00:00'),
       (8, 'Юлія', 'Лисенко', 'Тарасівна', '2025-09-01 00:00:00'),
       (9, 'Андрій', 'Савченко', 'Сергійович', '2025-09-01 00:00:00'),
       (10, 'Вікторія', 'Бондаренко', 'Миколаївна', '2025-09-01 00:00:00'),
       (11, 'Тарас', 'Мороз', 'Григорович', '2025-09-01 00:00:00'),
       (12, 'Катерина', 'Кравченко', 'Михайлівна', '2025-09-01 00:00:00'),
       (13, 'Віктор', 'Олійник', 'Степанович', '2025-09-01 00:00:00'),
       (14, 'Дарина', 'Павленко', 'Юріївна', '2025-09-01 00:00:00') ON CONFLICT DO NOTHING;

-- password: "admin"
INSERT INTO user_credential (user_id, password_hash, updated_at)
VALUES (1, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (2, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (3, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (4, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (5, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (6, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (7, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (8, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (9, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (10, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (11, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (12, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (13, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
       (14, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC',
        '2025-09-01 00:00:00') ON CONFLICT DO NOTHING;

-- Academic year mirror (normally populated via SQS events from syllabus).
-- Must match the academic year(s) seeded in syllabus.sql.
INSERT INTO academic_year_info (id, is_deleted)
VALUES (1, FALSE) ON CONFLICT DO NOTHING;

-- Attach every active user to every (non-deleted) academic year.
INSERT INTO academic_year_to_user (academic_year_id, user_id)
SELECT ay.id, u.id
FROM academic_year_info ay
         CROSS JOIN "user" u
WHERE ay.is_deleted = FALSE
  AND u.is_active = TRUE
  AND u.is_deleted = FALSE ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('"user"', 'id'), (SELECT MAX(id) FROM "user"));
SELECT setval(pg_get_serial_sequence('user_profile', 'id'), (SELECT MAX(id) FROM user_profile));