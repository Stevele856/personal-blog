package handler

import (
	"log"
	"net/http"

	"github.com/letrongvu/blog/internal/view"
)

func (h *Handler) About(w http.ResponseWriter, r *http.Request){
	if err := view.Render(w, "about.html", view.PageData{}); err != nil {
		log.Printf("about: render: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}