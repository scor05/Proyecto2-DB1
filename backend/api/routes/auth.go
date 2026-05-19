package routes

import (
	"encoding/json"
	"net/http"
	"time"

	database "proyecto2/backend/services"
)

type loginRequest struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

func LoginEmployee(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input loginRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		user, err := manager.Authenticate(r.Context(), input.Correo, input.Password)
		if err != nil {
			WriteDBError(w, err)
			return
		}

		token, err := manager.CreateSession(*user)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w, sessionCookie(token, 24*time.Hour))
		WriteJSON(w, http.StatusOK, user)
	}
}

func Logout(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(database.SessionCookieName); err == nil {
			manager.DeleteSession(cookie.Value)
		}

		http.SetCookie(w, expiredSessionCookie())
		w.WriteHeader(http.StatusNoContent)
	}
}

func CurrentSession(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := SessionUserFromRequest(manager, w, r)
		if !ok {
			return
		}

		WriteJSON(w, http.StatusOK, user)
	}
}

func sessionCookie(token string, duration time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     database.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(duration.Seconds()),
	}
}

func expiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     database.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
}
