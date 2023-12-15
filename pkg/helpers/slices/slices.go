package slices

func MapSlice[S any, T any](slice []S, mapper func(S) T) []T {
	result := make([]T, 0, len(slice))
	for i := range slice {
		result = append(result, mapper(slice[i]))
	}
	return result
}
