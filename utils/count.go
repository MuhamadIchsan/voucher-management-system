package utils

func CountPages(total int64, limit int64) int64 {
	return (total + limit - 1) / limit
}
