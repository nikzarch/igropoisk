package genre

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/get_genre_by_id.sql
var getGenreByIDSQL string

//go:embed queries/get_genre_by_name.sql
var getGenreByNameSQL string

//go:embed queries/add_genre.sql
var addGenreSQL string

type Repository interface {
	GetGenreById(ctx context.Context, id int) (*Genre, error)
	GetGenreByName(ctx context.Context, name string) (*Genre, error)
	AddGenre(ctx context.Context, name string) (*Genre, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (p *PostgresRepository) GetGenreById(ctx context.Context, id int) (*Genre, error) {
	var g Genre
	err := p.pool.QueryRow(ctx, getGenreByIDSQL, id).Scan(&g.ID, &g.Name)
	return &g, err
}

func (p *PostgresRepository) GetGenreByName(ctx context.Context, name string) (*Genre, error) {
	var g Genre
	err := p.pool.QueryRow(ctx, getGenreByNameSQL, name).Scan(&g.ID, &g.Name)
	return &g, err
}

func (p *PostgresRepository) AddGenre(ctx context.Context, name string) (*Genre, error) {
	r := p.pool.QueryRow(ctx, addGenreSQL, name)
	var g Genre
	err := r.Scan(&g.ID)
	g.Name = name
	return &g, err
}
