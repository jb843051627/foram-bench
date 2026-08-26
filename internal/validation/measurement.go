package validation

import (
	"fmt"
	"math"
)

func Finite(value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("measurement must be finite")
	}
	return nil
}

func Positive(value float64, name string) error {
	if err := Finite(value); err != nil {
		return err
	}
	if value <= 0 {
		return fmt.Errorf("%s must be positive", name)
	}
	return nil
}

func Range(value, min, max float64, name string) error {
	if err := Finite(value); err != nil {
		return err
	}
	if value < min || value > max {
		return fmt.Errorf("%s must be between %.3f and %.3f", name, min, max)
	}
	return nil
}

func Tolerance(actual, expected, tolerance float64) error {
	if err := Finite(actual); err != nil {
		return err
	}
	if err := Finite(expected); err != nil {
		return err
	}
	if tolerance < 0 {
		return fmt.Errorf("tolerance must not be negative")
	}
	if math.Abs(actual-expected) > tolerance {
		return fmt.Errorf("value %.4f differs from %.4f by more than %.4f", actual, expected, tolerance)
	}
	return nil
}

func Percent(value float64, name string) error {
	return Range(value, 0, 100, name)
}
