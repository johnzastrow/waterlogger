package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
	"waterlogger/internal/config"
	"waterlogger/internal/models"
)

// BackupData represents a complete database backup
type BackupData struct {
	Timestamp        time.Time              `json:"timestamp"`
	SourceDatabase   string                 `json:"source_database"`
	Users            []models.User          `json:"users"`
	UserPreferences  []models.UserPreferences `json:"user_preferences"`
	Pools            []models.Pool          `json:"pools"`
	Kits             []models.Kit           `json:"kits"`
	Samples          []models.Sample        `json:"samples"`
	Measurements     []models.Measurements  `json:"measurements"`
	Indices          []models.Indices       `json:"indices"`
}

// DatabaseMigrator handles database migrations between SQLite and MariaDB
type DatabaseMigrator struct {
	sourceDB *gorm.DB
	targetDB *gorm.DB
}

// NewDatabaseMigrator creates a new database migrator
func NewDatabaseMigrator(sourceDB, targetDB *gorm.DB) *DatabaseMigrator {
	return &DatabaseMigrator{
		sourceDB: sourceDB,
		targetDB: targetDB,
	}
}

// CreateBackup creates a complete backup of the database
func (dm *DatabaseMigrator) CreateBackup(backupPath string) error {
	log.Printf("Creating backup at %s", backupPath)
	
	backup := BackupData{
		Timestamp:      time.Now(),
		SourceDatabase: "unknown", // Will be set by caller
	}
	
	// Backup Users
	if err := dm.sourceDB.Find(&backup.Users).Error; err != nil {
		return fmt.Errorf("failed to backup users: %v", err)
	}
	
	// Backup UserPreferences
	if err := dm.sourceDB.Find(&backup.UserPreferences).Error; err != nil {
		return fmt.Errorf("failed to backup user preferences: %v", err)
	}
	
	// Backup Pools
	if err := dm.sourceDB.Find(&backup.Pools).Error; err != nil {
		return fmt.Errorf("failed to backup pools: %v", err)
	}
	
	// Backup Kits
	if err := dm.sourceDB.Find(&backup.Kits).Error; err != nil {
		return fmt.Errorf("failed to backup kits: %v", err)
	}
	
	// Backup Samples
	if err := dm.sourceDB.Find(&backup.Samples).Error; err != nil {
		return fmt.Errorf("failed to backup samples: %v", err)
	}
	
	// Backup Measurements
	if err := dm.sourceDB.Find(&backup.Measurements).Error; err != nil {
		return fmt.Errorf("failed to backup measurements: %v", err)
	}
	
	// Backup Indices
	if err := dm.sourceDB.Find(&backup.Indices).Error; err != nil {
		return fmt.Errorf("failed to backup indices: %v", err)
	}
	
	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}
	
	// Write backup to file
	file, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(backup); err != nil {
		return fmt.Errorf("failed to write backup data: %v", err)
	}
	
	log.Printf("Backup created successfully with %d users, %d pools, %d samples", 
		len(backup.Users), len(backup.Pools), len(backup.Samples))
	
	return nil
}

// RestoreFromBackup restores data from a backup file
func (dm *DatabaseMigrator) RestoreFromBackup(backupPath string) error {
	log.Printf("Restoring from backup at %s", backupPath)
	
	// Read backup file
	file, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %v", err)
	}
	defer file.Close()
	
	var backup BackupData
	if err := json.NewDecoder(file).Decode(&backup); err != nil {
		return fmt.Errorf("failed to decode backup data: %v", err)
	}
	
	// Ensure target database has the correct schema
	if err := dm.targetDB.AutoMigrate(
		&models.User{},
		&models.UserPreferences{},
		&models.Pool{},
		&models.Kit{},
		&models.Sample{},
		&models.Measurements{},
		&models.Indices{},
	); err != nil {
		return fmt.Errorf("failed to migrate target database schema: %v", err)
	}
	
	// Restore data in the correct order (respecting foreign key constraints)
	
	// 1. Users (no dependencies)
	if len(backup.Users) > 0 {
		if err := dm.targetDB.Create(&backup.Users).Error; err != nil {
			return fmt.Errorf("failed to restore users: %v", err)
		}
	}
	
	// 2. UserPreferences (depends on Users)
	if len(backup.UserPreferences) > 0 {
		if err := dm.targetDB.Create(&backup.UserPreferences).Error; err != nil {
			return fmt.Errorf("failed to restore user preferences: %v", err)
		}
	}
	
	// 3. Pools (depends on Users)
	if len(backup.Pools) > 0 {
		if err := dm.targetDB.Create(&backup.Pools).Error; err != nil {
			return fmt.Errorf("failed to restore pools: %v", err)
		}
	}
	
	// 4. Kits (no dependencies)
	if len(backup.Kits) > 0 {
		if err := dm.targetDB.Create(&backup.Kits).Error; err != nil {
			return fmt.Errorf("failed to restore kits: %v", err)
		}
	}
	
	// 5. Samples (depends on Pools, Users, Kits)
	if len(backup.Samples) > 0 {
		if err := dm.targetDB.Create(&backup.Samples).Error; err != nil {
			return fmt.Errorf("failed to restore samples: %v", err)
		}
	}
	
	// 6. Measurements (depends on Samples)
	if len(backup.Measurements) > 0 {
		if err := dm.targetDB.Create(&backup.Measurements).Error; err != nil {
			return fmt.Errorf("failed to restore measurements: %v", err)
		}
	}
	
	// 7. Indices (depends on Samples)
	if len(backup.Indices) > 0 {
		if err := dm.targetDB.Create(&backup.Indices).Error; err != nil {
			return fmt.Errorf("failed to restore indices: %v", err)
		}
	}
	
	log.Printf("Restore completed successfully")
	return nil
}

// MigrateDatabase migrates data from SQLite to MariaDB or vice versa
func (dm *DatabaseMigrator) MigrateDatabase(tempBackupPath string) error {
	log.Printf("Starting database migration")
	
	// Step 1: Create backup of source database
	if err := dm.CreateBackup(tempBackupPath); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}
	
	// Step 2: Restore to target database
	if err := dm.RestoreFromBackup(tempBackupPath); err != nil {
		return fmt.Errorf("failed to restore to target database: %v", err)
	}
	
	// Step 3: Clean up temporary backup file
	if err := os.Remove(tempBackupPath); err != nil {
		log.Printf("Warning: failed to remove temporary backup file: %v", err)
	}
	
	log.Printf("Database migration completed successfully")
	return nil
}

// MigrateSQLiteToMariaDB migrates data from SQLite to MariaDB
func MigrateSQLiteToMariaDB(cfg *config.Config) error {
	log.Printf("Migrating from SQLite to MariaDB")
	
	// Create SQLite connection
	sqliteConfig := &config.Config{
		Database: config.DatabaseConfig{
			Type:   "sqlite",
			SQLite: cfg.Database.SQLite,
		},
	}
	
	sqliteDB, err := NewDB(sqliteConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SQLite database: %v", err)
	}
	defer sqliteDB.Close()
	
	// Create MariaDB connection
	mariadbConfig := &config.Config{
		Database: config.DatabaseConfig{
			Type:    "mariadb",
			MariaDB: cfg.Database.MariaDB,
		},
	}
	
	mariadbDB, err := NewDB(mariadbConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to MariaDB database: %v", err)
	}
	defer mariadbDB.Close()
	
	// Perform migration
	migrator := NewDatabaseMigrator(sqliteDB.DB, mariadbDB.DB)
	tempBackupPath := fmt.Sprintf("temp_backup_%d.json", time.Now().Unix())
	
	return migrator.MigrateDatabase(tempBackupPath)
}

// MigrateMariaDBToSQLite migrates data from MariaDB to SQLite
func MigrateMariaDBToSQLite(cfg *config.Config) error {
	log.Printf("Migrating from MariaDB to SQLite")
	
	// Create MariaDB connection
	mariadbConfig := &config.Config{
		Database: config.DatabaseConfig{
			Type:    "mariadb",
			MariaDB: cfg.Database.MariaDB,
		},
	}
	
	mariadbDB, err := NewDB(mariadbConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to MariaDB database: %v", err)
	}
	defer mariadbDB.Close()
	
	// Create SQLite connection
	sqliteConfig := &config.Config{
		Database: config.DatabaseConfig{
			Type:   "sqlite",
			SQLite: cfg.Database.SQLite,
		},
	}
	
	sqliteDB, err := NewDB(sqliteConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SQLite database: %v", err)
	}
	defer sqliteDB.Close()
	
	// Perform migration
	migrator := NewDatabaseMigrator(mariadbDB.DB, sqliteDB.DB)
	tempBackupPath := fmt.Sprintf("temp_backup_%d.json", time.Now().Unix())
	
	return migrator.MigrateDatabase(tempBackupPath)
}

// ExportData exports database data to a backup file
func ExportData(db *gorm.DB, backupPath string, databaseType string) error {
	migrator := &DatabaseMigrator{sourceDB: db}
	
	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}
	
	return migrator.CreateBackup(backupPath)
}

// ImportData imports database data from a backup file
func ImportData(db *gorm.DB, backupPath string) error {
	migrator := &DatabaseMigrator{targetDB: db}
	return migrator.RestoreFromBackup(backupPath)
}