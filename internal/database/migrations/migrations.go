package migrations

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
	"waterlogger/internal/logging"
	"waterlogger/internal/models"
)

// Migration represents a database schema migration
type Migration interface {
	Version() string            // Unique version identifier (e.g., "001", "002")
	Name() string               // Human-readable name (e.g., "add_pool_dimensions")
	Up(db *gorm.DB) error      // Apply the migration
	Down(db *gorm.DB) error    // Rollback the migration
}

// MigrationRunner manages and executes database migrations
type MigrationRunner struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *gorm.DB) *MigrationRunner {
	return &MigrationRunner{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// Register adds a migration to the runner
func (mr *MigrationRunner) Register(migration Migration) {
	mr.migrations = append(mr.migrations, migration)
}

// Initialize creates the schema_migrations table if it doesn't exist
func (mr *MigrationRunner) Initialize() error {
	logging.Info().Msg("Initializing migration system")

	// Create schema_migrations table
	if err := mr.db.AutoMigrate(&models.SchemaMigration{}); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	logging.Info().Msg("Migration system initialized")
	return nil
}

// GetAppliedMigrations returns a map of applied migration versions
func (mr *MigrationRunner) GetAppliedMigrations() (map[string]bool, error) {
	var appliedMigrations []models.SchemaMigration
	if err := mr.db.Find(&appliedMigrations).Error; err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}

	applied := make(map[string]bool)
	for _, m := range appliedMigrations {
		applied[m.Version] = true
	}

	return applied, nil
}

// GetPendingMigrations returns migrations that haven't been applied yet
func (mr *MigrationRunner) GetPendingMigrations() ([]Migration, error) {
	applied, err := mr.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	pending := make([]Migration, 0)
	for _, migration := range mr.migrations {
		if !applied[migration.Version()] {
			pending = append(pending, migration)
		}
	}

	// Sort by version
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Version() < pending[j].Version()
	})

	return pending, nil
}

// RunPendingMigrations applies all pending migrations
func (mr *MigrationRunner) RunPendingMigrations() error {
	pending, err := mr.GetPendingMigrations()
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		logging.Info().Msg("No pending migrations to apply")
		return nil
	}

	logging.Info().Int("count", len(pending)).Msg("Found pending migrations")

	for _, migration := range pending {
		if err := mr.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version(), err)
		}
	}

	logging.Info().Msg("All pending migrations applied successfully")
	return nil
}

// applyMigration applies a single migration within a transaction
func (mr *MigrationRunner) applyMigration(migration Migration) error {
	logging.Info().
		Str("version", migration.Version()).
		Str("name", migration.Name()).
		Msg("Applying migration")

	// Run migration in a transaction
	err := mr.db.Transaction(func(tx *gorm.DB) error {
		// Apply the migration
		if err := migration.Up(tx); err != nil {
			return fmt.Errorf("migration Up() failed: %w", err)
		}

		// Record the migration
		record := models.SchemaMigration{
			Version:   migration.Version(),
			Name:      migration.Name(),
			AppliedAt: time.Now(),
		}

		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		return nil
	})

	if err != nil {
		logging.Error().
			Err(err).
			Str("version", migration.Version()).
			Str("name", migration.Name()).
			Msg("Migration failed")
		return err
	}

	logging.Info().
		Str("version", migration.Version()).
		Str("name", migration.Name()).
		Msg("Migration applied successfully")

	return nil
}

// Rollback rolls back the last applied migration
func (mr *MigrationRunner) Rollback() error {
	// Get the last applied migration
	var lastMigration models.SchemaMigration
	if err := mr.db.Order("applied_at DESC").First(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logging.Info().Msg("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to query last migration: %w", err)
	}

	// Find the migration in our registry
	var migrationToRollback Migration
	for _, m := range mr.migrations {
		if m.Version() == lastMigration.Version {
			migrationToRollback = m
			break
		}
	}

	if migrationToRollback == nil {
		return fmt.Errorf("migration %s not found in registry", lastMigration.Version)
	}

	logging.Info().
		Str("version", migrationToRollback.Version()).
		Str("name", migrationToRollback.Name()).
		Msg("Rolling back migration")

	// Run rollback in a transaction
	err := mr.db.Transaction(func(tx *gorm.DB) error {
		// Rollback the migration
		if err := migrationToRollback.Down(tx); err != nil {
			return fmt.Errorf("migration Down() failed: %w", err)
		}

		// Delete the migration record
		if err := tx.Delete(&lastMigration).Error; err != nil {
			return fmt.Errorf("failed to delete migration record: %w", err)
		}

		return nil
	})

	if err != nil {
		logging.Error().
			Err(err).
			Str("version", migrationToRollback.Version()).
			Str("name", migrationToRollback.Name()).
			Msg("Rollback failed")
		return err
	}

	logging.Info().
		Str("version", migrationToRollback.Version()).
		Str("name", migrationToRollback.Name()).
		Msg("Migration rolled back successfully")

	return nil
}

// GetMigrationStatus returns the status of all migrations
func (mr *MigrationRunner) GetMigrationStatus() ([]MigrationStatus, error) {
	applied, err := mr.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	status := make([]MigrationStatus, 0, len(mr.migrations))

	for _, migration := range mr.migrations {
		s := MigrationStatus{
			Version: migration.Version(),
			Name:    migration.Name(),
			Applied: applied[migration.Version()],
		}
		status = append(status, s)
	}

	// Sort by version
	sort.Slice(status, func(i, j int) bool {
		return status[i].Version < status[j].Version
	})

	return status, nil
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version string
	Name    string
	Applied bool
}
