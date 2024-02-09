// cmd/api/main.go
package main

import (
	"github.com/Kovalenkoyo81/weather/internal/repositories/memory"
	"github.com/Kovalenkoyo81/weather/internal/services"
	"github.com/Kovalenkoyo81/weather/internal/transport/rest"
)

func main() {
	repo := memory.NewRepository()
	service := services.New(repo)
	server := rest.NewServer(service)
	server.Run(":8080")
}
