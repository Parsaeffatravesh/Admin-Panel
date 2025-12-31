ALTER TABLE sessions
    ADD COLUMN IF NOT EXISTS rotated_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS replaced_by_token TEXT,
    ADD COLUMN IF NOT EXISTS revoked_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_sessions_revoked_at ON sessions(revoked_at);
