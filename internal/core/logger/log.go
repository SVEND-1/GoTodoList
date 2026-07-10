package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerContextKey struct {
}

var (
	key = loggerContextKey{}
)

type Log struct {
	*zap.Logger
	file *os.File
}

func ToContext(ctx context.Context, log *Log) context.Context {
	return context.WithValue(ctx, key, log)
}

func FromContext(ctx context.Context) *Log {
	log, ok := ctx.Value(key).(*Log)
	if !ok {
		panic("logger not found in context")
	}
	return log
}

func NewLog(config LogConfig) (*Log, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("error parsing log level: %v", err)
	}
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("error creating log directory: %v", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(config.Folder, fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)
	logger := zap.New(core, zap.AddCaller())

	return &Log{logger, logFile}, nil
}

func (l *Log) With(fields ...zap.Field) *Log {
	return &Log{l.Logger.With(fields...), l.file}
}

func (log *Log) Close() {
	if err := log.file.Close(); err != nil {
		fmt.Printf("error closing log file: %v", err)
	}
}
