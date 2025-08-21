package game

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"igropoisk_backend/internal/game/genre"
	"igropoisk_backend/internal/logger"
	"igropoisk_backend/internal/middleware"
	"strings"
)

type Service interface {
	AddGame(ctx context.Context, request AddGameRequest) error
	DeleteGameByID(ctx context.Context, id int) error
	GetGameByID(ctx context.Context, id int) (*Game, error)
	GetGameByName(ctx context.Context, name string) (*Game, error)
	GetAllGames(ctx context.Context) ([]Game, error)
	SearchGames(ctx context.Context, query string) ([]Game, error)
}
type service struct {
	gameRepo   Repository
	genreRepo  genre.Repository
	searchRepo SearchRepository
}

func NewService(gameRepo Repository, genreRepo genre.Repository, searchRepo SearchRepository) Service {
	return &service{gameRepo: gameRepo, genreRepo: genreRepo, searchRepo: searchRepo}
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
	if len(game.ImageURL) == 0 {
		return false, errors.New("game image url is empty")
	}
	return true, nil
}

func (s *service) AddGame(ctx context.Context, request AddGameRequest) error {
	genre, err := s.genreRepo.GetGenreByName(ctx, request.Genre)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			genre, err = s.genreRepo.AddGenre(ctx, request.Genre)
			if err != nil {
				logger.Logger.Error(
					"Failed to add genre",
					"genre_name", request.Genre,
					"user_id", ctx.Value(middleware.UserIDKey),
					"error", err,
				)
				return errors.New("failed to add a new genre")
			}
		} else {
			logger.Logger.Error(
				"Failed to find genre",
				"genre_name", request.Genre,
				"user_id", ctx.Value(middleware.UserIDKey),
				"error", err,
			)
			return errors.New("failed to find a genre")
		}
	}

	var game Game = Game{Name: request.Name, Description: request.Description, ImageURL: request.ImageURL, Genre: *genre}

	if valid, err := validateGame(game); !valid {
		return err
	}

	name := strings.TrimSpace(game.Name)
	if len(name) > 0 {
		name = strings.ToLower(name)
		name = strings.ToUpper(string(name[0])) + name[1:]
	}
	game.Name = name
	err = s.gameRepo.AddGame(ctx, &game)
	if err != nil {
		logger.Logger.Error(
			"Failed to add a new game",
			"game_id", game.ID,
			"game_name", game.Name,
			"user_id", ctx.Value(middleware.UserIDKey),
			"error", err,
		)
		return errors.New("failed to add a new game")
	}
	if err := s.searchRepo.IndexGame(ctx, &game); err != nil {
		logger.Logger.Warn("Failed to add a new game to search repo",
			"game_id", game.ID,
			"user_id", ctx.Value(middleware.UserIDKey),
			"error", err)
	}
	return nil
}

func (s *service) DeleteGameByID(ctx context.Context, id int) error {

	if id <= 0 {
		logger.Logger.Error("Invalid game id",
			"game_id", id,
			"user_id", ctx.Value(middleware.UserIDKey),
		)
		return errors.New("id must be greater than zero")
	}
	err := s.gameRepo.RemoveGameByID(ctx, id)
	if err != nil {
		logger.Logger.Error("Failed to remove a game",
			"game_id", id,
			"user_id", ctx.Value(middleware.UserIDKey),
			"error", err,
		)
		return errors.New("failed to remove a game")
	}

	if err := s.searchRepo.DeleteGame(ctx, id); err != nil {
		logger.Logger.Error("Failed to remove a game",
			"game_id", id,
			"user_id", ctx.Value(middleware.UserIDKey),
			"error", err)
	}
	return nil
}

func (s *service) GetGameByID(ctx context.Context, id int) (*Game, error) {
	if id <= 0 {
		logger.Logger.Error("Invalid game id",
			"game_id", id,
			"user_id", ctx.Value(middleware.UserIDKey),
		)
		return nil, errors.New("id must be greater than zero")
	}
	game, err := s.gameRepo.GetGameByID(ctx, id)
	if err != nil {
		logger.Logger.Error("Failed to get a game",
			"game_id", id,
			"user_id", ctx.Value(middleware.UserIDKey),
			"error", err)
		return nil, errors.New("failed to get a game")
	}
	return game, nil
}

func (s *service) GetGameByName(ctx context.Context, name string) (*Game, error) {
	if name == "" {
		logger.Logger.Error("Invalid game name",
			"game_name", name,
			"user_id", ctx.Value(middleware.UserIDKey))
	}
	game, err := s.gameRepo.GetGameByName(ctx, name)
	if err != nil {
		logger.Logger.Error("Failed to get a game",
			"game_name", name,
			"error", err)
		return nil, errors.New("failed to get a game")
	}
	return game, nil
}

func (s *service) GetAllGames(ctx context.Context) ([]Game, error) {
	games, err := s.gameRepo.GetAllGames(ctx)
	if err != nil {
		logger.Logger.Error("Failed to get all games",
			"user_id", ctx.Value(middleware.UserIDKey),
			"error", err)
		return nil, errors.New("failed to get all games")
	}
	return games, nil
}

func (s *service) SearchGames(ctx context.Context, query string) ([]Game, error) {
	games, err := s.searchRepo.SearchGames(ctx, query)
	if err != nil {
		logger.Logger.Error("Failed to search games",
			"user_id", ctx.Value(middleware.UserIDKey),
			"query", query,
			"error", err)
		return nil, errors.New("failed to search games")
	}
	return games, nil
}
