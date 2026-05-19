-- Demo seed data for user-management database
-- Apply after migrations with: make db-dump-up
--
-- Creates 5 teachers (user ids 2-6).
-- user id=1 (admin@admin.com) is already inserted by migration 000004.
-- All demo users share the same password as admin.

INSERT INTO "user" (id, email, role, is_active, is_deleted, created_at)
OVERRIDING SYSTEM VALUE
VALUES
    (2, 'petrenko.i@edu.local',   'teacher', TRUE, FALSE, '2025-01-01 00:00:00'),
    (3, 'kovalenko.m@edu.local',  'teacher', TRUE, FALSE, '2025-01-01 00:00:00'),
    (4, 'shevchenko.o@edu.local', 'teacher', TRUE, FALSE, '2025-01-01 00:00:00'),
    (5, 'bondarenko.n@edu.local', 'teacher', TRUE, FALSE, '2025-01-01 00:00:00'),
    (6, 'marchenko.d@edu.local',  'teacher', TRUE, FALSE, '2025-01-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO user_profile (user_id, first_name, last_name, patronymic, created_at)
VALUES
    (2, 'Іван',    'Петренко',   'Миколайович', '2025-01-01 00:00:00'),
    (3, 'Марія',   'Коваленко',  'Василівна',   '2025-01-01 00:00:00'),
    (4, 'Олексій', 'Шевченко',   'Іванович',    '2025-01-01 00:00:00'),
    (5, 'Наталія', 'Бондаренко', 'Петрівна',    '2025-01-01 00:00:00'),
    (6, 'Дмитро',  'Марченко',   'Олексійович', '2025-01-01 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO user_credential (user_id, password_hash, updated_at)
VALUES
    (2, '$2a$10$c19WRNx5srdpjShPjbjFiOKpUQkCIY/nf3LYuglm8B8o.w77vkGHe', '2025-01-01 00:00:00'),
    (3, '$2a$10$c19WRNx5srdpjShPjbjFiOKpUQkCIY/nf3LYuglm8B8o.w77vkGHe', '2025-01-01 00:00:00'),
    (4, '$2a$10$c19WRNx5srdpjShPjbjFiOKpUQkCIY/nf3LYuglm8B8o.w77vkGHe', '2025-01-01 00:00:00'),
    (5, '$2a$10$c19WRNx5srdpjShPjbjFiOKpUQkCIY/nf3LYuglm8B8o.w77vkGHe', '2025-01-01 00:00:00'),
    (6, '$2a$10$c19WRNx5srdpjShPjbjFiOKpUQkCIY/nf3LYuglm8B8o.w77vkGHe', '2025-01-01 00:00:00')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('"user"', 'id'),       (SELECT MAX(id) FROM "user"));
SELECT setval(pg_get_serial_sequence('user_profile', 'id'), (SELECT MAX(id) FROM user_profile));
