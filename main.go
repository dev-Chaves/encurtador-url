package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generate(n int) string {
	b := make([]rune, n) //criando slice com tamanho n
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))] //pra cada indice gera uma letra aleatória do letters
	}
	return string(b)
}

type UrlDTO struct {
	Url string `json:"url"`
}

type UrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func sendResponse(w http.ResponseWriter, data any, status int) {

	w.Header().Set("Content-Type", "application/json")

	bytes, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}

	w.WriteHeader(status)

	w.Write(bytes)

}

type Server struct {
	BaseUrl string
	Hash    map[string]string
}

func (s *Server) healthHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "Application is health")

}

func (s *Server) encurtUrl(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	var url UrlDTO

	decoder := json.NewDecoder(req.Body)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&url); err != nil {
		sendResponse(w, ErrorResponse{Error: "Invalid JSON"}, 400)
		return
	}

	urlHash := generate(6)

	s.Hash[urlHash] = url.Url

	response := s.BaseUrl + urlHash

	sendResponse(w, UrlResponse{ShortUrl: response}, 200)
}

func (s *Server) accessShortUrl(w http.ResponseWriter, req *http.Request) {

	path := req.PathValue("hash")

	url := s.Hash[path]

	fmt.Println(url)

	if url == "" {
		response := "URL not found"
		sendResponse(w, ErrorResponse{Error: response}, 404)
		return
	}

	http.Redirect(w, req, url, 301)

}

func main() {

	server := &Server{
		BaseUrl: "localhost:8080/url/",
		Hash:    make(map[string]string),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", server.healthHandler)

	mux.HandleFunc("POST /encurt", server.encurtUrl)

	mux.HandleFunc("GET /url/{hash}", server.accessShortUrl)

	fmt.Println("Server start on 8080")

	log.Fatal(http.ListenAndServe(":8080", mux))

}
