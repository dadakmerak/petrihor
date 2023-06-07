// package to sanitize input and output, field name, table name etc
package sanitize

import (
	"encoding/json"
	"html"
	"strings"

	"github.com/dadakmerak/petrihor/pkg/shared"
	"github.com/dadakmerak/petrihor/pkg/stringer"
)

// convert any snake_case to camelCase with sanitize
func ToCamelCase(word string) string {
	word = sanitizeHTML(word)
	word = allowedInput(word)
	return stringer.CamelCase(word)
}

// convert any camelCase to snake_case with sanitize
func ToSnakeCase(word string) string {
	word = sanitizeHTML(word)
	word = allowedInput(word)
	return stringer.SnakeCase(word)
}

func sanitizeHTML(str string) string {
	return html.EscapeString(str)
}

func unsanitizeHTML(str string) string {
	return html.UnescapeString(str)
}

// Replace special characters for field, table, schema
// except: alphabet, number, space, dash, undescore, dot (without regex because perform)
func allowedInput(str string) string {
	var result string
	str = strings.TrimSpace(str)
	for _, char := range str {
		switch {
		case (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == ' ' || char == '.' || char == '-' || char == '_' || char == '*' || char == ':':
			result += string(char)
		default:
			result += ""
		}
	}
	return result
}

func sliceContains(slices []string, look string) bool {
	for _, v := range slices {
		if v == look {
			return true
		}
	}
	return false
}

func QueryURIToMap(query string) (shared.Map, error) {
	paramMapped := make(shared.Map)
	if len(query) > 0 {
		err := json.Unmarshal([]byte(query), &paramMapped)
		if err != nil {
			return nil, err
		}
	}
	return paramMapped, nil
}

func QueryURIToInterface(query string, v any) error {
	err := json.Unmarshal([]byte(query), &v)
	if err != nil {
		return err
	}
	return nil
}

func RemoveDuplicates(slice []interface{}) []interface{} {
	exist := map[interface{}]bool{}
	result := []interface{}{}

	for _, v := range slice {
		if exist[v] {
			continue
		} else {
			exist[v] = true
			result = append(result, v)
		}
	}
	return result
}

func QueryURIToStrings(query string) ([]string, error) {
	var jsonData []string
	if query != "" {
		err := json.Unmarshal([]byte(query), &jsonData)
		if err != nil {
			return nil, err
		}
	}
	return jsonData, nil
}
