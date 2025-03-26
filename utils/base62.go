package utils

func IntToBase62(num int64) string {
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if num == 0 {
		return string(base62Chars[0])
	}

	base62 := ""
	for num > 0 {
		remainder := num % 62
		base62 = string(base62Chars[remainder]) + base62
		num = num / 62
	}

	// Add padding to the base62 string (6 characters)
	for len(base62 ) < 6 {
		base62 = "0" + base62
	}

	return base62
}