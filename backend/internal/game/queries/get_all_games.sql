SELECT
    g.id,
    g.name,
    g.avg_rating,
    g.reviews_count,
    g.description,
    g.image_url,
    ge.id AS genre_id,
    ge.name AS genre_name
FROM games g
         JOIN genres ge ON g.genre_id = ge.id