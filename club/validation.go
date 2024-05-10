package club

import (
	"strings"
	"unicode"
)

func ValidateClientName(name string) bool {
	return !strings.ContainsFunc(name, func(r rune) bool {
		return !(unicode.IsLower(r) && unicode.Is(unicode.Latin, r) || unicode.IsNumber(r) || r == '_')
	})
}
func ValidateTableNumber(tableNum, maxTableNum int) bool {
	return tableNum >= 1 && tableNum <= maxTableNum
}
