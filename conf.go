package lager

const (
	// cstRFC3339Nano = "2006-01-02 15:04:05.999999999Z08:00" // only placeholder.
	
	// default console separator.
	defaultConsoleSeparator = " | "
	
	// placeholder names for the log schema.
	logPlaceholderLoggerName = "@logger"
	logPlaceholderMessage    = "@message"
	logPlaceholderInstance   = "@instance"
	logPlaceholderVer        = "@ver"
	logPlaceholderAppID      = "@app_id"
	
	// PlaceholderOutputLevelEnvName the placeholder env names for the log schema.
	PlaceholderOutputLevelEnvName = "LOG_OUTPUT_LEVEL"
)

func GetLogPlaceholderLoggerName() string {
	return logPlaceholderLoggerName
}

func GetLogPlaceholderMessage() string {
	return logPlaceholderMessage
}
