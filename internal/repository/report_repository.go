package repository

import (
	"time"

	"github.com/nabil/book-store-system/internal/entity"
	"github.com/nabil/book-store-system/pkg/logger"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetSalesReport(startDate, endDate time.Time) ([]*SalesReportItem, error)
	GetTopBooks(limit int) ([]*TopBookItem, error)
}

type SalesReportItem struct {
	Date        time.Time
	TotalSales  float64
	TotalOrders int
}

type TopBookItem struct {
	Book      entity.Book
	TotalSold int
}

type reportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepositoryImpl{
		db: db,
	}
}

// GetSalesReport retrieves sales data for reporting
func (r *reportRepositoryImpl) GetSalesReport(startDate, endDate time.Time) ([]*SalesReportItem, error) {
	logger.Info("Generating sales report from %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	var results []*SalesReportItem

	err := r.db.Model(&entity.Order{}).
		Select("DATE(created_at) as date, SUM(total_price) as total_sales, COUNT(*) as total_orders").
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "completed").
		Group("DATE(created_at)").
		Order("date").
		Scan(&results).Error

	if err != nil {
		logger.Error("Failed to generate sales report: %v", err)
		return nil, err
	}

	logger.Info("Successfully generated sales report with %d entries", len(results))
	return results, nil
}

// GetTopBooks retrieves top selling books
func (r *reportRepositoryImpl) GetTopBooks(limit int) ([]*TopBookItem, error) {
	logger.Info("Generating top books report with limit: %d", limit)
	type BookWithSales struct {
		entity.Book
		TotalSold int `gorm:"column:total_sold"`
	}

	var bookResults []BookWithSales
	err := r.db.Table("order_items").
		Select("books.*, SUM(order_items.quantity) as total_sold").
		Joins("JOIN books ON books.id = order_items.book_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status = ?", "completed").
		Group("books.id").
		Order("total_sold DESC").
		Limit(limit).
		Scan(&bookResults).Error

	if err != nil {
		logger.Error("Failed to generate top books report: %v", err)
		return nil, err
	}

	// Convert to TopBookItem format
	var results []*TopBookItem
	for _, bookResult := range bookResults {
		results = append(results, &TopBookItem{
			Book:      bookResult.Book,
			TotalSold: bookResult.TotalSold,
		})
	}

	logger.Info("Successfully generated top books report with %d entries", len(results))
	return results, nil
}
