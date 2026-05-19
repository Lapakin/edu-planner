-- pg_notify has an 8 KB payload limit. Sending row_to_json(NEW) breaks when
-- the schedule JSON is large. Instead, send only the event_id (UUID) and let
-- the producer fetch the full row from the DB before forwarding to SQS.
CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS
$$
BEGIN
    PERFORM pg_notify('events_channel', NEW.event_id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
