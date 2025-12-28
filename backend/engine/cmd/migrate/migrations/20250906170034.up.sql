CREATE TABLE qualifying
(
    id            UUID               DEFAULT gen_random_uuid(),
    tournament_id UUID REFERENCES tournaments (id) ON DELETE CASCADE ON UPDATE CASCADE,
    player_id     UUID REFERENCES players (id) ON DELETE CASCADE ON UPDATE CASCADE,
    time          INT                DEFAULT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX qualifying_tournament_id_idx ON qualifying (tournament_id);