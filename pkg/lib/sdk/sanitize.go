package sdk

func Sanitize[T any](values []T, fn func(k int, v T) T) []T {
	sanitized := make([]T, 0)

	for key, v := range values {
		sanitized = append(sanitized, fn(key, v))
	}

	return sanitized
}
