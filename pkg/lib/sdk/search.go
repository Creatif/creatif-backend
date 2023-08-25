package sdk

func SearchMap[T comparable, F any](values map[T]F, search T) F {
	var def F
	for key, value := range values {
		if key == search {
			return value
		}
	}

	return def
}

func SearchFn[T any](values []T, fn func(idx int, val T) bool) (T, int, bool) {
	var temp T
	for idx, val := range values {
		if fn(idx, val) {
			return val, idx, true
		}
	}

	return temp, -1, false
}
