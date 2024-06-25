package config

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrModeInvalid = errors.New("invalid mode")
	ErrModeUnknown = errors.New("unknown mode")
)

type Mode int

const (
	ModeDev Mode = iota + 1
	ModeProd
)

func (m Mode) String() string {
	b, err := m.MarshalText()
	if err != nil {
		return ""
	}

	return string(b)
}

func (m Mode) MarshalText() ([]byte, error) {
	switch m {
	case ModeDev:
		return []byte("DEV"), nil
	case ModeProd:
		return []byte("PROD"), nil
	default:
		return nil, ErrModeInvalid
	}
}

func (m *Mode) UnmarshalText(text []byte) error {
	t := string(text)
	switch strings.ToUpper(t) {
	case "DEV":
		*m = ModeDev
		return nil
	case "PROD":
		*m = ModeProd
		return nil
	default:
		return fmt.Errorf("string '%s': %w", t, ErrModeUnknown)
	}
}
