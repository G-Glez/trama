CREATE TABLE users (
    id              TEXT PRIMARY KEY,
    username        TEXT NOT NULL UNIQUE,
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    default_faction TEXT,
    created_at      DATETIME NOT NULL,
    updated_at      DATETIME NOT NULL
);
