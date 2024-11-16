package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once     sync.Once
	instance *Logger
)

type Logger struct {
	zapLogger *zap.Logger
}

var LoggerInstance *Logger

func init() {
	once.Do(func() {
		logLevel := getLogLevelFromEnv()

		config := zap.Config{

			Level:            zap.NewAtomicLevelAt(logLevel),
			Development:      true,
			Encoding:         "json",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO, ERROR, etc.
				EncodeTime:     zapcore.ISO8601TimeEncoder,  // Human-readable time format
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}

		if os.Getenv("ENV") == "production" {
			config.Development = false
		}

		zapLogger, err := config.Build()
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}

		instance = &Logger{zapLogger: zapLogger}
	})

	LoggerInstance = instance

}

func getLogLevelFromEnv() zapcore.Level {
	fmt.Println("Getting log level from env")
	fmt.Println("LOG_LEVEL: ", os.Getenv("LOG_LEVEL"))
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "DEBUG"
	}
	switch logLevel {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.zapLogger.Panic(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.zapLogger.Sync()
}

// wrapper func's
func Info(msg string, fields ...zap.Field) {
	LoggerInstance.Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	LoggerInstance.Warn(msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	LoggerInstance.Debug(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	LoggerInstance.Error(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	LoggerInstance.Fatal(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	LoggerInstance.Panic(msg, fields...)
}

// sunc is not needed as of now. can be used in main.go
func Sync() error {
	return LoggerInstance.Sync()
}

func SetLogLevel(logLevel string) {
	os.Setenv("LOG_LEVEL", logLevel)
}
