package user

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/add_user.sql
var addUserSQL string

//go:embed queries/get_user_by_id.sql
var getUserByIDSQL string

//go:embed queries/get_user_by_name.sql
var getUserByNameSQL string

type Repository interface {
	AddUser(ctx context.Context, name, passwordHash string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByName(ctx context.Context, name string) (*User, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (p *PostgresRepository) AddUser(ctx context.Context, name, passwordHash string) (*User, error) {
	user := User{}
	err := p.pool.QueryRow(ctx, addUserSQL, name, passwordHash).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, fmt.Errorf("AddUser : %w", err)
	}
	return &user, nil
}

func (p *PostgresRepository) GetUserByID(ctx context.Context, id int) (*User, error) {
	user := &User{}
	err := p.pool.QueryRow(ctx, getUserByIDSQL, id).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID : %w", err)
	}
	return user, nil
}

func (p *PostgresRepository) GetUserByName(ctx context.Context, name string) (*User, error) {
	user := &User{}
	err := p.pool.QueryRow(ctx, getUserByNameSQL, name).Scan(&user.ID, &user.Name, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("GetUserByName : %w", err)
	}
	return user, nil
}
