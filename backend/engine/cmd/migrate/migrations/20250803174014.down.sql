-- Step 1: Add VARCHAR columns back
ALTER TABLE tournaments
    ADD COLUMN id_old VARCHAR(255);
ALTER TABLE groups
    ADD COLUMN id_old VARCHAR(255);
ALTER TABLE groups
    ADD COLUMN round_id_old VARCHAR(255);
ALTER TABLE rounds
    ADD COLUMN id_old VARCHAR(255);
ALTER TABLE rounds
    ADD COLUMN tournament_id_old VARCHAR(255);
ALTER TABLE players
    ADD COLUMN id_old VARCHAR(255);
ALTER TABLE players
    ADD COLUMN tournament_id_old VARCHAR(255);
ALTER TABLE matches
    ADD COLUMN id_old VARCHAR(255);
ALTER TABLE placements
    ADD COLUMN id_old VARCHAR(255);
ALTER TABLE placements
    ADD COLUMN match_id_old VARCHAR(255);
ALTER TABLE placements
    ADD COLUMN player_id_old VARCHAR(255);
ALTER TABLE users
    ADD COLUMN id_old VARCHAR(255);
ALTER TABLE player_groups
    ADD COLUMN player_id_old VARCHAR(255);
ALTER TABLE player_groups
    ADD COLUMN group_id_old VARCHAR(255);

-- Step 2: Convert UUIDs back to strings
UPDATE tournaments
SET id_old = id::TEXT;
UPDATE groups
SET id_old = id::TEXT;
UPDATE groups
SET round_id_old = round_id::TEXT;
UPDATE rounds
SET id_old = id::TEXT;
UPDATE rounds
SET tournament_id_old = tournament_id::TEXT;
UPDATE players
SET id_old = id::TEXT;
UPDATE players
SET tournament_id_old = tournament_id::TEXT;
UPDATE matches
SET id_old = id::TEXT;
UPDATE placements
SET id_old = id::TEXT;
UPDATE placements
SET match_id_old = match_id::TEXT;
UPDATE placements
SET player_id_old = player_id::TEXT;
UPDATE users
SET id_old = id::TEXT;
UPDATE player_groups
SET player_id_old = player_id::TEXT;
UPDATE player_groups
SET group_id_old = group_id::TEXT;

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

-- Step 4: Drop UUID columns
ALTER TABLE tournaments
    DROP COLUMN id CASCADE;
ALTER TABLE groups
    DROP COLUMN id CASCADE;
ALTER TABLE groups
    DROP COLUMN round_id CASCADE;
ALTER TABLE rounds
    DROP COLUMN id CASCADE;
ALTER TABLE rounds
    DROP COLUMN tournament_id CASCADE;
ALTER TABLE players
    DROP COLUMN id CASCADE;
ALTER TABLE players
    DROP COLUMN tournament_id CASCADE;
ALTER TABLE matches
    DROP COLUMN id CASCADE;
ALTER TABLE placements
    DROP COLUMN id CASCADE;
ALTER TABLE placements
    DROP COLUMN match_id CASCADE;
ALTER TABLE placements
    DROP COLUMN player_id CASCADE;
ALTER TABLE users
    DROP COLUMN id CASCADE;
ALTER TABLE player_groups
    DROP COLUMN player_id CASCADE;
ALTER TABLE player_groups
    DROP COLUMN group_id CASCADE;

-- Step 5: Rename VARCHAR columns back
ALTER TABLE tournaments
    RENAME COLUMN id_old TO id;
ALTER TABLE groups
    RENAME COLUMN id_old TO id;
ALTER TABLE groups
    RENAME COLUMN round_id_old TO round_id;
ALTER TABLE rounds
    RENAME COLUMN id_old TO id;
ALTER TABLE rounds
    RENAME COLUMN tournament_id_old TO tournament_id;
ALTER TABLE players
    RENAME COLUMN id_old TO id;
ALTER TABLE players
    RENAME COLUMN tournament_id_old TO tournament_id;
ALTER TABLE matches
    RENAME COLUMN id_old TO id;
ALTER TABLE placements
    RENAME COLUMN id_old TO id;
ALTER TABLE placements
    RENAME COLUMN match_id_old TO match_id;
ALTER TABLE placements
    RENAME COLUMN player_id_old TO player_id;
ALTER TABLE users
    RENAME COLUMN id_old TO id;
ALTER TABLE player_groups
    RENAME COLUMN player_id_old TO player_id;
ALTER TABLE player_groups
    RENAME COLUMN group_id_old TO group_id;

-- Step 6: Add primary key constraints back
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

-- Step 7: Add foreign key constraints back
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