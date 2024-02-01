package rest

import (
	"net/http"

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

// createFavorite добавляет новую закладку для пользователя
func (r *Rest) createFavorite(c *gin.Context) {
	token := c.Param("token") // Токен пользователя используется как идентификатор

	var favorite models.Favorite
	if err := c.BindJSON(&favorite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := r.service.SaveFavorite(token, favorite); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite saved successfully"})
}

// getFavorites возвращает список закладок пользователя
func (r *Rest) getFavorites(c *gin.Context) {
	token := c.Param("token") // Токен пользователя используется как идентификатор

	favorites, err := r.service.GetFavorites(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

// deleteFavorite удаляет указанную закладку пользователя
func (r *Rest) deleteFavorite(c *gin.Context) {
	token := c.Param("token") // Токен пользователя используется как идентификатор
	city := c.Param("city")   // Параметр из URL для идентификации города

	if err := r.service.DeleteFavorite(token, city); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite deleted successfully"})
}
