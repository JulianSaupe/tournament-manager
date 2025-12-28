-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Step 1: Add new UUID columns to all tables (including groups)
ALTER TABLE tournaments
    ADD COLUMN id_new UUID DEFAULT gen_random_uuid();
ALTER TABLE groups
    ADD COLUMN id_new UUID DEFAULT gen_random_uuid();
ALTER TABLE groups
    ADD COLUMN round_id_new UUID;
ALTER TABLE rounds
    ADD COLUMN id_new UUID DEFAULT gen_random_uuid();
ALTER TABLE rounds
    ADD COLUMN tournament_id_new UUID;
ALTER TABLE players
    ADD COLUMN id_new UUID DEFAULT gen_random_uuid();
ALTER TABLE players
    ADD COLUMN tournament_id_new UUID;
ALTER TABLE matches
    ADD COLUMN id_new UUID DEFAULT gen_random_uuid();
ALTER TABLE placements
    ADD COLUMN id_new UUID DEFAULT gen_random_uuid();
ALTER TABLE placements
    ADD COLUMN match_id_new UUID;
ALTER TABLE placements
    ADD COLUMN player_id_new UUID;
ALTER TABLE users
    ADD COLUMN id_new UUID DEFAULT gen_random_uuid();
ALTER TABLE player_groups
    ADD COLUMN player_id_new UUID;
ALTER TABLE player_groups
    ADD COLUMN group_id_new UUID;

-- Step 2: Convert existing string UUIDs to UUID type
UPDATE tournaments
SET id_new = id::UUID;
UPDATE groups
SET id_new = id::UUID;
UPDATE groups
SET round_id_new = round_id::UUID
WHERE round_id IS NOT NULL;
UPDATE rounds
SET id_new = id::UUID;
UPDATE rounds
SET tournament_id_new = tournament_id::UUID
WHERE tournament_id IS NOT NULL;
UPDATE players
SET id_new = id::UUID;
UPDATE players
SET tournament_id_new = tournament_id::UUID;
UPDATE matches
SET id_new = id::UUID;
UPDATE placements
SET id_new = id::UUID;
UPDATE placements
SET match_id_new = match_id::UUID
WHERE match_id IS NOT NULL;
UPDATE placements
SET player_id_new = player_id::UUID
WHERE player_id IS NOT NULL;
UPDATE users
SET id_new = id::UUID;
UPDATE player_groups
SET player_id_new = player_id::UUID
WHERE player_id IS NOT NULL;
UPDATE player_groups
SET group_id_new = group_id::UUID
WHERE group_id IS NOT NULL;

-- Step 3: Drop foreign key constraints
ALTER TABLE groups
    DROP CONSTRAINT IF EXISTS groups_round_id_fkey;
ALTER TABLE rounds
    DROP CONSTRAINT IF EXISTS rounds_tournament_id_fkey;
ALTER TABLE players
    DROP CONSTRAINT IF EXISTS players_tournament_id_fkey;
ALTER TABLE placements
    DROP CONSTRAINT IF EXISTS placements_match_id_fkey;
ALTER TABLE placements
    DROP CONSTRAINT IF EXISTS placements_player_id_fkey;
ALTER TABLE player_groups
    DROP CONSTRAINT IF EXISTS player_groups_player_id_fkey;
ALTER TABLE player_groups
    DROP CONSTRAINT IF EXISTS player_groups_group_id_fkey;

-- Step 4: Drop old columns
ALTER TABLE tournaments
    DROP COLUMN id CASCADE;
ALTER TABLE groups
    DROP COLUMN id CASCADE;
ALTER TABLE groups
    DROP COLUMN round_id CASCADE;
ALTER TABLE rounds
    DROP COLUMN id CASCADE,
    DROP COLUMN tournament_id CASCADE;
ALTER TABLE players
    DROP COLUMN id CASCADE,
    DROP COLUMN tournament_id CASCADE;
ALTER TABLE matches
    DROP COLUMN id CASCADE;
ALTER TABLE placements
    DROP COLUMN id CASCADE,
    DROP COLUMN match_id CASCADE,
    DROP COLUMN player_id CASCADE;
ALTER TABLE users
    DROP COLUMN id CASCADE;
ALTER TABLE player_groups
    DROP COLUMN player_id CASCADE;
ALTER TABLE player_groups
    DROP COLUMN group_id CASCADE;

-- Step 5: Rename new columns to original names
ALTER TABLE tournaments
    RENAME COLUMN id_new TO id;
ALTER TABLE groups
    RENAME COLUMN id_new TO id;
ALTER TABLE groups
    RENAME COLUMN round_id_new TO round_id;
ALTER TABLE rounds
    RENAME COLUMN id_new TO id;
ALTER TABLE rounds
    RENAME COLUMN tournament_id_new TO tournament_id;
ALTER TABLE players
    RENAME COLUMN id_new TO id;
ALTER TABLE players
    RENAME COLUMN tournament_id_new TO tournament_id;
ALTER TABLE matches
    RENAME COLUMN id_new TO id;
ALTER TABLE placements
    RENAME COLUMN id_new TO id;
ALTER TABLE placements
    RENAME COLUMN match_id_new TO match_id;
ALTER TABLE placements
    RENAME COLUMN player_id_new TO player_id;
ALTER TABLE users
    RENAME COLUMN id_new TO id;
ALTER TABLE player_groups
    RENAME COLUMN player_id_new TO player_id;
ALTER TABLE player_groups
    RENAME COLUMN group_id_new TO group_id;

-- Step 6: Add primary key constraints
ALTER TABLE tournaments
    ADD PRIMARY KEY (id);
ALTER TABLE groups
    ADD PRIMARY KEY (id);
ALTER TABLE rounds
    ADD PRIMARY KEY (id);
ALTER TABLE players
    ADD PRIMARY KEY (id);
ALTER TABLE matches
    ADD PRIMARY KEY (id);
ALTER TABLE placements
    ADD PRIMARY KEY (id);
ALTER TABLE users
    ADD PRIMARY KEY (id);

-- Step 7: Add foreign key constraints (in correct order)
ALTER TABLE rounds
    ADD CONSTRAINT rounds_tournament_id_fkey
        FOREIGN KEY (tournament_id) REFERENCES tournaments (id) ON DELETE CASCADE;
ALTER TABLE groups
    ADD CONSTRAINT groups_round_id_fkey
        FOREIGN KEY (round_id) REFERENCES rounds (id) ON DELETE CASCADE;
ALTER TABLE players
    ADD CONSTRAINT players_tournament_id_fkey
        FOREIGN KEY (tournament_id) REFERENCES tournaments (id) ON DELETE CASCADE;
ALTER TABLE placements
    ADD CONSTRAINT placements_match_id_fkey
        FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE CASCADE;
ALTER TABLE placements
    ADD CONSTRAINT placements_player_id_fkey
        FOREIGN KEY (player_id) REFERENCES players (id) ON DELETE CASCADE;
ALTER TABLE player_groups
    ADD CONSTRAINT player_groups_player_id_fkey
        FOREIGN KEY (player_id) REFERENCES players (id) ON DELETE CASCADE;
ALTER TABLE player_groups
    ADD CONSTRAINT player_groups_group_id_fkey
        FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE;