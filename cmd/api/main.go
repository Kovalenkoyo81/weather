// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kovalenkoyo81/weather/internal/config"
	"github.com/Kovalenkoyo81/weather/internal/repositories/leveldb"
	"github.com/Kovalenkoyo81/weather/internal/repositories/memory"
	"github.com/Kovalenkoyo81/weather/internal/services"
	"github.com/Kovalenkoyo81/weather/internal/transport/rest"
)

func main() {
	var repo services.UsersRepository
	if config.RepoIsLevelDB {
		dbRepo, err := leveldb.NewRepository(config.DbPath)
		if err != nil {
			log.Fatalf("Не удалось инициализировать репозиторий LevelDB: %v", err)
		}
		repo = dbRepo
		defer repo.Close()

	} else {
		repo = memory.NewRepository()
	}

	service := services.New(repo)
	router := rest.NewServer(service)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	println("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
