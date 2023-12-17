package validation

import "fmt"

func OneOfField[T comparable](field string, value T, possibleValues []T) error {
	for _, validV := range possibleValues {
		if validV == value {
			return nil
		}
	}
	return fmt.Errorf("%q has invalid value. Valid values are %v", field, possibleValues)
}
