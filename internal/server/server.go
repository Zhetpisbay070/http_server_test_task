package server

import (
	"net/http"
	"time"
)

// NewServer создает и настраивает HTTP-сервер
func NewServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/fetch", fetchHandler) // Регистрируем хендлер

	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
