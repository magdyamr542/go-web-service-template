package pointers

func DefaultIfNil[T any](input *T) T {
	if input != nil {
		return *input
	}
	return Default[T]()
}

func Default[T any]() T {
	var value T
	return value
}
