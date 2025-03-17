CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    update_at
    published_at TIMESTAMP,
    lease_until TIMESTAMP
);

CREATE TRIGGER on_update
    BEFORE UPDATE
    ON outbox
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();
