package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
	"waterlogger/internal/config"
	"waterlogger/internal/logging"
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
	Adjustments      []models.Adjustment    `json:"adjustments"`
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
	logging.Info().Str("path", backupPath).Msg("Creating database backup")
	
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
	
	// Backup Adjustments
	if err := dm.sourceDB.Find(&backup.Adjustments).Error; err != nil {
		return fmt.Errorf("failed to backup adjustments: %v", err)
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
	
	logging.Info().
		Int("users", len(backup.Users)).
		Int("pools", len(backup.Pools)).
		Int("samples", len(backup.Samples)).
		Int("adjustments", len(backup.Adjustments)).
		Msg("Backup created successfully")

	return nil
}

// cleanInvalidForeignKeys fixes invalid foreign key references in backup data
func cleanInvalidForeignKeys(backup *BackupData) error {
	logging.Info().Msg("Cleaning invalid foreign key references")

	// Build a map of valid user IDs
	validUserIDs := make(map[uint]bool)
	var defaultUserID uint = 1 // Fallback default

	for _, user := range backup.Users {
		validUserIDs[user.ID] = true
		if defaultUserID == 1 || user.ID < defaultUserID {
			defaultUserID = user.ID // Use the first/lowest valid user ID
		}
	}

	if len(validUserIDs) == 0 {
		logging.Warn().Msg("No users found in backup, cannot clean foreign keys")
		return fmt.Errorf("no users in backup data")
	}

	logging.Info().
		Uint("default_user_id", defaultUserID).
		Int("valid_users", len(validUserIDs)).
		Msg("Using default user ID for invalid references")

	fixedCount := 0

	// Fix UserPreferences
	for i := range backup.UserPreferences {
		if !validUserIDs[backup.UserPreferences[i].CreatedBy] && backup.UserPreferences[i].CreatedBy != 0 {
			logging.Warn().
				Uint("invalid_id", backup.UserPreferences[i].CreatedBy).
				Uint("fixed_to", defaultUserID).
				Msg("Fixed invalid created_by in user_preferences")
			backup.UserPreferences[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if !validUserIDs[backup.UserPreferences[i].UpdatedBy] && backup.UserPreferences[i].UpdatedBy != 0 {
			backup.UserPreferences[i].UpdatedBy = defaultUserID
			fixedCount++
		}
		// Fix zero values
		if backup.UserPreferences[i].CreatedBy == 0 {
			backup.UserPreferences[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if backup.UserPreferences[i].UpdatedBy == 0 {
			backup.UserPreferences[i].UpdatedBy = defaultUserID
			fixedCount++
		}
	}

	// Fix Pools
	for i := range backup.Pools {
		if !validUserIDs[backup.Pools[i].CreatedBy] && backup.Pools[i].CreatedBy != 0 {
			logging.Warn().
				Uint("invalid_id", backup.Pools[i].CreatedBy).
				Uint("fixed_to", defaultUserID).
				Str("pool", backup.Pools[i].Name).
				Msg("Fixed invalid created_by in pool")
			backup.Pools[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if !validUserIDs[backup.Pools[i].UpdatedBy] && backup.Pools[i].UpdatedBy != 0 {
			backup.Pools[i].UpdatedBy = defaultUserID
			fixedCount++
		}
		// Fix zero values
		if backup.Pools[i].CreatedBy == 0 {
			backup.Pools[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if backup.Pools[i].UpdatedBy == 0 {
			backup.Pools[i].UpdatedBy = defaultUserID
			fixedCount++
		}
	}

	// Build map of valid pool IDs
	validPoolIDs := make(map[uint]bool)
	var defaultPoolID uint = 0
	for _, pool := range backup.Pools {
		validPoolIDs[pool.ID] = true
		if defaultPoolID == 0 || pool.ID < defaultPoolID {
			defaultPoolID = pool.ID
		}
	}

	// Build map of valid kit IDs
	validKitIDs := make(map[uint]bool)
	var defaultKitID uint = 0
	for _, kit := range backup.Kits {
		validKitIDs[kit.ID] = true
		if defaultKitID == 0 || kit.ID < defaultKitID {
			defaultKitID = kit.ID
		}
	}

	// Fix Kits
	for i := range backup.Kits {
		if !validUserIDs[backup.Kits[i].CreatedBy] && backup.Kits[i].CreatedBy != 0 {
			backup.Kits[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if !validUserIDs[backup.Kits[i].UpdatedBy] && backup.Kits[i].UpdatedBy != 0 {
			backup.Kits[i].UpdatedBy = defaultUserID
			fixedCount++
		}
		if backup.Kits[i].CreatedBy == 0 {
			backup.Kits[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if backup.Kits[i].UpdatedBy == 0 {
			backup.Kits[i].UpdatedBy = defaultUserID
			fixedCount++
		}
	}

	// Fix Samples
	for i := range backup.Samples {
		// Fix created_by/updated_by
		if !validUserIDs[backup.Samples[i].CreatedBy] && backup.Samples[i].CreatedBy != 0 {
			backup.Samples[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if !validUserIDs[backup.Samples[i].UpdatedBy] && backup.Samples[i].UpdatedBy != 0 {
			backup.Samples[i].UpdatedBy = defaultUserID
			fixedCount++
		}
		if backup.Samples[i].CreatedBy == 0 {
			backup.Samples[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if backup.Samples[i].UpdatedBy == 0 {
			backup.Samples[i].UpdatedBy = defaultUserID
			fixedCount++
		}

		// Fix pool_id references
		if !validPoolIDs[backup.Samples[i].PoolID] {
			if defaultPoolID > 0 {
				logging.Warn().
					Uint("invalid_pool_id", backup.Samples[i].PoolID).
					Uint("fixed_to", defaultPoolID).
					Uint("sample_id", backup.Samples[i].ID).
					Msg("Fixed invalid pool_id in sample")
				backup.Samples[i].PoolID = defaultPoolID
				fixedCount++
			}
		}

		// Fix kit_id references
		if !validKitIDs[backup.Samples[i].KitID] {
			if defaultKitID > 0 {
				logging.Warn().
					Uint("invalid_kit_id", backup.Samples[i].KitID).
					Uint("fixed_to", defaultKitID).
					Uint("sample_id", backup.Samples[i].ID).
					Msg("Fixed invalid kit_id in sample")
				backup.Samples[i].KitID = defaultKitID
				fixedCount++
			}
		}
	}

	// Build map of valid sample IDs (needed for measurements and indices validation)
	validSampleIDs := make(map[uint]bool)
	for _, sample := range backup.Samples {
		validSampleIDs[sample.ID] = true
	}

	// Fix Measurements - remove any with invalid sample_id references
	validMeasurements := make([]models.Measurements, 0)
	for i := range backup.Measurements {
		// Check if sample_id is valid
		if !validSampleIDs[backup.Measurements[i].SampleID] {
			logging.Warn().
				Uint("invalid_sample_id", backup.Measurements[i].SampleID).
				Uint("measurement_id", backup.Measurements[i].ID).
				Msg("Removing measurement with invalid sample_id reference")
			fixedCount++
			continue // Skip this measurement
		}

		// Fix user references
		if !validUserIDs[backup.Measurements[i].CreatedBy] && backup.Measurements[i].CreatedBy != 0 {
			backup.Measurements[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if !validUserIDs[backup.Measurements[i].UpdatedBy] && backup.Measurements[i].UpdatedBy != 0 {
			backup.Measurements[i].UpdatedBy = defaultUserID
			fixedCount++
		}
		if backup.Measurements[i].CreatedBy == 0 {
			backup.Measurements[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if backup.Measurements[i].UpdatedBy == 0 {
			backup.Measurements[i].UpdatedBy = defaultUserID
			fixedCount++
		}

		validMeasurements = append(validMeasurements, backup.Measurements[i])
	}
	backup.Measurements = validMeasurements

	// Fix Indices - remove any with invalid sample_id references
	validIndices := make([]models.Indices, 0)
	for i := range backup.Indices {
		// Check if sample_id is valid
		if !validSampleIDs[backup.Indices[i].SampleID] {
			logging.Warn().
				Uint("invalid_sample_id", backup.Indices[i].SampleID).
				Uint("index_id", backup.Indices[i].ID).
				Msg("Removing index with invalid sample_id reference")
			fixedCount++
			continue // Skip this index
		}

		// Fix user references
		if !validUserIDs[backup.Indices[i].CreatedBy] && backup.Indices[i].CreatedBy != 0 {
			backup.Indices[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if !validUserIDs[backup.Indices[i].UpdatedBy] && backup.Indices[i].UpdatedBy != 0 {
			backup.Indices[i].UpdatedBy = defaultUserID
			fixedCount++
		}
		if backup.Indices[i].CreatedBy == 0 {
			backup.Indices[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if backup.Indices[i].UpdatedBy == 0 {
			backup.Indices[i].UpdatedBy = defaultUserID
			fixedCount++
		}

		validIndices = append(validIndices, backup.Indices[i])
	}
	backup.Indices = validIndices

	// Fix Adjustments
	for i := range backup.Adjustments {
		// Fix created_by/updated_by
		if !validUserIDs[backup.Adjustments[i].CreatedBy] && backup.Adjustments[i].CreatedBy != 0 {
			backup.Adjustments[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if !validUserIDs[backup.Adjustments[i].UpdatedBy] && backup.Adjustments[i].UpdatedBy != 0 {
			backup.Adjustments[i].UpdatedBy = defaultUserID
			fixedCount++
		}
		if backup.Adjustments[i].CreatedBy == 0 {
			backup.Adjustments[i].CreatedBy = defaultUserID
			fixedCount++
		}
		if backup.Adjustments[i].UpdatedBy == 0 {
			backup.Adjustments[i].UpdatedBy = defaultUserID
			fixedCount++
		}

		// Fix pool_id references
		if !validPoolIDs[backup.Adjustments[i].PoolID] {
			if defaultPoolID > 0 {
				logging.Warn().
					Uint("invalid_pool_id", backup.Adjustments[i].PoolID).
					Uint("fixed_to", defaultPoolID).
					Uint("adjustment_id", backup.Adjustments[i].ID).
					Msg("Fixed invalid pool_id in adjustment")
				backup.Adjustments[i].PoolID = defaultPoolID
				fixedCount++
			}
		}

		// Fix sample_id references (if not NULL)
		if backup.Adjustments[i].SampleID != nil {
			if !validSampleIDs[*backup.Adjustments[i].SampleID] {
				logging.Warn().
					Uint("invalid_sample_id", *backup.Adjustments[i].SampleID).
					Uint("adjustment_id", backup.Adjustments[i].ID).
					Msg("Clearing invalid sample_id in adjustment")
				backup.Adjustments[i].SampleID = nil // Set to NULL if invalid
				fixedCount++
			}
		}
	}

	logging.Info().
		Int("fixed_references", fixedCount).
		Msg("Completed cleaning invalid foreign key references")

	return nil
}

// RestoreFromBackup restores data from a backup file
func (dm *DatabaseMigrator) RestoreFromBackup(backupPath string) error {
	logging.Info().Str("path", backupPath).Msg("Restoring from backup")

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
		&models.Adjustment{},
	); err != nil {
		return fmt.Errorf("failed to migrate target database schema: %v", err)
	}

	// Clean invalid foreign key references before restoring
	if err := cleanInvalidForeignKeys(&backup); err != nil {
		return fmt.Errorf("failed to clean invalid foreign keys: %v", err)
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
	
	// 8. Adjustments (depends on Pools, optionally on Samples)
	if len(backup.Adjustments) > 0 {
		if err := dm.targetDB.Create(&backup.Adjustments).Error; err != nil {
			return fmt.Errorf("failed to restore adjustments: %v", err)
		}
	}
	
	logging.Info().Msg("Restore completed successfully")
	return nil
}

// MigrateDatabase migrates data from SQLite to MariaDB or vice versa
func (dm *DatabaseMigrator) MigrateDatabase(tempBackupPath string) error {
	logging.Info().Msg("Starting database migration")

	// Step 1: Create backup of source database
	if err := dm.CreateBackup(tempBackupPath); err != nil {
		logging.Error().Err(err).Msg("Failed to create backup")
		return fmt.Errorf("failed to create backup: %v", err)
	}

	// Step 2: Restore to target database
	if err := dm.RestoreFromBackup(tempBackupPath); err != nil {
		logging.Error().Err(err).Msg("Failed to restore to target database")
		return fmt.Errorf("failed to restore to target database: %v", err)
	}

	// Step 3: Clean up temporary backup file
	if err := os.Remove(tempBackupPath); err != nil {
		logging.Warn().Err(err).Str("path", tempBackupPath).Msg("Failed to remove temporary backup file")
	}

	logging.Info().Msg("Database migration completed successfully")
	return nil
}

// MigrateSQLiteToMariaDB migrates data from SQLite to MariaDB
func MigrateSQLiteToMariaDB(cfg *config.Config) error {
	logging.Info().Msg("Migrating from SQLite to MariaDB")
	
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
	logging.Info().Msg("Migrating from MariaDB to SQLite")
	
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