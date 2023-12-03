package services

import (
	"context"
	"thyra/internal/users/repositories"
	"thyra/internal/users/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.AuthRepository
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) AuthenticateUser(ctx context.Context, username, password string) (string, error) {
	userID, storedPassword, userType, err := s.repo.GetUserCredentials(username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		return "", err
	}

	return utils.GenerateJWTToken(userID, username, userType) // Assuming jwt.GenerateToken is implemented
}
