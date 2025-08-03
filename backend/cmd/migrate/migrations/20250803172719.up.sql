CREATE TABLE matches
(
    id       VARCHAR(255) PRIMARY KEY,
    group_id VARCHAR(255) references groups (id) ON DELETE CASCADE ON UPDATE CASCADE,
    map_name VARCHAR(255)
)