package gocache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCacheGetStringMissingKey(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	cache := NewCache(NewRueidisStore(c))

	_, err := cache.GetString(context.Background(), "bongo")
	require.ErrorIs(t, err, ErrMissingKey)
}

func TestCacheGetValidKey(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	cache := NewCache(NewRueidisStore(c))

	err := cache.PutString(context.Background(), "bongo", "bingo", time.Minute)
	require.Nil(t, err)

	val, err := cache.GetString(context.Background(), "bongo")
	require.Nil(t, err)
	require.Equal(t, "bingo", val)

	cache.Forget(context.Background(), "bongo")
}

func TestCacheGetStructMissingKey(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	cache := NewCache(NewRueidisStore(c))

	out := &demo{}
	err := cache.GetStruct(context.Background(), "bongo", out)
	require.ErrorIs(t, err, ErrMissingKey)
}

func TestCacheGetStructValidKey(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	cache := NewCache(NewRueidisStore(c))

	item := demo{Data: "bingo"}
	err := cache.PutStruct(context.Background(), "bongo", item, time.Minute)
	require.Nil(t, err)

	out := &demo{}
	err = cache.GetStruct(context.Background(), "bongo", out)
	require.Nil(t, err)
	require.Equal(t, "bingo", out.Data)
}

func TestRememebrStringCallsFuncWhenNotInCache(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	cache := NewCache(NewRueidisStore(c))

	called := false

	out, err := cache.RememberString(context.Background(), "bongo", time.Second*30, func() (string, error) {
		called = true
		return "apples", nil
	})
	require.Nil(t, err)
	require.Equal(t, "apples", out)
	require.True(t, called)
}

func TestRememberStringDoesntCallFuncWhenExistsInCache(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	cache := NewCache(NewRueidisStore(c))

	called := false

	cache.PutString(context.Background(), "bongo", "apples", time.Second*30)

	out, err := cache.RememberString(context.Background(), "bongo", time.Second*30, func() (string, error) {
		called = true
		return "apples", nil
	})
	require.Nil(t, err)
	require.Equal(t, "apples", out)
	require.False(t, called)
}

func TestRememberStructCallsFuncWhenNotInCache(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	cache := NewCache(NewRueidisStore(c))

	called := false

	var out demo
	err := cache.RememberStruct(context.Background(), "bongo", &out, time.Second*30, func() (any, error) {
		called = true
		return demo{Data: "oranges"}, nil
	})
	require.Nil(t, err)
	require.Equal(t, "oranges", out.Data)
	require.True(t, called)
}

func TestRememberStructDoesntCallFuncWhenInCache(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	cache := NewCache(NewRueidisStore(c))

	called := false

	require.Nil(t, cache.PutStruct(context.Background(), "bongo", demo{Data: "pears"}, time.Second*30))

	var out demo
	err := cache.RememberStruct(context.Background(), "bongo", &out, time.Second*30, func() (any, error) {
		called = true
		return demo{Data: "pears"}, nil
	})

	require.Nil(t, err)
	require.Equal(t, "pears", out.Data)
	require.False(t, called)
}
