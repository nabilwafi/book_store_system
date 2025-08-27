package grpc

import (
	"context"
	"github.com/nabil/book-store-system/internal/transport/dto"
	"github.com/nabil/book-store-system/internal/service"
	"github.com/nabil/book-store-system/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CategoryHandler handles gRPC requests for category operations
type CategoryHandler struct {
	proto.UnimplementedCategoryServiceServer
	categoryService service.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory handles category creation
func (h *CategoryHandler) CreateCategory(ctx context.Context, req *proto.CreateCategoryRequest) (*proto.CreateCategoryResponse, error) {
	// Validate request using DTO
	createDTO := &dto.CreateCategoryRequestDTO{
		Name:  req.Name,
		Token: req.Token,
	}
	
	if err := createDTO.ValidateCreateCategoryRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	category, err := h.categoryService.CreateCategory(req.Name, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create category: %v", err)
	}
	
	return &proto.CreateCategoryResponse{
		Success: true,
		Category: &proto.Category{
			Id:   uint32(category.ID),
			Name: category.Name,
		},
		Message: "Category created successfully",
	}, nil
}

// GetCategories retrieves all categories with pagination
func (h *CategoryHandler) GetCategories(ctx context.Context, req *proto.GetCategoriesRequest) (*proto.GetCategoriesResponse, error) {
	// Validate request using DTO
	getCategoriesDTO := &dto.GetCategoriesRequestDTO{
		Page:  req.Page,
		Limit: req.Limit,
	}
	
	if err := getCategoriesDTO.ValidateGetCategoriesRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	categories, total, err := h.categoryService.GetCategories(int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get categories: %v", err)
	}
	
	var protoCategories []*proto.Category
	for _, category := range categories {
		protoCategories = append(protoCategories, &proto.Category{
			Id:   uint32(category.ID),
			Name: category.Name,
		})
	}
	
	return &proto.GetCategoriesResponse{
		Success:    true,
		Message:    "Categories retrieved successfully",
		Categories: protoCategories,
		Total:      int32(total),
	}, nil
}

// GetCategory retrieves a category by ID
func (h *CategoryHandler) GetCategory(ctx context.Context, req *proto.GetCategoryRequest) (*proto.GetCategoryResponse, error) {
	// Validate request using DTO
	getCategoryDTO := &dto.GetCategoryRequestDTO{
		ID: req.Id,
	}
	
	if err := getCategoryDTO.ValidateGetCategoryRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	category, err := h.categoryService.GetCategory(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Category not found: %v", err)
	}
	
	return &proto.GetCategoryResponse{
		Success: true,
		Message: "Category retrieved successfully",
		Category: &proto.Category{
			Id:   uint32(category.ID),
			Name: category.Name,
		},
	}, nil
}

// UpdateCategory updates an existing category
func (h *CategoryHandler) UpdateCategory(ctx context.Context, req *proto.UpdateCategoryRequest) (*proto.UpdateCategoryResponse, error) {
	// Validate request using DTO
	updateDTO := &dto.UpdateCategoryRequestDTO{
		ID:    req.Id,
		Name:  req.Name,
		Token: req.Token,
	}
	
	if err := updateDTO.ValidateUpdateCategoryRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	category, err := h.categoryService.UpdateCategory(uint(req.Id), req.Name, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update category: %v", err)
	}
	
	return &proto.UpdateCategoryResponse{
		Success: true,
		Category: &proto.Category{
			Id:   uint32(category.ID),
			Name: category.Name,
		},
		Message: "Category updated successfully",
	}, nil
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*proto.DeleteCategoryResponse, error) {
	// Validate request using DTO
	deleteDTO := &dto.DeleteCategoryRequestDTO{
		ID:    req.Id,
		Token: req.Token,
	}
	
	if err := deleteDTO.ValidateDeleteCategoryRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}
	
	err := h.categoryService.DeleteCategory(uint(req.Id), req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete category: %v", err)
	}
	
	return &proto.DeleteCategoryResponse{
		Success: true,
		Message: "Category deleted successfully",
	}, nil
}