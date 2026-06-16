CREATE TABLE factions (
    id              TEXT PRIMARY KEY,
    edition_id      TEXT NOT NULL,
    game_system_id  TEXT NOT NULL,
    name            TEXT NOT NULL,
    created_at      DATETIME NOT NULL,
    updated_at      DATETIME NOT NULL,
    UNIQUE(edition_id, name)
);
