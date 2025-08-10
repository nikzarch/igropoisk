package game

import (
	"errors"
	"strings"
)

type Service interface {
	AddGame(Game) error
	DeleteGameById(id int) error
	GetGameById(id int) (*Game, error)
	GetGameByName(name string) (*Game, error)
	GetAllGames() ([]Game, error)
}
type service struct {
	repo Repository
}

func validateGame(game Game) (bool, error) {
	if strings.TrimSpace(game.Name) == "" {
		return false, errors.New("game name is empty")
	}
	if len([]rune(game.Name)) < 3 {
		return false, errors.New("game name must be at least 3 characters")
	}
	if len([]rune(game.Name)) > 100 {
		return false, errors.New("game name is too long")
	}
	return true, nil
}

func (s *service) AddGame(game Game) error {

	if valid, err := validateGame(game); !valid {
		return err
	}
	err := s.repo.AddGame(&game)
	return err
}

func (s *service) DeleteGameById(id int) error {
	if id <= 0 {
		return errors.New("id must be greater than zero")
	}
	err := s.repo.RemoveGameById(id)
	return err
}

func (s *service) GetGameById(id int) (*Game, error) {
	if id <= 0 {
		return nil, errors.New("id must be greater than zero")
	}
	game, err := s.repo.GetGameById(id)
	return game, err
}

func (s *service) GetGameByName(name string) (*Game, error) {
	game, err := s.repo.GetGameByName(name)
	return game, err
}

func (s *service) GetAllGames() ([]Game, error) {
	games, err := s.repo.GetAllGames()
	return games, err
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
