package lager

import (
	"go.uber.org/atomic"
)

type (
	// LevelEnablerFunc is a convenient way to implement LevelEnabler with
	// an anonymous function.
	//
	// It's particularly useful when splitting log output between different
	// outputs (e.g., standard error and standard out). For sample code, see the
	// package-level AdvancedConfiguration example.
	LevelEnablerFunc func(Level) bool
	
	// An AtomicLevel is an atomically changeable, dynamic logging level. It lets
	// you safely change the log level of a tree of loggers (the root logger and
	// any children created by adding context) at runtime.
	//
	// The AtomicLevel itself is an http.Handler that serves a JSON endpoint to
	// alter its level.
	//
	// AtomicLevels must be created with the NewAtomicLevel constructor to allocate
	// their internal atomic pointer.
	AtomicLevel struct {
		l *atomic.Int32
	}
)

// Enabled calls the wrapped function.
func (f LevelEnablerFunc) Enabled(lvl Level) bool { return f(lvl) }

// NewAtomicLevel creates an AtomicLevel with InfoLevel and above logging
// enabled.
func NewAtomicLevel() AtomicLevel {
	return AtomicLevel{
		l: atomic.NewInt32(int32(InfoLevel)),
	}
}

// NewAtomicLevelAt is a convenience function that creates an AtomicLevel
// and then calls SetLevel with the given level.
func NewAtomicLevelAt(l Level) AtomicLevel {
	a := NewAtomicLevel()
	a.SetLevel(l)
	return a
}

// Enabled implements the LevelEnabler interface, which allows the
// AtomicLevel to be used in place of traditional static levels.
func (lvl AtomicLevel) Enabled(l Level) bool {
	return lvl.Level().Enabled(l)
}

// Level returns the minimum enabled log level.
func (lvl AtomicLevel) Level() Level {
	return Level(int8(lvl.l.Load()))
}

// SetLevel alters the logging level.
func (lvl AtomicLevel) SetLevel(l Level) {
	lvl.l.Store(int32(l))
}

// String returns the string representation of the underlying Level.
func (lvl AtomicLevel) String() string {
	return lvl.Level().String()
}

// UnmarshalText unmarshals the text to an AtomicLevel. It uses the same text
// representations as the static Levels ("debug", "info", "warn",
// "error", "dpanic", "panic", and "fatal").
func (lvl *AtomicLevel) UnmarshalText(text []byte) error {
	if lvl.l == nil {
		lvl.l = &atomic.Int32{}
	}
	
	var l Level
	if err := l.UnmarshalText(text); err != nil {
		return err
	}
	
	lvl.SetLevel(l)
	return nil
}

// MarshalText marshals the AtomicLevel to a byte slice. It uses the same
// text representation as the static Levels ("debug", "info", "warn",
// "error", "dpanic", "panic", and "fatal").
func (lvl AtomicLevel) MarshalText() (text []byte, err error) {
	return lvl.Level().MarshalText()
}
