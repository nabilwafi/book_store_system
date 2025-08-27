package repository

import (
	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/pkg/logger"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
	GetAll(page, limit int) ([]*entity.User, int64, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// Create creates a new user
func (r *userRepositoryImpl) Create(user *entity.User) error {
	logger.Info("Creating new user with email: %s", user.Email)
	err := r.db.Create(user).Error
	if err != nil {
		logger.Error("Failed to create user: %v", err)
		return err
	}
	logger.Info("Successfully created user with ID: %d", user.ID)
	return nil
}

// GetByID gets a user by ID
func (r *userRepositoryImpl) GetByID(id uint) (*entity.User, error) {
	logger.Info("Fetching user by ID: %d", id)
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		logger.Error("Failed to fetch user by ID %d: %v", id, err)
		return nil, err
	}
	logger.Info("Successfully fetched user: %s", user.Email)
	return &user, nil
}

// GetByEmail gets a user by email
func (r *userRepositoryImpl) GetByEmail(email string) (*entity.User, error) {
	logger.Info("Fetching user by email: %s", email)
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		logger.Error("Failed to fetch user by email %s: %v", email, err)
		return nil, err
	}
	logger.Info("Successfully fetched user by email: %s", email)
	return &user, nil
}

// Update updates a user
func (r *userRepositoryImpl) Update(user *entity.User) error {
	logger.Info("Updating user with ID: %d", user.ID)
	err := r.db.Save(user).Error
	if err != nil {
		logger.Error("Failed to update user with ID %d: %v", user.ID, err)
		return err
	}
	logger.Info("Successfully updated user with ID: %d", user.ID)
	return nil
}

// Delete deletes a user
func (r *userRepositoryImpl) Delete(id uint) error {
	logger.Info("Deleting user with ID: %d", id)
	err := r.db.Delete(&entity.User{}, id).Error
	if err != nil {
		logger.Error("Failed to delete user with ID %d: %v", id, err)
		return err
	}
	logger.Info("Successfully deleted user with ID: %d", id)
	return nil
}

// GetAll gets all users with pagination
func (r *userRepositoryImpl) GetAll(page, limit int) ([]*entity.User, int64, error) {
	logger.Info("Fetching all users with pagination - page: %d, limit: %d", page, limit)
	var users []*entity.User
	var total int64

	// Count total records
	err := r.db.Model(&entity.User{}).Count(&total).Error
	if err != nil {
		logger.Error("Failed to count users: %v", err)
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * limit
	err = r.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		logger.Error("Failed to fetch users with pagination: %v", err)
		return nil, 0, err
	}

	logger.Info("Successfully fetched %d users out of %d total", len(users), total)
	return users, total, nil
}
