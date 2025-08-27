package repository

import (
	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/pkg/logger"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *entity.Order, items []*entity.OrderItem) error
	GetByID(id uint) (*entity.Order, error)
	GetByUserID(userID uint, page, limit int) ([]*entity.Order, int64, error)
	Update(order *entity.Order) error
	UpdateStatus(id uint, status string) error
	UpdatePaymentURL(id uint, paymentURL string) error
	GetAll(page, limit int) ([]*entity.Order, int64, error)
}

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepositoryImpl{
		db: db,
	}
}

// Create creates a new order with items
func (r *orderRepositoryImpl) Create(order *entity.Order, items []*entity.OrderItem) error {
	logger.Info("Creating new order for user ID: %d with %d items", order.UserID, len(items))
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			logger.Error("Failed to create order: %v", err)
			return err
		}

		// Set order ID for all items
		for _, item := range items {
			item.OrderID = order.ID
		}

		// Create order items
		if err := tx.Create(&items).Error; err != nil {
			logger.Error("Failed to create order items: %v", err)
			return err
		}

		logger.Info("Successfully created order with ID: %d", order.ID)

		return nil
	})
}

// GetByID retrieves an order by ID with items
func (r *orderRepositoryImpl) GetByID(id uint) (*entity.Order, error) {
	logger.Info("Fetching order by ID: %d", id)
	var order entity.Order
	err := r.db.Preload("User").Preload("OrderItems.Book").First(&order, id).Error
	if err != nil {
		logger.Error("Failed to fetch order by ID %d: %v", id, err)
		return nil, err
	}
	logger.Info("Successfully fetched order ID: %d for user ID: %d", order.ID, order.UserID)
	return &order, nil
}

// GetByUserID retrieves orders by user ID with pagination
func (r *orderRepositoryImpl) GetByUserID(userID uint, page, limit int) ([]*entity.Order, int64, error) {
	logger.Info("Fetching orders for user ID: %d - page: %d, limit: %d", userID, page, limit)
	var orders []*entity.Order
	var total int64

	query := r.db.Model(&entity.Order{}).Where("user_id = ?", userID)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		logger.Error("Failed to count orders for user ID %d: %v", userID, err)
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve orders with pagination
	err := query.Preload("OrderItems.Book").Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		logger.Error("Failed to fetch orders for user ID %d: %v", userID, err)
		return nil, 0, err
	}

	logger.Info("Successfully fetched %d orders for user ID %d out of %d total", len(orders), userID, total)
	return orders, total, nil
}

// Update updates an existing order
func (r *orderRepositoryImpl) Update(order *entity.Order) error {
	logger.Info("Updating order with ID: %d", order.ID)
	err := r.db.Save(order).Error
	if err != nil {
		logger.Error("Failed to update order with ID %d: %v", order.ID, err)
		return err
	}
	logger.Info("Successfully updated order with ID: %d", order.ID)
	return nil
}

// UpdateStatus updates order status
func (r *orderRepositoryImpl) UpdateStatus(id uint, status string) error {
	logger.Info("Updating order status for ID %d to: %s", id, status)
	err := r.db.Model(&entity.Order{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		logger.Error("Failed to update order status for ID %d: %v", id, err)
		return err
	}
	logger.Info("Successfully updated order status for ID %d to: %s", id, status)
	return nil
}

// UpdatePaymentURL updates order payment URL
func (r *orderRepositoryImpl) UpdatePaymentURL(id uint, paymentURL string) error {
	logger.Info("Updating payment URL for order ID %d", id)
	err := r.db.Model(&entity.Order{}).Where("id = ?", id).Update("payment_url", paymentURL).Error
	if err != nil {
		logger.Error("Failed to update payment URL for order ID %d: %v", id, err)
		return err
	}
	logger.Info("Successfully updated payment URL for order ID %d", id)
	return nil
}

// GetAll retrieves all orders with pagination (admin only)
func (r *orderRepositoryImpl) GetAll(page, limit int) ([]*entity.Order, int64, error) {
	logger.Info("Fetching all orders with pagination - page: %d, limit: %d", page, limit)
	var orders []*entity.Order
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Order{}).Count(&total).Error; err != nil {
		logger.Error("Failed to count all orders: %v", err)
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve orders with pagination
	err := r.db.Preload("User").Preload("OrderItems.Book").Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		logger.Error("Failed to fetch all orders with pagination: %v", err)
		return nil, 0, err
	}

	logger.Info("Successfully fetched %d orders out of %d total", len(orders), total)
	return orders, total, nil
}
