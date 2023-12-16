package gocache

import "errors"

var (
	ErrMissingKey = errors.New("key not found")
)
