package rest

import (
	"github.com/gin-gonic/gin"
)

type Rest struct {
	service services.Service
}

func NewServer(service services.Service) *gin.Engine {
	// конфигурация Gin и другие настройки...

	rest := Rest{service: service}

	r := gin.Default()
	r.POST("/users", rest.createUser)
	r.GET("/users/:name/exists", rest.userExists)

	return r
}
