package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	WithFields(fields ...zap.Field) Logger
}
type ZapLogger struct {
	*zap.Logger
}

func NewLogger() Logger {
	file, err := os.OpenFile("./logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(file),
		zapcore.InfoLevel,
	)

	stdout := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	core := zapcore.NewTee(fileCore, stdout)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	if err != nil {
		return nil
	}
	return &ZapLogger{Logger: zapLogger}
}

func (l *ZapLogger) WithFields(fields ...zap.Field) Logger {
	return &ZapLogger{Logger: l.Logger.With(fields...)}
}
