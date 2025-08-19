package game

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"igropoisk_backend/internal/game/genre"
)

//go:embed queries/add_game.sql
var addGameSQL string

//go:embed queries/remove_game.sql
var removeGameSQL string

//go:embed queries/get_game_by_id.sql
var getGameByIDSQL string

//go:embed queries/get_game_by_name.sql
var getGameByNameSQL string

//go:embed queries/get_all_games.sql
var getAllGamesSQL string

type Repository interface {
	AddGame(ctx context.Context, game *Game) error
	RemoveGameByID(ctx context.Context, id int) error
	GetGameByID(ctx context.Context, id int) (*Game, error)
	GetAllGames(ctx context.Context) ([]Game, error)
	GetGameByName(ctx context.Context, name string) (*Game, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (p *PostgresRepository) GetGameByID(ctx context.Context, id int) (*Game, error) {
	game := &Game{}
	genre := &genre.Genre{}
	err := p.pool.QueryRow(ctx, getGameByIDSQL, id).Scan(
		&game.ID,
		&game.Name,
		&game.AvgRating,
		&game.ReviewsCount,
		&game.Description,
		&game.ImageURL,
		&genre.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("GetGameByID: %w", err)
	}
	return game, nil
}

func (p *PostgresRepository) GetGameByName(ctx context.Context, name string) (*Game, error) {
	game := &Game{}
	genre := &genre.Genre{}
	err := p.pool.QueryRow(ctx, getGameByNameSQL, name).Scan(
		&game.ID,
		&game.Name,
		&game.AvgRating,
		&game.ReviewsCount,
		&game.Description,
		&game.ImageURL,
		&genre.ID,
	)
	game.Genre = *genre
	if err != nil {
		return nil, fmt.Errorf("GetGameByName: %w", err)
	}
	return game, nil
}

func (p *PostgresRepository) GetAllGames(ctx context.Context) ([]Game, error) {
	rows, err := p.pool.Query(ctx, getAllGamesSQL)
	if err != nil {
		return nil, fmt.Errorf("GetAllGames: %w", err)
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		var genre genre.Genre
		if err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.AvgRating,
			&game.ReviewsCount,
			&game.Description,
			&game.ImageURL,
			&genre.ID,
			&genre.Name,
		); err != nil {
			return nil, fmt.Errorf("GetAllGames Scan: %w", err)
		}
		game.Genre = genre
		games = append(games, game)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllGames rows: %w", err)
	}
	return games, nil
}
func (p *PostgresRepository) AddGame(ctx context.Context, game *Game) error {
	_, err := p.pool.Exec(ctx, addGameSQL, game.Name, game.Description, game.ImageURL, game.Genre.ID)
	if err != nil {
		return fmt.Errorf("AddGame: %w", err)
	}
	return nil
}

func (p *PostgresRepository) RemoveGameByID(ctx context.Context, id int) error {
	_, err := p.pool.Exec(ctx, removeGameSQL, id)
	if err != nil {
		return fmt.Errorf("RemoveGameByID: %w", err)
	}
	return nil
}
