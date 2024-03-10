package logger

type Logger interface {
	Fatal(string, error)
	Panic(string, error)
	Info(string)
	Debug(string)
	Warn(string)
}
