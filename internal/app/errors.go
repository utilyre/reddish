package app

import "errors"

var (
	ErrNoRecord = errors.New("no record")
	ErrExpired  = errors.New("expired")
)
