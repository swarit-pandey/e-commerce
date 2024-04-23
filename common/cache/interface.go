package cache

import (
	"context"
	"time"
)

// Repository defines a simple interface to interact with Aerospike database
type Repository[T any] interface {
	// Set will set a type T with TTL
	Set(ctx context.Context, value *T, exp time.Duration) error

	// Get will fetch the result against a given ID
	Get(ctx context.Context, id string) (T, error)
}
