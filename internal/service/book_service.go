package service

import (
	"errors"

	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/internal/repository"
	"github.com/nabil/book-store-system/pkg/logger"
	"github.com/nabil/book-store-system/pkg/middleware"
)

type BookService interface {
	CreateBook(title, author, imageBase64, token string, price float64, stock, year int, categoryID uint) (*entity.Book, error)
	GetBooks(page, limit int, search string) ([]*entity.Book, int64, error)
	GetBook(id uint) (*entity.Book, error)
	UpdateBook(id uint, title, author, imageBase64, token string, price float64, stock, year int, categoryID uint) (*entity.Book, error)
	DeleteBook(id uint, token string) error
	GetBooksByCategory(categoryID uint, page, limit int) ([]*entity.Book, int64, error)
	CheckBookAvailability(bookID uint, quantity int) (bool, error)
}

type bookServiceImpl struct {
	bookRepo     repository.BookRepository
	categoryRepo repository.CategoryRepository
	userRepo     repository.UserRepository
	auth         *middleware.AuthMiddleware
}

func NewBookService(bookRepo repository.BookRepository, categoryRepo repository.CategoryRepository, userRepo repository.UserRepository) BookService {
	auth := middleware.NewAuthMiddleware(userRepo)
	return &bookServiceImpl{
		bookRepo:     bookRepo,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
		auth:         auth,
	}
}

// CreateBook creates a new book (admin only)
func (s *bookServiceImpl) CreateBook(title, author, imageBase64, token string, price float64, stock, year int, categoryID uint) (*entity.Book, error) {
	logger.Info("Starting book creation", "title", title, "author", author, "categoryID", categoryID)

	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Book creation failed - invalid admin token", "title", title, "error", err)
		return nil, err
	}

	_, err = s.categoryRepo.GetByID(categoryID)
	if err != nil {
		logger.Error("Book creation failed - category not found", "title", title, "categoryID", categoryID, "error", err)
		return nil, errors.New("category not found")
	}

	// Create book
	book := &entity.Book{
		Title:       title,
		Author:      author,
		Price:       price,
		Stock:       stock,
		Year:        year,
		CategoryID:  categoryID,
		ImageBase64: imageBase64,
	}

	// Save book
	err = s.bookRepo.Create(book)
	if err != nil {
		logger.Error("Failed to create book", "title", title, "error", err)
		return nil, err
	}

	logger.Info("Book creation successful", "title", title, "bookID", book.ID, "categoryID", categoryID)
	return book, nil
}

// GetBooks retrieves books with pagination and search
func (s *bookServiceImpl) GetBooks(page, limit int, search string) ([]*entity.Book, int64, error) {
	logger.Info("Getting books", "page", page, "limit", limit, "search", search)

	books, total, err := s.bookRepo.GetAll(page, limit, search)
	if err != nil {
		logger.Error("Failed to get books", "page", page, "limit", limit, "search", search, "error", err)
		return nil, 0, err
	}

	logger.Info("Books retrieved successfully", "count", len(books), "total", total)
	return books, total, nil
}

// GetBook retrieves a book by ID
func (s *bookServiceImpl) GetBook(id uint) (*entity.Book, error) {
	logger.Info("Getting book by ID", "bookID", id)

	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get book", "bookID", id, "error", err)
		return nil, err
	}

	logger.Info("Book retrieved successfully", "bookID", id, "title", book.Title)
	return book, nil
}

// UpdateBook updates a book (admin only)
func (s *bookServiceImpl) UpdateBook(id uint, title, author, imageBase64, token string, price float64, stock, year int, categoryID uint) (*entity.Book, error) {
	logger.Info("Starting book update", "bookID", id, "title", title, "categoryID", categoryID)

	// Validate admin token
	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Book update failed - invalid admin token", "bookID", id, "error", err)
		return nil, err
	}

	// Check if book exists
	existingBook, err := s.bookRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get book for update", "bookID", id, "error", err)
		return nil, err
	}

	// Validate category exists
	_, err = s.categoryRepo.GetByID(categoryID)
	if err != nil {
		logger.Error("Book update failed - category not found", "bookID", id, "categoryID", categoryID, "error", err)
		return nil, errors.New("category not found")
	}

	// Update book fields
	existingBook.Title = title
	existingBook.Author = author
	existingBook.Price = price
	existingBook.Stock = stock
	existingBook.Year = year
	existingBook.CategoryID = categoryID
	existingBook.ImageBase64 = imageBase64

	// Update existing book
	err = s.bookRepo.Update(existingBook)
	if err != nil {
		logger.Error("Failed to update book", "bookID", id, "error", err)
		return nil, err
	}

	logger.Info("Book update successful", "bookID", id, "title", title)
	return existingBook, nil
}

// DeleteBook deletes a book (admin only)
func (s *bookServiceImpl) DeleteBook(id uint, token string) error {
	logger.Info("Starting book deletion", "bookID", id)

	// Validate admin token
	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Book deletion failed - invalid admin token", "bookID", id, "error", err)
		return err
	}

	// Check if book exists
	_, err = s.bookRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get book for deletion", "bookID", id, "error", err)
		return err
	}

	// Delete book
	err = s.bookRepo.Delete(id)
	if err != nil {
		logger.Error("Failed to delete book", "bookID", id, "error", err)
		return err
	}

	logger.Info("Book deletion successful", "bookID", id)
	return nil
}

// GetBooksByCategory retrieves books by category with pagination
func (s *bookServiceImpl) GetBooksByCategory(categoryID uint, page, limit int) ([]*entity.Book, int64, error) {
	logger.Info("Getting books by category", "categoryID", categoryID, "page", page, "limit", limit)

	_, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		logger.Error("Failed to get books by category - category not found", "categoryID", categoryID, "error", err)
		return nil, 0, errors.New("category not found")
	}

	books, total, err := s.bookRepo.GetByCategory(categoryID, page, limit)
	if err != nil {
		logger.Error("Failed to get books by category", "categoryID", categoryID, "error", err)
		return nil, 0, err
	}

	logger.Info("Successfully retrieved books by category", "categoryID", categoryID, "count", len(books), "total", total)
	return books, total, nil
}

// CheckBookAvailability checks if book is available for purchase
func (s *bookServiceImpl) CheckBookAvailability(bookID uint, quantity int) (bool, error) {
	logger.Info("Checking book availability", "bookID", bookID, "quantity", quantity)

	available, err := s.bookRepo.CheckStock(bookID, quantity)
	if err != nil {
		logger.Error("Failed to check book availability", "bookID", bookID, "quantity", quantity, "error", err)
		return false, err
	}

	logger.Info("Book availability check completed", "bookID", bookID, "quantity", quantity, "available", available)
	return available, nil
}
