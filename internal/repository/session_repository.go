package repository

import (
	"database/sql"
	"errors"
	"time"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionRepository interface {
	Create(userID int, token string, expiresAt time.Time) error
	GetUserIDByToken(token string) (int, error)
	Delete(token string) error
}

type sqliteSessionRepository struct {
	db *sql.DB
}

const sqliteTimeLayout = "2006-01-02T15:04:05.000Z"

func (r *sqliteSessionRepository) Create(userID int, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expiresAt.UTC().Format(sqliteTimeLayout),
	)

	return err
}

func (r *sqliteSessionRepository) GetUserIDByToken(token string) (int, error) {
	var uID int
	err := r.db.QueryRow(
		"SELECT user_id FROM sessions WHERE token = ? AND expires_at > ?",
		token, time.Now().UTC().Format(sqliteTimeLayout),
	).Scan(&uID)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrSessionNotFound
	}
	if err != nil {
		return 0, err
	}

	return uID, nil
}

func (r *sqliteSessionRepository) Delete(token string) error {
	result, err := r.db.Exec("DELETE FROM sessions WHERE token = ?",token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrSessionNotFound
	}

	return nil
}
