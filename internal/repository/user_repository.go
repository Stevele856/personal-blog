package repository

import (
	"database/sql"
	"errors"

	"github.com/letrongvu/blog/internal/model"
)

// Sentinel error
var ErrUserNotFound = errors.New("user not found")

// Interface
type UserRepository interface {
	GetByUsername(username string) (*model.User, error)
}

// Struct
type sqliteUserRepository struct {
	db *sql.DB
}

// Define sqliteUserRepository
func (r *sqliteUserRepository) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow("SELECT id, username, password_hash, created_at FROM users WHERE username = ?",
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Constructor
func NewUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}
