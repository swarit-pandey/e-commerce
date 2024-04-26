package cache

import "errors"

var (
	ErrCacheMiss        = errors.New("cache miss")
	ErrCacheInvalidated = errors.New("cache invalid")
)

var (
	ErrEntityNotFound = errors.New("requested entity is not found in repository layer")
)
