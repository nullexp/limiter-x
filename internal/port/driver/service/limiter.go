package service

import (
	"context"
)

// RateLimiter defines the interface for rate-limiting service operations
type RateLimiter interface {
	// RateLimit checks if a request is allowed for a specific user based on the rate limit
	RateLimit(ctx context.Context, userId string, limit int) (bool, error)

	// GetUserRateLimit fetches the current rate limit configuration for a specific user
	GetUserRateLimit(ctx context.Context, userId string) (*RateLimitModel, error)

	// UpdateUserRateLimit updates the rate limit for a specific user
	UpdateUserRateLimit(ctx context.Context, userId string, newLimit int) error
}

// RateLimitModel holds the rate limit configuration for a user
type RateLimitModel struct {
	UserId    string // The ID of the user
	Limit     int    // The rate limit for the user (e.g., 100 requests per second)
	Remaining int    // The remaining number of requests the user can make in the current window
	Window    string // The time window for the rate limit (e.g., "10 seconds")
}
