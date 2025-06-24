package helper

import "strings"

func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			if result.String()[len(result.String())-1] == 'I' && r == 'D' {
				result.WriteRune(r)
				continue
			}
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
