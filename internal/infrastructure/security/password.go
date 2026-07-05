package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

var (
	ErrInvalidPassword = errors.New("invalid password")
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrInvalidPassword
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckPassword compares a plain text password with a hashed password
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength validates password strength requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check for at least one number
	hasNumber := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
			break
		}
	}

	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}
