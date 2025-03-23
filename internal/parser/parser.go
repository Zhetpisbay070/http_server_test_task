package parser

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

// Result — структура для хранения ответа
type Result struct {
	URL  string `json:"url"`
	Body string `json:"body"`
}

// FetchAll запрашивает данные по списку URL
func FetchAll(ctx context.Context, urls []string, maxWorkers int) ([]Result, error) {
	results := make([]Result, 0, len(urls))
	resultChan := make(chan Result, len(urls)) // Канал для результатов
	errChan := make(chan error, 1)             // Канал для передачи первой ошибки
	sem := make(chan struct{}, maxWorkers)     // Ограничение на воркеры

	var wg sync.WaitGroup

	for _, url := range urls {
		select {
		case <-ctx.Done(): // Если контекст отменен — выходим
			return nil, ctx.Err()
		case sem <- struct{}{}: // Ограничиваем количество горутин
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				defer func() { <-sem }() // Освобождаем место в семафоре

				resp, err := fetchURL(ctx, url)
				if err != nil {
					select {
					case errChan <- err: // Отправляем первую ошибку
					default:
					}
					return
				}
				resultChan <- resp
			}(url)
		}
	}

	// Горутина для ожидания завершения всех воркеров
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Собираем результаты или выходим при ошибке
	for {
		select {
		case err := <-errChan:
			return nil, err
		case res, ok := <-resultChan:
			if !ok {
				return results, nil
			}
			results = append(results, res)
		}
	}
}

// fetchURL выполняет HTTP-запрос с таймаутом 1 секунда
func fetchURL(ctx context.Context, url string) (Result, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{}, err
	}

	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код
	if resp.StatusCode != http.StatusOK {
		return Result{}, errors.New("failed to fetch " + url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	return Result{URL: url, Body: string(body)}, nil
}
