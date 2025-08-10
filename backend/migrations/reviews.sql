CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    game_id INT REFERENCES games(id) ON DELETE CASCADE,
    rating NUMERIC(3,2) NOT NULL
    CHECK (rating >= 0 AND rating <= 10)
);