package dto

import (
	"github.com/nabil/book-store-system/pkg/helpers"
)

type CreateBookRequestDTO struct {
	Title       string  `json:"title" validate:"required,min=2,max=200"`
	Author      string  `json:"author" validate:"required,min=2,max=100"`
	ImageBase64 string  `json:"image_base64"`
	Price       float64 `json:"price" validate:"required,min=0.01"`
	Stock       int32   `json:"stock" validate:"required,min=0"`
	Year        int32   `json:"year" validate:"required"`
	CategoryID  uint32  `json:"category_id" validate:"required,min=1"`
	Token       string  `json:"token" validate:"required"`
}

// ValidateCreateBookRequest validates the CreateBookRequestDTO
func (c *CreateBookRequestDTO) ValidateCreateBookRequest() error {
	return helpers.ValidateStruct(c)
}

type UpdateBookRequestDTO struct {
	ID         uint32  `json:"id" validate:"required,min=1"`
	Title      string  `json:"title" validate:"required,min=2,max=200"`
	Author     string  `json:"author" validate:"required,min=2,max=100"`
	Price      float64 `json:"price" validate:"required,min=0.01"`
	Stock      int32   `json:"stock" validate:"required,min=0"`
	Year       int32   `json:"year" validate:"required"`
	CategoryID uint32  `json:"category_id" validate:"required,min=1"`
	Token      string  `json:"token" validate:"required"`
}

// ValidateUpdateBookRequest validates the UpdateBookRequestDTO
func (u *UpdateBookRequestDTO) ValidateUpdateBookRequest() error {
	return helpers.ValidateStruct(u)
}

type DeleteBookRequestDTO struct {
	ID    uint32 `json:"id" validate:"required,min=1"`
	Token string `json:"token" validate:"required"`
}

// ValidateDeleteBookRequest validates the DeleteBookRequestDTO
func (d *DeleteBookRequestDTO) ValidateDeleteBookRequest() error {
	return helpers.ValidateStruct(d)
}

type GetBooksRequestDTO struct {
	Page   int32  `json:"page" validate:"omitempty,min=1"`
	Limit  int32  `json:"limit" validate:"omitempty,min=1,max=100"`
	Search string `json:"search"`
}

// ValidateGetBooksRequest validates the GetBooksRequestDTO
func (g *GetBooksRequestDTO) ValidateGetBooksRequest() error {
	// Set default values if not provided
	if g.Page < 1 {
		g.Page = 1
	}
	if g.Limit < 1 {
		g.Limit = 10
	}
	return helpers.ValidateStruct(g)
}

type GetBookRequestDTO struct {
	ID uint32 `json:"id" validate:"required,min=1"`
}

// ValidateGetBookRequest validates the GetBookRequestDTO
func (g *GetBookRequestDTO) ValidateGetBookRequest() error {
	return helpers.ValidateStruct(g)
}

type GetBooksByCategoryRequestDTO struct {
	CategoryID uint32 `json:"category_id" validate:"required,min=1"`
	Page       int32  `json:"page" validate:"omitempty,min=1"`
	Limit      int32  `json:"limit" validate:"omitempty,min=1,max=100"`
}

// ValidateGetBooksByCategoryRequest validates the GetBooksByCategoryRequestDTO
func (g *GetBooksByCategoryRequestDTO) ValidateGetBooksByCategoryRequest() error {
	if g.Page < 1 {
		g.Page = 1
	}
	if g.Limit < 1 {
		g.Limit = 10
	}
	return helpers.ValidateStruct(g)
}
