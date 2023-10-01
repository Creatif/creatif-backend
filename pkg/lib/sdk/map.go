package sdk

func Map[T any, F any](values []T, fn func(idx int, value T) F) []F {
	t := make([]F, len(values))

	for i, item := range values {
		t[i] = fn(i, item)
	}

	return t
}
