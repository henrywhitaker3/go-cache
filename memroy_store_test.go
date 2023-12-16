package gocache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type demo struct {
	Data string
}

func TestMemGetStringMissingKey(t *testing.T) {
	store := NewMemoryStore()

	out, err := store.GetString(context.Background(), "bongo")

	assert.Equal(t, "", out)
	assert.Equal(t, ErrMissingKey, err)
}

func TestMemGetStringHit(t *testing.T) {
	store := NewMemoryStore()
	store.PutString(context.Background(), "bongo", "bingo", 0)

	out, err := store.GetString(context.Background(), "bongo")

	assert.Nil(t, err)
	assert.Equal(t, "bingo", out)
}

func TestMemPutString(t *testing.T) {
	store := NewMemoryStore()

	err := store.PutString(context.Background(), "bongo", "bongo", time.Second*30)

	assert.Nil(t, err)
}

func TestMemGetStructReturnsMissingKeyWhenNotInCache(t *testing.T) {
	store := NewMemoryStore()

	out := &demo{}

	err := store.GetStruct(context.Background(), "bongo", out)
	assert.ErrorIs(t, err, ErrMissingKey)
}

func TestMemGetStructReturnsStruct(t *testing.T) {
	store := NewMemoryStore()
	d := demo{Data: "bingo"}
	store.PutStruct(context.Background(), "bongo", d, 0)

	out := &demo{}

	err := store.GetStruct(context.Background(), "bongo", out)
	assert.Nil(t, err)
	assert.Equal(t, "bingo", out.Data)
}
