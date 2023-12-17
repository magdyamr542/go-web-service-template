package validation

import "fmt"

func MinLenField[T any](field string, slice []T, minLen int) error {
	if len(slice) < minLen {
		return fmt.Errorf("%q min length is %d", field, minLen)
	}
	return nil
}

func ValidItemsField[T any](field string, slice []T, validator func(T) error) error {
	for i, item := range slice {
		if err := validator(item); err != nil {
			return fmt.Errorf("%q invalid. item at index %d invalid. %v", field, i, err)
		}
	}
	return nil
}
