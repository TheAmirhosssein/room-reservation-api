package utils

func PageToOffset(PageNumber, pageSize int) int {
	return (PageNumber - 1) * pageSize
}
