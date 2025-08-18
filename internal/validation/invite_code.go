package validation

import (
	"regexp"
)

type InviteCodeValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors"`
}

func ValidateInviteCode(code string) InviteCodeValidationResult {
	var errors []string

	// Check if empty
	if code == "" {
		errors = append(errors, "Invite code is required")
		return InviteCodeValidationResult{IsValid: false, Errors: errors}
	}

	// Check length
	if len(code) < 1 || len(code) > 6 {
		errors = append(errors, "Custom invite code must be between 1 and 6 characters long")
		return InviteCodeValidationResult{IsValid: false, Errors: errors}
	}

	// Check for invalid characters
	validCodePattern := `^[A-Z0-9]+$`
	matched, _ := regexp.MatchString(validCodePattern, code)
	if !matched {
		errors = append(errors, "Invite code can only contain uppercase letters and numbers")
		return InviteCodeValidationResult{IsValid: false, Errors: errors}
	}

	return InviteCodeValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}
