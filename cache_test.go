package gocache

import (
	"context"
	"testing"

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
