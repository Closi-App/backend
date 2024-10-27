package logger

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
	WithField(key string, value interface{}) Logger

	Logger() interface{}
}
