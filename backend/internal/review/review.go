package review

import (
	"database/sql"
	"errors"
)

type Review struct {
	ID          int            `json:"id"`
	GameID      int            `json:"game_id"`
	UserID      int            `json:"user_id"`
	Rating      int            `json:"rating"`
	Description sql.NullString `json:"description"`
}

func NewReview(gameID, userID, rating int, description string) (*Review, error) {
	if rating <= 0 || rating > 10 {
		return nil, errors.New("Rating must be between 0 and 10")
	}
	review := &Review{GameID: gameID, UserID: userID, Rating: rating}
	review.Description = sql.NullString{String: description, Valid: true}
	return review, nil
}
