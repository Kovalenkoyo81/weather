// internal/transport/users.go
package rest

import (
	"net/http"

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

	c.Status(http.StatusOK)
}

func (r *Rest) userExists(c *gin.Context) {
	name := c.Param("name")
	exists, err := r.service.UserExists(c, name)
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
	token := c.Query("token")
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
	login := c.Param("login")

	favorites, err := r.service.UserExists(c, login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}
