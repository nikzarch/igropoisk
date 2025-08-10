package game

const MinReviews = 3

type Game struct {
	ID           int
	Name         string
	AvgRating    *float64 // may be nil
	ReviewsCount int
}

func (g *Game) Average() *float64 {
	if g.ReviewsCount < MinReviews {
		return nil
	}
	return g.AvgRating
}
