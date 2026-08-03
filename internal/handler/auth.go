package handler

import (
	"net/http"
	"time"

	"github.com/letrongvu/blog/internal/middleware"
	"github.com/letrongvu/blog/internal/view"
)

// Login/Logout
func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "login.html", view.PageData{})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	token, err := h.userService.Login(username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		view.Render(w, "login.html", view.PageData{Data: "Invalid username or password"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(middleware.SessionCookieName); err == nil {
		_ = h.userService.Logout(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
