package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"strings"

	"github.com/gin-gonic/gin"
)

func UserAuthorizator(c *gin.Context) (username string, err error) {
	// Проверка наличия заголовка Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}

	// Проверка формата заголовка (должен начинаться с "Bearer ")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("Invalid token format")
	}

	// Извлечение токена из заголовка
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Декодирование токена из формата Base64 в JSON
	decodedBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", errors.New("Invalid token format")
	}

	// Распаковка JSON и извлечение имени пользователя
	if err := json.Unmarshal(decodedBytes, &username); err != nil {
		return "", errors.New("Invalid token")
	}

	// Возврат имени пользователя в случае успеха
	return username, nil
}
