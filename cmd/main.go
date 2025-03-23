package main

import (
	"context"
	"http_server/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Создаем сервер
	srv := server.NewServer(":8080")

	// Канал для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Ожидаем сигнала остановки
	<-stop
	log.Println("Shutting down server...")

	// Создаем контекст с таймаутом для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
