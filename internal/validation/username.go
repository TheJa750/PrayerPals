package validation

import (
	"regexp"
	"strings"
)

type UsernameValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors"`
}

func ValidateUsername(username string) UsernameValidationResult {
	var errors []string

	// Trim whitespace
	username = strings.TrimSpace(username)

	// Check if empty
	if username == "" {
		errors = append(errors, "Username is required")
		return UsernameValidationResult{IsValid: false, Errors: errors}
	}

	// Minimum 2 characters
	if len(username) < 2 {
		errors = append(errors, "Username must be at least 2 characters long")
	}

	// Maximum 30 characters
	if len(username) > 25 {
		errors = append(errors, "Username must be no more than 25 characters long")
	}

	// Check for valid characters (letters, numbers, underscores, hyphens)
	validChars := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validChars.MatchString(username) {
		errors = append(errors, "Username can only contain letters, numbers, underscores, and hyphens")
	}

	return UsernameValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}
