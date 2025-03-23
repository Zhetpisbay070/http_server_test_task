package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestServer_Fetch(t *testing.T) {
	go func() {
		main() // Запускаем сервер
	}()
	time.Sleep(1 * time.Second) // Даем серверу время стартануть

	client := &http.Client{}
	reqBody, _ := json.Marshal(map[string]interface{}{
		"urls": []string{"https://example.com"},
	})

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/fetch", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
