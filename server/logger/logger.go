package logger

import (
	"chatting/config"
	"github.com/labstack/gommon/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var log Logger

const (
	logLevelDebug = "debug"
	logLevelInfo  = "info"
	logLevelWarn  = "warn"
	logLevelError = "error"
)

type Logger interface {
	Debug(i ...interface{})
	Debugf(format string, args ...interface{})
	Error(i ...interface{})
	Errorf(format string, args ...interface{})
	Info(i ...interface{})
	Infof(format string, args ...interface{})
}

func init() {
	zapConf := zap.NewProductionConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	zapConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logConfig := config.Config().GetStringMapString("logger")
	switch logConfig["level"] {
	case logLevelDebug:
		zapConf.Level.SetLevel(zapcore.DebugLevel)
		break
	case logLevelInfo:
		zapConf.Level.SetLevel(zapcore.InfoLevel)
		break
	case logLevelWarn:
		zapConf.Level.SetLevel(zapcore.WarnLevel)
		break
	case logLevelError:
		zapConf.Level.SetLevel(zapcore.ErrorLevel)
		break
	}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConf.EncoderConfig),
		zapcore.AddSync(color.New().Output()),
		zapConf.Level,
	), zap.WithCaller(true))

	log = logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
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
