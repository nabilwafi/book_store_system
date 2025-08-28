package service

import (
	"github.com/nabil/book-store-system/internal/repository"
	"github.com/nabil/book-store-system/pkg/logger"
	"github.com/nabil/book-store-system/pkg/middleware"
	"time"
)

type ReportService interface {
	GetSalesReport(startDate, endDate time.Time, token string) ([]*repository.SalesReportItem, float64, error)
	GetTopBooks(limit int, token string) ([]*repository.TopBookItem, error)
	GetBookPriceStatistics(token string) (*repository.BookPriceStatistics, error)
}

type reportServiceImpl struct {
	reportRepo repository.ReportRepository
	userRepo   repository.UserRepository
	auth       *middleware.AuthMiddleware
}

func NewReportService(reportRepo repository.ReportRepository, userRepo repository.UserRepository) ReportService {
	auth := middleware.NewAuthMiddleware(userRepo)
	return &reportServiceImpl{
		reportRepo: reportRepo,
		userRepo:   userRepo,
		auth:       auth,
	}
}

// GetSalesReport generates sales report for a date range (admin only)
func (s *reportServiceImpl) GetSalesReport(startDate, endDate time.Time, token string) ([]*repository.SalesReportItem, float64, error) {
	logger.Info("Generating sales report", "startDate", startDate.Format("2006-01-02"), "endDate", endDate.Format("2006-01-02"))
	
	// Validate admin token
	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Sales report generation failed - invalid admin token", "error", err)
		return nil, 0, err
	}
	
	// Get sales report
	reportItems, err := s.reportRepo.GetSalesReport(startDate, endDate)
	if err != nil {
		logger.Error("Failed to get sales report data", "startDate", startDate.Format("2006-01-02"), "endDate", endDate.Format("2006-01-02"), "error", err)
		return nil, 0, err
	}
	
	// Calculate total sales
	var totalSales float64
	for _, item := range reportItems {
		totalSales += item.TotalSales
	}
	
	logger.Info("Sales report generation successful", "itemCount", len(reportItems), "totalSales", totalSales)
	return reportItems, totalSales, nil
}

// GetTopBooks retrieves top selling books (admin only)
func (s *reportServiceImpl) GetTopBooks(limit int, token string) ([]*repository.TopBookItem, error) {
	logger.Info("Getting top selling books", "limit", limit)
	
	// Validate admin token
	_, err := s.auth.ValidateAdminToken(token)
	if err != nil {
		logger.Error("Top books retrieval failed - invalid admin token", "error", err)
		return nil, err
	}
	
	topBooks, err := s.reportRepo.GetTopBooks(limit)
	if err != nil {
		logger.Error("Failed to get top books data", "limit", limit, "error", err)
		return nil, err
	}
	
	logger.Info("Top books retrieval successful", "count", len(topBooks), "limit", limit)
	return topBooks, nil
}

// GetBookPriceStatistics retrieves book price statistics with admin authentication
func (s *reportServiceImpl) GetBookPriceStatistics(token string) (*repository.BookPriceStatistics, error) {
	logger.Info("Book price statistics request received")
	
	// Validate admin token
	if _, err := s.auth.ValidateAdminToken(token); err != nil {
		logger.Error("Admin token validation failed", "error", err)
		return nil, err
	}
	
	// Get book price statistics from repository
	stats, err := s.reportRepo.GetBookPriceStatistics()
	if err != nil {
		logger.Error("Failed to get book price statistics", "error", err)
		return nil, err
	}
	
	logger.Info("Book price statistics retrieval successful", "total_books", stats.TotalBooks)
	return stats, nil
}