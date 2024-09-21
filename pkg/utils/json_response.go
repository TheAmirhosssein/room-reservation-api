package utils

import (
	"math"
)

func GenerateListResponse(dbResult any, itemsCount, pageSize, currentPage int) map[string]any {
	pageCount := int(math.Ceil(float64(itemsCount) / float64(pageSize)))
	return map[string]any{
		"page_count":   pageCount,
		"current_page": currentPage,
		"result":       dbResult,
	}
}
