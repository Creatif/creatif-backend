package sdk

func ForEach[T any](values []T, fn func(key int, val T)) {
	for key, val := range values {
		fn(key, val)
	}
}
