package gocache

import (
	"context"
	"errors"
	"time"
)

type Cache struct {
	store Store
}

func NewCache(store Store) *Cache {
	return &Cache{
		store: store,
	}
}

func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	return c.store.GetString(ctx, key)
}

func (c *Cache) PutString(ctx context.Context, key string, data string, ttl time.Duration) error {
	return c.store.PutString(ctx, key, data, ttl)
}

func (c *Cache) RememberString(ctx context.Context, key string, ttl time.Duration, f CacheStringFunc) (string, error) {
	if val, err := c.GetString(ctx, key); err != nil {
		// Only return the error if it's not a missing key error (we want to run f when there's nothing in the cache)
		if !errors.Is(err, ErrMissingKey) {
			return "", err
		}
	} else {
		// When there's no error, we got a cache hit so exit
		return val, nil
	}

	item, err := f()
	if err != nil {
		return "", err
	}

	if err := c.PutString(ctx, key, item, ttl); err != nil {
		return "", err
	}
	return c.GetString(ctx, key)
}

func (c *Cache) GetStruct(ctx context.Context, key string, data any) error {
	return c.store.GetStruct(ctx, key, data)
}

func (c *Cache) PutStruct(ctx context.Context, key string, data any, ttl time.Duration) error {
	return c.store.PutStruct(ctx, key, data, ttl)
}

func (c *Cache) RememberStruct(ctx context.Context, key string, data any, ttl time.Duration, f CacheStructFunc) error {
	if err := c.GetStruct(ctx, key, data); err != nil {
		// Only return the error if it's not a missing key error (we want to run f when there's nothing in the cache)
		if !errors.Is(err, ErrMissingKey) {
			return err
		}
	} else {
		// When there's no error, we got a cache hit so exit
		return nil
	}

	item, err := f()
	if err != nil {
		return err
	}

	if err := c.PutStruct(ctx, key, item, ttl); err != nil {
		return err
	}
	return c.GetStruct(ctx, key, data)
}

func (c *Cache) Forget(ctx context.Context, key string) error {
	return c.store.Forget(ctx, key)
}
