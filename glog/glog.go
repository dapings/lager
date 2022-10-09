package glog

import (
	"fmt"
	"os"
	
	"go.uber.org/zap"
)

// Package glog exposes an API subset of the [glog](https://github.com/golang/glog) package.
// All logging state delivered to this package is shunted to the global [zap logger](https://github.com/uber-go/zap).

type (
	Level int32
	
	Verbose bool
)

func Flush() {
	_ = zap.L().Sync()
}

func V(Level) Verbose {
	return Verbose(zap.L().Core().Enabled(zap.DebugLevel))
}

func (v Verbose) Info(args ...interface{}) {
	zap.S().Debug(args...)
}

func (v Verbose) Infoln(args ...interface{}) {
	zap.S().Debug(fmt.Sprint(args...), '\n')
}

func (v Verbose) Infof(format string, args ...interface{}) {
	zap.S().Debugf(format, args...)
}

func Info(args ...interface{}) {
	zap.S().Info(args...)
}

func InfoDepth(depth int, args ...interface{}) {
	zap.S().Info(args...)
}

func Infoln(args ...interface{}) {
	zap.S().Info(fmt.Sprint(args...), '\n')
}

func Infof(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}

func Warning(args ...interface{}) {
	zap.S().Warn(args...)
}

func WarningDepth(depth int, args ...interface{}) {
	zap.S().Warn(args...)
}

func Warningln(args ...interface{}) {
	zap.S().Warn(fmt.Sprint(args...), '\n')
}

func Warningf(format string, args ...interface{}) {
	zap.S().Warnf(format, args...)
}

func Error(args ...interface{}) {
	zap.S().Error(args...)
}

func ErrorDepth(depth int, args ...interface{}) {
	zap.S().Error(args...)
}

func Errorln(args ...interface{}) {
	zap.S().Error(fmt.Sprint(args...), '\n')
}

func Errorf(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	zap.S().Error(args...)
	os.Exit(255)
}

func FatalDepth(depth int, args ...interface{}) {
	zap.S().Error(args...)
	os.Exit(255)
}

func Fatalln(args ...interface{}) {
	zap.S().Error(fmt.Sprint(args...), '\n')
	os.Exit(255)
}

func Fatalf(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
	os.Exit(255)
}

func Exit(args ...interface{}) {
	zap.S().Error(args...)
	os.Exit(1)
}

func ExitDepth(depth int, args ...interface{}) {
	zap.S().Error(fmt.Sprint(args...), '\n')
	os.Exit(1)
}

func Exitln(args ...interface{}) {
	zap.S().Error(fmt.Sprint(args...), '\n')
	os.Exit(1)
}

func Exitf(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
	os.Exit(1)
}
