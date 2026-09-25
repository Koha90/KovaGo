package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Name
		wantErr error
	}{
		{
			name:  "valid name",
			value: "Test Name",
			want:  "Test Name",
		},
		{
			name:  "trims spaces",
			value: "  Test Name  ",
			want:  "Test Name",
		},
		{
			name:    "empty name",
			value:   "",
			wantErr: ErrNameEmpty,
		},
		{
			name:    "spaces only",
			value:   "   ",
			wantErr: ErrNameEmpty,
		},
		{
			name:    "too short unicode name",
			value:   "Я",
			wantErr: ErrNameTooShort,
		},
		{
			name:  "valid unicode name",
			value: "Тестовый Магазин",
			want:  "Тестовый Магазин",
		},
		{
			name:    "too long name",
			value:   strings.Repeat("Я", maxNameLength+1),
			wantErr: ErrNameTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewName(tt.value)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf(
						"NewName() error = %v, want %v",
						err,
						tt.wantErr,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"NewName() unexpected error = %v",
					err,
				)
			}

			if got != tt.want {
				t.Errorf(
					"NewName() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestShopRename(t *testing.T) {
	createdAt := time.Date(
		2026,
		time.September,
		24,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	updatedAt := createdAt.Add(time.Hour)

	shop, err := NewShop(
		"Old Name",
		"old-name",
		"",
		createdAt,
	)
	if err != nil {
		t.Fatalf(
			"NewShop() error = %v", err,
		)
	}

	if err := shop.Rename(
		"New Name",
		updatedAt,
	); err != nil {
		t.Fatalf("Rename() error = %v", err)
	}

	if shop.name != "New Name" {
		t.Errorf(
			"Name = %q, want %q",
			shop.name,
			"New Name",
		)
	}

	if !shop.updatedAt.Equal(updatedAt) {
		t.Errorf(
			"UpdatedAt = %v, want %v",
			shop.updatedAt,
			updatedAt,
		)
	}

	if shop.slug != "old-name" {
		t.Errorf(
			"Slug = %q, want unchanged %q",
			shop.slug,
			"old-name",
		)
	}
}
