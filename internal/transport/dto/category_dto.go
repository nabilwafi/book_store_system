package dto

import (
	"github.com/nabil/book-store-system/pkg/helpers"
)

type CreateCategoryRequestDTO struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Token string `json:"token" validate:"required"`
}

// ValidateCreateCategoryRequest validates the CreateCategoryRequestDTO
func (c *CreateCategoryRequestDTO) ValidateCreateCategoryRequest() error {
	return helpers.ValidateStruct(c)
}

type UpdateCategoryRequestDTO struct {
	ID    uint32 `json:"id" validate:"required,min=1"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Token string `json:"token" validate:"required"`
}

// ValidateUpdateCategoryRequest validates the UpdateCategoryRequestDTO
func (u *UpdateCategoryRequestDTO) ValidateUpdateCategoryRequest() error {
	return helpers.ValidateStruct(u)
}

type DeleteCategoryRequestDTO struct {
	ID    uint32 `json:"id" validate:"required,min=1"`
	Token string `json:"token" validate:"required"`
}

// ValidateDeleteCategoryRequest validates the DeleteCategoryRequestDTO
func (d *DeleteCategoryRequestDTO) ValidateDeleteCategoryRequest() error {
	return helpers.ValidateStruct(d)
}

type GetCategoriesRequestDTO struct {
	Page  int32 `json:"page" validate:"omitempty,min=1"`
	Limit int32 `json:"limit" validate:"omitempty,min=1,max=100"`
}

// ValidateGetCategoriesRequest validates the GetCategoriesRequestDTO
func (g *GetCategoriesRequestDTO) ValidateGetCategoriesRequest() error {
	// Set default values if not provided
	if g.Page < 1 {
		g.Page = 1
	}
	if g.Limit < 1 {
		g.Limit = 10
	}
	return helpers.ValidateStruct(g)
}

type GetCategoryRequestDTO struct {
	ID uint32 `json:"id" validate:"required,min=1"`
}

// ValidateGetCategoryRequest validates the GetCategoryRequestDTO
func (g *GetCategoryRequestDTO) ValidateGetCategoryRequest() error {
	return helpers.ValidateStruct(g)
}