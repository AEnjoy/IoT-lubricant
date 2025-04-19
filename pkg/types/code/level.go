package code

type Level uint8

const (
	Unknown Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
	Panic
)

func (l Level) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

func (l Level) String() string {
	switch l {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	case Fatal:
		return "fatal"
	case Panic:
		return "panic"
	default:
		return "unknown"
	}
}
