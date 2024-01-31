package main

import (
	"github.com//kovalenkoyo81/weather//internal/transport/rest"
	"github.com//kovalenkoyo81/weather/internal/services"
	"github.com/kovalenkoyo81/weather/internal/repostories/memory"
)

func main() {
	repo := &memory.Repository{}
	service := services.New(repo)

	server := rest.NewServer(service)
	server.Run(":8080")
}
