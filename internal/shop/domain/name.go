package domain

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const (
	minNameLength = 2
	maxNameLength = 160
)

var (
	ErrNameEmpty    = errors.New("shop name is empty")
	ErrNameTooShort = errors.New("shop name is too short")
	ErrNameTooLong  = errors.New("shop name is too long")
)

// Name represents a valid shop name.
type Name string

// NewName creates a validated shop name.
func NewName(value string) (Name, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return "", ErrNameEmpty
	}

	length := utf8.RuneCountInString(value)

	if length < minNameLength {
		return "", ErrNameTooShort
	}

	if length > maxNameLength {
		return "", ErrNameTooLong
	}

	return Name(value), nil
}
