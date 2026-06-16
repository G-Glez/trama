-- v0001: Core domain tables (factions)

CREATE TABLE IF NOT EXISTS factions (
    id              TEXT PRIMARY KEY,
    edition_id      TEXT NOT NULL,
    game_system_id  TEXT NOT NULL,
    name            TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(edition_id, name)
);
