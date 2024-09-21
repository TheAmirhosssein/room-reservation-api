package utils_test

import (
	"testing"

	"github.com/TheAmirhosssein/room-reservation-api/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseQueryParamToInt(t *testing.T) {
	result := utils.ParseQueryParamToInt("10", 5)
	assert.Equal(t, result, 10)

	result = utils.ParseQueryParamToInt("", 5)
	assert.Equal(t, result, 5)
}
