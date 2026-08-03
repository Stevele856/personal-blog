package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/letrongvu/blog/internal/repository"
	"github.com/letrongvu/blog/internal/view"
)

func (h *Handler) Post(w http.ResponseWriter, r *http.Request){
	slug := r.PathValue("slug")

	post, err := h.postService.GetBySlug(slug)
	if errors.Is(err, repository.ErrPostNotFound){
		h.NotFound(w,r)
		return
	}

	if err != nil {
		log.Printf("post: get by slug: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := view.Render(w, "post.html",view.PageData{Data: post}); err != nil {
		log.Printf("post: render: %v", err)
		http.Error(w,"Internal server error", http.StatusInternalServerError)
	}
}