package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/readiness", handlerHealthCheck)
	mux.HandleFunc("GET /api/v1/err", handlerErrorCheck)
	corsMux := corsMiddleware(mux)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	fmt.Printf("Listening on port : %v", port)
	log.Fatal(server.ListenAndServe())
}
