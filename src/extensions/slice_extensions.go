package extensions

func Contains[T comparable](items []T, predicate func(item T) bool) bool {
	for _, item := range items {
		if predicate(item) == true {
			return true
		}
	}

	return false
}
