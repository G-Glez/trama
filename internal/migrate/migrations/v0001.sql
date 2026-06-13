-- v0001: Core domain tables (game_systems, editions, factions)

CREATE TABLE IF NOT EXISTS game_systems (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS editions (
    id             TEXT PRIMARY KEY,
    game_system_id TEXT NOT NULL REFERENCES game_systems(id),
    name           TEXT NOT NULL,
    version        TEXT NOT NULL,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS factions (
    id         TEXT PRIMARY KEY,
    edition_id TEXT NOT NULL REFERENCES editions(id),
    name       TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
