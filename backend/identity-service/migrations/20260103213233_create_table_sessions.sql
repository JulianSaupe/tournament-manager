CREATE TABLE sessions
(
    session_id      UUID      DEFAULT uuidv7() PRIMARY KEY,
    user_id         UUID      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at      TIMESTAMP NOT NULL,
    last_accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
