package service

import (
	"errors"

	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/internal/repository"
	"github.com/nabil/book-store-system/pkg/logger"
	"github.com/nabil/book-store-system/pkg/middleware"
	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(name, token string) (*entity.Category, error)
	GetCategories(page, limit int) ([]*entity.Category, int64, error)
	GetCategory(id uint) (*entity.Category, error)
	UpdateCategory(id uint, name, token string) (*entity.Category, error)
	DeleteCategory(id uint, token string) error
}

type categoryServiceImpl struct {
	categoryRepo repository.CategoryRepository
	userRepo     repository.UserRepository
	auth         *middleware.AuthMiddleware
}

func NewCategoryService(categoryRepo repository.CategoryRepository, userRepo repository.UserRepository) CategoryService {
	auth := middleware.NewAuthMiddleware(userRepo)
	return &categoryServiceImpl{
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
		auth:         auth,
	}
}

// CreateCategory creates a new category (admin only)
func (s *categoryServiceImpl) CreateCategory(name, token string) (*entity.Category, error) {
	logger.Info("Starting category creation", "name", name)

	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Category creation failed - invalid admin token", "name", name, "error", err)
		return nil, err
	}

	// Check if category already exists
	existingCategory, err := s.categoryRepo.GetByName(name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to check existing category", "name", name, "error", err)
		return nil, err
	}
	if existingCategory != nil {
		logger.Error("Category creation failed - name already exists", "name", name)
		return nil, errors.New("category with this name already exists")
	}

	category := &entity.Category{
		Name: name,
	}

	err = s.categoryRepo.Create(category)
	if err != nil {
		logger.Error("Failed to create category", "name", name, "error", err)
		return nil, err
	}

	logger.Info("Category creation successful", "name", name, "categoryID", category.ID)
	return category, nil
}

// GetCategories retrieves all categories with pagination
func (s *categoryServiceImpl) GetCategories(page, limit int) ([]*entity.Category, int64, error) {
	logger.Info("Getting categories", "page", page, "limit", limit)

	categories, total, err := s.categoryRepo.GetAll(page, limit)
	if err != nil {
		logger.Error("Failed to get categories", "page", page, "limit", limit, "error", err)
		return nil, 0, err
	}

	logger.Info("Categories retrieved successfully", "count", len(categories), "total", total)
	return categories, total, nil
}

// GetCategory retrieves a category by ID
func (s *categoryServiceImpl) GetCategory(id uint) (*entity.Category, error) {
	logger.Info("Getting category by ID", "categoryID", id)

	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get category", "categoryID", id, "error", err)
		return nil, err
	}

	logger.Info("Category retrieved successfully", "categoryID", id, "name", category.Name)
	return category, nil
}

// UpdateCategory updates an existing category (admin only)
func (s *categoryServiceImpl) UpdateCategory(id uint, name, token string) (*entity.Category, error) {
	logger.Info("Starting category update", "categoryID", id, "name", name)

	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Category update failed - invalid admin token", "categoryID", id, "error", err)
		return nil, err
	}

	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get category for update", "categoryID", id, "error", err)
		return nil, err
	}

	if name != category.Name {
		existingCategory, err := s.categoryRepo.GetByName(name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Failed to check category name availability", "name", name, "error", err)
			return nil, err
		}
		if existingCategory != nil {
			logger.Error("Category update failed - name already taken", "name", name, "categoryID", id)
			return nil, errors.New("category name is already taken")
		}
	}

	category.Name = name

	err = s.categoryRepo.Update(category)
	if err != nil {
		logger.Error("Failed to update category", "categoryID", id, "error", err)
		return nil, err
	}

	logger.Info("Category update successful", "categoryID", id, "name", name)
	return category, nil
}

// DeleteCategory deletes a category (admin only)
func (s *categoryServiceImpl) DeleteCategory(id uint, token string) error {
	logger.Info("Starting category deletion", "categoryID", id)

	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Category deletion failed - invalid admin token", "categoryID", id, "error", err)
		return err
	}

	_, err = s.categoryRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get category for deletion", "categoryID", id, "error", err)
		return err
	}

	err = s.categoryRepo.Delete(id)
	if err != nil {
		logger.Error("Failed to delete category", "categoryID", id, "error", err)
		return err
	}

	logger.Info("Category deletion successful", "categoryID", id)
	return nil
}
