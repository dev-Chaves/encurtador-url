package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func newServer() Server {
	return Server{
		Mu:      sync.Mutex{},
		BaseUrl: "localhost:8080",
		Hash:    make(map[string]string),
	}
}

func TestHealthCheck(t *testing.T) {

	server := newServer()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	rec := httptest.NewRecorder()

	server.healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Esperava status %d, recebeu %d", http.StatusOK, rec.Code)
	}

	esperado := sucessResponse{Message: "Application is Health"}

	var teste sucessResponse

	err := json.Unmarshal(rec.Body.Bytes(), &teste)

	if err != nil {
		t.Errorf("Cannot decode de JSON: %q", err)
	}

	if esperado != teste {
		t.Errorf("Esperava corpo: %q, recebeu %q", esperado, rec.Body.String())
	}

}
