package validation

import "fmt"

func MinLen[T any](slice []T, minLen int) error {
	if len(slice) < minLen {
		return fmt.Errorf("min length is %d", minLen)
	}
	return nil
}

func ValidItems[T any](slice []T, validator func(T) error) error {
	for i, item := range slice {
		if err := validator(item); err != nil {
			return fmt.Errorf("item at index %d invalid. %v", i, err)
		}
	}
	return nil
}
