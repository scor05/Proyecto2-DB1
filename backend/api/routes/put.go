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

func UpdateCategory(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := CategoryIDFromRequest(w, r)
		if !ok {
			return
		}

		var input database.CategoryWrite
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		category, err := manager.UpdateCategory(r.Context(), id, input)
		if err != nil {
			WriteDBError(w, err)
			return
		}
		WriteJSON(w, http.StatusOK, category)
	}
}

func UpdateProvider(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := ProviderIDFromRequest(w, r)
		if !ok {
			return
		}

		var input database.ProviderWrite
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		provider, err := manager.UpdateProvider(r.Context(), id, input)
		if err != nil {
			WriteDBError(w, err)
			return
		}
		WriteJSON(w, http.StatusOK, provider)
	}
}

func UpdateEmployee(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := EmployeeIDFromRequest(w, r)
		if !ok {
			return
		}

		var input database.EmployeeWrite
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		employee, err := manager.UpdateEmployee(r.Context(), id, input)
		if err != nil {
			WriteDBError(w, err)
			return
		}
		WriteJSON(w, http.StatusOK, employee)
	}
}

func UpdateClient(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := ClientIDFromRequest(w, r)
		if !ok {
			return
		}

		var input database.ClientWrite
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		client, err := manager.UpdateClient(r.Context(), id, input)
		if err != nil {
			WriteDBError(w, err)
			return
		}
		WriteJSON(w, http.StatusOK, client)
	}
}
