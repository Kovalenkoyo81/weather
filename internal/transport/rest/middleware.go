// internal/transport/middleware.go
package rest

import (
	"encoding/base64"
	"net/http"

	"github.com/Kovalenkoyo81/weather/internal/services"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(service *services.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверка наличия токена
		token := c.Param("token")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API token required"})
			return
		}

		// Декодируем токен из base64
		decodedBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		username := string(decodedBytes)

		// Проверяем, существует ли пользователь
		exists, err := service.UserExists(c, username)
		if err != nil || !exists {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}
