CREATE INDEX IF NOT EXISTS idx_outbox_created_unprocessed
    ON outbox (created_at)
    WHERE published_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_outbox_lease_until
    ON outbox (lease_until)
    WHERE published_at IS NULL AND lease_until IS NULL;