// internal/transport/users.go
package rest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "User name cannot be empty, format {name:user}"})
		return
	}

	err := r.service.CreateNewUser(c, user)
	if err != nil {
		// Проверяем, является ли ошибка "пользователь уже существует"
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		} else {
			// Для всех остальных ошибок возвращаем 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.Status(http.StatusCreated) // Возвращаем статус 201 Created для успешно созданного пользователя
}

func (r *Rest) userExists(c *gin.Context) {
	name := c.Param("name")
	exists, err := r.service.UserExists(name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
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
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	lang := config.Lang

	// Получаем город из query параметров запроса
	city := c.Query("city")

	if city == "" {
		favorites, err := r.service.GetFavorites(c, token)
		if err != nil || len(favorites) == 0 {
			// Если нет закладок, используем "rostov" как город по умолчанию
			city = config.DefaultCity
		} else {
			// Используем город из первой закладки
			city = favorites[0].City
		}
	}

	weatherData, err := r.service.GetCurrentWeather(city, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current weather"})
		return
	}

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
		return
	}

	city := c.Query("city")
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
