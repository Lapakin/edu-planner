CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS
$$
BEGIN
    PERFORM pg_notify('events_channel', row_to_json(NEW)::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
