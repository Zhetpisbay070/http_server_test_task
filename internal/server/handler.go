package server

import (
	"context"
	"encoding/json"
	"http_server/internal/parser"
	"net/http"
	"time"
)

const (
	maxURLs        = 20
	requestLimit   = 100              // Лимит на 100 входящих запросов
	workerLimit    = 4                // Одновременно не более 4 исходящих запросов
	requestTimeout = 10 * time.Second // Общий таймаут обработки запроса
)

var sem = make(chan struct{}, requestLimit) // Семафор для входящих запросов

type fetchRequest struct {
	URLs []string `json:"urls"`
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	// Ограничение на 100 одновременных запросов
	select {
	case sem <- struct{}{}:
		defer func() { <-sem }()
	default:
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// Проверяем метод
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON напрямую без io.ReadAll()
	var req fetchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Проверяем количество URL
	if len(req.URLs) == 0 || len(req.URLs) > maxURLs {
		http.Error(w, "URLs count must be between 1 and 20", http.StatusBadRequest)
		return
	}

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// Запускаем обработку URL в воркер-пуле
	results, err := parser.FetchAll(ctx, req.URLs, workerLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
