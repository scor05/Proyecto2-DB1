package routes

import (
	"encoding/json"
	"net/http"
	database "proyecto2/backend/db"
)

func UpdateProduct(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := ProductIDFromRequest(w, r)
		if !ok {
			return
		}

		var input database.ProductWrite
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		product, err := manager.Update(r.Context(), id, input)
		if err != nil {
			WriteDBError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, product)
	}
}

func UpdateCompra(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := CompraIDFromRequest(w, r)
		if !ok {
			return
		}

		var input database.CompraWrite
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		compra, err := manager.UpdateCompra(r.Context(), id, input)
		if err != nil {
			WriteDBError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, compra)
	}
}
