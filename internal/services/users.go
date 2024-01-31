package services

import (
	"context"

	"github.com/Kovalenkoyo81/weather/internal/models"
)

type UsersRepository interface {
	AddUser(user models.User)
	FindUser(name string) bool
}

type Service struct {
	repo UsersRepository
}

func New(repo UsersRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateNewUser(ctx context.Context, user models.User) error {
	s.repo.AddUser(user)
	return nil
}

func (s *Service) UserExists(ctx context.Context, name string) (bool, error) {
	return s.repo.FindUser(name), nil
}
