// Package domain contains the shop domain model and business rules.
package domain

import (
	"errors"
	"time"
)

var ErrTimeRequired = errors.New(
	"shop time is required",
)

// ID uniquely identifies a shop.
type ID int64

// Shop represents a store managed by KovaGo.
type Shop struct {
	id          ID
	name        Name
	slug        Slug
	description Description
	status      Status
	createdAt   time.Time
	updatedAt   time.Time
}

// NewShop creates a new shop.
func NewShop(
	name string,
	slug string,
	description string,
	now time.Time,
) (Shop, error) {
	validName, err := NewName(name)
	if err != nil {
		return Shop{}, err
	}

	validSlug, err := NewSlug(slug)
	if err != nil {
		return Shop{}, err
	}

	validDescription, err := NewDescription(description)
	if err != nil {
		return Shop{}, err
	}

	if err := validateTime(now); err != nil {
		return Shop{}, err
	}

	return Shop{
		name:        validName,
		slug:        validSlug,
		description: validDescription,
		status:      StatusActive,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Rename changes the shop name.
func (s *Shop) Rename(
	value string,
	now time.Time,
) error {
	name, err := NewName(value)
	if err != nil {
		return err
	}

	if s.name == name {
		return nil
	}

	if err := validateTime(now); err != nil {
		return err
	}

	s.name = name
	s.updatedAt = now

	return nil
}

// ChangeSlug changes the shop URL identifier.
func (s *Shop) ChangeSlug(
	value string,
	now time.Time,
) error {
	slug, err := NewSlug(value)
	if err != nil {
		return err
	}

	if s.slug == slug {
		return nil
	}

	if err := validateTime(now); err != nil {
		return err
	}

	s.slug = slug
	s.updatedAt = now

	return nil
}

// ChangeDescription changes the shop description.
func (s *Shop) ChangeDescription(
	value string,
	now time.Time,
) error {
	description, err := NewDescription(value)
	if err != nil {
		return err
	}

	if s.description == description {
		return nil
	}

	if err := validateTime(now); err != nil {
		return err
	}

	s.description = description
	s.updatedAt = now

	return nil
}

// Deactivate makes the shop unavailable without removing its data.
func (s *Shop) Deactivate(now time.Time) error {
	if !s.IsActive() {
		return nil
	}

	if err := validateTime(now); err != nil {
		return err
	}

	s.status = StatusInactive
	s.updatedAt = now

	return nil
}

// Activate makes the shop operational.
func (s *Shop) Activate(now time.Time) error {
	if s.IsActive() {
		return nil
	}

	if err := validateTime(now); err != nil {
		return err
	}

	s.status = StatusActive
	s.updatedAt = now

	return nil
}

// --- Accessors ---

func (s Shop) ID() ID {
	return s.id
}

func (s Shop) Name() Name {
	return s.name
}

func (s Shop) Slug() Slug {
	return s.slug
}

func (s Shop) Description() Description {
	return s.description
}

func (s Shop) Status() Status {
	return s.status
}

func (s Shop) IsActive() bool {
	return s.status == StatusActive
}

func (s Shop) CreatedAt() time.Time {
	return s.createdAt
}

func (s Shop) UpdatedAt() time.Time {
	return s.updatedAt
}

func validateTime(value time.Time) error {
	if value.IsZero() {
		return ErrTimeRequired
	}

	return nil
}
