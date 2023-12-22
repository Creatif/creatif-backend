package sdk

func Replace[T comparable](values []T, search T, replace T) []T {
	newT := make([]T, len(values))
	for i, val := range values {
		if val == search {
			newT[i] = replace
		} else {
			newT[i] = val
		}
	}

	return newT
}
