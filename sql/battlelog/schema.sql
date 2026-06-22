CREATE TABLE wh40k_11th_matches (
    id         TEXT PRIMARY KEY,
    player1_id TEXT NOT NULL,
    player2_id TEXT NOT NULL,
    points1    INTEGER NOT NULL CHECK(points1 >= 0 AND points1 <= 100),
    points2    INTEGER NOT NULL CHECK(points2 >= 0 AND points2 <= 100),
    winner_id  TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE teams (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE tournaments (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    player_id  TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE tournament_players (
    tournament_id TEXT NOT NULL,
    player_id     TEXT NOT NULL,
    team_id       TEXT,
    PRIMARY KEY (tournament_id, player_id)
);
