package handler

import (
	"log"
	"net/http"

	"github.com/letrongvu/blog/internal/view"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request){
	posts, err := h.postService.ListPublished()
	if err != nil {
		log.Printf("home: list published posts: %v", err)
		http.Error(w,"Internal server error",http.StatusInternalServerError)
		return
	}

	if err := view.Render(w, "home.html", view.PageData{Data: posts}); err != nil {
		log.Printf("home: render: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}