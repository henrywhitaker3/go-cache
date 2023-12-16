package gocache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisGetStringMissingKey(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	mock.ExpectGet("bongo").RedisNil()

	out, err := store.GetString(context.Background(), "bongo")

	assert.Equal(t, "", out)
	assert.Equal(t, ErrMissingKey, err)
}

func TestRedisGetStringHit(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	mock.ExpectGet("bongo").SetVal("bingo")

	out, err := store.GetString(context.Background(), "bongo")

	assert.Nil(t, err)
	assert.Equal(t, "bingo", out)
}

func TestRedisPutString(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	mock.ExpectSetEx("bongo", "bongo", time.Second*30).SetVal("bongo")

	err := store.PutString(context.Background(), "bongo", "bongo", time.Second*30)

	assert.Nil(t, err)
}

func TestRedisGetStructReturnsMissingKeyWhenNotInCache(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	mock.ExpectGet("bongo").RedisNil()

	out := &demo{}

	err := store.GetStruct(context.Background(), "bongo", out)
	assert.ErrorIs(t, err, ErrMissingKey)
}

func TestRedisGetStructReturnsStruct(t *testing.T) {
	c, mock := redismock.NewClientMock()
	store := NewRedisStore(c)

	d := demo{Data: "bingo"}

	b, err := json.Marshal(d)
	assert.Nil(t, err)

	mock.ExpectGet("bongo").SetVal(string(b))

	out := &demo{}

	err = store.GetStruct(context.Background(), "bongo", out)
	assert.Nil(t, err)

	assert.Equal(t, "bingo", out.Data)
}
