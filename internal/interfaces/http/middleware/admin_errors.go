package middleware

import "errors"

var (
	// ErrMissingUserID indicates user ID is missing from context
	ErrMissingUserID = errors.New("user ID not found in context")

	// ErrInvalidUserID indicates user ID format is invalid
	ErrInvalidUserID = errors.New("invalid user ID format")

	// ErrMissingEmail indicates email is missing from context
	ErrMissingEmail = errors.New("email not found in context")

	// ErrMissingRole indicates role is missing from context
	ErrMissingRole = errors.New("role not found in context")

	// ErrNotAdmin indicates user does not have admin role
	ErrNotAdmin = errors.New("admin role required")

	// ErrAdminSelfDemotion indicates admin trying to demote themselves
	ErrAdminSelfDemotion = errors.New("cannot demote yourself")
)
