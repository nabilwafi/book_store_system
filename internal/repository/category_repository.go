package repository

import (
	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/pkg/logger"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *entity.Category) error
	GetByID(id uint) (*entity.Category, error)
	GetByName(name string) (*entity.Category, error)
	Update(category *entity.Category) error
	Delete(id uint) error
	GetAll(page, limit int) ([]*entity.Category, int64, error)
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

// Create creates a new category
func (r *categoryRepositoryImpl) Create(category *entity.Category) error {
	logger.Info("Creating new category: %s", category.Name)
	err := r.db.Create(category).Error
	if err != nil {
		logger.Error("Failed to create category: %v", err)
		return err
	}
	logger.Info("Successfully created category with ID: %d", category.ID)
	return nil
}

// GetByID gets a category by ID
func (r *categoryRepositoryImpl) GetByID(id uint) (*entity.Category, error) {
	logger.Info("Fetching category by ID: %d", id)
	var category entity.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		logger.Error("Failed to fetch category by ID %d: %v", id, err)
		return nil, err
	}
	logger.Info("Successfully fetched category: %s", category.Name)
	return &category, nil
}

// GetByName gets a category by name
func (r *categoryRepositoryImpl) GetByName(name string) (*entity.Category, error) {
	logger.Info("Fetching category by name: %s", name)
	var category entity.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		logger.Error("Failed to fetch category by name %s: %v", name, err)
		return nil, err
	}
	logger.Info("Successfully fetched category by name: %s", name)
	return &category, nil
}

// Update updates a category
func (r *categoryRepositoryImpl) Update(category *entity.Category) error {
	logger.Info("Updating category with ID: %d", category.ID)
	err := r.db.Save(category).Error
	if err != nil {
		logger.Error("Failed to update category with ID %d: %v", category.ID, err)
		return err
	}
	logger.Info("Successfully updated category: %s", category.Name)
	return nil
}

// Delete deletes a category
func (r *categoryRepositoryImpl) Delete(id uint) error {
	logger.Info("Deleting category with ID: %d", id)
	err := r.db.Delete(&entity.Category{}, id).Error
	if err != nil {
		logger.Error("Failed to delete category with ID %d: %v", id, err)
		return err
	}
	logger.Info("Successfully deleted category with ID: %d", id)
	return nil
}

// GetAll gets all categories with pagination
func (r *categoryRepositoryImpl) GetAll(page, limit int) ([]*entity.Category, int64, error) {
	logger.Info("Fetching all categories with pagination - page: %d, limit: %d", page, limit)
	var categories []*entity.Category
	var total int64

	// Count total records
	err := r.db.Model(&entity.Category{}).Count(&total).Error
	if err != nil {
		logger.Error("Failed to count categories: %v", err)
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * limit
	err = r.db.Offset(offset).Limit(limit).Find(&categories).Error
	if err != nil {
		logger.Error("Failed to fetch categories with pagination: %v", err)
		return nil, 0, err
	}

	logger.Info("Successfully fetched %d categories out of %d total", len(categories), total)
	return categories, total, nil
}
