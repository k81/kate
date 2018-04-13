package log

type Level uint8

type LevelName string

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

func (lvl Level) String() string {
	ln := "UNKNOWN"
	switch lvl {
	case TraceLevel:
		ln = "TRACE"
	case DebugLevel:
		ln = "DEBUG"
	case InfoLevel:
		ln = "INFO"
	case WarningLevel:
		ln = "WARNING"
	case ErrorLevel:
		ln = "ERROR"
	case FatalLevel:
		ln = "FATAL"
	}
	return ln
}

func (ln LevelName) ToLevel() Level {
	lvl := InfoLevel
	switch ln {
	case "TRACE":
		lvl = TraceLevel
	case "DEBUG":
		lvl = DebugLevel
	case "INFO":
		lvl = InfoLevel
	case "WARNING":
		lvl = WarningLevel
	case "ERROR":
		lvl = ErrorLevel
	case "FATAL":
		lvl = FatalLevel
	}
	return lvl
}
