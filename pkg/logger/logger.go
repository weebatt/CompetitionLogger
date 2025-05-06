package logger

import (
	"context"
	"go.uber.org/zap"
)

const (
	key = "logger"
)

type Logger struct {
	logger *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, key, &Logger{logger: logger}), nil
}

func GetFromContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value(key).(*Logger)
	if !ok {
		return &Logger{zap.NewNop()}
	}
	return logger
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}
