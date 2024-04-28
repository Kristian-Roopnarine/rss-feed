package main

import "net/http"

var corsConfig = map[string]string{
	"Access-Control-ALlow-Origin":  "*",
	"Access-Control-Allow-Methods": "GET, POST, OPTIONS, PUT, DELETE",
	"Access-Control-Allow-Headers": "*",
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range corsConfig {
			w.Header().Set(key, value)
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
