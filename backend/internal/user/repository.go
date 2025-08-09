package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	AddUser(name, passwordHash string) (*User, error)
	GetUserById(id int) (*User, error)
	GetUserByName(name string) (*User, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) GetUserById(id int) (*User, error) {
	query := "SELECT id,name FROM USERS WHERE ID = $1"
	user := &User{}
	row := p.db.QueryRowContext(context.Background(), query, id)
	err := row.Scan(&user.Id, &user.Name)
	return user, err
}

func (p *postgresRepository) AddUser(name, passwordHash string) (*User, error) {
	query := "INSERT INTO users (name, password_hash) VALUES ($1, $2) RETURNING id, name"
	user := User{}
	err := p.db.QueryRowContext(context.Background(), query, name, passwordHash).
		Scan(&user.Id, &user.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, errors.New("user already exists")
			}
		}
		return nil, err
	}
	return &user, nil
}

func (p *postgresRepository) GetUserByName(name string) (*User, error) {
	query := "SELECT id, name, password_hash FROM users WHERE name = $1"
	user := &User{}
	row := p.db.QueryRowContext(context.Background(), query, name)
	err := row.Scan(&user.Id, &user.Name, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no such user")
		}
		return nil, err
	}
	return user, nil
}
