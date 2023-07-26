package helper

// Filter filters a collection.
func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// First finds the first item in the collection.
func First[T any](ss []T, test func(T) bool) (ret T) {
	for _, s := range ss {
		if test(s) {
			ret = s
			break
		}
	}
	return
}

// Has checks if the collection has the item.
func Has[T any](ss []T, test func(T) bool) (ret bool) {
	for _, s := range ss {
		if test(s) {
			ret = true
		}
	}
	return
}
