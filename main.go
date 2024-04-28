package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kristian-Roopnarine/rss/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("HOST")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	apiConfig := apiConfig{DB: dbQueries}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/readiness", handlerHealthCheck)
	mux.HandleFunc("GET /api/v1/err", handlerErrorCheck)
	mux.HandleFunc("POST /api/v1/users", apiConfig.handlerCreateUser)
	corsMux := corsMiddleware(mux)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	fmt.Printf("Listening on port : %v", port)
	log.Fatal(server.ListenAndServe())
}
