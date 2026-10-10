package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func newServer() Server {
	return Server{
		Mu:      sync.Mutex{},
		BaseUrl: "localhost:8080/",
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

func TestEncurtURL(t *testing.T) {

	server := newServer()

	url := UrlDTO{Url: "google.com"}

	body, err := json.Marshal(url)

	if err != nil {
		t.Errorf("Internal server error: %q", err)
	}

	req := httptest.NewRequest("POST", "/encurt", bytes.NewReader(body))

	rec := httptest.NewRecorder()

	server.encurtUrl(rec, req)

	var response UrlResponse

	err = json.Unmarshal(rec.Body.Bytes(), &response)

	if err != nil {
		t.Errorf("Error in parse JSON: %q", err)
	}

	hash, boo := strings.CutPrefix(response.ShortUrl, server.BaseUrl)

	if boo != true {
		t.Fatalf("The response dont have the hash")
	}

	valor, ok := server.Hash[hash]

	if ok != true || valor != url.Url {
		t.Fatalf("The value doesnt exist in the hash")
	}

	fmt.Printf("Hash: %q || Value: %q", hash, valor)

}
