package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitLogger(level string) {
	conf := zap.NewDevelopmentConfig()

	switch level {
	case "DEBUG":
		conf.Level.SetLevel(zapcore.DebugLevel)
	case "INFO":
		conf.Level.SetLevel(zapcore.InfoLevel)
	case "WARN":
	case "WARNING":
		conf.Level.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		conf.Level.SetLevel(zapcore.ErrorLevel)
	case "DPANIC":
		conf.Level.SetLevel(zapcore.DPanicLevel)
	case "PANIC":
		conf.Level.SetLevel(zapcore.PanicLevel)
	case "FATAL":
		conf.Level.SetLevel(zapcore.FatalLevel)
	default:
		conf.Level.SetLevel(zapcore.InfoLevel)
	}

	logger, err := conf.Build()
	if err != nil {
		zap.S().Panicf("error instantiating logger: %s", err)
	}
	Logger = logger.Sugar()
}
