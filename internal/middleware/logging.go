package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"waterlogger/internal/logging"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in header
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Set request ID in context and response header
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggingMiddleware logs HTTP requests with structured logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Get request ID from context
		requestID, _ := c.Get("request_id")
		reqID := ""
		if requestID != nil {
			reqID = requestID.(string)
		}

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get user ID if authenticated
		var userID uint
		if userIDVal, exists := c.Get("user_id"); exists {
			userID = userIDVal.(uint)
		}

		// Build full path
		if raw != "" {
			path = path + "?" + raw
		}

		// Log the request
		logger := logging.GetLogger().WithRequestID(reqID)
		if userID > 0 {
			logger = logger.WithUserID(userID)
		}

		event := logger.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("ip", clientIP).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent())

		// Add error if status is 4xx or 5xx
		if statusCode >= 400 {
			if len(c.Errors) > 0 {
				event.Str("error", c.Errors.String())
			}
		}

		// Different message based on status code
		if statusCode >= 500 {
			event.Msg("HTTP request - server error")
		} else if statusCode >= 400 {
			event.Msg("HTTP request - client error")
		} else {
			event.Msg("HTTP request")
		}
	}
}

// AuditLoggingMiddleware logs security-sensitive operations
func AuditLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only audit specific endpoints
		if shouldAudit(c.Request.Method, c.Request.URL.Path) {
			start := time.Now()

			// Get user ID if authenticated
			var userID uint
			if userIDVal, exists := c.Get("user_id"); exists {
				userID = userIDVal.(uint)
			}

			// Get request ID
			requestID, _ := c.Get("request_id")
			reqID := ""
			if requestID != nil {
				reqID = requestID.(string)
			}

			// Process request
			c.Next()

			// Log audit event after request is processed
			details := map[string]interface{}{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"ip":         c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
				"status":     c.Writer.Status(),
				"latency_ms": time.Since(start).Milliseconds(),
				"request_id": reqID,
			}

			// Add error details if present
			if len(c.Errors) > 0 {
				details["error"] = c.Errors.String()
			}

			// Determine audit event type
			event := getAuditEventType(c.Request.Method, c.Request.URL.Path)

			logging.AuditLog(event, userID, details)
		} else {
			c.Next()
		}
	}
}

// shouldAudit determines if a request should be audited
func shouldAudit(method, path string) bool {
	// Audit authentication operations
	if path == "/login" || path == "/logout" || path == "/api/auth/login" {
		return true
	}

	// Audit setup wizard
	if path == "/api/setup" {
		return true
	}

	// Audit user management operations
	if method != "GET" && (path == "/api/users" || contains(path, "/api/users/")) {
		return true
	}

	// Audit data modification operations
	if method == "POST" || method == "PUT" || method == "DELETE" {
		return true
	}

	// Audit export operations
	if contains(path, "/export") {
		return true
	}

	return false
}

// getAuditEventType returns a descriptive audit event type
func getAuditEventType(method, path string) string {
	switch {
	case path == "/login" || path == "/api/auth/login":
		return "user_login"
	case path == "/logout":
		return "user_logout"
	case path == "/api/setup":
		return "setup_wizard"
	case contains(path, "/api/users/") && method == "PUT":
		return "user_update"
	case contains(path, "/api/users/") && method == "DELETE":
		return "user_delete"
	case path == "/api/users" && method == "POST":
		return "user_create"
	case contains(path, "/api/pools/") && method == "DELETE":
		return "pool_delete"
	case contains(path, "/api/pools/") && method == "PUT":
		return "pool_update"
	case path == "/api/pools" && method == "POST":
		return "pool_create"
	case contains(path, "/api/samples/") && method == "DELETE":
		return "sample_delete"
	case path == "/api/samples" && method == "POST":
		return "sample_create"
	case contains(path, "/export"):
		return "data_export"
	default:
		return "data_modification"
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return hex.EncodeToString([]byte(time.Now().String()))[:16]
	}
	return hex.EncodeToString(b)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
