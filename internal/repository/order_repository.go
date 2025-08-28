package repository

import (
	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/pkg/logger"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *entity.Order, items []*entity.OrderItem) error
	CreateOrderTx(tx *gorm.DB, order *entity.Order, items []*entity.OrderItem) error
	GetByID(id uint) (*entity.Order, error)
	GetByUserID(userID uint, page, limit int) ([]*entity.Order, int64, error)
	Update(order *entity.Order) error
	UpdateStatus(id uint, status string) error
	UpdateStatusTx(tx *gorm.DB, id uint, status string) error
	UpdatePaymentURL(id uint, paymentURL string) error
	UpdatePaymentURLTx(tx *gorm.DB, id uint, paymentURL string) error
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
	logger.Infof("Creating new order for user ID: %d with %d items", order.UserID, len(items))
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			logger.Errorf("Failed to create order: %v", err)
			return err
		}

		// Set order ID for all items
		for _, item := range items {
			item.OrderID = order.ID
		}

		// Create order items
		if err := tx.Create(&items).Error; err != nil {
			logger.Errorf("Failed to create order items: %v", err)
			return err
		}

		logger.Infof("Successfully created order with ID: %d", order.ID)

		return nil
	})
}

// GetByID retrieves an order by ID with items
func (r *orderRepositoryImpl) GetByID(id uint) (*entity.Order, error) {
	logger.Infof("Fetching order by ID: %d", id)
	var order entity.Order
	err := r.db.Preload("User").Preload("OrderItems.Book").First(&order, id).Error
	if err != nil {
		logger.Errorf("Failed to fetch order by ID %d: %v", id, err)
		return nil, err
	}
	logger.Infof("Successfully fetched order ID: %d for user ID: %d", order.ID, order.UserID)
	return &order, nil
}

// GetByUserID retrieves orders by user ID with pagination
func (r *orderRepositoryImpl) GetByUserID(userID uint, page, limit int) ([]*entity.Order, int64, error) {
	logger.Infof("Fetching orders for user ID: %d - page: %d, limit: %d", userID, page, limit)
	var orders []*entity.Order
	var total int64

	query := r.db.Model(&entity.Order{}).Where("user_id = ?", userID)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		logger.Errorf("Failed to count orders for user ID %d: %v", userID, err)
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve orders with pagination
	err := query.Preload("OrderItems.Book").Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		logger.Errorf("Failed to fetch orders for user ID %d: %v", userID, err)
		return nil, 0, err
	}

	logger.Infof("Successfully fetched %d orders for user ID %d out of %d total", len(orders), userID, total)
	return orders, total, nil
}

// Update updates an existing order
func (r *orderRepositoryImpl) Update(order *entity.Order) error {
	logger.Infof("Updating order with ID: %d", order.ID)
	err := r.db.Save(order).Error
	if err != nil {
		logger.Errorf("Failed to update order with ID %d: %v", order.ID, err)
		return err
	}
	logger.Infof("Successfully updated order with ID: %d", order.ID)
	return nil
}

// UpdateStatus updates order status
func (r *orderRepositoryImpl) UpdateStatus(id uint, status string) error {
	logger.Infof("Updating order status for ID %d to: %s", id, status)
	err := r.db.Model(&entity.Order{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		logger.Errorf("Failed to update order status for ID %d: %v", id, err)
		return err
	}
	logger.Infof("Successfully updated order status for ID %d to: %s", id, status)
	return nil
}

// UpdatePaymentURL updates order payment URL
func (r *orderRepositoryImpl) UpdatePaymentURL(id uint, paymentURL string) error {
	logger.Infof("Updating payment URL for order ID %d", id)
	err := r.db.Model(&entity.Order{}).Where("id = ?", id).Update("payment_url", paymentURL).Error
	if err != nil {
		logger.Errorf("Failed to update payment URL for order ID %d: %v", id, err)
		return err
	}
	logger.Infof("Successfully updated payment URL for order ID %d", id)
	return nil
}

// UpdateStatusTx updates order status using external transaction
func (r *orderRepositoryImpl) UpdateStatusTx(tx *gorm.DB, id uint, status string) error {
	logger.Infof("Updating order status for ID %d to: %s with external transaction", id, status)
	err := tx.Model(&entity.Order{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		logger.Errorf("Failed to update order status for ID %d in transaction: %v", id, err)
		return err
	}
	logger.Infof("Successfully updated order status for ID %d to: %s in transaction", id, status)
	return nil
}

// UpdatePaymentURLTx updates order payment URL using external transaction
func (r *orderRepositoryImpl) UpdatePaymentURLTx(tx *gorm.DB, id uint, paymentURL string) error {
	logger.Infof("Updating payment URL for order ID %d with external transaction", id)
	err := tx.Model(&entity.Order{}).Where("id = ?", id).Update("payment_url", paymentURL).Error
	if err != nil {
		logger.Errorf("Failed to update payment URL for order ID %d in transaction: %v", id, err)
		return err
	}
	logger.Infof("Successfully updated payment URL for order ID %d in transaction", id)
	return nil
}

// CreateOrderTx creates a new order with items using external transaction
func (r *orderRepositoryImpl) CreateOrderTx(tx *gorm.DB, order *entity.Order, items []*entity.OrderItem) error {
	logger.Infof("Creating new order with external transaction for user ID: %d with %d items", order.UserID, len(items))

	// Create the order
	if err := tx.Create(order).Error; err != nil {
		logger.Errorf("Failed to create order in transaction: %v", err)
		return err
	}
	logger.Infof("Order created successfully with ID: %d", order.ID)

	// Set order ID for all items and create them
	for _, item := range items {
		item.OrderID = order.ID
	}

	if err := tx.Create(&items).Error; err != nil {
		logger.Errorf("Failed to create order items in transaction: %v", err)
		return err
	}
	logger.Infof("Order items created successfully for order ID: %d", order.ID)

	return nil
}

// GetAll retrieves all orders with pagination (admin only)
func (r *orderRepositoryImpl) GetAll(page, limit int) ([]*entity.Order, int64, error) {
	logger.Infof("Fetching all orders with pagination - page: %d, limit: %d", page, limit)
	var orders []*entity.Order
	var total int64

	// Count total records
	if err := r.db.Model(&entity.Order{}).Count(&total).Error; err != nil {
		logger.Errorf("Failed to count all orders: %v", err)
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Retrieve orders with pagination
	err := r.db.Preload("User").Preload("OrderItems.Book").Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		logger.Errorf("Failed to fetch all orders with pagination: %v", err)
		return nil, 0, err
	}

	logger.Infof("Successfully fetched %d orders out of %d total", len(orders), total)
	return orders, total, nil
}
