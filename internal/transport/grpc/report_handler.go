package grpc

import (
	"context"
	"github.com/nabil/book-store-system/internal/transport/dto"
	"github.com/nabil/book-store-system/internal/service"
	"github.com/nabil/book-store-system/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"time"
)

// ReportHandler handles gRPC requests for report operations
type ReportHandler struct {
	proto.UnimplementedReportServiceServer
	reportService service.ReportService
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GetSalesReport generates sales report for a date range (admin only)
func (h *ReportHandler) GetSalesReport(ctx context.Context, req *proto.GetSalesReportRequest) (*proto.GetSalesReportResponse, error) {
	// Validate request using DTO
	salesReportDTO := &dto.GetSalesReportRequestDTO{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Token:     req.Token,
	}
	
	if err := salesReportDTO.ValidateGetSalesReportRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid start date format: %v", err)
	}
	
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid end date format: %v", err)
	}
	
	// Get sales report
	reportItems, totalSales, err := h.reportService.GetSalesReport(startDate, endDate, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get sales report: %v", err)
	}
	
	// Convert to proto format
	var protoItems []*proto.SalesReportItem
	for _, item := range reportItems {
		protoItems = append(protoItems, &proto.SalesReportItem{
			Date:        item.Date.Format("2006-01-02"),
			TotalSales:  item.TotalSales,
			TotalOrders: int32(item.TotalOrders),
		})
	}
	
	return &proto.GetSalesReportResponse{
		Success:      true,
		Message:      "Sales report retrieved successfully",
		Report:       protoItems,
		TotalRevenue: totalSales,
	}, nil
}

// GetTopBooks retrieves top selling books (admin only)
func (h *ReportHandler) GetTopBooks(ctx context.Context, req *proto.GetTopBooksRequest) (*proto.GetTopBooksResponse, error) {
	// Validate request using DTO
	topBooksDTO := &dto.GetTopBooksRequestDTO{
		Limit: req.Limit,
		Token: req.Token,
	}
	
	if err := topBooksDTO.ValidateGetTopBooksRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	topBooks, err := h.reportService.GetTopBooks(int(req.Limit), req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get top books: %v", err)
	}
	
	// Convert to proto format
	var protoItems []*proto.TopBookItem
	for _, item := range topBooks {
		protoBook := &proto.Book{
			Id:          uint32(item.Book.ID),
			Title:       item.Book.Title,
			Author:      item.Book.Author,
			Price:       item.Book.Price,
			Stock:       int32(item.Book.Stock),
			Year:        int32(item.Book.Year),
			CategoryId:  uint32(item.Book.CategoryID),
			ImageBase64: item.Book.ImageBase64,
		}
		
		if item.Book.Category.ID != 0 {
			protoBook.Category = &proto.Category{
				Id:   uint32(item.Book.Category.ID),
				Name: item.Book.Category.Name,
			}
		}
		
		protoItems = append(protoItems, &proto.TopBookItem{
			Book:      protoBook,
			TotalSold: int32(item.TotalSold),
		})
	}
	
	return &proto.GetTopBooksResponse{
		Success: true,
		Message: "Top books retrieved successfully",
		Books:   protoItems,
	}, nil
}