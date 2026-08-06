package handler

import (
	"net/http"

	"github.com/letrongvu/blog/internal/middleware"
	"github.com/letrongvu/blog/internal/model"
	"github.com/letrongvu/blog/internal/services"
)


type Handler struct {
	postService *services.PostService
	userService *services.UserService
}

func New(postService *services.PostService, userService *services.UserService) *Handler{
	return &Handler{
		postService: postService,
		userService: userService,
	}
}

// Add helper to get CurrentUser (private helper)
func (h *Handler) currentUser(r *http.Request) *model.User{
	uID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		return nil
	}
	user, err := h.userService.GetByID(uID)
	if err != nil {
		return nil
	}
	return user
}