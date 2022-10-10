package lager

import (
	"bytes"
	"flag"
	"strings"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestLevelString(t *testing.T) {
	testCases := map[Level]string{
		Level(-42): "Level(-42)",
	}
	
	for l, v := range levelToString {
		testCases[l] = v
	}
	
	for lvl, stringLevel := range testCases {
		assert.Equal(t, stringLevel, lvl.String(), "unexpected lowercase level string.")
		assert.Equal(t, strings.ToUpper(stringLevel), lvl.CapitalString(), "unexpected all-caps level string.")
	}
}

func TestLevelText(t *testing.T) {
	testCases := []struct {
		text  string
		level Level
	}{
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"", InfoLevel}, // make the zero value useful
		{"warn", WarnLevel},
		{"error", ErrorLevel},
		{"none", NoneLevel},
		{"fatal", FatalLevel},
	}
	for _, tt := range testCases {
		if tt.text != "" {
			lvl := tt.level
			marshaled, err := lvl.MarshalText()
			assert.NoError(t, err, "unexpected error marshaling level %v to text.", &lvl)
			assert.Equal(t, tt.text, string(marshaled), "Marshaling level %v to text yielded unexpected result.", &lvl)
		}
		
		var unmarshalled Level
		err := unmarshalled.UnmarshalText([]byte(tt.text))
		assert.NoError(t, err, `unexpected error unmarshalling text %q to level.`, tt.text)
		assert.Equal(t, tt.level, unmarshalled, `Text %q unmarshalled to an unexpected level.`, tt.text)
	}
}

func TestCapitalLevelsParse(t *testing.T) {
	testCases := []struct {
		text  string
		level Level
	}{
		{"DEBUG", DebugLevel},
		{"INFO", InfoLevel},
		{"WARN", WarnLevel},
		{"ERROR", ErrorLevel},
		{"NONE", NoneLevel},
		{"FATAL", FatalLevel},
	}
	for _, tt := range testCases {
		var unmarshalled Level
		err := unmarshalled.UnmarshalText([]byte(tt.text))
		assert.NoError(t, err, `unexpected error unmarshalling text %q to level.`, tt.text)
		assert.Equal(t, tt.level, unmarshalled, `Text %q unmarshalled to an unexpected level.`, tt.text)
	}
}

func TestWeirdLevelsParse(t *testing.T) {
	testCases := []struct {
		text  string
		level Level
	}{
		// I guess...
		{"Debug", DebugLevel},
		{"Info", InfoLevel},
		{"Warn", WarnLevel},
		{"Error", ErrorLevel},
		{"None", NoneLevel},
		{"Fatal", FatalLevel},
		
		// What even is...
		{"DeBuG", DebugLevel},
		{"InFo", InfoLevel},
		{"WaRn", WarnLevel},
		{"ErRor", ErrorLevel},
		{"NoNe", NoneLevel},
		{"FaTaL", FatalLevel},
	}
	for _, tt := range testCases {
		var unmarshalled Level
		err := unmarshalled.UnmarshalText([]byte(tt.text))
		assert.NoError(t, err, `unexpected error unmarshalling text %q to level.`, tt.text)
		assert.Equal(t, tt.level, unmarshalled, `Text %q unmarshalled to an unexpected level.`, tt.text)
	}
}

func TestLevelNils(t *testing.T) {
	var l *Level
	
	// The String() method will not handle nil level properly.
	assert.Panics(t, func() {
		assert.Equal(t, "Level(nil)", l.String(), "unexpected result stringify nil *Level.")
	}, "Level(nil).String() should panic")
	
	assert.Panics(t, func() {
		_, _ = l.MarshalText()
	}, "expected to panic when marshalling a nil level.")
	
	err := l.UnmarshalText([]byte("debug"))
	assert.Equal(t, errUnmarshalNilLevel, err, "expected to error unmarshalling into a nil Level.")
}

func TestLevelUnmarshalUnknownText(t *testing.T) {
	var l Level
	err := l.UnmarshalText([]byte("foo"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unrecognized level", "expected unmarshalling arbitrary text to fail.")
}

func TestLevelAsFlagValue(t *testing.T) {
	var (
		buf bytes.Buffer
		lvl Level
	)
	fs := flag.NewFlagSet("levelTest", flag.ContinueOnError)
	fs.SetOutput(&buf)
	fs.Var(&lvl, "level", "log level")
	
	for _, expected := range []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, NoneLevel, FatalLevel} {
		assert.NoError(t, fs.Parse([]string{"-level", expected.String()}))
		assert.Equal(t, expected, lvl, "unexpected level after parsing flag.")
		assert.Equal(t, expected, lvl.Get(), "unexpected output using flag.Getter API.")
		assert.Empty(t, buf.String(), "unexpected error output parsing level flag.")
		buf.Reset()
	}
	
	assert.Error(t, fs.Parse([]string{"-level", "nope"}))
	assert.Equal(
		t,
		`invalid value "nope" for flag -level: unrecognized level: "nope"`,
		strings.Split(buf.String(), "\n")[0], // second line is help message
		"unexpected error output from invalid flag input.",
	)
}
