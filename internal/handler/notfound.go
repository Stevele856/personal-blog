package handler

import (
	"log"
	"net/http"

	"github.com/letrongvu/blog/internal/view"
)

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	if err := view.Render(w,"404.html", view.PageData{}); err != nil {
		log.Printf("not found: render: %v", err)
	}
}