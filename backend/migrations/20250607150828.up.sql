-- Migration: 20250607150828
-- Created at: 2025-06-07 15:08:28

CREATE TABLE IF NOT EXISTS tournaments
(
    id          VARCHAR(255) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    start_date  VARCHAR(255),
    end_date    VARCHAR(255),
    status      VARCHAR(50)  NOT NULL
);