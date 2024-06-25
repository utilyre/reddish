package domain

import (
	"errors"
	"fmt"
)

var (
	ErrEmpty = errors.New("empty")
)

type Key string

func NewKey(v string) (Key, error) {
	if len(v) == 0 {
		return "", fmt.Errorf("key: string '%s': %w", v, ErrEmpty)
	}

	return Key(v), nil
}

func (k Key) String() string {
	return string(k)
}

type Val []byte

func NewVal(v []byte) (Val, error) {
	return Val(v), nil
}

func (v Val) String() string {
	return string([]byte(v))
}
