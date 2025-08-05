package validation

import (
	"regexp"
	"strings"
	"unicode"
)

type PasswordValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors"`
}

func ValidatePassword(password string) PasswordValidationResult {
	var errors []string

	// Trim whitespace
	password = strings.TrimSpace(password)

	// Check if password is empty
	if password == "" {
		errors = append(errors, "Password cannot be empty")
		return PasswordValidationResult{IsValid: false, Errors: errors}
	}

	// Check length
	if len(password) < 8 {
		errors = append(errors, "Password must be longer than 8 characters")
	}

	// Check for uppercase letters
	hasUpper := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}

	// Check for lowercase letters
	hasLower := false
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}

	// Check for digits
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	if !hasDigit {
		errors = append(errors, "Password must contain at least one digit")
	}

	// Check for special characters
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]`).MatchString(password)
	if !hasSpecial {
		errors = append(errors, "Password must contain at least one special character (e.g., !@#$%^&*()_+)")
	}

	return PasswordValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}
