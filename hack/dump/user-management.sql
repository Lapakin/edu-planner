-- Seed data for user-management database
-- Based on dataset.sql (12 teachers from Department of Information Systems)
-- Apply after migrations with: make db-dump-up
--
-- User 1 = admin / head of department (Зеленцов Д. Г.)
-- Users 2–12 = department teachers
-- All users share the same password: admin  (bcrypt hash below)

INSERT INTO "user" (id, email, role, is_active, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (1,  'zelenczov.d@itis.edu.ua',     'admin',   TRUE, FALSE, '2025-09-01 00:00:00'),
    (2,  'solodka.n@itis.edu.ua',       'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (3,  'brychovskyy.o@itis.edu.ua',   'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (4,  'radul.o@itis.edu.ua',         'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (5,  'korotka.l@itis.edu.ua',       'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (6,  'lyashenko.o@itis.edu.ua',     'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (7,  'moskalenko.m@itis.edu.ua',    'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (8,  'serbulova.i@itis.edu.ua',     'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (9,  'ostashko.i@itis.edu.ua',      'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (10, 'anisimov.v@itis.edu.ua',      'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (11, 'shaptala.t@itis.edu.ua',      'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (12, 'naumenko.n@itis.edu.ua',      'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (13, 'ivanenko.p@itis.edu.ua',      'teacher', TRUE, FALSE, '2025-09-01 00:00:00'),
    (14, 'kovalenko.t@itis.edu.ua',     'teacher', TRUE, FALSE, '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO user_profile (user_id, first_name, last_name, patronymic, created_at)
VALUES
    (1,  'Дмитро',    'Зеленцов',     'Геннадійович',  '2025-09-01 00:00:00'),
    (2,  'Наталія',   'Солодка',      'Олександрівна', '2025-09-01 00:00:00'),
    (3,  'Олександр', 'Бричковський', 'Дмитрович',     '2025-09-01 00:00:00'),
    (4,  'Олександр', 'Радуль',       'Анатолійович',  '2025-09-01 00:00:00'),
    (5,  'Людмила',   'Коротка',      'Іванівна',      '2025-09-01 00:00:00'),
    (6,  'Олена',     'Ляшенко',      'Анатоліївна',   '2025-09-01 00:00:00'),
    (7,  'Марина',    'Москаленко',   'Леонідівна',    '2025-09-01 00:00:00'),
    (8,  'Ірина',     'Сербулова',    'Вікторівна',    '2025-09-01 00:00:00'),
    (9,  'Ірина',     'Осташко',      'Олександрівна', '2025-09-01 00:00:00'),
    (10, 'Володимир', 'Анісімов',     'Вікторович',    '2025-09-01 00:00:00'),
    (11, 'Тетяна',    'Шаптала',      'Михайлівна',    '2025-09-01 00:00:00'),
    (12, 'Наталія',   'Науменко',     'Юріївна',       '2025-09-01 00:00:00'),
    (13, 'Петро',     'Іваненко',     'Васильович',    '2025-09-01 00:00:00'),
    (14, 'Тетяна',    'Коваленко',    'Миколаївна',    '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

-- password: "admin"
INSERT INTO user_credential (user_id, password_hash, updated_at)
VALUES
    (1,  '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (2,  '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (3,  '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (4,  '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (5,  '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (6,  '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (7,  '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (8,  '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (9,  '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (10, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (11, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (12, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (13, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00'),
    (14, '$2b$10$K0ssVi2gb5aos/HoFNWC3uqWIHSQm.QQtUWpDlb48sGVDdN7gT3fC', '2025-09-01 00:00:00')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('"user"',       'id'), (SELECT MAX(id) FROM "user"));
SELECT setval(pg_get_serial_sequence('user_profile', 'id'), (SELECT MAX(id) FROM user_profile));
