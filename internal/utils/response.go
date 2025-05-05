package utils

import (
	"strings"
	"unicode"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// camelToSnakeCase converts a camelCase string to snake_case
func camelToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func ErrorResponse(message string, err interface{}) APIResponse {
	// Convert validation errors to a more structured format
	if errStr, ok := err.(string); ok && strings.Contains(errStr, "Field validation") {
		// Parse validation error string
		fields := make(map[string]string)
		lines := strings.Split(errStr, "\n")

		for _, line := range lines {
			if strings.Contains(line, "Field validation") {
				parts := strings.Split(line, "Error:")
				if len(parts) >= 2 {
					// Extract the field name properly
					keyParts := strings.Split(parts[0], "'")
					if len(keyParts) >= 2 {
						// The field path is in format 'StructName.FieldName'
						fieldPath := keyParts[1]
						// Extract only the field name, remove struct name if present
						fieldName := fieldPath
						if strings.Contains(fieldPath, ".") {
							pathParts := strings.Split(fieldPath, ".")
							fieldName = pathParts[len(pathParts)-1] // Get the last part which is the field name
						}

						// Get JSON tag format (already in snake_case)
						jsonFieldName := camelToSnakeCase(fieldName)

						tag := ""
						if strings.Contains(parts[1], "the '") {
							tagParts := strings.Split(parts[1], "the '")
							if len(tagParts) >= 2 {
								tag = strings.Split(tagParts[1], "'")[0]
							}
						} else {
							// Fallback if format is different
							tag = "invalid"
						}

						var msg string
						switch tag {
						case "required":
							msg = jsonFieldName + " is required"
						case "email":
							msg = jsonFieldName + " must be a valid email address"
						case "min":
							msg = jsonFieldName + " must be at least minimum length"
						case "max":
							msg = jsonFieldName + " must not exceed maximum length"
						default:
							msg = jsonFieldName + " is invalid"
						}

						fields[jsonFieldName] = msg
					}
				}
			}
		}

		// Convert to array format for response
		var errors []map[string]string
		for field, message := range fields {
			errors = append(errors, map[string]string{
				"field":   field,
				"message": message,
			})
		}

		return APIResponse{
			Success: false,
			Message: message,
			Error:   errors,
		}
	}

	// Handle regular errors
	return APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}
