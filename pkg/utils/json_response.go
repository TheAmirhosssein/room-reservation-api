package utils

import (
	"math"
)

type PaginatedResponse struct {
	PageCount   int `json:"page_count"`
	CurrentPage int `json:"current_page"`
	Result      any `json:"result"`
}

func GenerateListResponse(dbResult any, itemsCount, pageSize, currentPage int) PaginatedResponse {
	pageCount := int(math.Ceil(float64(itemsCount) / float64(pageSize)))
	return PaginatedResponse{
		PageCount:   pageCount,
		CurrentPage: itemsCount,
		Result:      dbResult,
	}
}
