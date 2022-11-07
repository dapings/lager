package lager

import (
	"io"
	
	"github.com/dapings/lager/experiments"
)

const (
	// DefaultScopeName default scope name for "@default".
	DefaultScopeName = "@default"
	// OverrideScopeName overriding to replace every scope.
	OverrideScopeName = "@all"
	// GrpcScopeName default gRPC scope name for "@grpc".
	GrpcScopeName = "@grpc"
	
	// DefaultOutputPath not use dev file, e.g. /dev/stdout, /dev/stderr, /dev/null.
	DefaultOutputPath = "stdout"
	// DefaultErrOutputPath not use dev file, e.g. /dev/stdout, /dev/stderr, /dev/null.
	DefaultErrOutputPath = "stderr"
)

const (
	scopeLevelSeparator = ":"
	logLevelSeparator   = ","
	
	// some default log level infos.
	defaultOutputLevel     = InfoLevel
	defaultStackTraceLevel = NoneLevel
	undefinedAppID         = ""
	
	// some default log rote infos.
	defaultRotationMaxAge     = 30
	defaultRotationMaxSize    = 100 * 1024 * 1024
	defaultRotationMaxBackups = 1000
)

type (
	CloseFunc func() error
	
	// Options the set of options supported by log kit.
	Options struct {
		// a list of file system paths to write the log data.
		// the special value: stdout, stderr, can be used to output the standard I/O stream, default: stdout.
		OutputPaths []string
		
		// a list of file system paths to write the error log data.
		// the special value: stdout, stderr, can be used to output the standard I/O stream, default: stderr.
		ErrOutputPaths []string
		
		// the rotating log file path, this file should be automatically rotated over time
		// based the RotationMaxSize, RotationMaxAge, RotationMaxBackups parameters to rotate, default not rotate.
		//
		// this path used as a foundational path. the log output is normally saved.
		// when the file is too big or too old, a rotation needs, and the file is renamed by appending a timestamp after the name.
		// once a renamed file has been created, the path resumes.
		RotateOutputPath string
		
		// the maximum log file size in megabytes before rotated.
		// default 100 megabytes.
		RotationMaxSize int
		
		// the maximum number of days to retain old log files, based the timestamp encoded in their filename.
		// default 30 days to remove the older log files.
		RotationMaxAge int
		
		// the maximum old log file number to retain. default at most 1000 log files.
		RotationMaxBackups int
		
		// whether the log is formatted as JSON.
		JSONEncoding bool
		
		// whether the log is formatted as XML.
		XMLEncoding bool
		
		// capture the grpc logs, default true.
		// not exposed by the CLI flags, mainly useful for testing.
		// even though grpc stack is closed, it hold on the logger to cases the data races.
		LogGrpc bool
		
		// a list of the specific io.Writer to write the log data.
		SpecificWriters []io.Writer
		
		// the application unique id
		appID string
		
		// can be separated by logLevelSeparator
		outputLevels     string
		stackTraceLevels string
		logCallers       string
		
		// experimental support
		// stackdriver
		useStackdriverFormat bool
		teeToStackdriver     bool
		stackdriverLogger    experiments.StackdriverLogger
		// stackdriverTargetProject string
		// stackdriverLogName       string
		
		// tee log to an UDS server
		teeToUDSServer bool
		udsSocketAddr  string
		udsServerPath  string
	}
)