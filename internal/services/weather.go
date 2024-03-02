// internal/services/weather.go
package services

import (
	"encoding/json"

	"github.com/Kovalenkoyo81/weather/internal/config"
	"github.com/Kovalenkoyo81/weather/internal/models"
	"github.com/go-resty/resty/v2"
)

// GetCurrentWeather делает запрос к API погоды и возвращает погодные условия для указанного города
func (s *Service) GetCurrentWeather(city, lang string) (*models.SimplifiedWeatherResponse, error) {
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"q":    city,
			"lang": lang,
			"key":  config.ApiKey,
		}).
		Get("https://api.weatherapi.com/v1/current.json")

	if err != nil {
		return nil, err
	}

	var weatherResponse models.СurrentWeatherResponse
	err = json.Unmarshal(resp.Body(), &weatherResponse)
	if err != nil {
		return nil, err
	}

	// Преобразуем полный ответ в сокращенный формат
	simplifiedResponse := SimplifyWeatherResponse(&weatherResponse)

	return simplifiedResponse, nil
}

// SimplifyWeatherResponse преобразует полный ответ о погоде в сокращенный формат.
func SimplifyWeatherResponse(fullResponse *models.СurrentWeatherResponse) *models.SimplifiedWeatherResponse {
	return &models.SimplifiedWeatherResponse{
		Temperature: fullResponse.Current.TempC,
		Description: fullResponse.Current.Condition.Text,
		WindSpeed:   fullResponse.Current.WindKph,
	}
}
