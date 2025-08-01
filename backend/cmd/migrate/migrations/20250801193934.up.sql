CREATE TABLE players
(
    id            VARCHAR(255) PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    tournament_id VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_players_tournament
        FOREIGN KEY (tournament_id)
            REFERENCES tournaments (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE
);

CREATE INDEX idx_players_tournament_id ON players (tournament_id);