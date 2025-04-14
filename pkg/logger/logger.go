package logger

import (
	"net/http"
	"os"
	"time"

	"go_byteeats/pkg/middleware"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// Config holds the logger configuration
type Config struct {
	Environment string
	LogLevel    string
	Pretty      bool
}

// Initialize sets up the logger with the given configuration
func Initialize(cfg Config) {
	// Set default timezone
	zerolog.TimeFieldFormat = time.RFC3339

	// Set global log level
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output writer
	var output zerolog.ConsoleWriter
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	// Create logger instance
	Logger = zerolog.New(output).
		With().
		Timestamp().
		Str("env", cfg.Environment).
		Logger()

	log.Logger = Logger
}

// Fields type for structured logging
type Fields map[string]interface{}

// Debug logs a debug message with optional fields
func Debug(msg string, fields ...Fields) {
	if len(fields) > 0 {
		Logger.Debug().Fields(fields[0]).Msg(msg)
	} else {
		Logger.Debug().Msg(msg)
	}
}

// Info logs an info message with optional fields
func Info(msg string, fields ...Fields) {
	if len(fields) > 0 {
		Logger.Info().Fields(fields[0]).Msg(msg)
	} else {
		Logger.Info().Msg(msg)
	}
}

// Warn logs a warning message with optional fields
func Warn(msg string, fields ...Fields) {
	if len(fields) > 0 {
		Logger.Warn().Fields(fields[0]).Msg(msg)
	} else {
		Logger.Warn().Msg(msg)
	}
}

// Error logs an error message with optional fields
func Error(err error, msg string, fields ...Fields) {
	if len(fields) > 0 {
		Logger.Error().Err(err).Fields(fields[0]).Msg(msg)
	} else {
		Logger.Error().Err(err).Msg(msg)
	}
}

// Fatal logs a fatal message with optional fields and exits
func Fatal(err error, msg string, fields ...Fields) {
	if len(fields) > 0 {
		Logger.Fatal().Err(err).Fields(fields[0]).Msg(msg)
	} else {
		Logger.Fatal().Err(err).Msg(msg)
	}
}

// HTTPMiddlewareLogger returns middleware for HTTP request logging
func HTTPMiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Process request
		next.ServeHTTP(ww, r)

		// Log request details
		Logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Int("status", ww.Status()).
			Str("latency", time.Since(start).String()).
			Str("user_agent", r.UserAgent()).
			Msg("HTTP Request")
	})
}
