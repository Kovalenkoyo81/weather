// internal/repositories/memory/memory.go

package memory

import (
	"errors"
	"sync"

	"github.com/Kovalenkoyo81/weather/internal/models"
)

type Repository struct {
	users        []models.User
	favoritesMap map[string][]models.Favorite
	tokensMap    map[string]string
	mu           sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		users:        []models.User{},
		favoritesMap: make(map[string][]models.Favorite),
		tokensMap:    make(map[string]string),
	}
}

// Close является заглушкой для закрытия ресурсов  из за того что в leveldb мы используем аналогичный
func (r *Repository) Close() error {
	// В случае памяти нет ресурсов, которые нужно освободить,
	// поэтому этот метод всегда возвращает nil в качестве ошибки.
	return nil
}

func (r *Repository) AddUser(user models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = append(r.users, user)
	return nil /// для унификации с leveldb

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
	r.mu.Lock()
	defer r.mu.Unlock()
	r.favoritesMap[userToken] = append(r.favoritesMap[userToken], favorite)
	return nil
}

func (r *Repository) GetFavorites(userToken string) ([]models.Favorite, error) {
	return r.favoritesMap[userToken], nil
}

func (r *Repository) DeleteFavorite(userToken, city string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
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
func (r *Repository) SaveToken(token string, username string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokensMap[token] = username
	return nil
}

// Метод для получения имени пользователя по токену
func (r *Repository) GetUserByToken(token string) (string, bool) {
	username, exists := r.tokensMap[token]
	return username, exists
}
