package routes

import (
	"net/http"
	database "proyecto2/backend/db"
)

func DestroyProduct(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := ProductIDFromRequest(w, r)
		if !ok {
			return
		}

		if err := manager.Destroy(r.Context(), id); err != nil {
			WriteDBError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DestroyCompra(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := CompraIDFromRequest(w, r)
		if !ok {
			return
		}

		if err := manager.DestroyCompra(r.Context(), id); err != nil {
			WriteDBError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
