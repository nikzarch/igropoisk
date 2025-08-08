package user

import (
	"context"
	"database/sql"
)

type Repository interface {
	AddUser(name, passwordHash string) (*User, error)
	GetUserById(id int) (*User, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) GetUserById(id int) (*User, error) {
	query := "SELECT id,name FROM USERS WHERE ID = $1"
	user := &User{}
	row := p.db.QueryRowContext(context.Background(), query, id)
	err := row.Scan(&user.Id, &user.Name)
	return user, err
}

func (p *PostgresRepository) AddUser(name, passwordHash string) (*User, error) {
	query := "INSERT INTO users (name, password_hash) VALUES ($1, $2) RETURNING id, name"
	user := User{}
	err := p.db.QueryRowContext(context.Background(), query, name, passwordHash).
		Scan(&user.Id, &user.Name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
