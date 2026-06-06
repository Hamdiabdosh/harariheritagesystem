CREATE TABLE status_history (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_type record_type NOT NULL,
    record_id   UUID NOT NULL,
    changed_by  UUID NOT NULL REFERENCES users(id),
    from_status VARCHAR(30),
    to_status   VARCHAR(30) NOT NULL,
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
