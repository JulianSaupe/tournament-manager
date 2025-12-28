CREATE TABLE groups
(
    id       VARCHAR(255) PRIMARY KEY,
    name     VARCHAR(255),
    round_id VARCHAR(255) REFERENCES rounds (id)
)