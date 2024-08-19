package gocache

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRueidisGetStringMissingKey(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	store := NewRueidisStore(c)

	out, err := store.GetString(context.Background(), "bongo")

	require.Equal(t, "", out)
	require.Equal(t, ErrMissingKey, err)
}

func TestRueidisGetStringHit(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	store := NewRueidisStore(c)

	cmd := c.B().Set().Key("bongo").Value("bingo").Build()
	require.Nil(t, c.Do(context.Background(), cmd).Error())

	out, err := store.GetString(context.Background(), "bongo")

	require.Nil(t, err)
	require.Equal(t, "bingo", out)
}

func TestRueidisPutString(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	store := NewRueidisStore(c)

	err := store.PutString(context.Background(), "bongo", "bongo", time.Second*30)

	require.Nil(t, err)

	cmd := c.B().Get().Key("bongo").Build()
	require.Nil(t, c.Do(context.Background(), cmd).Error())
}

func TestRueidisGetStructReturnsMissingKeyWhenNotInCache(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	store := NewRueidisStore(c)

	out := &demo{}

	err := store.GetStruct(context.Background(), "bongo", out)
	require.ErrorIs(t, err, ErrMissingKey)
}

func TestRueidisGetStructReturnsStruct(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	store := NewRueidisStore(c)

	d := demo{Data: "bingo"}

	b, err := json.Marshal(d)
	require.Nil(t, err)

	cmd := c.B().Set().Key("bongo").Value(string(b)).Build()
	require.Nil(t, c.Do(context.Background(), cmd).Error())

	out := &demo{}

	err = store.GetStruct(context.Background(), "bongo", out)
	require.Nil(t, err)

	require.Equal(t, "bingo", out.Data)
}

func TestRueidisPutStruct(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	store := NewRueidisStore(c)

	item := &demo{Data: "bingo"}

	store.PutStruct(context.Background(), "bongo", item, time.Second*30)

	cmd := c.B().Get().Key("bongo").Build()
	require.Nil(t, c.Do(context.Background(), cmd).Error())
}

func TestRueidisForgetKey(t *testing.T) {
	c, cancel := newRueidis(t)
	defer cancel()
	store := NewRueidisStore(c)

	cmd := c.B().Set().Key("bongo").Value("bingo").Build()
	require.Nil(t, c.Do(context.Background(), cmd).Error())

	store.Forget(context.Background(), "bongo")

	cmd = c.B().Get().Key("bongo").Build()
	require.NotNil(t, c.Do(context.Background(), cmd).Error())
}

func newRueidis(t *testing.T) (rueidis.Client, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	redis, err := redis.Run(ctx, "redis:latest")
	require.Nil(t, err)

	host, err := redis.Host(ctx)
	require.Nil(t, err)
	port, err := redis.MappedPort(ctx, nat.Port("6379/tcp"))
	require.Nil(t, err)
	conn := fmt.Sprintf("%s:%d", host, port.Int())

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{conn},
	})
	require.Nil(t, err)
	return client, func() {
		redis.Terminate(ctx)
		cancel()
	}
}
