package services

import (
	"crypto/rand"
	"encoding/base64"
	"sync"

	"proyecto2/backend/models"
)

const SessionCookieName = "pcfast_session"

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]models.AuthUser
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]models.AuthUser),
	}
}

func (s *SessionStore) Create(user models.AuthUser) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = user
	return token, nil
}

func (s *SessionStore) Get(token string) (*models.AuthUser, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.sessions[token]
	if !ok {
		return nil, false
	}
	return &user, true
}

func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

func randomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
