CREATE INDEX idx_immovable_registrar ON immovable_records(registrar_id);
CREATE INDEX idx_immovable_status ON immovable_records(status);
CREATE INDEX idx_immovable_woreda ON immovable_records(woreda);
CREATE INDEX idx_immovable_name ON immovable_records(name_amharic);

CREATE INDEX idx_movable_registrar ON movable_records(registrar_id);
CREATE INDEX idx_movable_status ON movable_records(status);

CREATE INDEX idx_photos_record ON record_photos(record_type, record_id);
CREATE INDEX idx_comments_record ON record_comments(record_type, record_id);
CREATE INDEX idx_history_record ON status_history(record_type, record_id);

CREATE INDEX idx_users_active ON users(is_active) WHERE deleted_at IS NULL;
