// internal/services/users.go
package services

import (
	"context"

	"github.com/Kovalenkoyo81/weather/internal/models"
)

type UsersRepository interface {
	AddUser(user models.User)
	FindUser(name string) bool
	SaveFavorite(userToken string, favorite models.Favorite) error
	GetFavorites(userToken string) ([]models.Favorite, error)
	DeleteFavorite(userToken, city string) error
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

func (s *Service) SaveFavorite(ctx context.Context, userToken string, favorite models.Favorite) error {
	return s.repo.SaveFavorite(userToken, favorite)
}

func (s *Service) GetFavorites(ctx context.Context, userToken string) ([]models.Favorite, error) {
	return s.repo.GetFavorites(userToken)
}

func (s *Service) DeleteFavorite(ctx context.Context, userToken, city string) error {
	return s.repo.DeleteFavorite(userToken, city)
}
