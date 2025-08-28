package grpc

import (
	"context"

	"github.com/nabil/book-store-system/internal/transport/dto"
	"github.com/nabil/book-store-system/internal/service"
	"github.com/nabil/book-store-system/pkg/helpers"
	"github.com/nabil/book-store-system/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderHandler handles gRPC requests for order operations
type OrderHandler struct {
	proto.UnimplementedOrderServiceServer
	orderService service.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder handles order creation
func (h *OrderHandler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	// Validate request using DTO
	var dtoItems []dto.OrderItemRequestDTO
	for _, item := range req.Items {
		dtoItems = append(dtoItems, dto.OrderItemRequestDTO{
			BookID:   item.BookId,
			Quantity: item.Quantity,
		})
	}
	
	createOrderDTO := &dto.CreateOrderRequestDTO{
		Items: dtoItems,
		Token: req.Token,
	}
	
	if err := createOrderDTO.ValidateCreateOrderRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	// Convert proto items to service items
	var items []service.OrderItem
	for _, item := range req.Items {
		items = append(items, service.OrderItem{
			BookID:   uint(item.BookId),
			Quantity: int(item.Quantity),
		})
	}
	
	order, err := h.orderService.CreateOrder(items, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create order: %v", err)
	}
	
	// Convert order items to proto
	var protoItems []*proto.OrderItem
	for _, item := range order.OrderItems {
		protoItem := &proto.OrderItem{
			Id:       uint32(item.ID),
			BookId:   uint32(item.BookID),
			Quantity: int32(item.Quantity),
			Price:    item.Price,
		}
		
		if item.Book.ID != 0 {
			protoItem.Book = &proto.Book{
				Id:          uint32(item.Book.ID),
				Title:       item.Book.Title,
				Author:      item.Book.Author,
				Price:       item.Book.Price,
				Stock:       int32(item.Book.Stock),
				Year:        int32(item.Book.Year),
				CategoryId:  uint32(item.Book.CategoryID),
				ImageBase64: item.Book.ImageBase64,
			}
		}
		
		protoItems = append(protoItems, protoItem)
	}
	
	return &proto.CreateOrderResponse{
		Success: true,
		Message: "Order created successfully",
		Order: &proto.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			Items:      protoItems,
		},
	}, nil
}

// GetOrders retrieves orders for a user with pagination
func (h *OrderHandler) GetOrders(ctx context.Context, req *proto.GetOrdersRequest) (*proto.GetOrdersResponse, error) {
	// Validate request using DTO
	getOrdersDTO := &dto.GetOrdersRequestDTO{
		Token: req.Token,
		Page:  req.Page,
		Limit: req.Limit,
	}
	
	if err := getOrdersDTO.ValidateGetOrdersRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	// Use validated DTO values instead of raw request values
	orders, total, err := h.orderService.GetOrders(req.Token, int(getOrdersDTO.Page), int(getOrdersDTO.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get orders: %v", err)
	}

	var protoOrders []*proto.Order
	for _, order := range orders {
		// Convert order items to proto
		var protoItems []*proto.OrderItem
		for _, item := range order.OrderItems {
			protoItem := &proto.OrderItem{
				Id:       uint32(item.ID),
				BookId:   uint32(item.BookID),
				Quantity: int32(item.Quantity),
				Price:    item.Price,
			}
			
			if item.Book.ID != 0 {
				protoItem.Book = &proto.Book{
					Id:          uint32(item.Book.ID),
					Title:       item.Book.Title,
					Author:      item.Book.Author,
					Price:       item.Book.Price,
					Stock:       int32(item.Book.Stock),
					Year:        int32(item.Book.Year),
					CategoryId:  uint32(item.Book.CategoryID),
					ImageBase64: item.Book.ImageBase64,
				}
			}
			
			protoItems = append(protoItems, protoItem)
		}
		
		protoOrder := &proto.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			Items:      protoItems,
		}
		
		protoOrders = append(protoOrders, protoOrder)
	}
	
	// Calculate pagination metadata using validated DTO values
	paginationMeta := helpers.CalculatePaginationMetadata(int(getOrdersDTO.Page), int(getOrdersDTO.Limit), total)

	return &proto.GetOrdersResponse{
		Success:      true,
		Message:      "Orders retrieved successfully",
		Orders:       protoOrders,
		Total:        int32(total),
		CurrentPage:  paginationMeta.CurrentPage,
		TotalPages:   paginationMeta.TotalPages,
		HasNext:      paginationMeta.HasNext,
		HasPrevious:  paginationMeta.HasPrevious,
	}, nil
}

// GetOrder retrieves a specific order by ID
func (h *OrderHandler) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.GetOrderResponse, error) {
	// Validate request using DTO
	getOrderDTO := &dto.GetOrderRequestDTO{
		ID:    req.Id,
		Token: req.Token,
	}
	
	if err := getOrderDTO.ValidateGetOrderRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	order, err := h.orderService.GetOrder(uint(req.Id), req.Token)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Order not found: %v", err)
	}
	
	// Convert order items to proto
	var protoItems []*proto.OrderItem
	for _, item := range order.OrderItems {
		protoItem := &proto.OrderItem{
			Id:       uint32(item.ID),
			BookId:   uint32(item.BookID),
			Quantity: int32(item.Quantity),
			Price:    item.Price,
		}
		
		if item.Book.ID != 0 {
			protoItem.Book = &proto.Book{
				Id:          uint32(item.Book.ID),
				Title:       item.Book.Title,
				Author:      item.Book.Author,
				Price:       item.Book.Price,
				Stock:       int32(item.Book.Stock),
				Year:        int32(item.Book.Year),
				CategoryId:  uint32(item.Book.CategoryID),
				ImageBase64: item.Book.ImageBase64,
			}
		}
		
		protoItems = append(protoItems, protoItem)
	}
	
	return &proto.GetOrderResponse{
		Success: true,
		Message: "Order retrieved successfully",
		Order: &proto.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			Items:       protoItems,
		},
	}, nil
}

// UpdateOrderStatus updates order status (admin only)
func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *proto.UpdateOrderStatusRequest) (*proto.UpdateOrderStatusResponse, error) {
	// Validate request using DTO
	updateStatusDTO := &dto.UpdateOrderStatusRequestDTO{
		ID:     req.Id,
		Status: req.Status,
		Token:  req.Token,
	}
	
	if err := updateStatusDTO.ValidateUpdateOrderStatusRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	order, err := h.orderService.UpdateOrderStatus(uint(req.Id), req.Status, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update order status: %v", err)
	}
	
	return &proto.UpdateOrderStatusResponse{
		Success: true,
		Message: "Order status updated successfully",
		Order: &proto.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
		},
	}, nil
}

// ProcessPayment processes payment for an order
func (h *OrderHandler) ProcessPayment(ctx context.Context, req *proto.ProcessPaymentRequest) (*proto.ProcessPaymentResponse, error) {
	// Validate request using DTO
	processPaymentDTO := &dto.ProcessPaymentRequestDTO{
		OrderID: req.OrderId,
		Token:   req.Token,
	}
	
	if err := processPaymentDTO.ValidateProcessPaymentRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	paymentURL, err := h.orderService.ProcessPayment(uint(req.OrderId), req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to process payment: %v", err)
	}
	
	return &proto.ProcessPaymentResponse{
		Success:    true,
		PaymentUrl: paymentURL,
		Message:    "Payment URL generated successfully",
	}, nil
}