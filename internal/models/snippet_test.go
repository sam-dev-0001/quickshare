package models

import (
	"testing"
	"time"
)

func TestSnippetIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiry    string
		createdAt time.Time
		want      bool
	}{
		{
			name:      "ExpiryNever",
			expiry:    ExpiryNever,
			createdAt: time.Now().Add(-10 * time.Hour),
			want:      false,
		},
		{
			name:      "Expiry10MinNotExpired",
			expiry:    Expiry10Min,
			createdAt: time.Now().Add(-5 * time.Minute),
			want:      false,
		},
		{
			name:      "Expiry10MinExpired",
			expiry:    Expiry10Min,
			createdAt: time.Now().Add(-11 * time.Minute),
			want:      true,
		},
		{
			name:      "Expiry1HourNotExpired",
			expiry:    Expiry1Hour,
			createdAt: time.Now().Add(-55 * time.Minute),
			want:      false,
		},
		{
			name:      "Expiry1HourExpired",
			expiry:    Expiry1Hour,
			createdAt: time.Now().Add(-65 * time.Minute),
			want:      true,
		},
		{
			name:      "Expiry1DayNotExpired",
			expiry:    Expiry1Day,
			createdAt: time.Now().Add(-23 * time.Hour),
			want:      false,
		},
		{
			name:      "Expiry1DayExpired",
			expiry:    Expiry1Day,
			createdAt: time.Now().Add(-25 * time.Hour),
			want:      true,
		},
		{
			name:      "Expiry1WeekNotExpired",
			expiry:    Expiry1Week,
			createdAt: time.Now().Add(-6 * 24 * time.Hour),
			want:      false,
		},
		{
			name:      "Expiry1WeekExpired",
			expiry:    Expiry1Week,
			createdAt: time.Now().Add(-8 * 24 * time.Hour),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Snippet{
				Expiry:    tt.expiry,
				CreatedAt: tt.createdAt,
			}
			if got := s.IsExpired(); got != tt.want {
				t.Errorf("Snippet.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnippetExpiresAt(t *testing.T) {
	now := time.Now()

	sNever := &Snippet{Expiry: ExpiryNever, CreatedAt: now}
	if !sNever.ExpiresAt().IsZero() {
		t.Errorf("Expected ExpiresAt for 'never' to be zero, got %v", sNever.ExpiresAt())
	}

	s10m := &Snippet{Expiry: Expiry10Min, CreatedAt: now}
	expected := now.Add(10 * time.Minute)
	if !s10m.ExpiresAt().Equal(expected) {
		t.Errorf("Expected ExpiresAt = %v, got %v", expected, s10m.ExpiresAt())
	}
}
