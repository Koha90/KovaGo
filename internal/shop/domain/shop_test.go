package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewShop(t *testing.T) {
	now := time.Date(
		2026,
		time.September,
		23,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	shop, err := NewShop(
		"Kova Store",
		"kova-store",
		"Store powered by KovaGo",
		now,
	)
	if err != nil {
		t.Fatalf("NewShop() error = %v", err)
	}

	if shop.Name() != "Kova Store" {
		t.Errorf(
			"Name = %q, want %q",
			shop.Name(),
			"Kova Store",
		)
	}

	if shop.Slug() != "kova-store" {
		t.Errorf(
			"Slug = %q, want %q",
			shop.Slug(),
			"kova-store",
		)
	}

	if shop.Description() != "Store powered by KovaGo" {
		t.Errorf(
			"Description = %q",
			shop.Description(),
		)
	}

	if !shop.CreatedAt().Equal(now) {
		t.Errorf(
			"CreatedAt = %v, want %v",
			shop.CreatedAt(),
			now,
		)
	}

	if !shop.UpdatedAt().Equal(now) {
		t.Errorf(
			"UpdatedAt = %v, want %v",
			shop.UpdatedAt(),
			now,
		)
	}
}

func TestShopDeactivate(t *testing.T) {
	createdAt := time.Date(
		2026,
		time.September,
		25,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	deactivatedAt := createdAt.Add(time.Hour)

	shop, err := NewShop(
		"Kova Store",
		"kova-store",
		"",
		createdAt,
	)
	if err != nil {
		t.Fatalf("NewShop() error = %v", err)
	}

	if !shop.IsActive() {
		t.Fatal("new shop must be active")
	}

	if err := shop.Deactivate(deactivatedAt); err != nil {
		t.Fatalf("Deactivate() error = %v", err)
	}

	if shop.IsActive() {
		t.Error("shop is active, want inactive")
	}

	if shop.Status() != StatusInactive {
		t.Errorf(
			"Status() = %q, want %q",
			shop.Status(),
			StatusInactive,
		)
	}

	if !shop.UpdatedAt().Equal(deactivatedAt) {
		t.Errorf(
			"UpdatedAt() = %v, want %v",
			shop.UpdatedAt(),
			deactivatedAt,
		)
	}
}

func TestShopDeactivateAlreadyInactive(t *testing.T) {
	createdAt := time.Date(
		2026,
		time.September,
		25,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	deactivatedAt := createdAt.Add(time.Hour)

	shop, err := NewShop(
		"Kova Store",
		"kova-store",
		"",
		createdAt,
	)
	if err != nil {
		t.Fatalf("NewShop() error = %v", err)
	}

	if !shop.IsActive() {
		t.Fatal("new shop must be active")
	}

	if err := shop.Deactivate(deactivatedAt); err != nil {
		t.Fatalf("Deactivate() error = %v", err)
	}

	updatedAt := shop.UpdatedAt()

	if shop.IsActive() {
		t.Error("shop is active, want inactive")
	}

	if shop.Status() != StatusInactive {
		t.Errorf(
			"Status() = %q, want %q",
			shop.Status(),
			StatusInactive,
		)
	}

	err = shop.Deactivate(
		updatedAt.Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("Deactivate() error = %v", err)
	}

	if !shop.UpdatedAt().Equal(updatedAt) {
		t.Error("UpdatedAt changed for inactive shop")
	}
}

func TestShopActivate(t *testing.T) {
	createdAt := time.Date(
		2026,
		time.September,
		25,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	deactivatedAt := createdAt.Add(time.Hour)
	activatedAt := deactivatedAt.Add(time.Hour)

	shop, err := NewShop(
		"Kova Store",
		"kova-store",
		"",
		createdAt,
	)
	if err != nil {
		t.Fatalf("NewShop() error = %v", err)
	}

	if err := shop.Deactivate(deactivatedAt); err != nil {
		t.Fatalf("Deactivate() error = %v", err)
	}

	if shop.IsActive() {
		t.Fatal("shop must be inactive after Deactivate()")
	}

	if err := shop.Activate(activatedAt); err != nil {
		t.Fatalf("Activate() error = %v", err)
	}

	if !shop.IsActive() {
		t.Error("shop is inactive, want active")
	}

	if shop.Status() != StatusActive {
		t.Errorf(
			"Status() = %q, want %q",
			shop.Status(),
			StatusActive,
		)
	}

	if !shop.UpdatedAt().Equal(activatedAt) {
		t.Errorf(
			"UpdatedAt() = %v, want %v",
			shop.UpdatedAt(),
			activatedAt,
		)
	}
}

func TestShopActivateAlreadyActive(t *testing.T) {
	createdAt := time.Date(
		2026,
		time.September,
		25,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	shop, err := NewShop(
		"Kova Store",
		"kova-store",
		"",
		createdAt,
	)
	if err != nil {
		t.Fatalf("NewShop() error = %v", err)
	}

	if !shop.IsActive() {
		t.Fatal("new shop must be active")
	}

	updatedAt := shop.UpdatedAt()

	if err := shop.Activate(
		createdAt.Add(time.Hour),
	); err != nil {
		t.Fatalf("Activate() error = %v", err)
	}

	if !shop.IsActive() {
		t.Error("shop is inactive, want active")
	}

	if shop.Status() != StatusActive {
		t.Errorf(
			"Status() = %q, want %q",
			shop.Status(),
			StatusActive,
		)
	}

	if !shop.UpdatedAt().Equal(updatedAt) {
		t.Errorf(
			"UpdatedAt() = %v, want unchanged %v",
			shop.UpdatedAt(),
			updatedAt,
		)
	}
}

func TestNewShopRequiresTime(t *testing.T) {
	_, err := NewShop(
		"Kova Store",
		"kova-store",
		"",
		time.Time{},
	)

	if !errors.Is(err, ErrTimeRequired) {
		t.Errorf(
			"NewShop() error = %v, want %v",
			err,
			ErrTimeRequired,
		)
	}
}
