package user

import (
	"golang.org/x/crypto/bcrypt"
	"igropoisk_backend/internal/auth"
)

type Service struct {
	repo *PostgresRepository
}

func NewService(repo *PostgresRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(name, password string) (token string, err error) {
	pass, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return "", hashErr
	}
	user, addErr := s.repo.AddUser(name, string(pass))
	if addErr != nil {
		return "", addErr
	}
	token, err = auth.GenerateToken(user.Id, user.Name)
	return token, nil
}
