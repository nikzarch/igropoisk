package user

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"igropoisk_backend/internal/auth"
	"igropoisk_backend/internal/logger"
)

type Service interface {
	Register(ctx context.Context, name, password string) (token string, err error)
	Login(ctx context.Context, name, password string) (token string, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, name, password string) (token string, err error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("Failed to hash password",
			"username", name,
			"error", err)
		return "", errors.New("failed to hash password")
	}
	user, err := s.repo.AddUser(ctx, name, string(passwordHash))
	if err != nil {
		logger.Logger.Error("Failed to add user",
			"username", name,
			"error", err)
		return "", errors.New("failed to add user")
	}
	token, err = auth.GenerateToken(user.ID, user.Name)
	if err != nil {
		logger.Logger.Error("Failed to generate token",
			"username", name,
			"error", err)
		return "", errors.New("failed to generate token")
	}
	return token, err
}

func (s *service) Login(ctx context.Context, name, password string) (token string, err error) {
	user, err := s.repo.GetUserByName(ctx, name)
	if err != nil {
		logger.Logger.Error("Failed to get user by name",
			"username", name,
			"error", err)
		return "", errors.New("failed to get user by name")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logger.Logger.Error("Invalid password",
			"username", name,
			"error", err)
		return "", errors.New("invalid username or password")
	}

	token, err = auth.GenerateToken(user.ID, user.Name)
	if err != nil {
		logger.Logger.Error("Failed to generate token",
			"username", name,
			"error", err)
		return "", errors.New("failed to generate token")
	}
	return token, nil
}
