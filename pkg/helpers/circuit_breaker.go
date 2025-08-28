package helpers

import (
	"time"

	"github.com/nabil/book-store-system/pkg/logger"
	"github.com/sony/gobreaker"
)

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Name           string
	MaxRequests    uint32
	Interval       time.Duration
	Timeout        time.Duration
	ReadyToTrip    func(counts gobreaker.Counts) bool
	OnStateChange  func(name string, from gobreaker.State, to gobreaker.State)
}

// DefaultPaymentCircuitBreakerConfig returns default configuration for payment circuit breaker
func DefaultPaymentCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Name:        "payment-service",
		MaxRequests: 3,
		Interval:    30 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Infof("Circuit breaker '%s' changed from %s to %s", name, from, to)
		},
	}
}

// NewCircuitBreaker creates a new circuit breaker with given configuration
func NewCircuitBreaker(config CircuitBreakerConfig) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:          config.Name,
		MaxRequests:   config.MaxRequests,
		Interval:      config.Interval,
		Timeout:       config.Timeout,
		ReadyToTrip:   config.ReadyToTrip,
		OnStateChange: config.OnStateChange,
	}

	return gobreaker.NewCircuitBreaker(settings)
}

// PaymentCircuitBreaker is a global circuit breaker instance for payment operations
var PaymentCircuitBreaker *gobreaker.CircuitBreaker

// InitPaymentCircuitBreaker initializes the payment circuit breaker
func InitPaymentCircuitBreaker() {
	config := DefaultPaymentCircuitBreakerConfig()
	PaymentCircuitBreaker = NewCircuitBreaker(config)
	logger.Infof("Payment circuit breaker initialized with name: %s", config.Name)
}

// ExecuteWithCircuitBreaker executes a function with circuit breaker protection
func ExecuteWithCircuitBreaker(cb *gobreaker.CircuitBreaker, fn func() (interface{}, error)) (interface{}, error) {
	return cb.Execute(fn)
}
