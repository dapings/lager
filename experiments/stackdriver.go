package experiments

import (
	"fmt"
	"time"
	
	"github.com/cockroachdb/errors"
	"github.com/dapings/lager"
	"go.uber.org/zap/zapcore"
)

type (
	// StackdriverLogger A strace driver logger.
	StackdriverLogger interface {
		Flush() error
		Log(entry loggingEntry)
	}
	
	// loggingEntry a logging entry.
	loggingEntry struct {
		// Timestamp is the time of the entry. If zero, the current time is used.
		Timestamp time.Time
		
		// Payload must be either a string, or something that marshals via the
		// encoding/json package to a JSON object (and not any other type of JSON value).
		Payload interface{}
		
		// Labels optionally specifies key/value labels for the log entry.
		// The StackdriverLogger.Log method takes ownership of this map.
		Labels map[string]string
		
		// LogName is the full log name, in the form "projects/{ProjectID}/logs/{LogID}".
		// It is set by the client when reading entries.
		// It is an error to set it when writing entries.
		LogName string
	}
	
	// stackdriverCore writes entries to a Logging API.
	stackdriverCore struct {
		logger       StackdriverLogger
		minimumLevel zapcore.Level
		fields       map[string]interface{}
	}
)

// Enabled impls zapcore.Core.
func (sdc *stackdriverCore) Enabled(l zapcore.Level) bool {
	return l >= sdc.minimumLevel
}

// With impls zapcore.Core.
func (sdc *stackdriverCore) With(fields []zapcore.Field) zapcore.Core {
	return &stackdriverCore{
		minimumLevel: sdc.minimumLevel,
		fields:       clone(sdc.fields, fields),
	}
}

// Check impls zapcore.Core.
func (sdc *stackdriverCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if sdc.Enabled(e.Level) {
		return ce.AddCore(e, sdc)
	}
	
	return ce
}

// Write impls zapcore.Core and
// writes a log entry to stackdriver.
func (sdc *stackdriverCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	payload := clone(sdc.fields, fields)
	payload[lager.GetLogPlaceholderLoggerName()] = e.LoggerName
	payload[lager.GetLogPlaceholderMessage()] = e.Message
	
	sdc.logger.Log(loggingEntry{
		Timestamp: e.Time,
		Payload:   payload,
	})
	
	return nil
}

// Sync impls zapcore.Core.
func (sdc *stackdriverCore) Sync() error {
	if err := sdc.logger.Flush(); err != nil {
		return errors.Wrap(err, "failed to Flush log messages")
	}
	
	return nil
}

// teeToStackdriver returns a zapcore.Core that writes the entries 
// to the provided core and the Stackdriver core.
func teeToStackdriver(baseCore zapcore.Core, logger StackdriverLogger) (zapcore.Core, lager.CloseFunc, error) {
	sdCore := &stackdriverCore{logger: logger}
	for l := zapcore.DebugLevel; l <= zapcore.FatalLevel; l++ {
		if baseCore.Enabled(l) {
			sdCore.minimumLevel = l
			break
		}
	}
	
	return zapcore.NewTee(baseCore, sdCore), func() error { return nil }, nil
}

// clone copy a new filed map.
func clone(original map[string]interface{}, newFields []zapcore.Field) map[string]interface{} {
	clone := make(map[string]interface{})
	
	for k, v := range original {
		clone[k] = v
	}
	
	for _, f := range newFields {
		switch f.Type {
		case zapcore.UintptrType:
			clone[f.Key] = uintptr(f.Integer)
		case zapcore.StringerType:
			clone[f.Key] = f.Interface.(fmt.Stringer).String()
		case zapcore.ErrorType:
			clone[f.Key] = f.Interface.(error).Error()
		case zapcore.ArrayMarshalerType, zapcore.ObjectMarshalerType, zapcore.BinaryType, zapcore.ReflectType:
			clone[f.Key] = f.Interface
		case zapcore.DurationType:
			clone[f.Key] = time.Duration(f.Integer).String()
		case zapcore.TimeType:
			clone[f.Key], _ = f.Interface.(time.Time)
		case zapcore.BoolType:
			clone[f.Key] = f.Integer == 1
		case zapcore.ByteStringType, zapcore.StringType:
			clone[f.Key] = f.String
		case zapcore.Complex64Type, zapcore.Complex128Type:
			clone[f.Key] = fmt.Sprint(f.Interface)
		case zapcore.Float64Type:
			clone[f.Key] = float64(f.Integer)
		case zapcore.Float32Type:
			clone[f.Key] = float32(f.Integer)
		case zapcore.Int64Type:
			clone[f.Key] = f.Integer
		case zapcore.Int32Type:
			clone[f.Key] = int32(f.Integer)
		case zapcore.Int16Type:
			clone[f.Key] = int16(f.Integer)
		case zapcore.Int8Type:
			clone[f.Key] = int8(f.Integer)
		case zapcore.Uint64Type:
			clone[f.Key] = uint64(f.Integer)
		case zapcore.Uint32Type:
			clone[f.Key] = uint32(f.Integer)
		case zapcore.Uint16Type:
			clone[f.Key] = uint16(f.Integer)
		case zapcore.Uint8Type:
			clone[f.Key] = uint8(f.Integer)
		case zapcore.SkipType:
			continue
		default:
			clone[f.Key] = f.Interface
		}
	}
	
	return clone
}
