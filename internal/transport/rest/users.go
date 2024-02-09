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
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := r.service.CreateNewUser(c, user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fmt.Println("User ", user, " created")
	c.Status(http.StatusOK)
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

func (r *Rest) createFavorite(c *gin.Context) {
	token := c.Param("token")
	var favorite models.Favorite
	if err := c.BindJSON(&favorite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := r.service.SaveFavorite(c, token, favorite); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite saved successfully"})
}

func (r *Rest) getFavorites(c *gin.Context) {
	token := c.Param("token")

	favorites, err := r.service.GetFavorites(c, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

func (r *Rest) deleteFavorite(c *gin.Context) {
	token := c.Param("token")
	city := c.Query("city")
	if err := r.service.DeleteFavorite(c, token, city); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite deleted successfully"})
}

func (r *Rest) handleCurrentWeather(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	//token := c.Query("token")
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

func (r *Rest) login(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

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
