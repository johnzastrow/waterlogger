package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps zerolog with additional context
type Logger struct {
	zerolog.Logger
	config *Config
}

// Config holds logging configuration
type Config struct {
	Level      string `yaml:"level"`       // debug, info, warn, error, fatal
	Format     string `yaml:"format"`      // json, console
	Output     string `yaml:"output"`      // stdout, file, both
	FilePath   string `yaml:"file_path"`   // path to log file
	MaxSize    int    `yaml:"max_size"`    // max size in MB before rotation
	MaxBackups int    `yaml:"max_backups"` // max number of old log files
	MaxAge     int    `yaml:"max_age"`     // max age in days
	Compress   bool   `yaml:"compress"`    // compress rotated files
}

var (
	// Global logger instance
	globalLogger *Logger
)

// DefaultConfig returns default logging configuration
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		Format:     "console",
		Output:     "stdout",
		FilePath:   "logs/waterlogger.log",
		MaxSize:    100, // 100 MB
		MaxBackups: 3,
		MaxAge:     28, // 28 days
		Compress:   true,
	}
}

// Initialize sets up the global logger
func Initialize(cfg *Config) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Apply defaults for empty fields
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.Format == "" {
		cfg.Format = "console"
	}
	if cfg.Output == "" {
		cfg.Output = "stdout"
	}
	if cfg.FilePath == "" {
		cfg.FilePath = "logs/waterlogger.log"
	}

	// Parse log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("invalid log level %s: %w", cfg.Level, err)
	}
	zerolog.SetGlobalLevel(level)

	// Set up output writers
	var writers []io.Writer

	// Console output
	if cfg.Output == "stdout" || cfg.Output == "both" {
		var consoleWriter io.Writer
		if cfg.Format == "console" {
			// Pretty console output for development
			consoleWriter = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
				NoColor:    false,
			}
		} else {
			consoleWriter = os.Stdout
		}
		writers = append(writers, consoleWriter)
	}

	// File output
	if cfg.Output == "file" || cfg.Output == "both" {
		// Create log directory if it doesn't exist
		logDir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		// Set up log rotation
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		writers = append(writers, fileWriter)
	}

	// Create multi-writer if needed
	var output io.Writer
	if len(writers) == 0 {
		return fmt.Errorf("no output configured for logger")
	} else if len(writers) == 1 {
		output = writers[0]
	} else {
		output = io.MultiWriter(writers...)
	}

	// Create zerolog logger
	zlog := zerolog.New(output).With().
		Timestamp().
		Logger()

	globalLogger = &Logger{
		Logger: zlog,
		config: cfg,
	}

	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		// Initialize with defaults if not already initialized
		_ = Initialize(DefaultConfig())
	}
	return globalLogger
}

// WithRequestID adds a request ID to the logger context
func (l *Logger) WithRequestID(requestID string) *Logger {
	newLogger := l.Logger.With().Str("request_id", requestID).Logger()
	return &Logger{
		Logger: newLogger,
		config: l.config,
	}
}

// WithUserID adds a user ID to the logger context
func (l *Logger) WithUserID(userID uint) *Logger {
	newLogger := l.Logger.With().Uint("user_id", userID).Logger()
	return &Logger{
		Logger: newLogger,
		config: l.config,
	}
}

// WithComponent adds a component name to the logger context
func (l *Logger) WithComponent(component string) *Logger {
	newLogger := l.Logger.With().Str("component", component).Logger()
	return &Logger{
		Logger: newLogger,
		config: l.config,
	}
}

// AuditLog logs security-sensitive operations
func (l *Logger) AuditLog(event string, userID uint, details map[string]interface{}) {
	evt := l.Info().
		Str("audit_event", event).
		Uint("user_id", userID).
		Time("timestamp", time.Now())

	for k, v := range details {
		evt = evt.Interface(k, v)
	}

	evt.Msg("Security audit event")
}

// DatabaseError logs database errors with context
func (l *Logger) DatabaseError(operation string, err error, query string) {
	l.Error().
		Err(err).
		Str("operation", operation).
		Str("query", sanitizeQuery(query)).
		Msg("Database operation failed")
}

// DatabaseSlowQuery logs slow database queries
func (l *Logger) DatabaseSlowQuery(query string, duration time.Duration, threshold time.Duration) {
	if duration > threshold {
		l.Warn().
			Str("query", sanitizeQuery(query)).
			Dur("duration", duration).
			Dur("threshold", threshold).
			Msg("Slow database query detected")
	}
}

// HTTPRequest logs HTTP request details
func (l *Logger) HTTPRequest(method, path, ip string, statusCode int, duration time.Duration) {
	evt := l.Info().
		Str("method", method).
		Str("path", path).
		Str("ip", ip).
		Int("status", statusCode).
		Dur("duration", duration)

	if statusCode >= 500 {
		evt.Msg("HTTP request - server error")
	} else if statusCode >= 400 {
		evt.Msg("HTTP request - client error")
	} else {
		evt.Msg("HTTP request")
	}
}

// sanitizeQuery removes sensitive data from SQL queries
func sanitizeQuery(query string) string {
	// Replace potential passwords or sensitive data in queries
	sensitive := []string{"password", "secret", "token", "key"}
	sanitized := query

	for _, word := range sensitive {
		if strings.Contains(strings.ToLower(query), word) {
			// This is a simple sanitization - in production, use more sophisticated methods
			sanitized = strings.ReplaceAll(sanitized, "'", "***")
			break
		}
	}

	// Limit query length for logging
	if len(sanitized) > 500 {
		sanitized = sanitized[:500] + "... (truncated)"
	}

	return sanitized
}

// Convenience functions that use the global logger

// Debug logs a debug message
func Debug() *zerolog.Event {
	return GetLogger().Debug()
}

// Info logs an info message
func Info() *zerolog.Event {
	return GetLogger().Info()
}

// Warn logs a warning message
func Warn() *zerolog.Event {
	return GetLogger().Warn()
}

// Error logs an error message
func Error() *zerolog.Event {
	return GetLogger().Error()
}

// Fatal logs a fatal message and exits
func Fatal() *zerolog.Event {
	return GetLogger().Fatal()
}

// WithRequestID returns a logger with request ID context
func WithRequestID(requestID string) *Logger {
	return GetLogger().WithRequestID(requestID)
}

// WithUserID returns a logger with user ID context
func WithUserID(userID uint) *Logger {
	return GetLogger().WithUserID(userID)
}

// WithComponent returns a logger with component context
func WithComponent(component string) *Logger {
	return GetLogger().WithComponent(component)
}

// AuditLog logs a security audit event
func AuditLog(event string, userID uint, details map[string]interface{}) {
	GetLogger().AuditLog(event, userID, details)
}

// DatabaseError logs a database error
func DatabaseError(operation string, err error, query string) {
	GetLogger().DatabaseError(operation, err, query)
}

// DatabaseSlowQuery logs a slow query
func DatabaseSlowQuery(query string, duration time.Duration, threshold time.Duration) {
	GetLogger().DatabaseSlowQuery(query, duration, threshold)
}

// HTTPRequest logs an HTTP request
func HTTPRequest(method, path, ip string, statusCode int, duration time.Duration) {
	GetLogger().HTTPRequest(method, path, ip, statusCode, duration)
}
