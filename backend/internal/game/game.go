package game

import "igropoisk_backend/internal/game/genre"

const MinReviews = 3

type Game struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	AvgRating    *float64    `json:"avg_rating"` // may be nil
	ReviewsCount int         `json:"reviews_count"`
	Description  string      `json:"description"`
	ImageURL     string      `json:"image_url"`
	Genre        genre.Genre `json:"genre"`
}

func (g *Game) Average() *float64 {
	if g.ReviewsCount < MinReviews || g.AvgRating == nil {
		return nil
	}
	return g.AvgRating
}
