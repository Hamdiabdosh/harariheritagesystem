CREATE TABLE record_comments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_type  record_type NOT NULL,
    record_id    UUID NOT NULL,
    author_id    UUID NOT NULL REFERENCES users(id),
    comment_text TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
