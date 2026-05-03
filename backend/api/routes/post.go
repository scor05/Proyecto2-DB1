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
