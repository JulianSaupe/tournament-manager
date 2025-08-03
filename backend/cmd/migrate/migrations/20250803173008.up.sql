CREATE TABLE player_groups
(
    player_id VARCHAR(255) references players(id) ON DELETE CASCADE ON UPDATE CASCADE,
    group_id VARCHAR(255) references groups(id) ON DELETE CASCADE ON UPDATE CASCADE
)