package log

type Logger interface {
	Log(*Entry)
}
