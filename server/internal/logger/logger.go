package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	LogFilePath string `mapstructure:"LOG_FILE_PATH" env-default:"./logs"`
}

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
}
type ZapLogger struct {
	*zap.Logger
}

func NewLogger(cfg LoggerConfig) (Logger, error) {
	_ = os.Mkdir(cfg.LogFilePath, 0755)
	file, err := os.OpenFile(cfg.LogFilePath+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
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

	return &ZapLogger{Logger: zapLogger}, nil
}

func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	z.Logger.Info(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...zap.Field) {
	z.Logger.Error(msg, fields...)
}

func (z *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{Logger: z.Logger.With(fields...)}
}
