package gocache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	memC *Cache = NewCache(NewMemoryStore())
)

func TestCacheGetStringMissingKey(t *testing.T) {
	_, err := memC.GetString(context.Background(), "bongo")
	assert.ErrorIs(t, err, ErrMissingKey)
}

func TestCacheGetValidKey(t *testing.T) {
	err := memC.PutString(context.Background(), "bongo", "bingo", 0)
	assert.Nil(t, err)

	val, err := memC.GetString(context.Background(), "bongo")
	assert.Nil(t, err)
	assert.Equal(t, "bingo", val)

	memC.Forget(context.Background(), "bongo")
}

func TestCacheGetStructMissingKey(t *testing.T) {
	out := &demo{}
	err := memC.GetStruct(context.Background(), "bongo", out)
	assert.ErrorIs(t, err, ErrMissingKey)
}

func TestCacheGetStructValidKey(t *testing.T) {
	item := demo{Data: "bingo"}
	err := memC.PutStruct(context.Background(), "bongo", item, 0)
	assert.Nil(t, err)

	out := &demo{}
	err = memC.GetStruct(context.Background(), "bongo", out)
	assert.Nil(t, err)
	assert.Equal(t, "bingo", out.Data)
}

func TestRememebrStringCallsFuncWhenNotInCache(t *testing.T) {
	cache := NewCache(memC)

	called := false

	out, err := cache.RememberString(context.Background(), "bongo", time.Second*30, func() (string, error) {
		called = true
		return "apples", nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "apples", out)
	assert.True(t, called)
}

func TestRememberStringDoesntCallFuncWhenExistsInCache(t *testing.T) {
	cache := NewCache(memC)

	called := false

	cache.PutString(context.Background(), "bongo", "apples", time.Second*30)

	out, err := cache.RememberString(context.Background(), "bongo", time.Second*30, func() (string, error) {
		called = true
		return "apples", nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "apples", out)
	assert.False(t, called)
}

func TestRememberStructCallsFuncWhenNotInCache(t *testing.T) {
	cache := NewCache(memC)

	called := false

	var out demo
	err := cache.RememberStruct(context.Background(), "bongo", &out, time.Second*30, func() (any, error) {
		called = true
		return demo{Data: "oranges"}, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "oranges", out.Data)
	assert.True(t, called)
}

func TestRememberStructDoesntCallFuncWhenInCache(t *testing.T) {
	cache := NewCache(memC)

	called := false

	assert.Nil(t, cache.PutStruct(context.Background(), "bongo", demo{Data: "pears"}, time.Second*30))

	var out demo
	err := cache.RememberStruct(context.Background(), "bongo", &out, time.Second*30, func() (any, error) {
		called = true
		return demo{Data: "pears"}, nil
	})

	assert.Nil(t, err)
	assert.Equal(t, "pears", out.Data)
	assert.False(t, called)
}
