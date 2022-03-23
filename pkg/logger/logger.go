package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	With(fields ...zapcore.Field) Logger
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Debug(msg string, fields ...zapcore.Field)
	Sync() error
}

type ZapLoggerAdapter struct {
	*zap.Logger
}

func NewLogger(name string, dev bool) (*ZapLoggerAdapter, error) {
	loggerConfig := zap.NewProductionConfig()
	if dev {
		loggerConfig = zap.NewDevelopmentConfig()
	}
	l, err := loggerConfig.Build(
		zap.Fields(defaultProcessFields(name)...),
		zap.AddCaller(),
	)
	if err != nil {
		return nil, err
	}
	return &ZapLoggerAdapter{Logger: l}, nil
}

func (l ZapLoggerAdapter) With(fields ...zapcore.Field) Logger {
	return &ZapLoggerAdapter{Logger: l.Logger.With(fields...)}
}

func defaultProcessFields(name string) []zap.Field {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	return []zap.Field{
		zap.Int("pid", os.Getpid()),
		zap.String("host", host),
		zap.String("app", name),
	}
}
