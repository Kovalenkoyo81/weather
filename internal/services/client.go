// internal/services/service.go
package services

import (
	"encoding/json"

	"github.com/Kovalenkoyo81/weather/internal/config"
	"github.com/Kovalenkoyo81/weather/internal/models"
	"github.com/go-resty/resty/v2"
)

// GetCurrentWeather делает запрос к API погоды и возвращает погодные условия для указанного города
func (s *Service) GetCurrentWeather(city, lang string) (*models.СurrentWeatherResponse, error) {

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

	return &weatherResponse, nil
}
