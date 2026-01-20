package validations

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode"
)

// ValidateLogin checks email and password for login
func ValidateLogin(email, password string) error {
	// Validate email
	if err := validateEmail(email); err != nil {
		return err
	}

	// Validate password
	password = strings.TrimSpace(password)
	if password == "" {
		return errors.New("password is required")
	}

	// Basic password validation for login (less strict than registration)
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if !hasNumber(password) || !hasLetter(password) {
		return fmt.Errorf("password must contain at least one letter and one number")
	}

	return nil
}

// ValidateRegister checks all required fields for user registration
func ValidateRegister(username, email, password string) error {
	// Trim whitespace
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	// Check empty fields
	if username == "" {
		return errors.New("username is required")
	}
	if password == "" {
		return errors.New("password is required")
	}

	// Validate username
	if err := validateUsername(username); err != nil {
		return err
	}

	// Validate email
	if err := validateEmail(email); err != nil {
		return err
	}

	// Validate password with strict requirements for registration
	if err := validatePasswordStrict(password); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateUser checks fields for user update
func ValidateUpdateUser(username, email string) error {
	// At least one field should be provided
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if username == "" && email == "" {
		return errors.New("at least one field (username or email) must be provided for update")
	}

	// Validate username if provided
	if username != "" {
		if err := validateUsername(username); err != nil {
			return err
		}
	}

	// Validate email if provided
	if email != "" {
		if err := validateEmail(email); err != nil {
			return err
		}
	}

	return nil
}

// ValidateUserID checks if a user ID is valid UUID format
func ValidateUserID(userID string) error {
	userID = strings.TrimSpace(userID)

	if userID == "" {
		return errors.New("user ID is required")
	}

	// Basic UUID format check (36 characters with hyphens)
	if len(userID) != 36 {
		return fmt.Errorf("invalid user ID format")
	}

	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}

	_, err := mail.ParseAddress(email)

	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 || !strings.Contains(parts[1], ".") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// validateUsername checks username format and length
func validateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return fmt.Errorf("username must not exceed 50 characters")
	}

	if !isValidUsername(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

// validatePasswordStrict enforces strict password requirements for registration
func validatePasswordStrict(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(password) > 100 {
		return fmt.Errorf("password must not exceed 100 characters")
	}
	if !hasNumber(password) {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasLetter(password) {
		return fmt.Errorf("password must contain at least one letter")
	}
	if !hasSpecialChar(password) {
		return fmt.Errorf("password must contain at least one special character (@, #, $, %, etc.)")
	}

	return nil
}

// isValidUsername checks if username contains only allowed characters
func isValidUsername(username string) bool {
	for _, c := range username {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_' ||
			c == '-') {
			return false
		}
	}
	return true
}

// hasNumber checks if string contains at least one digit
func hasNumber(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

// hasLetter checks if string contains at least one letter
func hasLetter(s string) bool {
	for _, c := range s {
		if unicode.IsLetter(c) {
			return true
		}
	}
	return false
}

// hasSpecialChar checks if string contains at least one special character
func hasSpecialChar(s string) bool {
	specialChars := "@#$%^&*()_+-=[]{}|;:,.<>?/"
	for _, c := range s {
		if strings.ContainsRune(specialChars, c) {
			return true
		}
	}
	return false
}
