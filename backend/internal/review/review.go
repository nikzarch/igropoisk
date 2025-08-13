package review

import (
	"database/sql"
	"errors"
	"igropoisk_backend/internal/user"
)

type Review struct {
	ID          int            `json:"id"`
	GameId      int            `json:"game_id"`
	User        *user.User     `json:"user"`
	Score       int            `json:"score"`
	Description sql.NullString `json:"description"`
}

func NewReview(gameId int, user *user.User, score int, description string) (*Review, error) {
	if score <= 0 || score > 10 {
		return nil, errors.New("Score must be between 0 and 10")
	}
	review := &Review{GameId: gameId, User: user, Score: score}
	review.Description = sql.NullString{String: description, Valid: true}
	return review, nil
}
