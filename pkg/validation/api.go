package validation

import "fmt"

type ValidationFunc func() error

func Validate(funcs ...ValidationFunc) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// Fielded adds the name of the field to the possible error of the validation func.
func Fielded(field string, err error) error {
	if err != nil {
		return fmt.Errorf("%q has invalid value. %s", field, err)
	}
	return nil
}
