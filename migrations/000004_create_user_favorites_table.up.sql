CREATE TABLE user_favorites (
    user_id    UUID        NOT NULL REFERENCES users(id),
    gym_id     UUID        NOT NULL REFERENCES gyms(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, gym_id)
);
