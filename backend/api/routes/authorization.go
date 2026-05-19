package routes

import (
	"errors"
	"net/http"

	database "proyecto2/backend/services"
)

var errForbidden = errors.New("forbidden")

func RequireRoles(manager *database.Manager, roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := SessionUserFromRequest(manager, w, r)
			if !ok {
				return
			}

			if _, ok := allowed[user.Rol]; !ok {
				WriteError(w, http.StatusForbidden, errForbidden)
				return
			}

			next(w, r)
		}
	}
}

func RequireSession(manager *database.Manager) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if _, ok := SessionUserFromRequest(manager, w, r); !ok {
				return
			}
			next(w, r)
		}
	}
}

func SessionUserFromRequest(manager *database.Manager, w http.ResponseWriter, r *http.Request) (*database.AuthUser, bool) {
	cookie, err := r.Cookie(database.SessionCookieName)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, database.ErrInvalidCredentials)
		return nil, false
	}

	user, ok := manager.SessionUser(cookie.Value)
	if !ok {
		http.SetCookie(w, expiredSessionCookie())
		WriteError(w, http.StatusUnauthorized, database.ErrInvalidCredentials)
		return nil, false
	}

	return user, true
}
