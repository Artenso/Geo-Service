package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Artenso/proxy-service/docs"
	"github.com/Artenso/proxy-service/internal/app"
	"github.com/Artenso/proxy-service/internal/logger"
)

// @title Geo-Service
// @version 1.0
// @description This service helps you to get full addres from its parts or coordinates

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatalf("failed to create app: %s", err)
	}

	// Создание канала для получения сигналов остановки
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		err := a.Run()
		if err != nil {
			logger.Fatalf("failed to run app: %s", err)
		}
	}()

	// Ожидание сигнала остановки
	<-sigChan

	// Создание контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Остановка сервера с использованием graceful shutdown
	log.Println("Shutting down server...")

	if err := a.Stop(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
