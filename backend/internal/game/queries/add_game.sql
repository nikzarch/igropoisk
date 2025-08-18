INSERT INTO games (name,description,image_url,genre_id) VALUES ($1,$2,$3,$4) RETURNING id;
