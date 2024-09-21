package utils_test

import (
	"reflect"
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
)

func TestGenerateListResponse(t *testing.T) {
	dbResult := []string{"item1", "item2", "item3"}
	itemsCount := 10
	pageSize := 3
	currentPage := 2

	expected := map[string]any{
		"page_count":   4,
		"current_page": 2,
		"result":       dbResult,
	}

	result := utils.GenerateListResponse(dbResult, itemsCount, pageSize, currentPage)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	itemsCount = 9
	pageSize = 3
	currentPage = 1

	expected = map[string]any{
		"page_count":   3,
		"current_page": 1,
		"result":       dbResult,
	}

	result = utils.GenerateListResponse(dbResult, itemsCount, pageSize, currentPage)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	itemsCount = 2
	pageSize = 10
	currentPage = 1

	expected = map[string]any{
		"page_count":   1,
		"current_page": 1,
		"result":       dbResult,
	}

	result = utils.GenerateListResponse(dbResult, itemsCount, pageSize, currentPage)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
