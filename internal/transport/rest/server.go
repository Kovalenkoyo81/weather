package rest

import (
	"github.com/Kovalenkoyo81/weather/internal/config"
	"github.com/Kovalenkoyo81/weather/internal/services"
	"github.com/gin-gonic/gin"
)

type Rest struct {
	service *services.Service
}

func NewServer(service *services.Service) *gin.Engine {
	if config.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	rest := &Rest{service: service}

	// Открытые маршруты
	r.GET("/users/:name/exists", rest.userExists)
	r.POST("/users", rest.createUser)
	r.POST("/login", rest.login)

	// Создание группы для защищенных маршрутов с применением мидлвара аутентификации
	authorized := r.Group("/")
	authorized.Use(TokenAuthMiddleware(service))
	{

		authorized.GET("/weather/current", rest.handleCurrentWeather)
		authorized.POST("/favorites", rest.createFavorite)
		authorized.GET("/favorites", rest.getFavorites)
		authorized.DELETE("/favorites", rest.deleteFavorite)
	}

	return r
}
