package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// The global logger instance
var Logger zerolog.Logger

func init() {
	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}

// Info logs an info message
func Info() *zerolog.Event {
	return Logger.Info()
}

// Fatal logs a fatal message and exits the program
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}
