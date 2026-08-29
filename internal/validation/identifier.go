package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var identifierPattern = regexp.MustCompile(`^[A-Z][A-Z0-9-]{2,31}$`)

func Identifier(value string) error {
	if !identifierPattern.MatchString(strings.TrimSpace(value)) {
		return fmt.Errorf("identifier %q must be uppercase and use letters, digits, hyphens", value)
	}
	return nil
}

func SiteCode(value string) error {
	value = strings.TrimSpace(value)
	if len(value) < 3 || len(value) > 16 || strings.ContainsAny(value, " /\\") {
		return fmt.Errorf("invalid site code %q", value)
	}
	return nil
}

func Material(value string) error {
	value = strings.TrimSpace(value)
	if len([]rune(value)) < 3 || len([]rune(value)) > 240 {
		return fmt.Errorf("material description length is outside 3..240")
	}
	return nil
}

func Operator(value string) error {
	value = strings.TrimSpace(value)
	if value == "" || len([]rune(value)) > 80 {
		return fmt.Errorf("operator is empty or too long")
	}
	return nil
}
