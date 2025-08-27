package dto

import (
	"errors"
	"time"

	"github.com/nabil/book-store-system/pkg/helpers"
)

// GetSalesReportRequestDTO represents the data transfer object for getting sales report
type GetSalesReportRequestDTO struct {
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
	Token     string `json:"token" validate:"required"`
}

// ValidateGetSalesReportRequest validates the GetSalesReportRequestDTO
func (g *GetSalesReportRequestDTO) ValidateGetSalesReportRequest() error {
	// First validate using struct tags
	if err := helpers.ValidateStruct(g); err != nil {
		return err
	}

	// Validate date format (YYYY-MM-DD)
	startDate, err := time.Parse("2006-01-02", g.StartDate)
	if err != nil {
		return errors.New("start_date must be in YYYY-MM-DD format")
	}

	endDate, err := time.Parse("2006-01-02", g.EndDate)
	if err != nil {
		return errors.New("end_date must be in YYYY-MM-DD format")
	}

	// Validate that start_date is not after end_date
	if startDate.After(endDate) {
		return errors.New("start_date cannot be after end_date")
	}

	// Validate that dates are not in the future
	now := time.Now()
	if startDate.After(now) {
		return errors.New("start_date cannot be in the future")
	}
	if endDate.After(now) {
		return errors.New("end_date cannot be in the future")
	}

	// Validate date range (max 1 year)
	maxRange := 365 * 24 * time.Hour
	if endDate.Sub(startDate) > maxRange {
		return errors.New("date range cannot exceed 1 year")
	}

	return nil
}

// GetTopBooksRequestDTO represents the data transfer object for getting top books
type GetTopBooksRequestDTO struct {
	Limit int32  `json:"limit" validate:"omitempty,min=1,max=100"`
	Token string `json:"token" validate:"required"`
}

// ValidateGetTopBooksRequest validates the GetTopBooksRequestDTO
func (g *GetTopBooksRequestDTO) ValidateGetTopBooksRequest() error {
	// Set default value if not provided
	if g.Limit < 1 {
		g.Limit = 10
	}

	return helpers.ValidateStruct(g)
}
