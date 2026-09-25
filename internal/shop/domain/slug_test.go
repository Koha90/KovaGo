package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewSlug(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Slug
		wantErr error
	}{
		{
			name:  "valid slug",
			value: "test-slug",
			want:  "test-slug",
		},
		{
			name:  "trims spaces",
			value: "  test-slug  ",
			want:  "test-slug",
		},
		{
			name:    "empty slug",
			value:   "",
			wantErr: ErrSlugEmpty,
		},
		{
			name:    "spaces only",
			value:   "   ",
			wantErr: ErrSlugEmpty,
		},
		{
			name:    "too short slug",
			value:   "t",
			wantErr: ErrSlugInvalid,
		},
		{
			name:    "too long slug",
			value:   strings.Repeat("t", maxSlugLength+1),
			wantErr: ErrSlugInvalid,
		},
		{
			name:    "uppercase",
			value:   "Kova-Go",
			wantErr: ErrSlugInvalid,
		},
		{
			name:    "starts with hyphen",
			value:   "-kova-go",
			wantErr: ErrSlugInvalid,
		},
		{
			name:    "ends with hyphen",
			value:   "kova-go-",
			wantErr: ErrSlugInvalid,
		},
		{
			name:    "double hyphen",
			value:   "kova--go",
			wantErr: ErrSlugInvalid,
		},
		{
			name:    "underscore",
			value:   "kova_go",
			wantErr: ErrSlugInvalid,
		},
		{
			name:    "unicode",
			value:   "коваго",
			wantErr: ErrSlugInvalid,
		},
		{
			name:  "letters and digits",
			value: "kovago-2026",
			want:  "kovago-2026",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSlug(tt.value)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf(
						"NewSlug() error = %v, want %v",
						err,
						tt.wantErr,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"NewSlug() unexpected error = %v",
					err,
				)
			}

			if got != tt.want {
				t.Errorf(
					"NewSlug() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}
