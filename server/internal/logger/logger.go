package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
}
type ZapLogger struct {
	*zap.Logger
}

func NewLogger() (Logger, error) {
	path := `./logs`                                                                     // ./var/log/app | ./logs
	_ = os.Mkdir(path, 0755)                                                             //TODO: вынести директорию в var/..
	file, err := os.OpenFile(path+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) //TODO: path вынести в real-time cfg
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
