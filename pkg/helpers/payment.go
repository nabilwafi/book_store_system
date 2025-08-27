package helpers

import (
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/nabil/book-store-system/config"
	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/pkg/logger"
)

// InitMidtrans initializes Midtrans client
func InitMidtrans() {
	cfg := config.LoadConfig()
	midtrans.ServerKey = cfg.MidtransKey
	midtrans.Environment = midtrans.Sandbox // Change to midtrans.Production in production
}

// CreatePayment creates a payment transaction in Midtrans
func CreatePayment(order *entity.Order, items []entity.OrderItem) (string, error) {
	// Add item details and calculate total
	var itemDetails []midtrans.ItemDetails
	var calculatedTotal int64

	for _, item := range items {
		itemPrice := int64(item.Price)
		itemQty := int32(item.Quantity)
		itemTotal := itemPrice * int64(itemQty)

		itemDetails = append(itemDetails, midtrans.ItemDetails{
			ID:    strconv.Itoa(int(item.BookID)),
			Price: itemPrice,
			Qty:   itemQty,
			Name:  item.Book.Title,
		})

		calculatedTotal += itemTotal
	}

	// Create transaction details using calculated total to ensure consistency
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(int(order.ID)),
			GrossAmt: calculatedTotal,
		},
		BcaVa: &snap.BcaVa{
			VaNumber: "123123123123",
		},
		Items: &itemDetails,
	}

	// Create Snap transaction
	resp, err := snap.CreateTransaction(req)
	if err != nil {
		logger.Errorf("Failed to create payment transaction: %v", err)
		return "", err
	}

	logger.Infof("Payment transaction created for order %d", order.ID)
	return resp.RedirectURL, nil
}

// SimulatePayment simulates payment completion (for testing)
func SimulatePayment(orderID string) bool {
	// In real implementation, this would verify with Midtrans
	// For now, we'll just simulate success
	logger.Infof("Payment simulation completed for order %s", orderID)
	return true
}
