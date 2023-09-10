package sdk

func ArrEqual[T comparable](values1 []T, values2 []T) bool {
	for _, t := range values1 {
		if !Includes(values2, t) {
			return false
		}
	}

	return true
}
