CREATE TABLE academic_year_to_user
(
    academic_year_id BIGINT    NOT NULL REFERENCES academic_year_info (id),
    user_id          BIGINT    NOT NULL REFERENCES "user" (id),
    created_at       TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (academic_year_id, user_id)
);
