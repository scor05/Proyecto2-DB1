package routes

import (
	"net/http"
	database "proyecto2/backend/db"
)

func IndexProducts(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := manager.Index(r.Context())
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(w, http.StatusOK, products)
	}
}

func ShowProduct(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := ProductIDFromRequest(w, r)
		if !ok {
			return
		}

		product, err := manager.Show(r.Context(), id)
		if err != nil {
			WriteDBError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, product)
	}
}

func ListCategories(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := manager.Categories(r.Context())
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}
		WriteJSON(w, http.StatusOK, categories)
	}
}

func ListProviders(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providers, err := manager.Providers(r.Context())
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}
		WriteJSON(w, http.StatusOK, providers)
	}
}
