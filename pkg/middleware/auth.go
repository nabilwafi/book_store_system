package middleware

import (
	"errors"
	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/internal/repository"
	"github.com/nabil/book-store-system/pkg/helpers"
)

// AuthMiddleware handles authentication and authorization
type AuthMiddleware struct {
	userRepo repository.UserRepository
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
	}
}

// ValidateUserToken validates if token belongs to a valid user
func (m *AuthMiddleware) ValidateUserToken(token string) (*entity.User, error) {
	// Validate token
	claims, err := helpers.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	user, err := m.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ValidateAdminToken validates if token belongs to admin user
func (m *AuthMiddleware) ValidateAdminToken(token string) (*entity.User, error) {
	// Validate token
	claims, err := helpers.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	user, err := m.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is admin
	if user.Role != "admin" {
		return nil, errors.New("access denied: admin role required")
	}

	return user, nil
}

// ValidateToken validates JWT token and returns user info (general purpose)
func (m *AuthMiddleware) ValidateToken(token string) (*entity.User, error) {
	return m.ValidateUserToken(token)
}