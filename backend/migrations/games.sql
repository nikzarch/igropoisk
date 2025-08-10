CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    avg_rating NUMERIC(3,2) DEFAULT NULL,
    reviews_count INT DEFAULT 0
);


