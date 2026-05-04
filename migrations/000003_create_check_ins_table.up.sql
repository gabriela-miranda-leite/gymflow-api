CREATE TYPE occupancy_level AS ENUM ('empty', 'moderate', 'busy', 'packed');

CREATE TABLE check_ins (
    id         UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID            NOT NULL REFERENCES users(id),
    gym_id     UUID            NOT NULL REFERENCES gyms(id),
    occupancy  occupancy_level NOT NULL,
    created_at TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
