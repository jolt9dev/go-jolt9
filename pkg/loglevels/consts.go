package loglevels

const (
	TRACE_NAME   = "trace"
	TRACE_VALUE  = 8
	DEBUG_NAME   = "debug"
	DEBUG_VALUE  = 7
	INFO_NAME    = "info"
	INFO_VALUE   = 6
	NOTICE_NAME  = "notice"
	NOTICE_VALUE = 5
	WARN_NAME    = "warn"
	WARN_VALUE   = 4
	ERROR_NAME   = "error"
	ERROR_VALUE  = 3
	FATAL_NAME   = "fatal"
	FATAL_VALUE  = 2
)

func GetLogLevelName(level int) string {
	switch level {
	case TRACE_VALUE:
		return TRACE_NAME
	case DEBUG_VALUE:
		return DEBUG_NAME
	case INFO_VALUE:
		return INFO_NAME
	case NOTICE_VALUE:
		return NOTICE_NAME
	case WARN_VALUE:
		return WARN_NAME
	case ERROR_VALUE:
		return ERROR_NAME
	case FATAL_VALUE:
		return FATAL_NAME
	default:
		return "unknown"
	}
}

func GetLogLevelValue(name string) int {
	switch name {
	case TRACE_NAME:
		return TRACE_VALUE
	case DEBUG_NAME:
		return DEBUG_VALUE
	case INFO_NAME:
		fallthrough
	case "information":
		return INFO_VALUE
	case NOTICE_NAME:
		return NOTICE_VALUE
	case WARN_NAME:
		fallthrough
	case "warning":
		return WARN_VALUE
	case ERROR_NAME:
		return ERROR_VALUE
	case "critical":
		fallthrough
	case FATAL_NAME:
		return FATAL_VALUE
	default:
		return 0
	}
}
