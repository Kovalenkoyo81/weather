// internal/services/users.go
package services

import (
	"context"
	"errors"

	"github.com/Kovalenkoyo81/weather/internal/models"
)

type UsersRepository interface {
	AddUser(user models.User)
	FindUser(name string) bool
	SaveFavorite(userToken string, favorite models.Favorite) error
	GetFavorites(userToken string) ([]models.Favorite, error)
	DeleteFavorite(userToken, city string) error
	SaveToken(token string, username string)
	GetUserByToken(token string) (string, bool)
}

type Service struct {
	repo UsersRepository
}

func New(repo UsersRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateNewUser(ctx context.Context, user models.User) error {
	exists, err := s.UserExists(user.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	s.repo.AddUser(user)
	return nil
}

func (s *Service) UserExists(name string) (bool, error) {
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

func (s *Service) SaveToken(token string, username string) error {
	s.repo.SaveToken(token, username)
	return nil
}

func (s *Service) GetUserByToken(token string) (string, bool) {
	return s.repo.GetUserByToken(token)
}
