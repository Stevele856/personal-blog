package handler

import "github.com/letrongvu/blog/internal/services"


type Handler struct {
	postService *services.PostService
}

func New(postService *services.PostService) *Handler{
	return &Handler{postService: postService}
}