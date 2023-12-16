package gocache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

var (
	_ Store = &RedisStore{}
)

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func (r *RedisStore) GetString(ctx context.Context, key string) (string, error) {
	output := r.client.Get(ctx, key)
	if output.Err() != nil {
		if output.Err() == redis.Nil {
			return "", ErrMissingKey
		}
		return "", fmt.Errorf("cahce GET err: %w", output.Err())
	}
	return output.Val(), nil
}

func (r *RedisStore) PutString(ctx context.Context, key string, data string, ttl time.Duration) error {
	out := r.client.SetEx(ctx, key, data, ttl)
	return out.Err()
}

func (r *RedisStore) GetStruct(ctx context.Context, key string, data any) error {
	output := r.client.Get(ctx, key)
	if output.Err() != nil {
		if output.Err() == redis.Nil {
			return ErrMissingKey
		}
		return fmt.Errorf("cahce GET err: %w", output.Err())
	}

	b, err := output.Bytes()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, data); err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) PutStruct(ctx context.Context, key string, data any, ttl time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	out := r.client.SetEx(ctx, key, string(b), ttl)
	return out.Err()
}

func (r *RedisStore) Forget(ctx context.Context, key string) error {
	out := r.client.Del(ctx, key)
	return out.Err()
}
