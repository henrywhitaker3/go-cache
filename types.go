package gocache

import (
	"context"
	"time"
)

type Store interface {
	GetString(ctx context.Context, key string) (string, error)
	PutString(ctx context.Context, key string, data string, ttl time.Duration) error

	GetStruct(ctx context.Context, key string, data any) error
	PutStruct(ctx context.Context, key string, data any, ttl time.Duration) error

	Forget(ctx context.Context, key string) error
}

type CacheStructFunc func() (any, error)
type CacheStringFunc func() (string, error)
