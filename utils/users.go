package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// ExtractUsernameFromToken декодирует токен, содержащий имя пользователя в формате JSON, закодированное в Base64,
// и возвращает извлеченное имя пользователя.
func ExtractUsernameFromToken(token string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", errors.New("invalid token format")
	}

	var username string
	if err := json.Unmarshal(decodedBytes, &username); err != nil {
		return "", errors.New("invalid token")
	}

	return username, nil
}
