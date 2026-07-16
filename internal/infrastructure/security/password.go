package security

import (
	"errors"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

const defaultBcryptCost = 12

var (
	ErrInvalidPassword = errors.New("invalid password")
)

// getBcryptCost returns the bcrypt cost, using a lower value in test environment
func getBcryptCost() int {
	// Check if we're in test environment
	if testCost := os.Getenv("BCRYPT_COST"); testCost != "" {
		if cost, err := strconv.Atoi(testCost); err == nil && cost >= bcrypt.MinCost && cost <= bcrypt.MaxCost {
			return cost
		}
	}
	return defaultBcryptCost
}

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrInvalidPassword
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), getBcryptCost())
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
