// // internal/transport/middleware.go
// package rest

// import (
// 	"encoding/base64"
// 	"net/http"
// 	"strings"

// 	"github.com/Kovalenkoyo81/weather/internal/services"
// 	"github.com/gin-gonic/gin"
// )

// func TokenAuthMiddleware(service *services.Service) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
// 			return
// 		}

// 		token := strings.TrimPrefix(authHeader, "Bearer ")
// 		decodedBytes, err := base64.StdEncoding.DecodeString(token)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
// 			return
// 		}
// 		username := string(decodedBytes)

// 		exists, err := service.UserExists(c, username)
// 		if err != nil || !exists {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			return
// 		}

// 		// Можете добавить имя пользователя в контекст, если это необходимо для последующих обработчиков
// 		c.Set("username", username)

//			c.Next()
//		}
//	}
package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Kovalenkoyo81/weather/internal/services"
	"github.com/Kovalenkoyo81/weather/internal/utils"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(service *services.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		username, err := utils.ExtractUsernameFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		var username string
		if err := json.Unmarshal(decodedBytes, &username); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		exists, err := service.UserExists(username)
		if err != nil || !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
