package validation

type ValidationFunc func() error

func Validate(funcs ...ValidationFunc) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
