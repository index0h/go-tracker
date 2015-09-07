package share

type LoggerInterface interface {
	Debug(args ...interface{}) LoggerInterface
	Info(args ...interface{}) LoggerInterface
	Warn(args ...interface{}) LoggerInterface
	Error(args ...interface{}) LoggerInterface
	Fatal(args ...interface{}) LoggerInterface
}
