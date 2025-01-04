CREATE TABLE IF NOT EXISTS users(
    id           uuid primary key default gen_random_uuid(),
    email        TEXT    NOT NULL UNIQUE,
    password    bytea    NOT NULL,
    totp_secret bytea NOT NULL UNIQUE
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps(
    id     uuid primary key default gen_random_uuid(),
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);