package validation

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func Max[T constraints.Integer](value T, maxValue T, includeMax bool) error {
	if (includeMax && value > maxValue) || (!includeMax && value >= maxValue) {
		return fmt.Errorf("value %d bigger than max value %d", value, maxValue)
	}
	return nil
}

func Min[T constraints.Integer](value T, minValue T, includeMin bool) error {
	if (includeMin && value < minValue) || (!includeMin && value <= minValue) {
		return fmt.Errorf("value %d smaller than min value %d", value, minValue)
	}
	return nil
}
