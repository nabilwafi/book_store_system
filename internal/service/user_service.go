package service

import (
	"errors"

	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/internal/repository"
	"github.com/nabil/book-store-system/pkg/helpers"
	"github.com/nabil/book-store-system/pkg/logger"
	"github.com/nabil/book-store-system/pkg/middleware"
	"gorm.io/gorm"
)

type UserService interface {
	Register(name, email, password string) (*entity.User, error)
	Login(email, password string) (string, *entity.User, error)
	GetProfile(token string) (*entity.User, error)
	UpdateProfile(userID uint, name, email string) (*entity.User, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
}

type userServiceImpl struct {
	userRepo repository.UserRepository
	auth     *middleware.AuthMiddleware
}

func NewUserService(userRepo repository.UserRepository) UserService {
	auth := middleware.NewAuthMiddleware(userRepo)
	return &userServiceImpl{
		userRepo: userRepo,
		auth:     auth,
	}
}

// Register creates a new user account
func (s *userServiceImpl) Register(name, email, password string) (*entity.User, error) {
	logger.Info("Starting user registration", "email", email, "name", name)

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to check existing user", "email", email, "error", err)
		return nil, err
	}
	if existingUser != nil {
		logger.Error("User registration failed - email already exists", "email", email)
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		logger.Error("Failed to hash password", "email", email, "error", err)
		return nil, err
	}

	// Create user
	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     "user", // Default role
	}

	err = s.userRepo.Create(user)
	if err != nil {
		logger.Error("Failed to create user", "email", email, "error", err)
		return nil, err
	}

	logger.Info("User registration successful", "email", email, "userID", user.ID)
	return user, nil
}

// Login authenticates a user and returns a token
func (s *userServiceImpl) Login(email, password string) (string, *entity.User, error) {
	logger.Info("Starting user login", "email", email)

	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		logger.Error("Login failed - user not found", "email", email, "error", err)
		return "", nil, errors.New("invalid email or password")
	}

	// Verify password
	if !helpers.CheckPassword(password, user.Password) {
		logger.Error("Login failed - invalid password", "email", email)
		return "", nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := helpers.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		logger.Error("Failed to generate token", "email", email, "error", err)
		return "", nil, err
	}

	logger.Info("User login successful", "email", email, "userID", user.ID)
	return token, user, nil
}

// GetProfile retrieves user profile by token
func (s *userServiceImpl) GetProfile(token string) (*entity.User, error) {
	logger.Info("Getting user profile")

	user, err := s.auth.ValidateUserToken(token)
	if err != nil {
		logger.Error("Failed to get user profile", "error", err)
		return nil, err
	}

	logger.Info("User profile retrieved successfully", "userID", user.ID, "email", user.Email)
	return user, nil
}

// UpdateProfile updates user profile
func (s *userServiceImpl) UpdateProfile(userID uint, name, email string) (*entity.User, error) {
	logger.Info("Starting profile update", "userID", userID, "name", name, "email", email)

	// Get existing user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		logger.Error("Failed to get user for profile update", "userID", userID, "error", err)
		return nil, err
	}

	// Check if email is already taken by another user
	if email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Failed to check email availability", "email", email, "error", err)
			return nil, err
		}
		if existingUser != nil {
			logger.Error("Profile update failed - email already taken", "email", email, "userID", userID)
			return nil, errors.New("email is already taken")
		}
	}

	// Update user fields
	user.Name = name
	user.Email = email

	// Save updated user
	err = s.userRepo.Update(user)
	if err != nil {
		logger.Error("Failed to update user profile", "userID", userID, "error", err)
		return nil, err
	}

	logger.Info("Profile update successful", "userID", userID, "email", email)
	return user, nil
}

// ChangePassword changes user password
func (s *userServiceImpl) ChangePassword(userID uint, oldPassword, newPassword string) error {
	logger.Info("Starting password change", "userID", userID)

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		logger.Error("Failed to get user for password change", "userID", userID, "error", err)
		return err
	}

	// Verify old password
	if !helpers.CheckPassword(oldPassword, user.Password) {
		logger.Error("Password change failed - invalid old password", "userID", userID)
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := helpers.HashPassword(newPassword)
	if err != nil {
		logger.Error("Failed to hash new password", "userID", userID, "error", err)
		return err
	}

	// Update password
	user.Password = hashedPassword
	err = s.userRepo.Update(user)
	if err != nil {
		logger.Error("Failed to update password", "userID", userID, "error", err)
		return err
	}

	logger.Info("Password change successful", "userID", userID)
	return nil
}
