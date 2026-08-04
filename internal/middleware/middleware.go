package middleware

import (
	"context"
	"log"
	"net/http"
)

type ctxKey int

const userIDKey ctxKey = iota

// Get UserID from context 
func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}

const SessionCookieName = "session"

type SessionValidator interface {
	ValidateSession(token string) (int, error)
}

func RequireAuth(validator SessionValidator) func(http.Handler) http.Handler{
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				http.Redirect(w,r,"/login", http.StatusSeeOther)
				return
			}

			userID, err := validator.ValidateSession(cookie.Value)
			if err != nil {
				http.Redirect(w,r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			h.ServeHTTP(w,r.WithContext(ctx))
		})
	}
}

func Recover(h http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w,"Internal Server Error", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w,r)
	})
}