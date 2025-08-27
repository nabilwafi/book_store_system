package dto

import (
	"errors"

	"github.com/nabil/book-store-system/pkg/helpers"
)

// OrderItemRequestDTO represents the data transfer object for order item
type OrderItemRequestDTO struct {
	BookID   uint32 `json:"book_id" validate:"required,min=1"`
	Quantity int32  `json:"quantity" validate:"required,min=1"`
}

// ValidateOrderItemRequest validates the OrderItemRequestDTO
func (o *OrderItemRequestDTO) ValidateOrderItemRequest() error {
	return helpers.ValidateStruct(o)
}

// CreateOrderRequestDTO represents the data transfer object for creating order
type CreateOrderRequestDTO struct {
	Items []OrderItemRequestDTO `json:"items" validate:"required,min=1,dive"`
	Token string                `json:"token" validate:"required"`
}

// ValidateCreateOrderRequest validates the CreateOrderRequestDTO
func (c *CreateOrderRequestDTO) ValidateCreateOrderRequest() error {
	// Check for duplicate book IDs
	bookIDs := make(map[uint32]bool)
	for _, item := range c.Items {
		if bookIDs[item.BookID] {
			return errors.New("duplicate book ID found in order items")
		}
		bookIDs[item.BookID] = true
	}
	return helpers.ValidateStruct(c)
}

// GetOrdersRequestDTO represents the data transfer object for getting orders with pagination
type GetOrdersRequestDTO struct {
	Token string `json:"token" validate:"required"`
	Page  int32  `json:"page" validate:"omitempty,min=1"`
	Limit int32  `json:"limit" validate:"omitempty,min=1,max=100"`
}

// ValidateGetOrdersRequest validates the GetOrdersRequestDTO
func (g *GetOrdersRequestDTO) ValidateGetOrdersRequest() error {
	// Set default values if not provided
	if g.Page < 1 {
		g.Page = 1
	}
	if g.Limit < 1 {
		g.Limit = 10
	}
	return helpers.ValidateStruct(g)
}

// GetOrderRequestDTO represents the data transfer object for getting single order
type GetOrderRequestDTO struct {
	ID    uint32 `json:"id" validate:"required,min=1"`
	Token string `json:"token" validate:"required"`
}

// ValidateGetOrderRequest validates the GetOrderRequestDTO
func (g *GetOrderRequestDTO) ValidateGetOrderRequest() error {
	return helpers.ValidateStruct(g)
}

// UpdateOrderStatusRequestDTO represents the data transfer object for updating order status
type UpdateOrderStatusRequestDTO struct {
	ID     uint32 `json:"id" validate:"required,min=1"`
	Status string `json:"status"`
	Token  string `json:"token" validate:"required"`
}

// ValidateUpdateOrderStatusRequest validates the UpdateOrderStatusRequestDTO
func (u *UpdateOrderStatusRequestDTO) ValidateUpdateOrderStatusRequest() error {
	return helpers.ValidateStruct(u)
}

// ProcessPaymentRequestDTO represents the data transfer object for processing payment
type ProcessPaymentRequestDTO struct {
	OrderID uint32 `json:"order_id" validate:"required,min=1"`
	Token   string `json:"token" validate:"required"`
}

// ValidateProcessPaymentRequest validates the ProcessPaymentRequestDTO
func (p *ProcessPaymentRequestDTO) ValidateProcessPaymentRequest() error {
	return helpers.ValidateStruct(p)
}
