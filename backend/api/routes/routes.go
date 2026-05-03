package routes

import (
	"net/http"

	database "proyecto2/backend/db"
)

func RegisterRoutes(mux *http.ServeMux, manager *database.Manager) {
	mux.HandleFunc("GET /api/health", Health(manager))
	mux.HandleFunc("POST /api/login", LoginEmployee(manager))
	mux.HandleFunc("GET /api/productos", IndexProducts(manager))
	mux.HandleFunc("GET /api/productos/{id}", ShowProduct(manager))
	mux.HandleFunc("PUT /api/productos/{id}", UpdateProduct(manager))
	mux.HandleFunc("PATCH /api/productos/{id}", PatchProduct(manager))
	mux.HandleFunc("DELETE /api/productos/{id}", DestroyProduct(manager))
	mux.HandleFunc("GET /api/categorias", ListCategories(manager))
	mux.HandleFunc("GET /api/proveedores", ListProviders(manager))
}

func Health(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		databaseStatus := "ok"

		if err := manager.Ping(r.Context()); err != nil {
			status = http.StatusServiceUnavailable
			databaseStatus = "error"
		}

		WriteJSON(w, status, map[string]string{
			"status":   "ok",
			"backend":  "ok",
			"database": databaseStatus,
		})
	}
}
