// internal/repositories/memory/memory.go

package memory

import (
	"errors"

	"github.com/Kovalenkoyo81/weather/internal/models"
)

type Repository struct {
	users        []models.User
	favoritesMap map[string][]models.Favorite
	tokensMap    map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		users:        []models.User{},
		favoritesMap: make(map[string][]models.Favorite),
		tokensMap:    make(map[string]string),
	}
}

func (r *Repository) AddUser(user models.User) {
	r.users = append(r.users, user)
}

func (r *Repository) FindUser(name string) bool {
	for _, u := range r.users {
		if u.Name == name {
			return true
		}
	}
	return false
}

// Методы для работы с закладками
func (r *Repository) SaveFavorite(userToken string, favorite models.Favorite) error {
	r.favoritesMap[userToken] = append(r.favoritesMap[userToken], favorite)
	return nil
}

func (r *Repository) GetFavorites(userToken string) ([]models.Favorite, error) {
	return r.favoritesMap[userToken], nil
}

func (r *Repository) DeleteFavorite(userToken, city string) error {
	favorites, exists := r.favoritesMap[userToken]
	if !exists {
		return errors.New("no favorites found")
	}

	for i, f := range favorites {
		if f.City == city {
			r.favoritesMap[userToken] = append(favorites[:i], favorites[i+1:]...)
			return nil
		}
	}

	return errors.New("favorite not found")
}

// Метод для сохранения токена пользователя
func (r *Repository) SaveToken(token string, username string) {
	r.tokensMap[token] = username
}

// Метод для получения имени пользователя по токену
func (r *Repository) GetUserByToken(token string) (string, bool) {
	username, exists := r.tokensMap[token]
	return username, exists
}
