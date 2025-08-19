package review

import (
	"context"
	"errors"
	"fmt"
	"igropoisk_backend/internal/game"
	"igropoisk_backend/internal/logger"
	"igropoisk_backend/internal/middleware"
)

type Service interface {
	AddReview(ctx context.Context, request AddReviewRequest) error
	GetReviewsByGameID(ctx context.Context, id int) ([]Review, error)
	GetReviewByID(ctx context.Context, id int) (*Review, error)
	RemoveReview(ctx context.Context, id int) error
}

type service struct {
	repo        Repository
	GameService game.Service
}

func NewService(repo Repository, gameService game.Service) Service {
	return &service{repo, gameService}
}

func (s *service) AddReview(ctx context.Context, request AddReviewRequest) error {
	exists, err := s.repo.IsGameReviewedByUserID(ctx, request.User.ID, request.GameID)
	if err != nil {
		logger.Logger.Error("Failed to check if review exists",
			"game_id", request.GameID,
			"user_id", request.User.ID,
			"error", err)
		return errors.New("failed to check if review exists")
	}
	if exists {
		return fmt.Errorf("review on this game by user %s already exists", request.User.Name)
	}
	review, err := NewReview(request.GameID, request.User.ID, request.Rating, request.Content)
	if err != nil {
		logger.Logger.Error("Failed to create review",
			"game_id", request.GameID,
			"user_id", request.User.ID,
			"error", err)
		return errors.New("failed to create review")
	}
	err = s.repo.AddReview(ctx, *review)
	if err != nil {
		logger.Logger.Error("Failed to add review",
			"game_id", request.GameID,
			"user_id", request.User.ID,
			"error", err)
		return errors.New("failed to add review")
	}
	return nil
}

func (s *service) GetReviewsByGameID(ctx context.Context, id int) ([]Review, error) {
	reviews, err := s.repo.GetReviewsByGameID(ctx, id)
	if err != nil {
		logger.Logger.Error("Failed to get reviews",
			"game_id", id,
			"user_id", ctx.Value(middleware.UserIDKey),
			"error", err)
		return nil, errors.New("failed to get reviews")
	}
	return reviews, nil
}

func (s *service) GetReviewByID(ctx context.Context, id int) (*Review, error) {
	review, err := s.repo.GetReviewByID(ctx, id)
	if err != nil {
		logger.Logger.Error("Failed to get review",
			"game_id", id,
			"user_id", ctx.Value(middleware.UserIDKey),
			"error", err)
		return nil, errors.New("failed to get review")
	}
	return review, nil
}

func (s *service) RemoveReview(ctx context.Context, id int) error {
	err := s.repo.RemoveReviewByID(ctx, id)
	if err != nil {
		logger.Logger.Error("Failed to remove review",
			"game_id", id,
			"user_id", ctx.Value(middleware.UserIDKey),
			"error", err)
		return errors.New("failed to remove review")
	}
	return nil
}
