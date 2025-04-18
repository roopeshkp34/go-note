package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitStore(s *sessions.CookieStore) {
	Store = s
}

func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session")
		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
