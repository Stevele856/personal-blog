package handler

import (
	"log"
	"net/http"

	"github.com/letrongvu/blog/internal/view"
)

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request){
	posts, err := h.postService.ListAll()
	if err != nil {
		log.Printf("dashboard: list all posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := view.Render(w, "dashboard.html", view.PageData{Data: posts}); err != nil {
		log.Printf("dashboard: render: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

