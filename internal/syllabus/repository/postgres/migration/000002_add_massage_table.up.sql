CREATE TABLE IF NOT EXISTS massage
(
    id          BIGSERIAL   PRIMARY KEY,
    event_id    UUID        NOT NULL,
    action_type VARCHAR     NOT NULL,
    object_type VARCHAR     NOT NULL,
    object      JSONB       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS
$$
BEGIN
    PERFORM pg_notify('events_channel', row_to_json(NEW)::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER events_notify_trigger
    AFTER INSERT
    ON massage
    FOR EACH ROW
EXECUTE FUNCTION notify_event();
