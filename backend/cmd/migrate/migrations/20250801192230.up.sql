CREATE TABLE tournaments
(
    id           VARCHAR(255) PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    start_date   VARCHAR(255),
    end_date     VARCHAR(255),
    status       VARCHAR(50)  NOT NULL,
    player_count INT DEFAULT 0
)