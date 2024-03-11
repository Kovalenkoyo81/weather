package leveldb

import (
	"encoding/json"
	"errors"

	"github.com/Kovalenkoyo81/weather/internal/models"
	"github.com/syndtr/goleveldb/leveldb"
)

type Repository struct {
	db *leveldb.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) AddUser(user models.User) error {
	// Преобразуем пользователя в JSON
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Используем имя пользователя в качестве ключа
	key := []byte(user.Name)

	// Сохраняем данные пользователя в базе данных LevelDB
	err = r.db.Put(key, userData, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindUser(name string) bool {
	// Преобразуем имя пользователя в байтовый массив ключа
	key := []byte(name)

	// Проверяем существование пользователя в базе данных LevelDB
	_, err := r.db.Get(key, nil)
	if err != nil {
		// Если ошибка "не найдено", пользователь не существует
		if errors.Is(err, leveldb.ErrNotFound) {
			return false
		}
		// В случае другой ошибки, возникает ошибка при поиске пользователя
		return false
	}

	// Пользователь найден

	return true
}

func (r *Repository) SaveFavorite(userToken string, favorite models.Favorite) error {
	// Получаем текущий список избранных пользователя
	currentFavorites, err := r.GetFavorites(userToken)
	if err != nil {
		return err
	}

	// Добавляем новое избранное в список
	currentFavorites = append(currentFavorites, favorite)

	// Преобразуем обновленный список избранных в JSON
	updatedFavoritesData, err := json.Marshal(currentFavorites)
	if err != nil {
		return err
	}

	// Преобразуем токен пользователя в байтовый массив ключа
	key := []byte(userToken)

	// Сохраняем обновленный список избранных пользователя в базе данных LevelDB
	err = r.db.Put(key, updatedFavoritesData, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetFavorites(userToken string) ([]models.Favorite, error) {
	// Преобразуем токен пользователя в байтовый массив ключа
	key := []byte(userToken)

	// Получаем список избранных пользователя из базы данных LevelDB
	favoritesData, err := r.db.Get(key, nil)
	if err != nil {
		// Если ошибка "не найдено", возвращаем пустой список избранных
		if errors.Is(err, leveldb.ErrNotFound) {
			return []models.Favorite{}, nil
		}
		return nil, err
	}

	// Распаковываем список избранных из JSON
	var favorites []models.Favorite
	err = json.Unmarshal(favoritesData, &favorites)
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (r *Repository) DeleteFavorite(userToken, city string) error {
	// Получаем текущий список избранных пользователя
	currentFavorites, err := r.GetFavorites(userToken)
	if err != nil {
		return err
	}

	// Удаляем указанный город из списка избранных
	var updatedFavorites []models.Favorite
	for _, favorite := range currentFavorites {
		if favorite.City != city {
			updatedFavorites = append(updatedFavorites, favorite)
		}
	}

	// Преобразуем обновленный список избранных в JSON
	updatedFavoritesData, err := json.Marshal(updatedFavorites)
	if err != nil {
		return err
	}

	// Преобразуем токен пользователя в байтовый массив ключа
	key := []byte(userToken)

	// Сохраняем обновленный список избранных пользователя в базе данных LevelDB
	err = r.db.Put(key, updatedFavoritesData, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SaveToken(token string, username string) error {
	// Преобразуем токен и имя пользователя в байтовые массивы
	tokenBytes := []byte(token)
	usernameBytes := []byte(username)

	// Сохраняем токен в базе данных LevelDB с именем пользователя в качестве значения
	err := r.db.Put(tokenBytes, usernameBytes, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByToken(token string) (string, bool) {
	// Преобразуем токен в байтовый массив ключа
	tokenBytes := []byte(token)

	// Получаем имя пользователя из базы данных LevelDB по токену
	usernameBytes, err := r.db.Get(tokenBytes, nil)
	if err != nil {
		// Если ошибка "не найдено", пользователь не существует
		if errors.Is(err, leveldb.ErrNotFound) {
			return "", false
		}
		return "", false
	}

	// Пользователь найден, возвращаем имя пользователя в виде строки
	return string(usernameBytes), true
}
