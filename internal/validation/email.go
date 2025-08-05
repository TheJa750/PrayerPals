package validation

import (
	"net/mail"
	"regexp"
	"strings"
)

type EmailValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors"`
}

func ValidateEmail(email string) EmailValidationResult {
	var errors []string

	// Trim whitespace
	email = strings.TrimSpace(email)

	// Check if email is empty
	if email == "" {
		errors = append(errors, "Email cannot be empty")
		return EmailValidationResult{IsValid: false, Errors: errors}
	}

	// Check length
	if len(email) > 254 {
		errors = append(errors, "Email cannot exceed 254 characters")
	}

	// Validate email format
	_, err := mail.ParseAddress(email)
	if err != nil {
		errors = append(errors, "Invalid email format")
	}

	// Additional checks
	if strings.Count(email, "@") != 1 {
		errors = append(errors, "Email must contain exactly one '@' symbol")
	}

	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		localPart := parts[0]
		domainPart := parts[1]

		// Check local part length
		if len(localPart) == 0 {
			errors = append(errors, "Email must have content before '@'")
		}
		if len(localPart) > 64 {
			errors = append(errors, "Local part of email is too long")
		}

		if len(domainPart) == 0 {
			errors = append(errors, "Email must have a domain")
		}

		// Check for valid domain format
		domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
		if !domainRegex.MatchString(domainPart) {
			errors = append(errors, "Invalid domain format")
		}
	}

	return EmailValidationResult{IsValid: len(errors) == 0, Errors: errors}
}
