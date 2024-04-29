package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	mux.HandleFunc("GET /api/v1/users", apiConfig.handlerGetUser)

	mux.HandleFunc("POST /api/v1/feeds", apiConfig.handlerCreateFeed)
	mux.HandleFunc("GET /api/v1/feeds", apiConfig.handlerGetAllFeeds)

	mux.HandleFunc("POST /api/v1/feed_follows", apiConfig.handlerCreateFeedFollow)
	mux.HandleFunc("DELETE /api/v1/feed_follows/{feedFollowID}", apiConfig.handlerDeleteFeedFollow)
	mux.HandleFunc("GET /api/v1/feed_follows", apiConfig.handlerGetUserFeedFollows)

	mux.HandleFunc("GET /api/v1/posts", apiConfig.handlerGetUserPosts)
	corsMux := corsMiddleware(mux)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	concurrencyResults := 10
	go Scraper(dbQueries, concurrencyResults, time.Minute)
	fmt.Printf("Listening on port : %v", port)
	log.Fatal(server.ListenAndServe())
}
