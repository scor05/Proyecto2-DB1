package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"proyecto2/backend/api/routes"
	database "proyecto2/backend/db"
)

func enableCORS(next http.Handler) http.Handler {
	origin := getenv("FRONTEND_ORIGIN", "http://localhost:5173")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	manager, err := database.NewManager()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer manager.Close()

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, manager)

	fmt.Println("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", enableCORS(mux)); err != nil {
		log.Fatal("Error listening to port:", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
