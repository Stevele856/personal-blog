package repository

import (
	"database/sql"
	"errors"

	"github.com/letrongvu/blog/internal/model"
)

var ErrPostNotFound = errors.New("post not found")

type PostRepository interface {
	ListPublished(publishedOnly bool) ([]model.Post, error)
	// Get all posts (published + draft)
	ListAll() ([]model.Post, error)
	GetBySlug(slug string) (*model.Post, error)
	Create(p *model.Post) error
	Update(p *model.Post) error
	Delete(id int) error
	//080326 Add GetByID - edit form need an post ID to Populating data into a post_form.html
	GetByID(id int) (*model.Post, error)
}

type sqlitePostRepository struct {
	db *sql.DB
}

func (r *sqlitePostRepository) ListPublished(publishedOnly bool) ([]model.Post, error) {
	rows, err := r.db.Query("SELECT id, title, slug, content, published, created_at, updated_at FROM posts WHERE published = ?", publishedOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	// Interate each rows
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *sqlitePostRepository) ListAll() ([]model.Post, error) {
	rows, err := r.db.Query("SELECT id, title, slug, content, published, created_at, updated_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, err
}

func (r *sqlitePostRepository) GetBySlug(slug string) (*model.Post, error) {
	var p model.Post
	err := r.db.QueryRow("SELECT id, title, slug, content, published, created_at, updated_at FROM posts WHERE slug = ?",
		slug,
	).Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPostNotFound
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *sqlitePostRepository) Create(p *model.Post) error {
	result, err := r.db.Exec(
		"INSERT INTO posts (title, slug, content, published) VALUES (?, ?, ?, ?)",
		p.Title, p.Slug, p.Content, p.Published,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)

	return nil
}

func (r *sqlitePostRepository) Update(p *model.Post) error {
	result, err := r.db.Exec(
		"UPDATE posts SET title = ?, slug = ?, content = ?, published = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id =?",
		p.Title, p.Slug, p.Content, p.Published, p.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrPostNotFound
	}

	return nil
}

func (r *sqlitePostRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM posts WHERE id = ?", id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrPostNotFound
	}

	return nil
}

func (r *sqlitePostRepository) GetByID(id int) (*model.Post, error) {
	var p model.Post
	err := r.db.QueryRow("SELECT id, title, slug, content, published, created_at, updated_at FROM posts WHERE id = ?",
		id,
	).Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.Published, &p.CreatedAt, &p.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPostNotFound
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &sqlitePostRepository{db: db}
}
