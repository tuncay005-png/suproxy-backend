package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/infrastructure/security"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

func TestPasswordService_HashPassword(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Success_ValidPassword", "SecurePassword123", false},
		{"Success_LongPassword", "VeryLongPasswordWith123Numbers!@#$%", false},
		{"Success_SpecialChars", "P@ssw0rd!@#$%^&*()", false},
		{"Failure_EmptyPassword", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := security.HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash)
				assert.True(t, len(hash) > 50) // BCrypt hashes are long
			}
		})
	}
}

func TestPasswordService_CheckPassword(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	t.Run("Success_CorrectPassword", func(t *testing.T) {
		password := "SecurePassword123"
		hash, err := security.HashPassword(password)
		require.NoError(t, err)

		err = security.CheckPassword(hash, password)
		assert.NoError(t, err)
	})

	t.Run("Failure_WrongPassword", func(t *testing.T) {
		password := "SecurePassword123"
		hash, err := security.HashPassword(password)
		require.NoError(t, err)

		err = security.CheckPassword(hash, "WrongPassword123")
		assert.Error(t, err)
	})

	t.Run("Failure_EmptyPassword", func(t *testing.T) {
		password := "SecurePassword123"
		hash, err := security.HashPassword(password)
		require.NoError(t, err)

		err = security.CheckPassword(hash, "")
		assert.Error(t, err)
	})

	t.Run("Failure_InvalidHash", func(t *testing.T) {
		err := security.CheckPassword("invalid-hash", "password")
		assert.Error(t, err)
	})
}

func TestPasswordService_HashConsistency(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	t.Run("DifferentHashesForSamePassword", func(t *testing.T) {
		password := "SecurePassword123"

		hash1, err := security.HashPassword(password)
		require.NoError(t, err)

		hash2, err := security.HashPassword(password)
		require.NoError(t, err)

		// Hashes should be different due to salt
		assert.NotEqual(t, hash1, hash2)

		// But both should validate correctly
		assert.NoError(t, security.CheckPassword(hash1, password))
		assert.NoError(t, security.CheckPassword(hash2, password))
	})
}

func TestPasswordService_ValidatePasswordStrength(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{"Success_ValidPassword", "Password123", false, ""},
		{"Success_LongPassword", "VeryLongPassword123", false, ""},
		{"Success_WithSpecialChars", "P@ssw0rd123", false, ""},
		{"Failure_TooShort", "Pass1", true, "at least 8 characters"},
		{"Failure_NoNumber", "PasswordOnly", true, "contain at least one number"},
		{"Failure_Empty", "", true, "at least 8 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := security.ValidatePasswordStrength(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordService_HashAndValidateFlow(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	t.Run("CompleteFlow_ValidPassword", func(t *testing.T) {
		password := "SecurePassword123"

		// Validate strength
		err := security.ValidatePasswordStrength(password)
		require.NoError(t, err)

		// Hash password
		hash, err := security.HashPassword(password)
		require.NoError(t, err)

		// Check password
		err = security.CheckPassword(hash, password)
		assert.NoError(t, err)
	})

	t.Run("CompleteFlow_WeakPassword", func(t *testing.T) {
		password := "weak"

		// Validate strength should fail
		err := security.ValidatePasswordStrength(password)
		assert.Error(t, err)
	})
}

func TestPasswordService_PasswordCaseSensitivity(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	t.Run("CaseSensitive", func(t *testing.T) {
		password := "SecurePassword123"
		hash, err := security.HashPassword(password)
		require.NoError(t, err)

		// Correct case should work
		err = security.CheckPassword(hash, password)
		assert.NoError(t, err)

		// Different case should fail
		err = security.CheckPassword(hash, "securepassword123")
		assert.Error(t, err)

		err = security.CheckPassword(hash, "SECUREPASSWORD123")
		assert.Error(t, err)
	})
}

func TestPasswordService_SpecialCharacters(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	specialPasswords := []string{
		"Pass123!@#$%^&*()",
		"Пароль123",      // Cyrillic
		"密码Password123",  // Chinese
		"パスワード123",       // Japanese
		"🔒Password123🔑",  // Emoji
		"Tab\tPass\n123", // Whitespace chars
	}

	for _, password := range specialPasswords {
		t.Run(password, func(t *testing.T) {
			hash, err := security.HashPassword(password)
			require.NoError(t, err)

			err = security.CheckPassword(hash, password)
			assert.NoError(t, err)
		})
	}
}

func TestPasswordService_LongPasswords(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	tests := []struct {
		name   string
		length int
	}{
		{"Length_50", 50},
		{"Length_100", 100},
		{"Length_200", 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate password with numbers
			password := ""
			for i := 0; i < tt.length-3; i++ {
				password += "a"
			}
			password += "123" // Add numbers for validation

			hash, err := security.HashPassword(password)
			require.NoError(t, err)

			err = security.CheckPassword(hash, password)
			assert.NoError(t, err)
		})
	}
}
