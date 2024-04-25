package cache

import (
	"context"
	"time"
)

// Repository defines a simple interface to interact with Aerospike database
type Repository[T any] interface {
	// Set will set a type T with TTL
	Set(ctx context.Context, key any, value T, exp time.Duration) error

	// Get will fetch the result against a given ID
	Get(ctx context.Context, key any) (T, error)

	// Delete will delete the object for the given ID
	Delete(ctx context.Context, key any) error
}
