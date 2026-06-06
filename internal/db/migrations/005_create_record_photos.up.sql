CREATE TABLE record_photos (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_type     record_type NOT NULL,
    record_id       UUID NOT NULL,
    file_path       VARCHAR(500) NOT NULL,
    file_name       VARCHAR(255),
    file_size_bytes INTEGER,
    uploaded_by     UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
