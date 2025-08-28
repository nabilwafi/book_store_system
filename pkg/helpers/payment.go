package helpers

import (
	"errors"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/nabil/book-store-system/config"
	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/pkg/logger"
	"github.com/sony/gobreaker"
)

func InitMidtrans() {
	cfg := config.LoadConfig()
	midtrans.ServerKey = cfg.MidtransKey
	midtrans.Environment = midtrans.Sandbox

	InitPaymentCircuitBreaker()
}

func CreatePayment(order *entity.Order, items []entity.OrderItem) (string, error) {
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

	// Execute payment creation with circuit breaker protection
	result, err := ExecuteWithCircuitBreaker(PaymentCircuitBreaker, func() (interface{}, error) {
		resp, err := snap.CreateTransaction(req)
		if err != nil {
			logger.Errorf("Failed to create payment transaction: %v", err)
			return nil, err
		}
		return resp, nil
	})

	if err != nil {
		// Check if error is from circuit breaker
		if err == gobreaker.ErrOpenState {
			logger.Errorf("Payment service circuit breaker is open - service temporarily unavailable")
			return "", errors.New("payment service temporarily unavailable, please try again later")
		}

		if err == gobreaker.ErrTooManyRequests {
			logger.Errorf("Payment service circuit breaker - too many requests")
			return "", errors.New("payment service is busy, please try again later")
		}

		return "", err
	}

	resp, ok := result.(*snap.Response)
	if !ok {
		logger.Errorf("Invalid response type from payment service")
		return "", errors.New("invalid response from payment service")
	}

	logger.Infof("Payment transaction created for order %d", order.ID)
	return resp.RedirectURL, nil
}

func SimulatePayment(orderID string) bool {
	logger.Infof("Payment simulation completed for order %s", orderID)
	return true
}
