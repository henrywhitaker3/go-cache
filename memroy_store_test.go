package gocache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type demo struct {
	Data string
}

func TestMemGetStringMissingKey(t *testing.T) {
	store := NewMemoryStore()

	out, err := store.GetString(context.Background(), "bongo")

	require.Equal(t, "", out)
	require.Equal(t, ErrMissingKey, err)
}

func TestMemGetStringHit(t *testing.T) {
	store := NewMemoryStore()
	store.PutString(context.Background(), "bongo", "bingo", 0)

	out, err := store.GetString(context.Background(), "bongo")

	require.Nil(t, err)
	require.Equal(t, "bingo", out)
}

func TestMemPutString(t *testing.T) {
	store := NewMemoryStore()

	err := store.PutString(context.Background(), "bongo", "bongo", time.Second*30)

	require.Nil(t, err)
}

func TestMemGetStructReturnsMissingKeyWhenNotInCache(t *testing.T) {
	store := NewMemoryStore()

	out := &demo{}

	err := store.GetStruct(context.Background(), "bongo", out)
	require.ErrorIs(t, err, ErrMissingKey)
}

func TestMemGetStructReturnsStruct(t *testing.T) {
	store := NewMemoryStore()
	d := demo{Data: "bingo"}
	store.PutStruct(context.Background(), "bongo", d, 0)

	out := &demo{}

	err := store.GetStruct(context.Background(), "bongo", out)
	require.Nil(t, err)
	require.Equal(t, "bingo", out.Data)
}

func TestMemForgetKey(t *testing.T) {
	store := NewMemoryStore()
	store.PutString(context.Background(), "bongo", "bingo", time.Second*30)

	g, err := store.GetString(context.Background(), "bongo")
	require.Nil(t, err)
	require.Equal(t, "bingo", g)

	store.Forget(context.Background(), "bongo")

	_, err = store.GetString(context.Background(), "bongo")
	require.ErrorIs(t, err, ErrMissingKey)
}
