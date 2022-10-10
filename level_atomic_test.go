package lager

import (
	"sync"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestLevelEnablerFunc(t *testing.T) {
	enabled := LevelEnablerFunc(func(l Level) bool { return l == InfoLevel })
	testCases := []struct {
		level   Level
		enabled bool
	}{
		{DebugLevel, false},
		{InfoLevel, true},
		{WarnLevel, false},
		{ErrorLevel, false},
		{NoneLevel, false},
		{FatalLevel, false},
	}
	for _, tt := range testCases {
		assert.Equal(t, tt.enabled, enabled.Enabled(tt.level), "unexpected result applying LevelEnablerFunc to %s", tt.level)
	}
}

func TestNewAtomicLevel(t *testing.T) {
	lvl := NewAtomicLevel()
	assert.Equal(t, InfoLevel, lvl.Level(), "unexpected initial level.")
	lvl.SetLevel(ErrorLevel)
	assert.Equal(t, ErrorLevel, lvl.Level(), "unexpected level after SetLevel.")
	lvl = NewAtomicLevelAt(WarnLevel)
	assert.Equal(t, WarnLevel, lvl.Level(), "unexpected level after SetLevel.")
}

func TestAtomicLevelMutation(t *testing.T) {
	lvl := NewAtomicLevel()
	lvl.SetLevel(WarnLevel)
	// Trigger races for non-atomic level mutations.
	proceed := make(chan struct{})
	wg := &sync.WaitGroup{}
	runConcurrently(10, 100, wg, func() {
		<-proceed
		assert.Equal(t, WarnLevel, lvl.Level())
	})
	runConcurrently(10, 100, wg, func() {
		<-proceed
		lvl.SetLevel(WarnLevel)
	})
	close(proceed)
	wg.Wait()
}

func TestAtomicLevelText(t *testing.T) {
	testCases := []struct {
		text   string
		expect Level
		err    bool
	}{
		{"debug", DebugLevel, false},
		{"info", InfoLevel, false},
		{"", InfoLevel, false},
		{"warn", WarnLevel, false},
		{"error", ErrorLevel, false},
		{"none", NoneLevel, false},
		{"fatal", FatalLevel, false},
		{"foobar", NoneLevel, true},
	}
	
	for _, tt := range testCases {
		var lvl AtomicLevel
		// Test both initial unmarshalling and overwriting existing value.
		for i := 0; i < 2; i++ {
			if tt.err {
				assert.Error(t, lvl.UnmarshalText([]byte(tt.text)), "expected unmarshalling %q to fail.", tt.text)
			} else {
				assert.NoError(t, lvl.UnmarshalText([]byte(tt.text)), "expected unmarshalling %q to succeed.", tt.text)
			}
			assert.Equal(t, tt.expect, lvl.Level(), "unexpected level after unmarshalling.")
			lvl.SetLevel(NoneLevel)
		}
		
		// Test marshalling
		if tt.text != "" && !tt.err {
			lvl.SetLevel(tt.expect)
			marshaled, err := lvl.MarshalText()
			assert.NoError(t, err, `unexpected error marshalling level "%v" to text.`, tt.expect)
			assert.Equal(t, tt.text, string(marshaled), "expected marshaled text to match")
			assert.Equal(t, tt.text, lvl.String(), "expected Stringer call to match")
		}
	}
}
