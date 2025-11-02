# Waterlogger Logging Implementation

**Date:** 2025-11-01
**Version:** 1.3.1
**Status:** ✅ Complete

---

## Overview

Waterlogger now implements comprehensive structured logging using **zerolog** with log rotation, multiple output destinations, configurable log levels, request tracing, and security audit logging.

## Features Implemented

### ✅ **1. Structured Logging with Zerolog**
- **Library:** `github.com/rs/zerolog`
- **Benefits:**
  - Zero-allocation JSON logger
  - Fast performance
  - Structured key-value logging
  - Pretty console output for development
  - JSON output for production

### ✅ **2. Log Rotation with Lumberjack**
- **Library:** `gopkg.in/natefinch/lumberjack.v2`
- **Configuration:**
  - Max file size: 100MB (configurable)
  - Max backups: 3 files (configurable)
  - Max age: 28 days (configurable)
  - Compression: Enabled (configurable)

### ✅ **3. Log Levels**
Supports all standard log levels:
- **DEBUG** - Detailed diagnostic information
- **INFO** - General informational messages
- **WARN** - Warning messages for potentially harmful situations
- **ERROR** - Error events that might still allow the application to continue
- **FATAL** - Severe errors that cause application termination

### ✅ **4. Multiple Output Destinations**
- **stdout** - Console output only
- **file** - Log file output only
- **both** - Console and file output simultaneously (recommended)

### ✅ **5. Output Formats**
- **console** - Pretty formatted output with colors (development)
- **json** - Structured JSON output (production)

### ✅ **6. Request ID Middleware**
- Generates unique request ID for each HTTP request
- Propagates request ID through entire request lifecycle
- Enables request tracing across multiple log entries
- Returns request ID in response header: `X-Request-ID`

### ✅ **7. Security Audit Logging**
Automatically logs security-sensitive operations:
- User login/logout
- User creation/modification/deletion
- Setup wizard completion
- Data exports
- All POST/PUT/DELETE operations

Audit logs include:
- User ID (if authenticated)
- Request method and path
- Client IP address
- User agent
- Response status code
- Request latency
- Request ID for tracing

### ✅ **8. Custom GORM Logger**
Integrates structured logging with GORM database operations:

**Features:**
- Logs all SQL queries (at DEBUG level)
- Tracks query execution time
- Detects slow queries (threshold: 200ms)
- **Comprehensive database error tracking:**
  - Duplicate key errors
  - Foreign key constraint violations
  - NOT NULL constraint violations
  - UNIQUE constraint violations
  - CHECK constraint violations
  - Record not found errors
  - Invalid transaction errors
  - Missing WHERE clause errors
- Sanitizes sensitive data in query logs
- Truncates long queries for readability

**Error Detection:**
```go
// Automatically logs and categorizes database errors
- Duplicate entry errors
- Foreign key violations
- Constraint violations
- Record not found
- Invalid data
- Primary key required
- And more...
```

### ✅ **9. HTTP Request Logging**
Every HTTP request is logged with:
- Request method (GET, POST, PUT, DELETE)
- Request path
- Client IP address
- Response status code
- Request latency
- User agent
- Request ID
- User ID (if authenticated)

Color-coded by status:
- **2xx** - Success (INFO level)
- **4xx** - Client error (INFO level)
- **5xx** - Server error (INFO level with error context)

### ✅ **10. Production-Ready Logger Configuration (v1.3.1+)**
- **Removed caller information logging** for production safety and portability
  - Caller information (`.Caller()`) was embedding development file paths in logs
  - Binary can now be deployed from any location (e.g., `/opt`) without exposing source code structure
  - Request tracing via Request IDs remains fully functional for troubleshooting
  - Improves security by not disclosing source code paths in logs

---

## Configuration

### Location
`config.yaml` - Logging section

### Example Configuration

```yaml
logging:
    level: info           # debug, info, warn, error, fatal
    format: console       # json, console
    output: both          # stdout, file, both
    file_path: logs/waterlogger.log
    max_size: 100         # MB
    max_backups: 3        # number of old log files to keep
    max_age: 28           # days
    compress: true        # compress old log files
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `level` | string | `info` | Minimum log level to output |
| `format` | string | `console` | Output format (console or json) |
| `output` | string | `both` | Where to send logs |
| `file_path` | string | `logs/waterlogger.log` | Path to log file |
| `max_size` | int | 100 | Max size in MB before rotation |
| `max_backups` | int | 3 | Max number of old log files |
| `max_age` | int | 28 | Max age of log files in days |
| `compress` | bool | true | Compress rotated log files |

### Environment-Specific Configurations

**Development:**
```yaml
logging:
    level: debug
    format: console
    output: both
```

**Production:**
```yaml
logging:
    level: info
    format: json
    output: file
    file_path: /var/log/waterlogger/waterlogger.log
```

---

## Usage Examples

### Basic Logging

```go
import "waterlogger/internal/logging"

// Simple messages
logging.Info().Msg("Server starting")
logging.Error().Msg("Failed to connect")
logging.Warn().Msg("Using default configuration")
logging.Debug().Msg("Request received")

// With context
logging.Info().
    Str("host", "localhost").
    Int("port", 2342).
    Msg("Starting server")

// With error
logging.Error().
    Err(err).
    Str("operation", "database_connect").
    Msg("Database connection failed")
```

### Request Context Logging

```go
// In HTTP handlers
logger := logging.WithRequestID(c.GetString("request_id"))
logger.Info().Str("username", username).Msg("User logged in")
```

### User Context Logging

```go
// With user ID
logger := logging.WithUserID(userID)
logger.Info().Msg("User preferences updated")
```

### Component-Specific Logging

```go
// For specific components
logger := logging.WithComponent("database")
logger.Info().Msg("Running migrations")
```

### Database Error Logging

```go
if err := db.Create(&user).Error; err != nil {
    logger.Error().
        Err(err).
        Str("username", user.Username).
        Msg("Failed to create user")
    // GORM logger will automatically log detailed database error
}
```

### Audit Logging

```go
logging.AuditLog("user_login", userID, map[string]interface{}{
    "username": username,
    "ip": clientIP,
    "success": true,
})
```

### Database Operation Logging

```go
// Automatically handled by custom GORM logger
db.Where("id = ?", id).First(&user)
// Logs: SQL query, execution time, result
```

---

## Log File Structure

### Directory Layout
```
waterlogger/
├── logs/
│   ├── waterlogger.log           # Current log file
│   ├── waterlogger-2025-10-25.log.gz  # Rotated log (compressed)
│   ├── waterlogger-2025-10-24.log.gz  # Older rotated log
│   └── waterlogger-2025-10-23.log.gz  # Oldest rotated log
```

### Log Rotation Behavior
1. When `waterlogger.log` reaches 100MB, it's renamed to `waterlogger-YYYY-MM-DD.log`
2. A new `waterlogger.log` file is created
3. Old log files are compressed to `.gz` format
4. Only 3 backup files are kept
5. Files older than 28 days are automatically deleted

---

## Log Output Examples

### Console Format (Development)

```
2025-10-26T14:23:45-05:00 INF Waterlogger starting...
2025-10-26T14:23:45-05:00 INF Build information build_date=2025-10-26 build_time=14:23:45 version=1.0.0
2025-10-26T14:23:45-05:00 INF Connecting to SQLite database path=waterlogger.db
2025-10-26T14:23:45-05:00 INF Successfully connected to SQLite database
2025-10-26T14:23:45-05:00 INF Running database migrations
2025-10-26T14:23:45-05:00 INF Database migrations completed successfully
2025-10-26T14:23:45-05:00 INF Starting Waterlogger server address=0.0.0.0:2342 host=0.0.0.0 port=2342
2025-10-26T14:23:45-05:00 INF Open your browser to access Waterlogger url=http://0.0.0.0:2342
```

### JSON Format (Production)

```json
{"level":"info","time":"2025-10-26T14:23:45-05:00","caller":"main.go:111","message":"Waterlogger starting..."}
{"level":"info","version":"1.0.0","build_time":"14:23:45","build_date":"2025-10-26","time":"2025-10-26T14:23:45-05:00","caller":"main.go:112","message":"Build information"}
{"level":"info","path":"waterlogger.db","time":"2025-10-26T14:23:45-05:00","caller":"database.go:35","message":"Connecting to SQLite database"}
{"level":"info","time":"2025-10-26T14:23:45-05:00","caller":"database.go:41","message":"Successfully connected to SQLite database"}
{"level":"info","address":"0.0.0.0:2342","host":"0.0.0.0","port":2342,"time":"2025-10-26T14:23:45-05:00","caller":"main.go:216","message":"Starting Waterlogger server"}
```

### HTTP Request Log Example

```
2025-10-26T14:24:01-05:00 INF HTTP request method=POST path=/api/samples ip=127.0.0.1 status=201 latency=24ms request_id=a1b2c3d4e5f6g7h8 user_id=1 user_agent=Mozilla/5.0
```

### Database Error Log Example

```
2025-10-26T14:24:15-05:00 ERR Database operation failed component=gorm error="UNIQUE constraint failed: users.username" sql="INSERT INTO users (username,email,password) VALUES (?,?,?)" rows=0 elapsed=2ms message="Duplicate key error"
```

### Audit Log Example

```
2025-10-26T14:25:30-05:00 INF Security audit event audit_event=user_login user_id=1 timestamp=2025-10-26T14:25:30-05:00 username=admin ip=127.0.0.1 success=true message="Security audit event"
```

---

## Middleware Stack

### Order of Middleware (Important!)

1. **Recovery Middleware** - Gin's built-in panic recovery
2. **Request ID Middleware** - Generates unique request ID
3. **Logging Middleware** - Logs HTTP requests
4. **Audit Logging Middleware** - Logs security events
5. **CORS Middleware** - Handles cross-origin requests
6. **Require Setup Middleware** - Checks if setup is complete
7. **Auth Middleware** - Handles authentication

This order ensures:
- Every request gets a request ID
- All requests are logged
- Security events are captured
- Authentication happens after logging is set up

---

## Performance Considerations

### Zerolog Performance
- **Zero allocation** for most log operations
- **Extremely fast** - designed for high-performance applications
- **Minimal overhead** - typically <100ns per log operation

### Log Levels in Production
- Set to `info` or `warn` in production to reduce log volume
- Use `debug` only for troubleshooting
- `error` and `fatal` are always logged regardless of level

### File I/O
- Log rotation happens asynchronously
- Compression is done in background
- Minimal impact on application performance

---

## Troubleshooting

### Logs Not Appearing

**Problem:** No logs being written
**Solution:**
1. Check `config.yaml` logging section
2. Verify log directory exists and is writable: `mkdir -p logs`
3. Check file permissions: `chmod 755 logs`
4. Verify log level allows the message level

### Log Files Not Rotating

**Problem:** Single large log file
**Solution:**
1. Check `max_size` configuration
2. Verify `lumberjack` is imported correctly
3. Check disk space availability

### Missing Request IDs

**Problem:** Request ID not in logs
**Solution:**
1. Verify `RequestIDMiddleware` is registered
2. Ensure middleware order is correct (RequestID must be first)
3. Check that logger is using `WithRequestID()`

### Database Errors Not Logged

**Problem:** Database errors not appearing in logs
**Solution:**
1. Verify custom GORM logger is configured in `database.NewDB()`
2. Check log level (should be `info` or lower)
3. Ensure GORM logger level is set correctly

### Slow Query Detection Not Working

**Problem:** Slow queries not being logged
**Solution:**
1. Check slow query threshold (default: 200ms)
2. Verify log level is `warn` or lower
3. Confirm GORM logger is properly initialized

---

## Migration from Old Logging

### Changes Made

1. **Removed:** Standard library `log` package
2. **Removed:** All `fmt.Printf` debug statements
3. **Added:** Structured logging with context
4. **Added:** Request tracing
5. **Added:** Security audit trail
6. **Added:** Database error tracking

### Breaking Changes

**None** - All changes are internal implementation. External API remains the same.

---

## Dependencies Added

```go
require (
    github.com/rs/zerolog v1.33.0
    gopkg.in/natefinch/lumberjack.v2 v2.2.1
)
```

To install:
```bash
go mod tidy
```

---

## Best Practices

### 1. Always Use Structured Logging
✅ **Good:**
```go
logging.Info().
    Str("username", username).
    Int("user_id", userID).
    Msg("User logged in")
```

❌ **Bad:**
```go
log.Printf("User %s (ID: %d) logged in", username, userID)
```

### 2. Include Context in Logs
✅ **Good:**
```go
logger := logging.WithRequestID(requestID).WithUserID(userID)
logger.Info().Msg("Processing request")
```

❌ **Bad:**
```go
logging.Info().Msg("Processing request")
```

### 3. Use Appropriate Log Levels
- `Debug` - Development and troubleshooting only
- `Info` - Normal operation events
- `Warn` - Concerning but not critical
- `Error` - Errors that need attention
- `Fatal` - Critical errors (application exits)

### 4. Don't Log Sensitive Data
❌ **Never log:**
- Passwords (even hashed)
- API keys
- Session tokens
- Credit card numbers
- Personal information (unless required for audit)

✅ **Do log:**
- Usernames
- User IDs
- Operation types
- Timestamps
- Error messages (sanitized)

### 5. Use Audit Logging for Security Events
```go
logging.AuditLog("password_change", userID, map[string]interface{}{
    "username": username,
    "ip": clientIP,
    "timestamp": time.Now(),
})
```

### 6. Handle Database Errors Properly
```go
if err := db.Create(&user).Error; err != nil {
    msg, details := logging.ParseGormError(err)
    logger.Error().
        Err(err).
        Str("error_type", details["error_type"].(string)).
        Msg(msg)
    // Handle error appropriately
}
```

---

## Testing

### Development Testing

```bash
# Run with debug logging
# Edit config.yaml: level: debug

./waterlogger

# Check logs directory
ls -lh logs/

# Tail logs in real-time
tail -f logs/waterlogger.log
```

### Production Testing

```bash
# Run with JSON output
# Edit config.yaml:
#   level: info
#   format: json
#   output: file

./waterlogger

# Parse JSON logs
cat logs/waterlogger.log | jq .

# Filter by level
cat logs/waterlogger.log | jq 'select(.level=="error")'

# Filter by component
cat logs/waterlogger.log | jq 'select(.component=="database")'
```

### Log Rotation Testing

```bash
# Force rotation (make logs/waterlogger.log > 100MB)
for i in {1..1000000}; do echo "Test log entry $i" >> logs/waterlogger.log; done

# Check rotation occurred
ls -lh logs/
```

---

## Monitoring and Alerting

### Integration with Log Aggregation Tools

**Elasticsearch/Logstash/Kibana (ELK):**
```bash
# Ship JSON logs to Logstash
filebeat -e -c filebeat.yml
```

**Grafana Loki:**
```yaml
# promtail.yml
clients:
  - url: http://loki:3100/loki/api/v1/push
```

**Datadog:**
```bash
# Configure Datadog agent
dd-agent log collect
```

### Alert Rules (Example)

**High Error Rate:**
```
Alert if error count > 100 in 5 minutes
Query: level:"error" | count > 100
```

**Slow Database Queries:**
```
Alert if slow query detected
Query: "Slow SQL query detected" | count > 10
```

**Failed Login Attempts:**
```
Alert if failed logins > 10 from same IP
Query: audit_event:"user_login" AND success:false | count_by_ip > 10
```

---

## Future Enhancements

Potential improvements for future versions:

1. **Log Sampling** - Sample high-volume logs to reduce storage
2. **Log Filtering** - Filter out health check requests
3. **Custom Log Formatters** - Support custom log formats
4. **Remote Logging** - Send logs to remote syslog servers
5. **Log Metrics** - Export log metrics to Prometheus
6. **Trace Integration** - Integrate with OpenTelemetry for distributed tracing
7. **Log Alerts** - Built-in alerting for critical log patterns

---

## Summary

Waterlogger now has **enterprise-grade logging** with:

✅ Structured logging (JSON/Console)
✅ Multiple log levels
✅ Log rotation and compression
✅ Request tracing (Request IDs)
✅ Security audit logging
✅ Database error tracking
✅ HTTP request logging
✅ Performance monitoring (slow queries)
✅ Configurable output destinations
✅ Zero-allocation performance

All logging is **production-ready** and follows **best practices** for observability and security.

---

## Support

For questions or issues:
- Check this documentation
- Review `config.yaml` settings
- Check log file permissions
- Verify disk space
- Test with `level: debug` for troubleshooting

---

**Implementation Date:** 2025-10-26
**Implemented By:** Claude Code AI Assistant
**Status:** ✅ Complete and Production Ready
