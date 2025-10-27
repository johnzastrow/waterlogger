package logging

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger is a custom GORM logger that uses zerolog
type GormLogger struct {
	logger                *Logger
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
	LogLevel              gormlogger.LogLevel
}

// NewGormLogger creates a new GORM logger
func NewGormLogger(logger *Logger, slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		logger:                logger,
		SlowThreshold:         slowThreshold,
		SkipErrRecordNotFound: true,
		LogLevel:              gormlogger.Info,
	}
}

// SetLogLevel sets the log level for GORM logger
func (l *GormLogger) SetLogLevel(level string) {
	switch level {
	case "debug":
		l.LogLevel = gormlogger.Info
	case "info":
		l.LogLevel = gormlogger.Warn
	case "warn", "error", "fatal":
		l.LogLevel = gormlogger.Error
	default:
		l.LogLevel = gormlogger.Info
	}
}

// LogMode sets the log mode
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info messages
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.logger.Info().
			Str("component", "gorm").
			Msgf(msg, data...)
	}
}

// Warn logs warning messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger.Warn().
			Str("component", "gorm").
			Msgf(msg, data...)
	}
}

// Error logs error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger.Error().
			Str("component", "gorm").
			Msgf(msg, data...)
	}
}

// Trace logs SQL queries with execution time
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// Build log event
	var event *zerolog.Event

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.SkipErrRecordNotFound):
		// Log errors (except ErrRecordNotFound if configured to skip)
		event = l.logger.Error().
			Err(err).
			Str("component", "gorm").
			Str("sql", sanitizeQuery(sql)).
			Int64("rows", rows).
			Dur("elapsed", elapsed)

		// Add specific error context
		if errors.Is(err, gorm.ErrRecordNotFound) {
			event.Msg("Record not found")
		} else if errors.Is(err, gorm.ErrInvalidTransaction) {
			event.Msg("Invalid transaction")
		} else if errors.Is(err, gorm.ErrNotImplemented) {
			event.Msg("Not implemented")
		} else if errors.Is(err, gorm.ErrMissingWhereClause) {
			event.Msg("Missing WHERE clause")
		} else if errors.Is(err, gorm.ErrUnsupportedRelation) {
			event.Msg("Unsupported relation")
		} else if errors.Is(err, gorm.ErrPrimaryKeyRequired) {
			event.Msg("Primary key required")
		} else if errors.Is(err, gorm.ErrModelValueRequired) {
			event.Msg("Model value required")
		} else if errors.Is(err, gorm.ErrInvalidData) {
			event.Msg("Invalid data")
		} else if errors.Is(err, gorm.ErrInvalidField) {
			event.Msg("Invalid field")
		} else if errors.Is(err, gorm.ErrEmptySlice) {
			event.Msg("Empty slice")
		} else if errors.Is(err, gorm.ErrDryRunModeUnsupported) {
			event.Msg("Dry run mode unsupported")
		} else if isDatabaseConstraintError(err) {
			// Specific handling for database constraint violations
			event.Msg("Database constraint violation")
		} else if isDuplicateKeyError(err) {
			// Specific handling for duplicate key errors
			event.Msg("Duplicate key error")
		} else if isForeignKeyError(err) {
			// Specific handling for foreign key constraint errors
			event.Msg("Foreign key constraint violation")
		} else {
			event.Msg("Database query error")
		}

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		// Log slow queries
		event = l.logger.Warn().
			Str("component", "gorm").
			Str("sql", sanitizeQuery(sql)).
			Int64("rows", rows).
			Dur("elapsed", elapsed).
			Dur("threshold", l.SlowThreshold)
		event.Msg("Slow SQL query detected")

	case l.LogLevel == gormlogger.Info:
		// Log all queries at debug level
		event = l.logger.Debug().
			Str("component", "gorm").
			Str("sql", sanitizeQuery(sql)).
			Int64("rows", rows).
			Dur("elapsed", elapsed)
		event.Msg("SQL query executed")
	}
}

// isDatabaseConstraintError checks if error is a constraint violation
func isDatabaseConstraintError(err error) bool {
	errStr := err.Error()
	// Check for common constraint error patterns
	constraintPatterns := []string{
		"constraint",
		"violates",
		"CHECK constraint",
		"UNIQUE constraint",
		"NOT NULL constraint",
	}

	for _, pattern := range constraintPatterns {
		if contains(errStr, pattern) {
			return true
		}
	}
	return false
}

// isDuplicateKeyError checks if error is a duplicate key error
func isDuplicateKeyError(err error) bool {
	errStr := err.Error()
	duplicatePatterns := []string{
		"duplicate key",
		"UNIQUE constraint failed",
		"Duplicate entry",
		"already exists",
	}

	for _, pattern := range duplicatePatterns {
		if contains(errStr, pattern) {
			return true
		}
	}
	return false
}

// isForeignKeyError checks if error is a foreign key constraint error
func isForeignKeyError(err error) bool {
	errStr := err.Error()
	fkPatterns := []string{
		"foreign key",
		"FOREIGN KEY constraint",
		"violates foreign key constraint",
		"Cannot add or update a child row",
		"Cannot delete or update a parent row",
	}

	for _, pattern := range fkPatterns {
		if contains(errStr, pattern) {
			return true
		}
	}
	return false
}

// contains checks if string contains substring (case-insensitive)
func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
		 len(str) > len(substr) &&
		 (hasSubstring(str, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if equalFold(s[i:i+len(substr)], substr) {
			return true
		}
	}
	return false
}

func equalFold(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		c1 := s1[i]
		c2 := s2[i]
		if c1 >= 'A' && c1 <= 'Z' {
			c1 += 'a' - 'A'
		}
		if c2 >= 'A' && c2 <= 'Z' {
			c2 += 'a' - 'A'
		}
		if c1 != c2 {
			return false
		}
	}
	return true
}

// ParseGormError extracts meaningful information from GORM errors
func ParseGormError(err error) (string, map[string]interface{}) {
	if err == nil {
		return "", nil
	}

	details := make(map[string]interface{})
	var message string

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		message = "Record not found"
		details["error_type"] = "not_found"

	case errors.Is(err, gorm.ErrInvalidTransaction):
		message = "Invalid transaction"
		details["error_type"] = "invalid_transaction"

	case errors.Is(err, gorm.ErrNotImplemented):
		message = "Operation not implemented"
		details["error_type"] = "not_implemented"

	case errors.Is(err, gorm.ErrMissingWhereClause):
		message = "Missing WHERE clause in query"
		details["error_type"] = "missing_where"

	case errors.Is(err, gorm.ErrUnsupportedRelation):
		message = "Unsupported relation"
		details["error_type"] = "unsupported_relation"

	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		message = "Primary key required"
		details["error_type"] = "primary_key_required"

	case isDuplicateKeyError(err):
		message = "Duplicate entry - record already exists"
		details["error_type"] = "duplicate_key"

	case isForeignKeyError(err):
		message = "Foreign key constraint violation"
		details["error_type"] = "foreign_key_violation"

	case isDatabaseConstraintError(err):
		message = "Database constraint violation"
		details["error_type"] = "constraint_violation"

	default:
		message = fmt.Sprintf("Database error: %v", err)
		details["error_type"] = "unknown"
	}

	details["raw_error"] = err.Error()
	return message, details
}
