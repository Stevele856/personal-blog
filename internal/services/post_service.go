package services

import (
	"errors"
	"strings"

	"github.com/letrongvu/blog/internal/model"
	"github.com/letrongvu/blog/internal/repository"
)

var ErrInvalidPost = errors.New("title and content are required")

type PostService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (s *PostService) ListPublished() ([]model.Post, error) {
	return s.postRepo.ListPublished(true)
}

func (s *PostService) ListAll() ([]model.Post, error) {
	return s.postRepo.ListAll()
}

func (s *PostService) GetBySlug(slug string) (*model.Post, error) {
	return s.postRepo.GetBySlug(slug)
}

func (s *PostService) Create(title, content string, published bool) (*model.Post, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || content == "" {
		return nil, ErrInvalidPost
	}

	p := &model.Post{
		Title:     title,
		Slug:      slugify(title),
		Content:   content,
		Published: published,
	}

	if err := s.postRepo.Create(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PostService) Update(id int, title, content string, published bool) error {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || content == "" {
		return ErrInvalidPost
	}

	return s.postRepo.Update(&model.Post{
		ID: id,
		Title: title,
		Slug: slugify(title),
		Content: content,
		Published: published,
	})
}

func (s *PostService) Delete(id int) error {
	return s.postRepo.Delete(id)
}

func (s *PostService) GetByID(id int) (*model.Post, error){
	return s.postRepo.GetByID(id)
}

// Define slugify
func slugify(title string) string {
	var s strings.Builder
	lastDash := false
	for _, l := range strings.ToLower(title) {
		switch {
		case l >= 'a' && l <= 'z' || l >= '0' && l <= '9':
			s.WriteRune(l)
			lastDash = false
		default:
			if !lastDash {
				s.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(s.String(), "-")
}
