package review

import (
	"database/sql"
	"fmt"
	"igropoisk_backend/internal/user"
)

type Repository interface {
	AddReview(Review) error
	RemoveReviewById(id int) error
	GetReviewById(id int) (*Review, error)
	GetReviewsByGameId(id int) ([]*Review, error)
	IsGameReviewedByUserId(userId, gameId int) (bool, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) AddReview(review Review) error {
	query := "INSERT INTO REVIEWS(game_id,rating,description,user_id) VALUES ($1,$2,$3,$4);"
	_, err := p.db.Exec(query, review.GameId, review.Score, review.Description, review.User.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) RemoveReviewById(id int) error {
	query := "DELETE FROM reviews WHERE id = $1"
	_, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("RemoveReviewById: %w", err)
	}
	return nil
}

func (p *PostgresRepository) GetReviewById(id int) (*Review, error) {
	query := "SELECT id,game_id,description,user_id FROM REVIEWS WHERE id = $1"
	row := p.db.QueryRow(query, id)
	review := Review{User: &user.User{}}
	err := row.Scan(&review.Id, &review.GameId, &review.Score, &review.Description, &review.User.Id)
	if err != nil {
		return nil, fmt.Errorf("GetReviewById: %w", err)
	}
	return &review, nil
}

func (p *PostgresRepository) GetReviewsByGameId(id int) ([]*Review, error) {
	query := "SELECT id,game_id,rating,description,user_id FROM REVIEWS WHERE game_id = $1"
	rows, err := p.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("GetReviewsByGame: %w", err)
	}
	reviews := make([]*Review, 0)
	for rows.Next() {
		review := Review{User: &user.User{}}
		err := rows.Scan(&review.Id, &review.GameId, &review.Score, &review.Description, &review.User.Id)
		if err != nil {
			return nil, fmt.Errorf("GetReviewsByGame: %w", err)
		}
		reviews = append(reviews, &review)

	}
	return reviews, nil
}

func (p *PostgresRepository) IsGameReviewedByUserId(userId, gameId int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM reviews WHERE user_id = $1 AND game_id = $2)"
	var exists bool
	err := p.db.QueryRow(query, userId, gameId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("GetReviewsByUserId: %w", err)
	}
	return exists, nil
}
