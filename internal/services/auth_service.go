package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/letrongvu/blog/internal/model"
	"github.com/letrongvu/blog/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type UserService struct {
	userRepo repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewUserService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *UserService{
	return &UserService{
		userRepo: userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *UserService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := &model.User{
		Username: username,
		PasswordHash: string(hashedPassword),
	}
	return s.userRepo.Create(u)
}

func (s *UserService) Login(username, password string) (string, error) {
	u, err := s.userRepo.GetByUsername(username)

	if err != nil {
		return "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate random tokens
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(b)

	// Caculate expires and save session
	expiresAt := time.Now().UTC().Add(24 * time.Hour)
	if err := s.sessionRepo.Create(u.ID, token, expiresAt); err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) ValidateSession(token string) (int, error) {
	uID, err := s.sessionRepo.GetUserIDByToken(token)

	if err != nil {
		return 0, err
	}

	return uID, nil

	
}

func (s *UserService) Logout(token string) error {
    return s.sessionRepo.Delete(token)
}
