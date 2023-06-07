package stringer

import (
	"strings"

	str "github.com/stoewer/go-strcase"
)

func Replace(baseString, newString string) string {
	dsn := strings.Replace(baseString, "tenantdb", strings.ToLower(newString), 1)
	return dsn
}

func ReplaceBracketsAndTrimComma(input string, replaceWith string) string {
	input = strings.Replace(input, "[", replaceWith, -1)
	input = strings.Replace(input, "]", replaceWith, -1)
	input = strings.TrimSuffix(input, ",")
	return input
}

func SnakeCase(input string) string {
	return strings.ReplaceAll(strings.TrimSpace(str.SnakeCase(input)), ".", "")
}

func CamelCase(input string) string {
	return strings.ReplaceAll(strings.TrimSpace(str.LowerCamelCase(input)), ".", "")
}

func RemoveCurlyBraces(s string) string {
	s = strings.ReplaceAll(s, "{{", "")
	s = strings.ReplaceAll(s, "}}", "")
	return s
}
