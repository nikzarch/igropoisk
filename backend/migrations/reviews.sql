CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    game_id INT REFERENCES games(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id),
    rating INT NOT NULL,
    description TEXT
    CHECK (rating >= 0 AND rating <= 10)
);