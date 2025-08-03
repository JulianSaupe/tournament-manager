CREATE TABLE placements
(
    id        VARCHAR(255) PRIMARY KEY,
    match_id  VARCHAR(255) references matches (id) ON DELETE CASCADE ON UPDATE CASCADE,
    player_id VARCHAR(255) references players (id) ON DELETE SET NULL ON UPDATE CASCADE,
    placement INT
)