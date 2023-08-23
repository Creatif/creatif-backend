package sdk

func Search[T comparable, F any](values map[T]F, search T) F {
	var def F
	for key, value := range values {
		if key == search {
			return value
		}
	}

	return def
}
