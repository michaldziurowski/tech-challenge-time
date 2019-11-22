package main

import (
	"log"
	"net/http"

	"github.com/michaldziurowski/tech-challenge-time/server/timetracking/infrastructure"
)

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		h.ServeHTTP(w, r)
	}
}

func main() {
	timeTrackingHandler := infrastructure.HttpHandler()

	log.Println("The timetracking server is ON : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(timeTrackingHandler)))
}
