package main

import (
	"bufio"
	"flag"
	"fmt"
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
	"waterlogger/internal/database/migrations"
	"waterlogger/internal/handlers"
	"waterlogger/internal/logging"
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
	var migrationStatus bool
	var migrationRollback bool

	flag.StringVar(&configPath, "config", "config.yaml", "Path to configuration file")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.BoolVar(&migrateToMariaDB, "migrate-to-mariadb", false, "Migrate data from SQLite to MariaDB")
	flag.BoolVar(&migrateToSQLite, "migrate-to-sqlite", false, "Migrate data from MariaDB to SQLite")
	flag.StringVar(&exportData, "export", "", "Export database data to backup file")
	flag.StringVar(&importData, "import", "", "Import database data from backup file")
	flag.StringVar(&resetPassword, "reset-password", "", "Reset password for specified username")
	flag.BoolVar(&migrationStatus, "migration-status", false, "Show database migration status")
	flag.BoolVar(&migrationRollback, "migration-rollback", false, "Rollback the last applied migration")
	flag.Parse()

	if showVersion {
		fmt.Println("Waterlogger v1.5.0")
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
		fmt.Println()
		fmt.Println("Database Migration:")
		fmt.Println("  -migrate-to-mariadb      Migrate data from SQLite to MariaDB")
		fmt.Println("  -migrate-to-sqlite       Migrate data from MariaDB to SQLite")
		fmt.Println()
		fmt.Println("Schema Migration:")
		fmt.Println("  -migration-status        Show database schema migration status")
		fmt.Println("  -migration-rollback      Rollback the last applied schema migration")
		fmt.Println()
		fmt.Println("Data Management:")
		fmt.Println("  -export string           Export database data to backup file")
		fmt.Println("  -import string           Import database data from backup file")
		fmt.Println("  -reset-password string   Reset password for specified username")
		fmt.Println()
		fmt.Println("For more information, visit: https://github.com/johnzastrow/waterlogger")
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist, create default
			fmt.Printf("Config file not found: %v\n", err)
			fmt.Println("Creating default configuration...")
			cfg = config.Default()
			if err := cfg.Save(configPath); err != nil {
				fmt.Printf("Failed to save default config: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Config file exists but has errors, don't overwrite
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize logging
	logConfig := &logging.Config{
		Level:      cfg.Logging.Level,
		Format:     cfg.Logging.Format,
		Output:     cfg.Logging.Output,
		FilePath:   cfg.Logging.FilePath,
		MaxSize:    cfg.Logging.MaxSize,
		MaxBackups: cfg.Logging.MaxBackups,
		MaxAge:     cfg.Logging.MaxAge,
		Compress:   cfg.Logging.Compress,
	}
	if err := logging.Initialize(logConfig); err != nil {
		fmt.Printf("Failed to initialize logging: %v\n", err)
		os.Exit(1)
	}

	logging.Info().Msg("Waterlogger starting...")
	logging.Info().Str("version", "1.3.0").Str("build_time", BuildTime).Str("build_date", BuildDate).Msg("Build information")

	// Initialize database
	db, err := database.NewDB(cfg)
	if err != nil {
		logging.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()
	
	// Handle migration commands
	if migrateToMariaDB {
		logging.Info().Msg("Starting migration from SQLite to MariaDB")
		if err := database.MigrateSQLiteToMariaDB(cfg); err != nil {
			logging.Fatal().Err(err).Msg("Migration failed")
		}
		logging.Info().Msg("Migration completed successfully")
		os.Exit(0)
	}

	if migrateToSQLite {
		logging.Info().Msg("Starting migration from MariaDB to SQLite")
		if err := database.MigrateMariaDBToSQLite(cfg); err != nil {
			logging.Fatal().Err(err).Msg("Migration failed")
		}
		logging.Info().Msg("Migration completed successfully")
		os.Exit(0)
	}

	if exportData != "" {
		logging.Info().Str("file", exportData).Msg("Exporting database data")
		if err := database.ExportData(db.DB, exportData, cfg.Database.Type); err != nil {
			logging.Fatal().Err(err).Str("file", exportData).Msg("Export failed")
		}
		logging.Info().Str("file", exportData).Msg("Export completed successfully")
		os.Exit(0)
	}

	if importData != "" {
		logging.Info().Str("file", importData).Msg("Importing database data")
		if err := database.ImportData(db.DB, importData); err != nil {
			logging.Fatal().Err(err).Str("file", importData).Msg("Import failed")
		}
		logging.Info().Str("file", importData).Msg("Import completed successfully")
		os.Exit(0)
	}

	if resetPassword != "" {
		logging.Info().Str("username", resetPassword).Msg("Resetting user password")
		if err := resetUserPassword(db.DB, resetPassword); err != nil {
			logging.Fatal().Err(err).Str("username", resetPassword).Msg("Password reset failed")
		}
		logging.Info().Str("username", resetPassword).Msg("Password reset completed successfully")
		os.Exit(0)
	}

	// Handle schema migration commands
	if migrationStatus {
		logging.Info().Msg("Showing migration status")
		migrationRunner := migrations.GetMigrationRunner(db.DB)
		if err := migrationRunner.Initialize(); err != nil {
			logging.Fatal().Err(err).Msg("Failed to initialize migration system")
		}

		status, err := migrationRunner.GetMigrationStatus()
		if err != nil {
			logging.Fatal().Err(err).Msg("Failed to get migration status")
		}

		fmt.Println("\nDatabase Migration Status:")
		fmt.Println("==========================")
		for _, s := range status {
			appliedStr := "[ ]"
			if s.Applied {
				appliedStr = "[✓]"
			}
			fmt.Printf("%s %s - %s\n", appliedStr, s.Version, s.Name)
		}
		fmt.Println()
		os.Exit(0)
	}

	if migrationRollback {
		logging.Info().Msg("Rolling back last migration")
		migrationRunner := migrations.GetMigrationRunner(db.DB)
		if err := migrationRunner.Initialize(); err != nil {
			logging.Fatal().Err(err).Msg("Failed to initialize migration system")
		}

		if err := migrationRunner.Rollback(); err != nil {
			logging.Fatal().Err(err).Msg("Rollback failed")
		}
		logging.Info().Msg("Rollback completed successfully")
		os.Exit(0)
	}

	// Create default admin user if needed
	if err := db.CreateDefaultAdminUser(); err != nil {
		logging.Error().Err(err).Msg("Failed to create default admin user")
	}

	// Initialize Gin router
	if cfg.App.Name == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router without default middleware (we'll add our own logging)
	router := gin.New()

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Add our custom logging middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.AuditLoggingMiddleware())

	// Load HTML templates
	templatesPattern := filepath.Join("web", "templates", "*.html")
	router.LoadHTMLGlob(templatesPattern)
	logging.Info().Str("pattern", templatesPattern).Msg("Loaded HTML templates")

	// Add build info and version to template context
	router.Use(func(c *gin.Context) {
		c.Set("BuildTime", BuildTime)
		c.Set("BuildDate", BuildDate)
		c.Set("Version", cfg.App.Version)
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
	logging.Info().
		Str("address", addr).
		Str("host", cfg.Server.Host).
		Int("port", cfg.Server.Port).
		Msg("Starting Waterlogger server")
	logging.Info().Str("url", fmt.Sprintf("http://%s", addr)).Msg("Open your browser to access Waterlogger")

	if err := router.Run(addr); err != nil {
		logging.Fatal().Err(err).Str("address", addr).Msg("Failed to start server")
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
	router.GET("/adjustments", h.AdjustmentsPage)
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
		api.GET("/settings/migrations", h.GetMigrationStatus)
		api.POST("/settings/backup", h.ExportDatabaseBackup)
		api.POST("/settings/import", h.ImportDatabaseBackup)

		// Unit conversion
		api.POST("/convert", h.ConvertUnits)
		
		// Volume calculation
		api.POST("/volume/calculate", h.CalculateVolume)
		api.GET("/volume/info", h.GetVolumeInfo)
		
		// Chemical adjustments
		api.POST("/adjustments/calculate", h.CalculateAdjustments)
		api.GET("/adjustments/target-ranges", h.GetTargetRanges)
		api.POST("/adjustments", h.SaveAdjustment)
		api.GET("/adjustments", h.GetAdjustments)
		api.GET("/adjustments/:id", h.GetAdjustment)
		api.DELETE("/adjustments/:id", h.DeleteAdjustment)
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
			logging.Error().Str("username", username).Msg("User not found")
			return fmt.Errorf("user '%s' not found", username)
		}
		logging.Error().Err(err).Str("username", username).Msg("Database error looking up user")
		return fmt.Errorf("database error: %v", err)
	}

	logging.Info().Str("username", user.Username).Str("email", user.Email).Msg("Found user for password reset")
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