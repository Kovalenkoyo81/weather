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
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/Kovalenkoyo81/weather/internal/services"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(service *services.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("TokenAuthMiddleware: No Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println(token)
		decodedBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			fmt.Println("TokenAuthMiddleware: Error decoding token", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		username := string(decodedBytes)

		fmt.Printf("TokenAuthMiddleware: Decoded username: %s\n", username)

		exists, err := service.UserExists(username)
		if err != nil || !exists {
			fmt.Printf("TokenAuthMiddleware: User not found or error: %v, exists: %t\n", err, exists)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Можете добавить имя пользователя в контекст, если это необходимо для последующих обработчиков
		c.Set("username", username)

		fmt.Println("TokenAuthMiddleware: User authenticated successfully")

		c.Next()
	}
}
