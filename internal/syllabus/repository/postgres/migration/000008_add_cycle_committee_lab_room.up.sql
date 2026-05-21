CREATE TABLE IF NOT EXISTS cycle_committee_lab_room (
    id                  bigserial PRIMARY KEY,
    academic_year_id    BIGINT REFERENCES academic_year(id),
    cycle_committee_id  BIGINT NOT NULL,
    room_id             BIGINT NOT NULL,
    created_at          TIMESTAMP NOT NULL,
    modified_at         TIMESTAMP,
    UNIQUE (cycle_committee_id, room_id)
);
