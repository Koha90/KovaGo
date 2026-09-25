package domain

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const maxDescriptionLength = 5_000

var ErrDescriptionTooLong = errors.New(
	"shop description is too long",
)

// Description represents a shop description.
type Description string

// NewDescription creates a validated shop description.
func NewDescription(value string) (Description, error) {
	value = strings.TrimSpace(value)

	if utf8.RuneCountInString(value) > maxDescriptionLength {
		return "", ErrDescriptionTooLong
	}

	return Description(value), nil
}
