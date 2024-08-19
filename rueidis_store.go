package gocache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/rueidis"
)

type RueidisStore struct {
	client rueidis.Client
}

func NewRueidisStore(c rueidis.Client) *RueidisStore {
	return &RueidisStore{client: c}
}

func (r *RueidisStore) GetString(ctx context.Context, key string) (string, error) {
	cmd := r.client.B().Get().Key(key).Build()
	res := r.client.Do(ctx, cmd)
	if err := res.Error(); err != nil {
		if errors.Is(err, rueidis.Nil) {
			return "", ErrMissingKey
		}
		return "", err
	}
	str, err := res.ToString()
	if err != nil {
		return "", err
	}
	return str, nil
}

func (r *RueidisStore) PutString(ctx context.Context, key string, data string, ttl time.Duration) error {
	cmd := r.client.B().Set().Key(key).Value(data).Ex(ttl).Build()
	res := r.client.Do(ctx, cmd)
	return res.Error()
}

func (r *RueidisStore) GetStruct(ctx context.Context, key string, data any) error {
	cmd := r.client.B().Get().Key(key).Build()
	res := r.client.Do(ctx, cmd)
	if err := res.Error(); err != nil {
		if errors.Is(err, rueidis.Nil) {
			return ErrMissingKey
		}
		return err
	}
	by, err := res.AsBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(by, data)
}

func (r *RueidisStore) PutStruct(ctx context.Context, key string, data any, ttl time.Duration) error {
	by, err := json.Marshal(data)
	if err != nil {
		return err
	}
	cmd := r.client.B().Set().Key(key).Value(string(by)).Ex(ttl).Build()
	res := r.client.Do(ctx, cmd)
	return res.Error()
}

func (r *RueidisStore) Forget(ctx context.Context, key string) error {
	cmd := r.client.B().Del().Key(key).Build()
	res := r.client.Do(ctx, cmd)
	return res.Error()
}

var _ Store = &RueidisStore{}
