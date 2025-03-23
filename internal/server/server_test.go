package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchHandler_ValidRequest(t *testing.T) {
	reqBody, _ := json.Marshal(fetchRequest{
		URLs: []string{"https://example.com"},
	})

	req, err := http.NewRequest(http.MethodPost, "/fetch", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fetchHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

func TestFetchHandler_TooManyURLs(t *testing.T) {
	// Превышаем лимит (21 URL)
	urls := make([]string, 21)
	for i := range urls {
		urls[i] = "https://example.com"
	}
	reqBody, _ := json.Marshal(fetchRequest{URLs: urls})

	req, _ := http.NewRequest(http.MethodPost, "/fetch", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fetchHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %d", rr.Code)
	}
}

func TestFetchHandler_WrongMethod(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/fetch", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fetchHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status MethodNotAllowed, got %d", rr.Code)
	}
}
