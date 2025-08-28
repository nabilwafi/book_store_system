package grpc

import (
	"context"

	"github.com/nabil/book-store-system/internal/service"
	"github.com/nabil/book-store-system/internal/transport/dto"
	"github.com/nabil/book-store-system/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler handles gRPC requests for user operations
type UserHandler struct {
	proto.UnimplementedUserServiceServer
	userService service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles user registration
func (h *UserHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	// Validate request using DTO
	registerDTO := &dto.RegisterRequestDTO{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := registerDTO.ValidateRegisterRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	user, err := h.userService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
	}

	return &proto.RegisterResponse{
		Success: true,
		User: &proto.User{
			Id:    uint32(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Message: "User registered successfully",
	}, nil
}

// Login handles user authentication
func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Validate request using DTO
	loginDTO := &dto.LoginRequestDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := loginDTO.ValidateLoginRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	token, user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Login failed: %v", err)
	}

	return &proto.LoginResponse{
		Success: true,
		Token:   token,
		User: &proto.User{
			Id:    uint32(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Message: "Login successful",
	}, nil
}

// GetProfile retrieves user profile
func (h *UserHandler) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.GetProfileResponse, error) {
	// Validate request using DTO
	getProfileDTO := &dto.GetProfileRequestDTO{
		Token: req.Token,
	}

	if err := getProfileDTO.ValidateGetProfileRequest(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %v", err)
	}

	user, err := h.userService.GetProfile(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Failed to get profile: %v", err)
	}

	return &proto.GetProfileResponse{
		User: &proto.User{
			Id:    uint32(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
