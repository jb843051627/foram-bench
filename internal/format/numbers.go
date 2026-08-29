package format

import (
	"fmt"
	"math"
)

func Decimal(value float64, places int) string {
	if places < 0 {
		places = 0
	}
	if places > 8 {
		places = 8
	}
	return fmt.Sprintf("%.*f", places, value)
}

func Percent(value float64) string {
	return fmt.Sprintf("%s%%", Decimal(value, 1))
}

func Significant(value float64) string {
	if value == 0 || math.IsNaN(value) {
		return "0"
	}
	return fmt.Sprintf("%.4g", value)
}

func Range(value, uncertainty float64) string {
	if uncertainty <= 0 {
		return Significant(value)
	}
	return fmt.Sprintf("%s ± %s", Significant(value), Significant(uncertainty))
}
