package driven

import (
	"context"
	"errors"
	"time"
)

type (
	// Connecter defines an interface for connecting to a cache store.
	Connecter interface {
		Connect() error
	}

	// Disconnecter defines an interface for disconnecting from a cache store.
	Disconnecter interface {
		Disconnect() error
	}

	// RawSetter defines an interface for setting values in a cache store.
	RawSetter interface {
		Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	}

	// RawFetcher defines an interface for fetching values from a cache store.
	RawFetcher interface {
		Fetch(ctx context.Context, key string) ([]byte, error)
	}

	// Deleter defines an interface for deleting values from a cache store.
	Deleter interface {
		Delete(ctx context.Context, key string) error
	}

	// Cache combines all the raw cache operations into a single interface.
	Cache interface {
		RawSetter
		RawFetcher
		Deleter
		Connecter
		Disconnecter
	}
)

var ErrCacheMissed = errors.New("cache missed")
