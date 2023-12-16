package gocache

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type MemoryStore struct {
	stringStore map[string]string
	structStore map[string]any

	stringMutex *sync.RWMutex
	structMutex *sync.RWMutex
}

var (
	_ Store = &MemoryStore{}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		stringStore: make(map[string]string),
		structStore: make(map[string]any),

		stringMutex: &sync.RWMutex{},
		structMutex: &sync.RWMutex{},
	}
}

func (m *MemoryStore) GetString(ctx context.Context, key string) (string, error) {
	m.stringMutex.RLock()
	defer m.stringMutex.RUnlock()

	if val, ok := m.stringStore[key]; ok {
		return val, nil
	} else {
		return "", ErrMissingKey
	}
}

func (m *MemoryStore) PutString(ctx context.Context, key string, data string, ttl time.Duration) error {
	m.stringMutex.Lock()
	defer m.stringMutex.Unlock()

	m.stringStore[key] = data
	return nil
}

func (m *MemoryStore) GetStruct(ctx context.Context, key string, data any) error {
	m.structMutex.RLock()
	defer m.structMutex.RUnlock()

	if val, ok := m.structStore[key]; ok {
		v := reflect.ValueOf(data).Elem()
		v.Set(reflect.ValueOf(val))
		return nil
	} else {
		return ErrMissingKey
	}
}

func (m *MemoryStore) PutStruct(ctx context.Context, key string, data any, ttl time.Duration) error {
	m.structMutex.Lock()
	defer m.structMutex.Unlock()

	m.structStore[key] = data
	return nil
}
