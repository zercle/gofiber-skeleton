package security

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeInput removes potentially dangerous characters from user input
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// HTML escape to prevent XSS
	input = html.EscapeString(input)

	return input
}

// SanitizeHTML removes HTML tags from input (allows only plain text)
func SanitizeHTML(input string) string {
	// Remove all HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

// SanitizeSQL prevents basic SQL injection patterns (use prepared statements as primary defense)
func SanitizeSQL(input string) string {
	// Remove common SQL injection patterns
	dangerous := []string{
		"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
		"DROP", "INSERT", "DELETE", "UPDATE", "CREATE", "ALTER",
		"EXEC", "EXECUTE", "UNION", "SELECT",
	}

	sanitized := input
	for _, pattern := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, pattern, "")
		sanitized = strings.ReplaceAll(sanitized, strings.ToLower(pattern), "")
	}

	return sanitized
}

// SanitizeFilename removes dangerous characters from filenames
func SanitizeFilename(filename string) string {
	// Remove path traversal attempts
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")

	// Remove null bytes
	filename = strings.ReplaceAll(filename, "\x00", "")

	// Only allow alphanumeric, dash, underscore, and dot
	re := regexp.MustCompile(`[^a-zA-Z0-9\-_\.]`)
	filename = re.ReplaceAllString(filename, "_")

	return filename
}

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
	// Simple regex for email validation
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// ValidateURL checks if URL is safe
func ValidateURL(url string) bool {
	// Check for allowed protocols
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}

	// Block localhost and private IPs in production
	localhost := []string{"localhost", "127.0.0.1", "0.0.0.0", "::1"}
	for _, host := range localhost {
		if strings.Contains(url, host) {
			return false
		}
	}

	return true
}

// StripScriptTags removes <script> tags and their content
func StripScriptTags(input string) string {
	re := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	return re.ReplaceAllString(input, "")
}

// StripEventHandlers removes inline event handlers (onclick, onerror, etc.)
func StripEventHandlers(input string) string {
	re := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*["'][^"']*["']`)
	return re.ReplaceAllString(input, "")
}

// SanitizeUserInput performs comprehensive sanitization for user-generated content
func SanitizeUserInput(input string) string {
	// Strip script tags
	input = StripScriptTags(input)

	// Strip event handlers
	input = StripEventHandlers(input)

	// HTML escape
	input = html.EscapeString(input)

	// Trim whitespace
	input = strings.TrimSpace(input)

	return input
}

// IsSafePassword checks password strength
func IsSafePassword(password string) bool {
	// Minimum length
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase, one lowercase, and one digit
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasDigit
}

// SanitizeMap sanitizes all string values in a map
func SanitizeMap(data map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})
	for key, value := range data {
		if strValue, ok := value.(string); ok {
			sanitized[key] = SanitizeInput(strValue)
		} else {
			sanitized[key] = value
		}
	}
	return sanitized
}
