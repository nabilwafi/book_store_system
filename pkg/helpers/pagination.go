package helpers

import "math"

// PaginationMetadata represents pagination information
type PaginationMetadata struct {
	CurrentPage  int32 `json:"current_page"`
	TotalPages   int32 `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

func CalculatePaginationMetadata(page, limit int, total int64) PaginationMetadata {
	// Ensure minimum values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Calculate total pages
	totalPages := int32(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	currentPage := int32(page)

	// Calculate has_next and has_previous
	hasNext := currentPage < totalPages
	hasPrevious := currentPage > 1

	return PaginationMetadata{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		HasNext:      hasNext,
		HasPrevious:  hasPrevious,
	}
}