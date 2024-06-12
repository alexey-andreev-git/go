package logger

type Logger interface {
	Fatal(string, error)
	Panic(string, error)
	Error(string, error)
	Info(string)
	Debug(string)
	Warn(string)
}
