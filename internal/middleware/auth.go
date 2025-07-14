package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"waterlogger/internal/models"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for setup wizard and static files
		if strings.HasPrefix(c.Request.URL.Path, "/setup") || 
		   strings.HasPrefix(c.Request.URL.Path, "/static") ||
		   strings.HasPrefix(c.Request.URL.Path, "/login") ||
		   strings.HasPrefix(c.Request.URL.Path, "/api/setup") {
			c.Next()
			return
		}

		// Check for session or token
		userID, exists := c.Get("user_id")
		if !exists {
			// Check session cookie
			sessionCookie, err := c.Cookie("session")
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, "/login")
				c.Abort()
				return
			}

			// Validate session (simplified - in production use proper session store)
			if userIDStr, err := validateSession(sessionCookie); err == nil {
				if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
					userID = uint(uid)
					c.Set("user_id", userID)
				}
			}
		}

		if userID == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		// Add user_id to context for GORM hooks
		ctx := context.WithValue(c.Request.Context(), "user_id", userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePassword(password string) []string {
	var errors []string
	
	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long")
	}
	
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}
	if !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}
	if !hasDigit {
		errors = append(errors, "Password must contain at least one digit")
	}
	if !hasSpecial {
		errors = append(errors, "Password must contain at least one special character")
	}
	
	return errors
}

// Simplified session validation (in production, use proper session store)
func validateSession(sessionID string) (string, error) {
	// This is a placeholder - implement proper session validation
	// For now, just return the session as user ID
	return sessionID, nil
}

func CreateSession(userID uint) string {
	// This is a placeholder - implement proper session creation
	// For now, just return the user ID as string
	return strconv.FormatUint(uint64(userID), 10)
}

func RequireSetup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if setup is required
		var count int64
		if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}

		if count == 0 {
			// No users exist, setup is required
			if !strings.HasPrefix(c.Request.URL.Path, "/setup") {
				c.Redirect(http.StatusTemporaryRedirect, "/setup")
				c.Abort()
				return
			}
		} else {
			// Users exist, setup is not allowed
			if strings.HasPrefix(c.Request.URL.Path, "/setup") {
				c.Redirect(http.StatusTemporaryRedirect, "/")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}