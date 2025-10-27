package migrations

import (
	"gorm.io/gorm"
	"waterlogger/internal/logging"
)

// GetMigrationRunner creates and registers all migrations
func GetMigrationRunner(db *gorm.DB) *MigrationRunner {
	runner := NewMigrationRunner(db)

	// Register migrations in order
	runner.Register(&Migration001InitialSchema{})
	runner.Register(&Migration002AddTestColumn{})

	logging.Debug().Int("count", len(runner.migrations)).Msg("Registered migrations")

	return runner
}
