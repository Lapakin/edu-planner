CREATE TABLE "user"
(
    id          BIGSERIAL PRIMARY KEY,
    email       VARCHAR   NOT NULL,
    role        VARCHAR   NOT NULL DEFAULT 'user',
    is_active   BOOLEAN   NOT NULL DEFAULT FALSE,
    is_deleted  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL,
    modified_at TIMESTAMP
);

CREATE TABLE user_profile
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT    REFERENCES "user" (id),
    first_name  VARCHAR   NOT NULL,
    last_name   VARCHAR   NOT NULL,
    patronymic  VARCHAR,
    created_at  TIMESTAMP NOT NULL,
    modified_at TIMESTAMP
);

CREATE TABLE user_credential
(
    user_id       BIGINT    PRIMARY KEY REFERENCES "user" (id),
    password_hash VARCHAR   NOT NULL,
    updated_at    TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE refresh_token
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT    NOT NULL REFERENCES "user" (id),
    token_hash VARCHAR   NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN   NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE academic_year_info
(
    id         BIGINT  PRIMARY KEY,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE invite_token
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT    NOT NULL REFERENCES "user" (id),
    token      VARCHAR   NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at    TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
