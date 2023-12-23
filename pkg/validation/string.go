package validation

import (
	"errors"
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
