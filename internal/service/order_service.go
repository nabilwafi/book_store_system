package service

import (
	"errors"
	"fmt"

	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/internal/repository"
	"github.com/nabil/book-store-system/pkg/helpers"
	"github.com/nabil/book-store-system/pkg/logger"
	"github.com/nabil/book-store-system/pkg/middleware"
)

// OrderItem represents an item in an order request
type OrderItem struct {
	BookID   uint `json:"book_id"`
	Quantity int  `json:"quantity"`
}

type OrderService interface {
	CreateOrder(items []OrderItem, token string) (*entity.Order, error)
	GetOrders(token string, page, limit int) ([]*entity.Order, int64, error) 
	GetOrder(id uint, token string) (*entity.Order, error)
	UpdateOrderStatus(id uint, status, token string) (*entity.Order, error)
	ProcessPayment(orderID uint, token string) (string, error)
	GetAllOrders(token string, page, limit int) ([]*entity.Order, int64, error)
}

// orderServiceImpl implements the OrderService interface
type orderServiceImpl struct {
	orderRepo repository.OrderRepository
	bookRepo  repository.BookRepository
	userRepo  repository.UserRepository
	auth      *middleware.AuthMiddleware
}

func NewOrderService(orderRepo repository.OrderRepository, bookRepo repository.BookRepository, userRepo repository.UserRepository) OrderService {
	auth := middleware.NewAuthMiddleware(userRepo)
	return &orderServiceImpl{
		orderRepo: orderRepo,
		bookRepo:  bookRepo,
		userRepo:  userRepo,
		auth:      auth,
	}
}

// CreateOrder creates a new order with items
func (s *orderServiceImpl) CreateOrder(items []OrderItem, token string) (*entity.Order, error) {
	logger.Info("Starting order creation", "itemCount", len(items))
	
	// Validate user token
	user, err := s.auth.ValidateUserToken(token)
	if err != nil {
		logger.Error("Order creation failed - invalid user token", "error", err)
		return nil, err
	}

	if len(items) == 0 {
		logger.Error("Order creation failed - no items provided", "userID", user.ID)
		return nil, errors.New("order must contain at least one item")
	}

	var totalAmount float64
	var orderItems []*entity.OrderItem

	// Validate items and calculate total
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}

		// Get book
		book, err := s.bookRepo.GetByID(item.BookID)
		if err != nil {
			logger.Error("Order creation failed - book not found", "userID", user.ID, "bookID", item.BookID, "error", err)
			return nil, fmt.Errorf("book with ID %d not found", item.BookID)
		}

		// Check stock
		available, err := s.bookRepo.CheckStock(item.BookID, item.Quantity)
		if err != nil {
			logger.Error("Order creation failed - stock check error", "userID", user.ID, "bookID", item.BookID, "error", err)
			return nil, err
		}
		if !available {
			logger.Error("Order creation failed - insufficient stock", "userID", user.ID, "bookID", item.BookID, "bookTitle", book.Title, "requestedQuantity", item.Quantity)
			return nil, fmt.Errorf("insufficient stock for book: %s", book.Title)
		}

		// Calculate item total
		itemTotal := book.Price * float64(item.Quantity)
		totalAmount += itemTotal

		// Create order item
		orderItem := &entity.OrderItem{
			BookID:   item.BookID,
			Quantity: item.Quantity,
			Price:    book.Price,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Create order
	order := &entity.Order{
		UserID:     user.ID,
		TotalPrice: totalAmount,
		Status:     "pending",
	}

	err = s.orderRepo.Create(order, orderItems)
	if err != nil {
		logger.Error("Failed to create order", "userID", user.ID, "totalAmount", totalAmount, "error", err)
		return nil, err
	}

	// Update book stocks
	for _, item := range items {
		book, _ := s.bookRepo.GetByID(item.BookID)
		newStock := book.Stock - item.Quantity
		err := s.bookRepo.UpdateStock(item.BookID, newStock)
		if err != nil {
			logger.Error("Failed to update book stock after order creation", "orderID", order.ID, "bookID", item.BookID, "error", err)
		}
	}

	// Load order with relations
	order, err = s.orderRepo.GetByID(order.ID)
	if err != nil {
		logger.Error("Failed to load created order", "orderID", order.ID, "error", err)
		return nil, err
	}

	logger.Info("Order creation successful", "orderID", order.ID, "userID", user.ID, "totalAmount", totalAmount, "itemCount", len(orderItems))
	return order, nil
}

// GetOrders retrieves orders for a user with pagination
func (s *orderServiceImpl) GetOrders(token string, page, limit int) ([]*entity.Order, int64, error) {
	logger.Info("Getting user orders", "page", page, "limit", limit)
	
	// Validate user token
	user, err := s.auth.ValidateUserToken(token)
	if err != nil {
		logger.Error("Failed to get orders - invalid user token", "error", err)
		return nil, 0, err
	}

	orders, total, err := s.orderRepo.GetByUserID(user.ID, page, limit)
	if err != nil {
		logger.Error("Failed to get user orders", "userID", user.ID, "error", err)
		return nil, 0, err
	}
	
	logger.Info("Successfully retrieved user orders", "userID", user.ID, "count", len(orders), "total", total)
	return orders, total, nil
}

// GetOrder retrieves a specific order by ID
func (s *orderServiceImpl) GetOrder(id uint, token string) (*entity.Order, error) {
	logger.Info("Getting order by ID", "orderID", id)
	
	// Validate user token
	user, err := s.auth.ValidateUserToken(token)
	if err != nil {
		logger.Error("Failed to get order - invalid user token", "orderID", id, "error", err)
		return nil, err
	}

	// Get order
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get order", "orderID", id, "userID", user.ID, "error", err)
		return nil, err
	}

	// Check if user owns the order or is admin
	if order.UserID != user.ID && user.Role != "admin" {
		logger.Error("Order access denied", "orderID", id, "userID", user.ID, "orderOwnerID", order.UserID, "userRole", user.Role)
		return nil, errors.New("access denied")
	}

	logger.Info("Successfully retrieved order", "orderID", id, "userID", user.ID)
	return order, nil
}

// UpdateOrderStatus updates order status (admin only)
func (s *orderServiceImpl) UpdateOrderStatus(id uint, status, token string) (*entity.Order, error) {
	logger.Info("Starting order status update", "orderID", id, "newStatus", status)
	
	// Validate admin token
	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Order status update failed - invalid admin token", "orderID", id, "error", err)
		return nil, err
	}

	// Validate status
	validStatuses := []string{"pending", "processing", "shipped", "completed", "cancelled"}
	validStatus := false
	for _, validStat := range validStatuses {
		if status == validStat {
			validStatus = true
			break
		}
	}
	if !validStatus {
		logger.Error("Order status update failed - invalid status", "orderID", id, "status", status, "validStatuses", validStatuses)
		return nil, errors.New("invalid status")
	}

	// Update status
	err = s.orderRepo.UpdateStatus(id, status)
	if err != nil {
		logger.Error("Failed to update order status", "orderID", id, "status", status, "error", err)
		return nil, err
	}

	// Get updated order
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get updated order", "orderID", id, "error", err)
		return nil, err
	}

	logger.Info("Order status update successful", "orderID", id, "newStatus", status)
	return order, nil
}

// ProcessPayment processes payment for an order
func (s *orderServiceImpl) ProcessPayment(orderID uint, token string) (string, error) {
	logger.Info("Starting payment processing", "orderID", orderID)
	
	// Validate user token
	user, err := s.auth.ValidateUserToken(token)
	if err != nil {
		logger.Error("Payment processing failed - invalid user token", "orderID", orderID, "error", err)
		return "", err
	}

	// Get order
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		logger.Error("Payment processing failed - order not found", "orderID", orderID, "userID", user.ID, "error", err)
		return "", err
	}

	// Check if user owns the order
	if order.UserID != user.ID {
		logger.Error("Payment processing failed - access denied", "orderID", orderID, "userID", user.ID, "orderOwnerID", order.UserID)
		return "", errors.New("access denied")
	}

	// Check if order is pending
	if order.Status != "pending" {
		logger.Error("Payment processing failed - order not in pending status", "orderID", orderID, "currentStatus", order.Status)
		return "", errors.New("order is not in pending status")
	}

	// Create payment using Midtrans
	paymentURL, err := helpers.CreatePayment(order, order.OrderItems)
	if err != nil {
		logger.Error("Payment processing failed - payment creation error", "orderID", orderID, "error", err)
		return "", err
	}

	// Update order payment URL
	err = s.orderRepo.UpdatePaymentURL(orderID, paymentURL)
	if err != nil {
		logger.Error("Failed to update payment URL", "orderID", orderID, "error", err)
		return "", err
	}

	// Update order status to processing
	err = s.orderRepo.UpdateStatus(orderID, "processing")
	if err != nil {
		logger.Error("Failed to update order status after payment creation", "orderID", orderID, "error", err)
		return "", err
	}

	logger.Info("Payment processing successful", "orderID", orderID, "userID", user.ID, "paymentURL", paymentURL)
	return paymentURL, nil
}

// GetAllOrders retrieves all orders with pagination (admin only)
func (s *orderServiceImpl) GetAllOrders(token string, page, limit int) ([]*entity.Order, int64, error) {
	logger.Info("Getting all orders (admin)", "page", page, "limit", limit)
	
	// Validate admin token
	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Failed to get all orders - invalid admin token", "error", err)
		return nil, 0, err
	}

	orders, total, err := s.orderRepo.GetAll(page, limit)
	if err != nil {
		logger.Error("Failed to get all orders", "error", err)
		return nil, 0, err
	}
	
	logger.Info("Successfully retrieved all orders", "count", len(orders), "total", total)
	return orders, total, nil
}
