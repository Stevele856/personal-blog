package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/letrongvu/blog/internal/repository"
	"github.com/letrongvu/blog/internal/view"
)

// CRUD post
func (h *Handler) NewPostForm(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "post_form.html", view.PageData{})
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	published := r.FormValue("published") == "on"

	if _, err := h.postService.Create(title, content, published); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		view.Render(w, "post_form.html", view.PageData{Error: err.Error()})
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) EditPostForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.NotFound(w, r)
		return
	}

	post, err := h.postService.GetByID(id)
	if errors.Is(err, repository.ErrPostNotFound) {
		h.NotFound(w, r)
		return
	}

	if err != nil {
		log.Printf("edit post form: get by ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	view.Render(w, "post_form.html", view.PageData{Data: post})
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.NotFound(w, r)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	published := r.FormValue("published") == "on"

	if err := h.postService.Update(id, title, content, published); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		view.Render(w, "post_form.html", view.PageData{Error: err.Error()})
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.NotFound(w, r)
		return
	}

	if err := h.postService.Delete(id); err != nil {
		log.Printf("delete post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
