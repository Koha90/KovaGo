package domain

import (
	"errors"
	"regexp"
	"strings"
)

const (
	minSlugLength = 2
	maxSlugLength = 120
)

var (
	ErrSlugEmpty   = errors.New("shop slug is empty")
	ErrSlugInvalid = errors.New("shop slug is invalid")
)

var slugPattern = regexp.MustCompile(
	`^[a-z0-9]+(?:-[a-z0-9]+)*$`,
)

// Slug represents a valid shop URL identifier.
type Slug string

// NewSlug creates a validated shop slug.
func NewSlug(value string) (Slug, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return "", ErrSlugEmpty
	}

	if len(value) < minSlugLength ||
		len(value) > maxSlugLength {
		return "", ErrSlugInvalid
	}

	if !slugPattern.MatchString(value) {
		return "", ErrSlugInvalid
	}

	return Slug(value), nil
}
