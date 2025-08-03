CREATE TABLE rounds
(
    id                       VARCHAR(255) PRIMARY KEY,
    name                     VARCHAR(255) NOT NULL,
    tournament_id            VARCHAR(255) references tournaments (id) ON DELETE CASCADE ON UPDATE CASCADE,
    match_count              INT          NOT NULL,
    player_count             INT          NOT NULL,
    player_advancement_count INT          NOT NULL,
    group_size               INT          NOT NULL,
    concurrent_group_count   INT          NOT NULL,
    created_at               TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rounds_tournament_id ON rounds (tournament_id);