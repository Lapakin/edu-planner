ALTER TABLE schedule_template_setting ADD COLUMN IF NOT EXISTS hours_per_lesson FLOAT NOT NULL DEFAULT 2.0;
