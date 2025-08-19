SELECT
    game.id,
    game.name,
    game.avg_rating,
    game.reviews_count,
    game.description,
    game.image_url,
    ge.id AS genre_id,
    ge.name AS genre_name
FROM games game
         JOIN genres ge ON game.genre_id = ge.id