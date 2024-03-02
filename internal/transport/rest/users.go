// internal/transport/users.go
package rest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Kovalenkoyo81/weather/internal/config"
	"github.com/Kovalenkoyo81/weather/internal/models"

	"github.com/gin-gonic/gin"
)

func (r *Rest) createUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Проверка, что имя пользователя не пустое
	if user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User name cannot be empty"})
		return
	}

	// Проверка на существование пользователя
	exists, err := r.service.UserExists(user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error checking user existence"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Создание нового пользователя, поскольку он не существует
	if err := r.service.CreateNewUser(context.Background(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error creating user"})
		return
	}

	c.Status(http.StatusCreated) // Возвращаем статус 201 Created для успешно созданного пользователя
}

func (r *Rest) userExists(c *gin.Context) {
	name := c.Param("name")
	exists, err := r.service.UserExists(name)
	if err != nil {
		// Если произошла ошибка при проверке существования пользователя, возвращаем 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
		return
	}

	if !exists {
		// Если пользователь не найден, возвращаем 404 Not Found
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Если пользователь существует, возвращаем 200 OK и статус существования
	c.JSON(http.StatusOK, gin.H{"exists": true})
}

func (r *Rest) login(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	fmt.Println(loginRequest.User)
	// Проверка существования пользователя
	exists, err := r.service.UserExists(loginRequest.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not exist"})
		return
	}

	// Генерация токена
	userData := loginRequest.User
	userDataJson, _ := json.Marshal(userData)
	token := base64.StdEncoding.EncodeToString(userDataJson)

	// Сохранение токена в репозитории
	r.service.SaveToken(token, userData)

	// Возврат токена пользователю
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *Rest) handleCurrentWeather(c *gin.Context) {

	username, ok := GetUserFromContext(c)
	if !ok {
		// GetUserFromContext сам отправляет соответствующий ответ клиенту,
		// поэтому здесь нам просто нужно прервать выполнение функции.
		return
	}

	// Получаем город из query параметров запроса
	city := c.Query("city")

	if city == "" {
		// Если город не указан, получаем список избранных городов пользователя
		favorites, err := r.service.GetFavorites(c, username)
		if err != nil || len(favorites) == 0 {
			// Если у пользователя нет избранных городов, используем город по умолчанию
			city = config.DefaultCity
		} else {
			// Используем город из первой закладки
			city = favorites[0].City
		}
	}

	// Далее выполняем запрос к API погоды для получения данных по указанному городу
	weatherData, err := r.service.GetCurrentWeather(city, config.Lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current weather"})
		return
	}

	// Отправляем полученные данные о погоде клиенту
	c.JSON(http.StatusOK, weatherData)
}

func (r *Rest) createFavorite(c *gin.Context) {

	username, ok := GetUserFromContext(c)
	if !ok {
		return
	}

	var favorite models.Favorite
	if err := c.BindJSON(&favorite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Использование извлеченного имени пользователя для сохранения избранного
	if err := r.service.SaveFavorite(c, username, favorite); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite saved successfully"})
}

func (r *Rest) getFavorites(c *gin.Context) {
	username, ok := GetUserFromContext(c)
	if !ok {
		return
	}

	favorites, err := r.service.GetFavorites(c, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

func (r *Rest) deleteFavorite(c *gin.Context) {
	username, ok := GetUserFromContext(c)
	if !ok {
		// Если пользователь не найден в контексте, ответ уже отправлен функцией GetUserFromContext
		return
	}
	city := c.Param("city")
	if err := r.service.DeleteFavorite(c, username, city); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite deleted successfully"})
}

// GetUserFromContext извлекает имя пользователя из контекста запроса.
func GetUserFromContext(c *gin.Context) (string, bool) {
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found in context"})
		return "", false
	}
	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username in context is not a string"})
		return "", false
	}
	return username, true
}
