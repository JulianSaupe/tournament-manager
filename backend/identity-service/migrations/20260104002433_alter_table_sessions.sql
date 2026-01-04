ALTER TABLE sessions
    ALTER COLUMN created_at TYPE timestamptz;
ALTER TABLE sessions
    ALTER COLUMN expires_at TYPE timestamptz;
ALTER TABLE sessions
    ALTER COLUMN last_accessed_at TYPE timestamptz;