package gocache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisGetStringMissingKey(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	mock.ExpectGet("bongo").RedisNil()

	out, err := store.GetString(context.Background(), "bongo")

	require.Equal(t, "", out)
	require.Equal(t, ErrMissingKey, err)
}

func TestRedisGetStringHit(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	mock.ExpectGet("bongo").SetVal("bingo")

	out, err := store.GetString(context.Background(), "bongo")

	require.Nil(t, err)
	require.Equal(t, "bingo", out)
}

func TestRedisPutString(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	mock.ExpectSetEx("bongo", "bongo", time.Second*30).SetVal("bongo")

	err := store.PutString(context.Background(), "bongo", "bongo", time.Second*30)

	require.Nil(t, err)
}

func TestRedisGetStructReturnsMissingKeyWhenNotInCache(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	mock.ExpectGet("bongo").RedisNil()

	out := &demo{}

	err := store.GetStruct(context.Background(), "bongo", out)
	require.ErrorIs(t, err, ErrMissingKey)
}

func TestRedisGetStructReturnsStruct(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	d := demo{Data: "bingo"}

	b, err := json.Marshal(d)
	require.Nil(t, err)

	mock.ExpectGet("bongo").SetVal(string(b))

	out := &demo{}

	err = store.GetStruct(context.Background(), "bongo", out)
	require.Nil(t, err)

	require.Equal(t, "bingo", out.Data)
}

func TestRedisPutStruct(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	item := &demo{Data: "bingo"}

	mock.ExpectSetEx("bongo", item, time.Second*30).SetVal("OK")

	store.PutStruct(context.Background(), "bongo", item, time.Second*30)
}

func TestRedisForgetKey(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	mock.ExpectDel("bongo")

	store.Forget(context.Background(), "bongo")
}
