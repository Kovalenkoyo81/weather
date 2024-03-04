package rest

import (
	"time"

	"github.com/Kovalenkoyo81/weather/internal/config"
	"github.com/Kovalenkoyo81/weather/internal/services"
	"github.com/gin-contrib/cors"
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
	//Настройка CORS
	//r.Use(cors.Default())
	//Настраиваем CORS для разрешения предварительных запросов и запросов со специальными заголовками
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // Разрешить все домены или указать конкретные
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Кэшировать предварительные запросы в течение 12 часов
	}
	r.Use(cors.New(config))

	// Открытые маршруты
	r.GET("/users/:name/exists", rest.userExists) //проверка существования пользователя
	r.POST("/users", rest.createUser)             //создание пользователя
	r.POST("/login", rest.login)                  // вход пользовтеля

	// Создание группы для защищенных маршрутов с применением мидлвара аутентификации
	authorized := r.Group("/")
	authorized.Use(tokenAuthMiddleware(service))
	{

		authorized.GET("/weather/current", rest.handleCurrentWeather) // получение текущей погоды
		authorized.POST("/favorites", rest.createFavorite)            //создание закладки
		authorized.GET("/favorites", rest.getFavorites)               //получение списка закладок
		authorized.DELETE("/favorites/:city", rest.deleteFavorite)    //удаление закладки
	}

	return r
}
