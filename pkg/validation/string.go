package validation

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errEmptyStr = errors.New("empty string")
)

func RequiredStr(value string) error {
	if strings.TrimSpace(value) == "" {
		return errEmptyStr
	}
	return nil
}

func RequiredStrField(field, value string) error {
	if err := RequiredStr(value); err != nil {
		return fmt.Errorf("%q is required but is empty", field)
	}
	return nil
}
