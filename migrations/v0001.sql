-- v0001: Core domain tables

CREATE TABLE IF NOT EXISTS factions (
    id              TEXT PRIMARY KEY,
    edition_id      TEXT NOT NULL,
    game_system_id  TEXT NOT NULL,
    name            TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(edition_id, name)
);

CREATE TABLE IF NOT EXISTS users (
    id              TEXT PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE,
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    default_faction TEXT,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wh40k_11th_matches (
    id         TEXT PRIMARY KEY,
    player1_id TEXT NOT NULL,
    player2_id TEXT NOT NULL,
    points1    INTEGER NOT NULL CHECK(points1 >= 0 AND points1 <= 100),
    points2    INTEGER NOT NULL CHECK(points2 >= 0 AND points2 <= 100),
    winner_id  TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS teams (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tournaments (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    player_id  TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tournament_players (
    tournament_id TEXT NOT NULL,
    player_id     TEXT NOT NULL,
    team_id       TEXT,
    PRIMARY KEY (tournament_id, player_id)
);
