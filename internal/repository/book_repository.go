package repository

import (
	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/pkg/logger"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *entity.Book) error
	GetByID(id uint) (*entity.Book, error)
	Update(book *entity.Book) error
	Delete(id uint) error
	GetAll(page, limit int, search string) ([]*entity.Book, int64, error)
	GetByCategory(categoryID uint, page, limit int) ([]*entity.Book, int64, error)
	UpdateStock(id uint, stock int) error
	UpdateStockTx(tx *gorm.DB, id uint, stock int) error
	CheckStock(id uint, quantity int) (bool, error)
}

type bookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepositoryImpl{
		db: db,
	}
}

// Create creates a new book
func (r *bookRepositoryImpl) Create(book *entity.Book) error {
	logger.Infof("Creating new book: %s", book.Title)
	err := r.db.Create(book).Error
	if err != nil {
		logger.Errorf("Failed to create book: %v", err)
		return err
	}
	logger.Infof("Successfully created book with ID: %d", book.ID)
	return nil
}

// GetByID retrieves a book by ID with category
func (r *bookRepositoryImpl) GetByID(id uint) (*entity.Book, error) {
	logger.Infof("Fetching book by ID: %d", id)
	var book entity.Book
	err := r.db.Preload("Category").First(&book, id).Error
	if err != nil {
		logger.Errorf("Failed to fetch book by ID %d: %v", id, err)
		return nil, err
	}
	logger.Infof("Successfully fetched book: %s", book.Title)
	return &book, nil
}

// Update updates an existing book
func (r *bookRepositoryImpl) Update(book *entity.Book) error {
	logger.Infof("Updating book with ID: %d", book.ID)
	err := r.db.Save(book).Error
	if err != nil {
		logger.Errorf("Failed to update book with ID %d: %v", book.ID, err)
		return err
	}
	logger.Infof("Successfully updated book: %s", book.Title)
	return nil
}

// Delete deletes a book by ID
func (r *bookRepositoryImpl) Delete(id uint) error {
	logger.Infof("Deleting book with ID: %d", id)
	err := r.db.Delete(&entity.Book{}, id).Error
	if err != nil {
		logger.Errorf("Failed to delete book with ID %d: %v", id, err)
		return err
	}
	logger.Infof("Successfully deleted book with ID: %d", id)
	return nil
}

// GetAll retrieves all books with pagination and optional search
func (r *bookRepositoryImpl) GetAll(page, limit int, search string) ([]*entity.Book, int64, error) {
	logger.Infof("Fetching all books - page: %d, limit: %d, search: %s", page, limit, search)
	var books []*entity.Book
	var total int64

	query := r.db.Model(&entity.Book{}).Preload("Category")

	// Apply search filter if provided
	if search != "" {
		query = query.Where("title LIKE ? OR author LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		logger.Errorf("Failed to count books: %v", err)
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve books with pagination
	err := query.Offset(offset).Limit(limit).Find(&books).Error
	if err != nil {
		logger.Errorf("Failed to fetch books with pagination: %v", err)
		return nil, 0, err
	}

	logger.Infof("Successfully fetched %d books out of %d total", len(books), total)
	return books, total, nil
}

// GetByCategory retrieves books by category with pagination
func (r *bookRepositoryImpl) GetByCategory(categoryID uint, page, limit int) ([]*entity.Book, int64, error) {
	logger.Infof("Fetching books by category ID: %d - page: %d, limit: %d", categoryID, page, limit)
	var books []*entity.Book
	var total int64

	query := r.db.Model(&entity.Book{}).Preload("Category").Where("category_id = ?", categoryID)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		logger.Errorf("Failed to count books by category %d: %v", categoryID, err)
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve books with pagination
	err := query.Offset(offset).Limit(limit).Find(&books).Error
	if err != nil {
		logger.Errorf("Failed to fetch books by category %d: %v", categoryID, err)
		return nil, 0, err
	}

	logger.Infof("Successfully fetched %d books from category %d out of %d total", len(books), categoryID, total)
	return books, total, nil
}

// UpdateStock updates book stock
func (r *bookRepositoryImpl) UpdateStock(id uint, stock int) error {
	logger.Infof("Updating stock for book ID %d to %d", id, stock)
	err := r.db.Model(&entity.Book{}).Where("id = ?", id).Update("stock", stock).Error
	if err != nil {
		logger.Errorf("Failed to update stock for book ID %d: %v", id, err)
		return err
	}
	logger.Infof("Successfully updated stock for book ID %d", id)
	return nil
}

// UpdateStockTx updates book stock using external transaction
// This function allows for more flexible transaction management
func (r *bookRepositoryImpl) UpdateStockTx(tx *gorm.DB, id uint, stock int) error {
	logger.Infof("Updating stock for book ID %d to %d with external transaction", id, stock)
	err := tx.Model(&entity.Book{}).Where("id = ?", id).Update("stock", stock).Error
	if err != nil {
		logger.Errorf("Failed to update stock for book ID %d in transaction: %v", id, err)
		return err
	}
	logger.Infof("Successfully updated stock for book ID %d in transaction", id)
	return nil
}

// CheckStock checks if book has sufficient stock
func (r *bookRepositoryImpl) CheckStock(id uint, quantity int) (bool, error) {
	logger.Infof("Checking stock for book ID %d, required quantity: %d", id, quantity)
	var book entity.Book
	err := r.db.Select("stock").First(&book, id).Error
	if err != nil {
		logger.Errorf("Failed to check stock for book ID %d: %v", id, err)
		return false, err
	}
	hasStock := book.Stock >= quantity
	logger.Infof("Stock check for book ID %d: current stock %d, required %d, sufficient: %t", id, book.Stock, quantity, hasStock)
	return hasStock, nil
}
