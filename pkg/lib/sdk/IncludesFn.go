package sdk

func IncludesFn[T comparable](values []T, fn func(item T) bool) bool {
	left, right := 0, len(values)-1

	for i := left; i <= right; i++ {
		if left == right {
			return fn(values[left])
		}

		if fn(values[left]) {
			return true
		}

		if fn(values[right]) {
			return true
		}

		left++
		right--
	}

	return false
}
