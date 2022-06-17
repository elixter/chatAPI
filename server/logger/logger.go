package logger

var log Logger

type Logger interface {
	Debug(i ...interface{})
	Debugf(format string, args ...interface{})
	Error(i ...interface{})
	Errorf(format string, args ...interface{})
	Info(i ...interface{})
	Infof(format string, args ...interface{})
	Print(i ...interface{})
	Printf(format string, args ...interface{})
}

func SetLogger(logger Logger) {
	log = logger
}

func Debug(i ...interface{}) {
	log.Debug(i)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func Error(i ...interface{}) {
	log.Error(i)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}

func Info(i ...interface{}) {
	log.Info(i)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args)
}

func Print(i ...interface{}) {
	log.Print(i)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args)
}
