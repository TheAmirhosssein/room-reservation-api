package utils

import "strconv"

func ParseQueryParamToInt(queryParam string, defaultNumber int) int {
	intPageNumber, _ := strconv.ParseInt(queryParam, 10, 64)
	if intPageNumber == 0 {
		intPageNumber = int64(defaultNumber)
	}
	return int(intPageNumber)
}
