CREATE OR REPLACE FUNCTION update_game_rating() RETURNS trigger AS $$
BEGIN
UPDATE games
SET
    reviews_count = sub.count,
    avg_rating = CASE WHEN sub.count >= 3 THEN sub.avg ELSE NULL END
    FROM (
        SELECT game_id, COUNT(*) AS count, AVG(rating) AS avg
        FROM reviews
        WHERE game_id = NEW.game_id
        GROUP BY game_id
    ) AS sub
WHERE games.id = sub.game_id;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_rating
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_game_rating();
