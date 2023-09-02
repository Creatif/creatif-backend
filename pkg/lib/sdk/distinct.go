package sdk

func Distinct[T comparable](values []T) []T {
	distinct := make([]T, 0)
	for _, val := range values {
		if !Includes(distinct, val) {
			distinct = append(distinct, val)
		}
	}

	return distinct
}
