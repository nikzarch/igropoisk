package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"igropoisk_backend/internal/auth"
)

type Service interface {
	Register(name, password string) (token string, err error)
	Login(name, password string) (token string, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(name, password string) (token string, err error) {
	pass, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return "", hashErr
	}
	user, addErr := s.repo.AddUser(name, string(pass))
	if addErr != nil {
		return "", addErr
	}
	token, err = auth.GenerateToken(user.Id, user.Name)
	return token, err
}

func (s *service) Login(name, password string) (token string, err error) {
	user, err := s.repo.GetUserByName(name)
	if err != nil {
		return "", err
	}
	passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return "", hashErr
	}
	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.PasswordHash)) != nil {
		return "", errors.New("invalid username or password")
	}

	token, err = auth.GenerateToken(user.Id, user.Name)
	return token, err
}
