package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	teste := 0

	mux := http.NewServeMux()

	healthHandler := func(w http.ResponseWriter, req *http.Request) {

		fmt.Fprintf(w, "Application is health %d", teste)

		teste++
	}

	mux.HandleFunc("/health", healthHandler)

	fmt.Println("Server start on 8080")

	log.Fatal(http.ListenAndServe(":8080", mux))

}
