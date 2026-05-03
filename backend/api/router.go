package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"proyecto2/backend/db"
)

func enableCORS(next http.Handler) http.Handler {
	origin := getenv("FRONTEND_ORIGIN", "http://localhost:5173")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	manager, err := db.NewManager()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer manager.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		databaseStatus := "ok"

		if err := manager.Ping(r.Context()); err != nil {
			status = http.StatusServiceUnavailable
			databaseStatus = "error"
		}

		writeJSON(w, status, map[string]string{
			"status":   "ok",
			"backend":  "ok",
			"database": databaseStatus,
		})
	})
	mux.HandleFunc("GET /api/productos", indexProducts(manager))
	mux.HandleFunc("GET /api/productos/{id}", showProduct(manager))
	mux.HandleFunc("PUT /api/productos/{id}", updateProduct(manager))
	mux.HandleFunc("PATCH /api/productos/{id}", patchProduct(manager))
	mux.HandleFunc("DELETE /api/productos/{id}", destroyProduct(manager))

	fmt.Println("Server listening on port 8080")
	if err := http.ListenAndServe(":8080", enableCORS(mux)); err != nil {
		log.Fatal("Error listening to port:", err)
	}
}

func indexProducts(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := manager.Index(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusOK, products)
	}
}

func showProduct(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := productIDFromRequest(w, r)
		if !ok {
			return
		}

		product, err := manager.Show(r.Context(), id)
		if err != nil {
			writeDBError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, product)
	}
}

func updateProduct(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := productIDFromRequest(w, r)
		if !ok {
			return
		}

		var input database.ProductWrite
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		product, err := manager.Update(r.Context(), id, input)
		if err != nil {
			writeDBError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, product)
	}
}

func patchProduct(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := productIDFromRequest(w, r)
		if !ok {
			return
		}

		var input database.ProductPatch
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		product, err := manager.Patch(r.Context(), id, input)
		if err != nil {
			writeDBError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, product)
	}
}

func destroyProduct(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := productIDFromRequest(w, r)
		if !ok {
			return
		}

		if err := manager.Destroy(r.Context(), id); err != nil {
			writeDBError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func productIDFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid product id"))
		return 0, false
	}
	return id, true
}

func writeDBError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, database.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, database.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusConflict, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Println("error writing JSON response:", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
