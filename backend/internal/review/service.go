package review

import (
	"fmt"
	"igropoisk_backend/internal/game"
)

type Service interface {
	AddReview(request AddReviewRequest) error
	GetReviewsByGameId(id int) ([]*Review, error)
	GetReviewById(id int) (*Review, error)
	RemoveReview(id int) error
}

type service struct {
	repo        Repository
	GameService game.Service
}

func NewService(repo Repository, gameService game.Service) Service {
	return &service{repo, gameService}
}

func (s *service) AddReview(request AddReviewRequest) error {
	exists, err := s.repo.IsGameReviewedByUserId(request.User.Id, request.GameId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("review on this game by user %s already exists", request.User.Name)
	}
	review, err := NewReview(request.GameId, &request.User, request.Score, request.Content)
	if err != nil {
		return fmt.Errorf("AddReview: %w", err)
	}
	err = s.repo.AddReview(*review)
	if err != nil {
		return fmt.Errorf("AddReview: %w", err)
	}
	return nil
}

func (s *service) GetReviewsByGameId(id int) ([]*Review, error) {
	return s.repo.GetReviewsByGameId(id)
}

func (s *service) GetReviewById(id int) (*Review, error) {
	return s.repo.GetReviewById(id)
}

func (s *service) RemoveReview(id int) error {
	return s.repo.RemoveReviewById(id)
}
