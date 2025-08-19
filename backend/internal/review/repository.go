package review

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/add_review.sql
var addReviewSQL string

//go:embed queries/remove_review_by_id.sql
var removeReviewByIDSQL string

//go:embed queries/get_review_by_id.sql
var getReviewByIDSQL string

//go:embed queries/get_reviews_by_game_id.sql
var getReviewByGameIDSQL string

//go:embed queries/is_game_reviewed_by_user_id.sql
var isGameReviwedByUserIDSQL string

type Repository interface {
	AddReview(ctx context.Context, review Review) error
	RemoveReviewByID(ctx context.Context, id int) error
	GetReviewByID(ctx context.Context, id int) (*Review, error)
	GetReviewsByGameID(ctx context.Context, id int) ([]Review, error)
	IsGameReviewedByUserID(ctx context.Context, userId, gameId int) (bool, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (p *PostgresRepository) AddReview(ctx context.Context, review Review) error {
	_, err := p.pool.Exec(ctx, addReviewSQL, review.GameID, review.UserID, review.Rating, review.Description)
	if err != nil {
		return fmt.Errorf("AddReview: %w", err)
	}
	return nil
}

func (p *PostgresRepository) RemoveReviewByID(ctx context.Context, id int) error {
	_, err := p.pool.Exec(ctx, removeReviewByIDSQL, id)
	if err != nil {
		return fmt.Errorf("RemoveReviewByID: %w", err)
	}
	return nil
}

func (p *PostgresRepository) GetReviewByID(ctx context.Context, id int) (*Review, error) {
	row := p.pool.QueryRow(ctx, getReviewByIDSQL, id)
	review := Review{}
	err := row.Scan(&review.ID, &review.GameID, &review.UserID, &review.Rating, &review.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("GetReviewByID: %w", err)
	}
	return &review, nil
}

func (p *PostgresRepository) GetReviewsByGameID(ctx context.Context, id int) ([]Review, error) {
	rows, err := p.pool.Query(ctx, getReviewByGameIDSQL, id)
	if err != nil {
		return nil, fmt.Errorf("GetReviewsByGameID: %w", err)
	}
	defer rows.Close()
	var reviews []Review
	for rows.Next() {
		review := Review{}
		if err := rows.Scan(&review.ID, &review.GameID, &review.UserID, &review.Rating, &review.Description); err != nil {
			return nil, fmt.Errorf("GetReviewsByGameID: %w", err)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("GetReviewsByGameID rows: %w", err)
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (p *PostgresRepository) IsGameReviewedByUserID(ctx context.Context, userId, gameId int) (bool, error) {
	var exists bool
	err := p.pool.QueryRow(ctx, isGameReviwedByUserIDSQL, userId, gameId).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("IsGameReviewedByUserID: %w", err)
	}
	return exists, nil
}
