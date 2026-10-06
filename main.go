package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type url struct {
	hash map[string]string
}

func generate(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {

	teste := 0

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

	mux.HandleFunc("GET /health", healthHandler)

	mux.HandleFunc("POST /health", postHealthHandler)

	fmt.Println("Server start on 8080")

	log.Fatal(http.ListenAndServe(":8080", mux))

}
