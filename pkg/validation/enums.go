package validation

import "fmt"

func OneOfField[T comparable](field string, value T, possibleValues []T) error {
	for _, validV := range possibleValues {
		if validV == value {
			return nil
		}
	}
	return fmt.Errorf("%s has invalid value. values values are %v", field, possibleValues)
}
