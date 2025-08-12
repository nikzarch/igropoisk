package review

import (
	"database/sql"
	"errors"
	"igropoisk_backend/internal/user"
)

type Review struct {
	ID          int
	GameId      int
	User        *user.User
	Score       int
	Description sql.NullString
}

func NewReview(gameId int, user *user.User, score int, description string) (*Review, error) {
	if score <= 0 || score > 10 {
		return nil, errors.New("Score must be between 0 and 10")
	}
	review := &Review{GameId: gameId, User: user, Score: score}
	review.Description = sql.NullString{String: description, Valid: true}
	return review, nil
}
