CREATE TABLE game_systems (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE editions (
    id             TEXT PRIMARY KEY,
    game_system_id TEXT NOT NULL,
    name           TEXT NOT NULL,
    version        TEXT NOT NULL,
    created_at     DATETIME NOT NULL,
    updated_at     DATETIME NOT NULL
);

CREATE TABLE factions (
    id         TEXT PRIMARY KEY,
    edition_id TEXT NOT NULL,
    name       TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
