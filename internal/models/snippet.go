package models

import (
	"time"
)

// Expiration options
const (
	ExpiryNever    = "never"
	Expiry10Min    = "10m"
	Expiry1Hour    = "1h"
	Expiry1Day     = "1d"
	Expiry1Week    = "1w"
)

// Snippet represents a shared text or code snippet.
type Snippet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Language  string    `json:"language"`
	Expiry    string    `json:"expiry"`
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired checks if a snippet has exceeded its expiration period.
func (s *Snippet) IsExpired() bool {
	if s.Expiry == ExpiryNever || s.Expiry == "" {
		return false
	}

	var duration time.Duration
	switch s.Expiry {
	case Expiry10Min:
		duration = 10 * time.Minute
	case Expiry1Hour:
		duration = time.Hour
	case Expiry1Day:
		duration = 24 * time.Hour
	case Expiry1Week:
		duration = 7 * 24 * time.Hour
	default:
		return false
	}

	return time.Since(s.CreatedAt) > duration
}

// ExpiresAt calculates the absolute expiration time of a snippet.
func (s *Snippet) ExpiresAt() time.Time {
	if s.Expiry == ExpiryNever || s.Expiry == "" {
		return time.Time{}
	}

	var duration time.Duration
	switch s.Expiry {
	case Expiry10Min:
		duration = 10 * time.Minute
	case Expiry1Hour:
		duration = time.Hour
	case Expiry1Day:
		duration = 24 * time.Hour
	case Expiry1Week:
		duration = 7 * 24 * time.Hour
	default:
		return time.Time{}
	}

	return s.CreatedAt.Add(duration)
}
