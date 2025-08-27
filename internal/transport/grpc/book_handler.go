package grpc

import (
	"context"

	"github.com/nabil/book-store-system/internal/service"
	"github.com/nabil/book-store-system/internal/transport/dto"
	"github.com/nabil/book-store-system/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BookHandler handles gRPC requests for book operations
type BookHandler struct {
	proto.UnimplementedBookServiceServer
	bookService service.BookService
}

// NewBookHandler creates a new BookHandler
func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

// CreateBook handles book creation
func (h *BookHandler) CreateBook(ctx context.Context, req *proto.CreateBookRequest) (*proto.CreateBookResponse, error) {
	// Validate request using DTO
	createDTO := &dto.CreateBookRequestDTO{
		Title:       req.Title,
		Author:      req.Author,
		ImageBase64: req.ImageBase64,
		Token:       req.Token,
		Price:       req.Price,
		Stock:       req.Stock,
		Year:        req.Year,
		CategoryID:  req.CategoryId,
	}

	if err := createDTO.ValidateCreateBookRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	book, err := h.bookService.CreateBook(
		req.Title,
		req.Author,
		req.ImageBase64,
		req.Token,
		req.Price,
		int(req.Stock),
		int(req.Year),
		uint(req.CategoryId),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create book: %v", err)
	}

	return &proto.CreateBookResponse{
		Success: true,
		Book: &proto.Book{
			Id:          uint32(book.ID),
			Title:       book.Title,
			Author:      book.Author,
			Price:       book.Price,
			Stock:       int32(book.Stock),
			Year:        int32(book.Year),
			CategoryId:  uint32(book.CategoryID),
			ImageBase64: book.ImageBase64,
			Category: &proto.Category{
				Id:   uint32(book.Category.ID),
				Name: book.Category.Name,
			},
		},
		Message: "Book created successfully",
	}, nil
}

// GetBooks retrieves all books with pagination and optional search
func (h *BookHandler) GetBooks(ctx context.Context, req *proto.GetBooksRequest) (*proto.GetBooksResponse, error) {
	// Validate request using DTO
	getBooksDTO := &dto.GetBooksRequestDTO{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	}

	if err := getBooksDTO.ValidateGetBooksRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	books, total, err := h.bookService.GetBooks(int(req.Page), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get books: %v", err)
	}

	var protoBooks []*proto.Book
	for _, book := range books {
		protoBook := &proto.Book{
			Id:          uint32(book.ID),
			Title:       book.Title,
			Author:      book.Author,
			Price:       book.Price,
			Stock:       int32(book.Stock),
			Year:        int32(book.Year),
			CategoryId:  uint32(book.CategoryID),
			ImageBase64: book.ImageBase64,
		}

		if book.Category.ID != 0 {
			protoBook.Category = &proto.Category{
				Id:   uint32(book.Category.ID),
				Name: book.Category.Name,
			}
		}

		protoBooks = append(protoBooks, protoBook)
	}

	return &proto.GetBooksResponse{
		Success: true,
		Message: "Books retrieved successfully",
		Books:   protoBooks,
		Total:   int32(total),
	}, nil
}

// GetBook retrieves a book by ID
func (h *BookHandler) GetBook(ctx context.Context, req *proto.GetBookRequest) (*proto.GetBookResponse, error) {
	// Validate request using DTO
	getBookDTO := &dto.GetBookRequestDTO{
		ID: req.Id,
	}

	if err := getBookDTO.ValidateGetBookRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	book, err := h.bookService.GetBook(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Book not found: %v", err)
	}

	protoBook := &proto.Book{
		Id:          uint32(book.ID),
		Title:       book.Title,
		Author:      book.Author,
		Price:       book.Price,
		Stock:       int32(book.Stock),
		Year:        int32(book.Year),
		CategoryId:  uint32(book.CategoryID),
		ImageBase64: book.ImageBase64,
	}

	if book.Category.ID != 0 {
		protoBook.Category = &proto.Category{
			Id:   uint32(book.Category.ID),
			Name: book.Category.Name,
		}
	}

	return &proto.GetBookResponse{
		Success: true,
		Message: "Book retrieved successfully",
		Book:    protoBook,
	}, nil
}

// UpdateBook updates an existing book
func (h *BookHandler) UpdateBook(ctx context.Context, req *proto.UpdateBookRequest) (*proto.UpdateBookResponse, error) {
	// Validate request using DTO
	updateDTO := &dto.UpdateBookRequestDTO{
		ID:         req.Id,
		Title:      req.Title,
		Author:     req.Author,
		Token:      req.Token,
		Price:      req.Price,
		Stock:      req.Stock,
		Year:       req.Year,
		CategoryID: req.CategoryId,
	}

	if err := updateDTO.ValidateUpdateBookRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	book, err := h.bookService.UpdateBook(
		uint(req.Id),
		req.Title,
		req.Author,
		req.ImageBase64,
		req.Token,
		req.Price,
		int(req.Stock),
		int(req.Year),
		uint(req.CategoryId),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update book: %v", err)
	}

	return &proto.UpdateBookResponse{
		Success: true,
		Book: &proto.Book{
			Id:          uint32(book.ID),
			Title:       book.Title,
			Author:      book.Author,
			Price:       book.Price,
			Stock:       int32(book.Stock),
			Year:        int32(book.Year),
			CategoryId:  uint32(book.CategoryID),
			ImageBase64: book.ImageBase64,
			Category: &proto.Category{
				Id:   uint32(book.Category.ID),
				Name: book.Category.Name,
			},
		},
		Message: "Book updated successfully",
	}, nil
}

// DeleteBook deletes a book
func (h *BookHandler) DeleteBook(ctx context.Context, req *proto.DeleteBookRequest) (*proto.DeleteBookResponse, error) {
	// Validate request using DTO
	deleteDTO := &dto.DeleteBookRequestDTO{
		ID:    req.Id,
		Token: req.Token,
	}

	if err := deleteDTO.ValidateDeleteBookRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	err := h.bookService.DeleteBook(uint(req.Id), req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete book: %v", err)
	}

	return &proto.DeleteBookResponse{
		Success: true,
		Message: "Book deleted successfully",
	}, nil
}

// GetBooksByCategory retrieves books by category with pagination
func (h *BookHandler) GetBooksByCategory(ctx context.Context, req *proto.GetBooksByCategoryRequest) (*proto.GetBooksByCategoryResponse, error) {
	// Validate request using DTO
	getBooksByCategoryDTO := &dto.GetBooksByCategoryRequestDTO{
		CategoryID: req.CategoryId,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	if err := getBooksByCategoryDTO.ValidateGetBooksByCategoryRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	books, total, err := h.bookService.GetBooksByCategory(uint(req.CategoryId), int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get books by category: %v", err)
	}

	var protoBooks []*proto.Book
	for _, book := range books {
		protoBook := &proto.Book{
			Id:          uint32(book.ID),
			Title:       book.Title,
			Author:      book.Author,
			Price:       book.Price,
			Stock:       int32(book.Stock),
			Year:        int32(book.Year),
			CategoryId:  uint32(book.CategoryID),
			ImageBase64: book.ImageBase64,
		}

		if book.Category.ID != 0 {
			protoBook.Category = &proto.Category{
				Id:   uint32(book.Category.ID),
				Name: book.Category.Name,
			}
		}

		protoBooks = append(protoBooks, protoBook)
	}

	return &proto.GetBooksByCategoryResponse{
		Success: true,
		Message: "Books retrieved successfully",
		Books:   protoBooks,
		Total:   int32(total),
	}, nil
}
