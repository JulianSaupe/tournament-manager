-- Migration: create_table_tournaments
-- Created at: 2025-06-07 16:35:12

CREATE TABLE IF NOT EXISTS tournaments
(
    id          VARCHAR(255) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    start_date  VARCHAR(255),
    end_date    VARCHAR(255),
    status      VARCHAR(50)  NOT NULL
);