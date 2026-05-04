CREATE TABLE gyms (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT         NOT NULL,
    address       TEXT         NOT NULL,
    latitude      NUMERIC(9,6) NOT NULL,
    longitude     NUMERIC(9,6) NOT NULL,
    rating        NUMERIC(3,2) NOT NULL DEFAULT 0,
    review_count  INTEGER      NOT NULL DEFAULT 0,
    tags          TEXT[],
    opening_hours TEXT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
