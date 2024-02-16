package rest

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Kovalenkoyo81/weather/internal/services"
	"github.com/gin-gonic/gin"
)

func tokenAuthMiddleware(service *services.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлечение токена из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Декодирование токена и извлечение имени пользователя
		decodedBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		var username string
		if err := json.Unmarshal(decodedBytes, &username); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Проверка существования пользователя
		exists, err := service.UserExists(username)
		if err != nil || !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Установка имени пользователя в контекст для использования в последующих обработчиках
		c.Set("username", username)
		//username, _ := c.Get("username").(string)
		c.Next()
	}
}
