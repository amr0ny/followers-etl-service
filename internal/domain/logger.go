package domain

type Logger interface {
	Info(format string, a ...any)
	Error(format string, a ...any)
	Debug(format string, a ...any)
}
