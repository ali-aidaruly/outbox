INSERT INTO outbox
    (id, type, payload, created_at, published_at, lease_until)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'type-1', '{"some-data-1": "data-1"}'::jsonb, '2023-01-02T15:00:00Z', NULL, NULL),
    ('00000000-0000-0000-0000-000000000002', 'type-1', '{"some-data-2": "data-2"}'::jsonb, '2023-01-02T16:00:00Z', NULL, NULL),
    ('00000000-0000-0000-0000-000000000003', 'type-3', '{"some-data-3": "data-3"}'::jsonb, '2023-01-02T17:00:00Z', NULL, NULL),
    ('00000000-0000-0000-0000-000000000004', 'default-type', '{"default-data": "data"}'::jsonb, '2023-01-02T18:00:00Z', NULL, '2023-01-02T18:00:00Z'),
    ('00000000-0000-0000-0000-000000000005', 'default-type', '{"default-data": "data"}'::jsonb, '2023-01-02T19:00:00Z', NULL, '2023-01-02T18:00:00Z'),
    ('00000000-0000-0000-0000-000000000006', 'default-type', '{"default-data": "data"}'::jsonb, '2023-01-02T20:00:00Z', NULL, '2023-01-02T18:00:00Z'),
    ('00000000-0000-0000-0000-000000000007', 'default-type', '{"default-data": "data"}'::jsonb, '2023-01-02T18:00:00Z', NULL, '2090-01-02T18:00:00Z'),
    ('00000000-0000-0000-0000-000000000008', 'default-type', '{"default-data": "data"}'::jsonb, '2023-01-02T18:00:00Z', NULL, '2090-01-02T18:00:00Z');