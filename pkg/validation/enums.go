package validation

import "fmt"

func OneOf[T comparable](value T, possibleValues []T) error {
	for _, validV := range possibleValues {
		if validV == value {
			return nil
		}
	}
	return fmt.Errorf("valid values are %v", possibleValues)
}
