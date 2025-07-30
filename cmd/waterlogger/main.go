package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"golang.org/x/term"
	"gorm.io/gorm"
	"waterlogger/internal/config"
	"waterlogger/internal/database"
	"waterlogger/internal/handlers"
	"waterlogger/internal/middleware"
	"waterlogger/internal/models"
)

// Build information - set at compile time
var (
	BuildTime = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Parse command line flags
	var configPath string
	var showVersion bool
	var showHelp bool
	var migrateToMariaDB bool
	var migrateToSQLite bool
	var exportData string
	var importData string
	var resetPassword string
	
	flag.StringVar(&configPath, "config", "config.yaml", "Path to configuration file")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.BoolVar(&migrateToMariaDB, "migrate-to-mariadb", false, "Migrate data from SQLite to MariaDB")
	flag.BoolVar(&migrateToSQLite, "migrate-to-sqlite", false, "Migrate data from MariaDB to SQLite")
	flag.StringVar(&exportData, "export", "", "Export database data to backup file")
	flag.StringVar(&importData, "import", "", "Import database data from backup file")
	flag.StringVar(&resetPassword, "reset-password", "", "Reset password for specified username")
	flag.Parse()

	if showVersion {
		fmt.Println("Waterlogger v1.0.0")
		os.Exit(0)
	}

	if showHelp {
		fmt.Println("Waterlogger - Pool and Hot Tub Water Management System")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  waterlogger [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -config string           Path to configuration file (default: config.yaml)")
		fmt.Println("  -version                 Show version information")
		fmt.Println("  -help                    Show this help message")
		fmt.Println("  -migrate-to-mariadb      Migrate data from SQLite to MariaDB")
		fmt.Println("  -migrate-to-sqlite       Migrate data from MariaDB to SQLite")
		fmt.Println("  -export string           Export database data to backup file")
		fmt.Println("  -import string           Import database data from backup file")
		fmt.Println("  -reset-password string   Reset password for specified username")
		fmt.Println()
		fmt.Println("For more information, visit: https://github.com/your-org/waterlogger")
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist, create default
			log.Printf("Config file not found: %v", err)
			log.Println("Creating default configuration...")
			cfg = config.Default()
			if err := cfg.Save(configPath); err != nil {
				log.Fatalf("Failed to save default config: %v", err)
			}
		} else {
			// Config file exists but has errors, don't overwrite
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	// Initialize database
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	
	// Handle migration commands
	if migrateToMariaDB {
		log.Println("Starting migration from SQLite to MariaDB...")
		if err := database.MigrateSQLiteToMariaDB(cfg); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration completed successfully!")
		os.Exit(0)
	}
	
	if migrateToSQLite {
		log.Println("Starting migration from MariaDB to SQLite...")
		if err := database.MigrateMariaDBToSQLite(cfg); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration completed successfully!")
		os.Exit(0)
	}
	
	if exportData != "" {
		log.Printf("Exporting database data to %s...", exportData)
		if err := database.ExportData(db.DB, exportData, cfg.Database.Type); err != nil {
			log.Fatalf("Export failed: %v", err)
		}
		log.Println("Export completed successfully!")
		os.Exit(0)
	}
	
	if importData != "" {
		log.Printf("Importing database data from %s...", importData)
		if err := database.ImportData(db.DB, importData); err != nil {
			log.Fatalf("Import failed: %v", err)
		}
		log.Println("Import completed successfully!")
		os.Exit(0)
	}
	
	if resetPassword != "" {
		log.Printf("Resetting password for user: %s", resetPassword)
		if err := resetUserPassword(db.DB, resetPassword); err != nil {
			log.Fatalf("Password reset failed: %v", err)
		}
		log.Println("Password reset completed successfully!")
		os.Exit(0)
	}

	// Create default admin user if needed
	if err := db.CreateDefaultAdminUser(); err != nil {
		log.Printf("Failed to create default admin user: %v", err)
	}

	// Initialize Gin router
	if cfg.App.Name == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.Default()

	// Load HTML templates
	templatesPattern := filepath.Join("web", "templates", "*.html")
	router.LoadHTMLGlob(templatesPattern)

	// Add build info to template context
	router.Use(func(c *gin.Context) {
		c.Set("BuildTime", BuildTime)
		c.Set("BuildDate", BuildDate)
		c.Next()
	})

	// Serve static files
	router.Static("/static", "./web/static")

	// Setup middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequireSetup(db.DB))
	router.Use(middleware.AuthMiddleware(db.DB))

	// Initialize handlers
	h := handlers.NewHandlers(db.DB, cfg)

	// Setup routes
	setupRoutes(router, h)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting Waterlogger server on %s", addr)
	log.Printf("Open your browser to: http://%s", addr)
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(router *gin.Engine, h *handlers.Handlers) {
	// Setup wizard routes
	router.GET("/setup", h.SetupWizardPage)
	router.POST("/api/setup", h.SetupWizardAPI)

	// Auth routes
	router.GET("/login", h.LoginPage)
	router.POST("/api/login", h.LoginAPI)
	router.POST("/api/logout", h.LogoutAPI)

	// Main application routes
	router.GET("/", h.Dashboard)
	router.GET("/pools", h.PoolsPage)
	router.GET("/kits", h.KitsPage)
	router.GET("/samples", h.SamplesPage)
	router.GET("/export", h.ExportPage)
	router.GET("/settings", h.SettingsPage)

	// API routes
	api := router.Group("/api")
	{
		// Users
		api.GET("/users", h.GetUsers)
		api.POST("/users", h.CreateUser)
		api.PUT("/users/:id", h.UpdateUser)
		api.DELETE("/users/:id", h.DeleteUser)

		// Pools
		api.GET("/pools", h.GetPools)
		api.POST("/pools", h.CreatePool)
		api.PUT("/pools/:id", h.UpdatePool)
		api.DELETE("/pools/:id", h.DeletePool)

		// Kits
		api.GET("/kits", h.GetKits)
		api.POST("/kits", h.CreateKit)
		api.PUT("/kits/:id", h.UpdateKit)
		api.DELETE("/kits/:id", h.DeleteKit)

		// Samples
		api.GET("/samples", h.GetSamples)
		api.POST("/samples", h.CreateSample)
		api.PUT("/samples/:id", h.UpdateSample)
		api.DELETE("/samples/:id", h.DeleteSample)


		// Export
		api.GET("/export", h.ExportBackup)
		api.GET("/export/excel", h.ExportExcel)
		api.GET("/export/markdown", h.ExportMarkdown)

		// Settings
		api.GET("/settings", h.GetSettings)
		api.POST("/settings", h.UpdateSettings)
		
		// Unit conversion
		api.POST("/convert", h.ConvertUnits)
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Page Not Found",
		})
	})
}

// resetUserPassword resets the password for a specified user
func resetUserPassword(db *gorm.DB, username string) error {
	// Find the user
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user '%s' not found", username)
		}
		return fmt.Errorf("database error: %v", err)
	}

	fmt.Printf("Found user: %s (%s)\n", user.Username, user.Email)
	
	var newPassword string
	var err error
	
	// Check if input is being piped or redirected
	if !term.IsTerminal(int(syscall.Stdin)) {
		// Input is being piped - read directly
		reader := bufio.NewReader(os.Stdin)
		newPassword, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read password from input: %v", err)
		}
		newPassword = strings.TrimSpace(newPassword)
		fmt.Println("Password read from input")
	} else {
		// Interactive mode - get password with confirmation
		newPassword, err = getPasswordFromInput("Enter new password: ")
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}
		
		confirmPassword, err := getPasswordFromInput("Confirm new password: ")
		if err != nil {
			return fmt.Errorf("failed to read password confirmation: %v", err)
		}
		
		if newPassword != confirmPassword {
			return fmt.Errorf("passwords do not match")
		}
	}
	
	// Validate password
	if errors := middleware.ValidatePassword(newPassword); len(errors) > 0 {
		return fmt.Errorf("password validation failed: %s", strings.Join(errors, ", "))
	}
	
	// Hash the new password
	hashedPassword, err := middleware.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	
	// Update the user's password
	if err := db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}
	
	fmt.Printf("Password successfully reset for user: %s\n", user.Username)
	return nil
}

// getPasswordFromInput securely reads a password from stdin
func getPasswordFromInput(prompt string) (string, error) {
	fmt.Print(prompt)
	
	// Try to read from terminal with hidden input
	if term.IsTerminal(int(syscall.Stdin)) {
		password, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // Add newline after hidden input
		return string(password), err
	}
	
	// Fallback to regular input if not a terminal
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(password), nil
}