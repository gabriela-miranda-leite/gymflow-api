CREATE TABLE users (
    id                 UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name               TEXT        NOT NULL,
    email              TEXT        NOT NULL UNIQUE,
    phone              TEXT,
    password_hash      TEXT        NOT NULL,
    ideal_time_enabled BOOLEAN     NOT NULL DEFAULT false,
    occupancy_limit    TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
