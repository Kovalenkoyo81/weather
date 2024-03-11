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
	repo := memory.NewRepository()
	if config.RepoIsLevelDB {
		repo, err := leveldb.NewRepository(config.DbPath)
		if err != nil {
			log.Fatalf("Не удалось инициализировать репозиторий LevelDB: %v", err)
		}
		defer repo.Close()

	}

	service := services.New(repo)
	router := rest.NewServer(service) // Получаем *gin.Engine

	server := &http.Server{
		Addr:    ":8080",
		Handler: router, // Используем *gin.Engine как обработчик
	}

	go func() {
		// ListenAndServe всегда возвращает ошибку. ErrServerClosed возвращается при Graceful Shutdown
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	println("Server started")

	// Настройка Graceful Shutdown
	quit := make(chan os.Signal, 1)
	// Перехватываем SIGINT и SIGTERM (Ctrl+C и команды завершения)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Ожидаем сигнала
	log.Println("Shutting down server...")

	// Плавное завершение работы сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	// Закрываем репозиторий
	err := repo.Close()
	if err != nil {
		// Обработка ошибки
	}
	log.Println("Server exiting")
}
