package user

import (
	"errors"
	"regexp"
	"strings"
)

// Email value object
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return Email{}, errors.New("invalid email format")
	}
	return Email{value: email}, nil
}

func (e Email) String() string {
	return e.value
}

// Password value object
type Password struct {
	hash string
}

func NewPassword(hash string) (Password, error) {
	if hash == "" {
		return Password{}, errors.New("password hash cannot be empty")
	}
	return Password{hash: hash}, nil
}

func (p Password) Hash() string {
	return p.hash
}

// Profile value object
type Profile struct {
	FirstName string
	LastName  string
	Phone     string
	Avatar    string
}

func NewProfile(firstName, lastName, phone, avatar string) Profile {
	return Profile{
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Phone:     strings.TrimSpace(phone),
		Avatar:    strings.TrimSpace(avatar),
	}
}

func (p Profile) FullName() string {
	return strings.TrimSpace(p.FirstName + " " + p.LastName)
}

// Status enum
type Status string

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusSuspended Status = "suspended"
	StatusDeleted   Status = "deleted"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusSuspended, StatusDeleted:
		return true
	}
	return false
}

// Role enum
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin:
		return true
	}
	return false
}
