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
	rest := Rest{service: service}

	// Применение middleware к защищенным маршрутам
	authorized := r.Group("/users/:token")
	authorized.Use(TokenAuthMiddleware(service))
	{
		//authorized.POST("/favorites", rest.createFavorite)
		//authorized.GET("/favorites", rest.getFavorites)
		//	authorized.DELETE("/favorites/:city", rest.deleteFavorite)
	}

	// Открытые маршруты
	r.POST("/users", rest.createUser)
	r.GET("/users/:name/exists", rest.userExists)

	return r
}
