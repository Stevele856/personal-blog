package handler

import "github.com/letrongvu/blog/internal/services"


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