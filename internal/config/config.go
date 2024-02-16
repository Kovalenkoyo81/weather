// internal/config/config.go

package config

import (
	"log"
	"os"
)

var DebugMode bool = false

const Lang string = "ru"
const DefaultCity = "rostov"

var ApiKey string

func init() {
	// Попытка получить значение переменной окружения API_KEY
	ApiKey = os.Getenv("API_KEY")
	if ApiKey == "" {
		// Если переменная окружения не найдена, логируем ошибку и завершаем работу приложения
		log.Fatal("API_KEY environment variable for weather service is not set. Terminating application.")
	}

}
