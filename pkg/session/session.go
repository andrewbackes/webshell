package session

import (
	"github.com/pborman/uuid"
	"net/http"
)

const (
	cookieName = "session_token"
)

type Store struct {
	ids map[string]struct{}
}

func New() *Store {
	return &Store{
		ids: make(map[string]struct{}),
	}
}

func (s *Store) Exists(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		return false
	}
	if _, exists := s.ids[c.Value]; exists {
		return true
	}
	return false
}

func (s *Store) NewSession(w http.ResponseWriter) {
	token := uuid.New()
	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Path:  "/",
		Value: token,
	})
	s.ids[token] = struct{}{}
}

func (s *Store) Middleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Exists(r) {
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}
		next(w, r)
	}
}
