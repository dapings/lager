package lager

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
)

type (
	// Level an enum of all supported log levels.
	Level int
)

// Enable logging level: fatal, error, warning, info, debug.
const (
	// NoneLevel disable logging.
	NoneLevel Level = iota

	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var (
	levelToString = map[Level]string{
		DebugLevel: "debug",
		InfoLevel:  "info",
		WarnLevel:  "warn",
		ErrorLevel: "error",
		FatalLevel: "fatal",
		NoneLevel:  "none",
	}

	stringToLevel = map[string]Level{
		"debug": DebugLevel,
		"info":  InfoLevel,
		"warn":  WarnLevel,
		"error": ErrorLevel,
		"fatal": FatalLevel,
		"none":  NoneLevel,
	}

	errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level")
)

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	if v, ok := levelToString[l]; ok {
		return v
	}

	return fmt.Sprintf("Level(%d)", l)
}

// CapitalString returns an all-caps ASCII representation of the log level.
func (l Level) CapitalString() string {
	return strings.ToUpper(l.String())
}

// MarshalText marshals the Level to text. Note that the text representation
// drops the -Level suffix (see example).
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText unmarshal text to a level. Like MarshalText, UnmarshalText
// expects the text representation of a Level to drop the -Level suffix.
//
// In particular, this makes it easy to configure logging levels using YAML,
// TOML, or JSON files.
func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return errors.Errorf("unrecognized level: %q", text)
	}
	return nil
}

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO", "": // make the zero value useful
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	case "none", "NONE":
		*l = NoneLevel
	default:
		return false
	}
	return true
}

// Set sets the level for the flag.Value interface.
func (l *Level) Set(s string) error {
	return l.UnmarshalText([]byte(s))
}

// Get gets the level for the flag.Getter interface.
func (l *Level) Get() interface{} {
	return *l
}

// Enabled returns true if the given level is at or above this level.
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}
