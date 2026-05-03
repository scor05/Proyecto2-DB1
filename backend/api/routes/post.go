package routes

import (
	"encoding/json"
	"net/http"
	database "proyecto2/backend/db"
)

type loginRequest struct {
	Correo string `json:"correo"`
}

func LoginEmployee(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input loginRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		employee, err := manager.LoginEmployee(r.Context(), input.Correo)
		if err != nil {
			WriteDBError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, employee)
	}
}

func CreateProduct(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input database.ProductWrite
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		product, err := manager.Create(r.Context(), input)
		if err != nil {
			WriteDBError(w, err)
			return
		}

		WriteJSON(w, http.StatusCreated, product)
	}
}

func CreateCompra(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input database.CompraWrite
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		compra, err := manager.CreateCompra(r.Context(), input)
		if err != nil {
			WriteDBError(w, err)
			return
		}

		WriteJSON(w, http.StatusCreated, compra)
	}
}
