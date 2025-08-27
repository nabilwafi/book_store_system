package dto

import (
	"github.com/nabil/book-store-system/pkg/helpers"
)

// RegisterRequestDTO represents the data transfer object for user registration
type RegisterRequestDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// ValidateRegisterRequest validates the RegisterRequestDTO
func (r *RegisterRequestDTO) ValidateRegisterRequest() error {
	return helpers.ValidateStruct(r)
}

// LoginRequestDTO represents the data transfer object for user login
type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ValidateLoginRequest validates the LoginRequestDTO
func (l *LoginRequestDTO) ValidateLoginRequest() error {
	return helpers.ValidateStruct(l)
}

// GetProfileRequestDTO represents the data transfer object for getting user profile
type GetProfileRequestDTO struct {
	Token string `json:"token" validate:"required"`
}

// ValidateGetProfileRequest validates the GetProfileRequestDTO
func (g *GetProfileRequestDTO) ValidateGetProfileRequest() error {
	return helpers.ValidateStruct(g)
}