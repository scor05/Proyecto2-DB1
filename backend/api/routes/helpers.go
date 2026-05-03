package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	database "proyecto2/backend/db"
	"strconv"
)

func ProductIDFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product id"))
		return 0, false
	}
	return id, true
}

func CompraIDFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid compra id"))
		return 0, false
	}
	return id, true
}

func CategoryIDFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid category id"))
		return 0, false
	}
	return id, true
}

func ProviderIDFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid provider id"))
		return 0, false
	}
	return id, true
}

func WriteDBError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, database.ErrNotFound):
		WriteError(w, http.StatusNotFound, err)
	case errors.Is(err, database.ErrInvalidInput):
		WriteError(w, http.StatusBadRequest, err)
	default:
		WriteError(w, http.StatusConflict, err)
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Println("error writing JSON response:", err)
	}
}
