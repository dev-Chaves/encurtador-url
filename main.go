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

func main() {

	BaseUrl := "localhost:8080/"

	teste := 0

	hash := make(map[string]string)

	mux := http.NewServeMux()

	healthHandler := func(w http.ResponseWriter, req *http.Request) {

		if req.Method == "POST" {
			fmt.Println("É um GET em ... ")
		}

		x := generate(10)

		fmt.Fprintf(w, "Application is health %d \n", teste)

		fmt.Fprintf(w, "Valor gerado: %v", x)

		fmt.Fprintf(w, "Rota: %v", req.URL.Path)

		teste++
	}

	postHealthHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "POST no Health \n")
	}

	encurtUrl := func(w http.ResponseWriter, req *http.Request) {

		defer req.Body.Close()

		var url UrlDTO

		decoder := json.NewDecoder(req.Body)

		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&url); err != nil {
			http.Error(w, "JSON Invalid: "+err.Error(), http.StatusBadRequest)
			return
		}

		urlHash := generate(6)

		hash[urlHash] = url.Url

		fmt.Println(url.Url)

		fmt.Println(hash)

		response := BaseUrl + urlHash

		w.Write([]byte(response))

	}

	mux.HandleFunc("GET /health", healthHandler)

	mux.HandleFunc("POST /health", postHealthHandler)

	mux.HandleFunc("POST /encurt", encurtUrl)

	fmt.Println("Server start on 8080")

	log.Fatal(http.ListenAndServe(":8080", mux))

}
