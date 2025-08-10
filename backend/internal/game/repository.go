package game

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository interface {
	AddGame(game *Game) error
	RemoveGameById(id int) error
	GetGameById(id int) (*Game, error)
	GetAllGames() ([]Game, error)
	GetGameByName(name string) (*Game, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) GetGameById(id int) (*Game, error) {
	query := "SELECT id, name, avg_rating, reviews_count FROM games WHERE id = $1"
	game := &Game{}
	err := p.db.QueryRow(query, id).Scan(&game.ID, &game.Name, &game.AvgRating, &game.ReviewsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no such game")
		}
		return nil, fmt.Errorf("GetGameById: %w", err)
	}

	return game, nil
}

func (p *PostgresRepository) GetGameByName(name string) (*Game, error) {
	query := "SELECT id, name, avg_rating, reviews_count FROM games WHERE name = $1"
	game := &Game{}
	err := p.db.QueryRow(query, name).Scan(&game.ID, &game.Name, &game.AvgRating, &game.ReviewsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no such game")
		}
		return nil, fmt.Errorf("GetGameByName: %w", err)
	}

	return game, nil
}

func (p *PostgresRepository) GetAllGames() ([]Game, error) {
	query := "SELECT id, name, avg_rating, reviews_count FROM games"
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllGames: %w", err)
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game

		if err := rows.Scan(&game.ID, &game.Name, &game.AvgRating, &game.ReviewsCount); err != nil {
			return nil, fmt.Errorf("GetAllGames Scan: %w", err)
		}

		games = append(games, game)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllGames rows: %w", err)
	}
	return games, nil
}

func (p *PostgresRepository) AddGame(game *Game) error {
	query := `INSERT INTO games (name) VALUES ($1) RETURNING id`
	err := p.db.QueryRow(query, game.Name).Scan(&game.ID)
	return err
}

func (p *PostgresRepository) RemoveGameById(id int) error {
	query := "DELETE FROM games WHERE id = $1"
	_, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("RemoveGameById: %w", err)
	}
	return nil
}
