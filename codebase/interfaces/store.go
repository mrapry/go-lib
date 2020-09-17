package interfaces

import (
	"context"
	"time"
)

// Store abstract interface
type Store interface {
	Get(ctx context.Context, key string) (string, error)
	GetKeys(ctx context.Context, pattern string) ([]string, error)
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
}
