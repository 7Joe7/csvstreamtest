package logger

import (
	"io"
	"os"
	"strings"
	"sync/atomic"

	"github.com/rs/zerolog"
)

// LogMode defines the logging mode.
type LogMode uint8

// LodMode enum to define the different logging modes
const (
	Pretty LogMode = iota
	JSON
)

var gLevel = new(uint32)

// NewZeroLog creates a new zerolog logger
func NewZeroLog(writer io.Writer, mode LogMode) *zerolog.Logger {
	var zl zerolog.Logger

	// allow to override log mode with env
	envMode := os.Getenv("LOG_MODE")
	if envMode != "" {
		mode = parseLogMode(envMode)
	}
	switch mode {
	case JSON:
		zl = zerolog.New(writer)
	case Pretty:
		zl = zerolog.New(writer).Output(zerolog.ConsoleWriter{Out: writer}).With().Timestamp().Logger()
	}

	zl = zl.Hook(FileNameHook{
		pretty: mode == Pretty,
	})

	return &zl
}

func parseLogMode(logMode string) LogMode {
	switch strings.ToUpper(logMode) {
	case "PRETTY":
		return Pretty
	case "JSON":
		return JSON
	default:
		return Pretty
	}
}

// SetGlobalLevel sets the global level for all loggers
func SetGlobalLevel(lvl string) {
	zlLevel := parseLevel(lvl)
	atomic.StoreUint32(gLevel, uint32(zlLevel))
	zerolog.SetGlobalLevel(zlLevel)
}

// parseLevel parses a level from string to log level
func parseLevel(level string) zerolog.Level {
	switch strings.ToUpper(level) {
	case "FATAL":
		return zerolog.FatalLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "WARNING":
		return zerolog.WarnLevel
	case "INFO":
		return zerolog.InfoLevel
	case "DEBUG":
		return zerolog.DebugLevel
	default:
		return zerolog.DebugLevel
	}
}
