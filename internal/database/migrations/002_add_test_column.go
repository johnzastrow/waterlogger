package migrations

import (
	"gorm.io/gorm"
	"waterlogger/internal/logging"
)

// Migration002AddTestColumn adds a test column to demonstrate the migration system
type Migration002AddTestColumn struct{}

func (m *Migration002AddTestColumn) Version() string {
	return "002"
}

func (m *Migration002AddTestColumn) Name() string {
	return "add_test_column_to_pools"
}

// Up adds a test_column to the pools table
func (m *Migration002AddTestColumn) Up(db *gorm.DB) error {
	logging.Info().Msg("Adding test_column to pools table")

	// Add column - works for both SQLite and MariaDB
	err := db.Exec("ALTER TABLE pools ADD COLUMN test_column VARCHAR(255)").Error
	if err != nil {
		// Check if column already exists (migration re-run safety)
		var count int64
		checkErr := db.Raw("SELECT COUNT(*) FROM pragma_table_info('pools') WHERE name='test_column'").Scan(&count).Error
		if checkErr != nil {
			// Try MariaDB syntax
			checkErr = db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name='pools' AND column_name='test_column'").Scan(&count).Error
		}

		if checkErr == nil && count > 0 {
			logging.Info().Msg("test_column already exists, skipping")
			return nil
		}
		return err
	}

	logging.Info().Msg("test_column added successfully")
	return nil
}

// Down removes the test_column from the pools table
func (m *Migration002AddTestColumn) Down(db *gorm.DB) error {
	logging.Info().Msg("Removing test_column from pools table")

	// Drop column - works for both SQLite and MariaDB
	// Note: SQLite requires table recreation for column removal in older versions
	err := db.Exec("ALTER TABLE pools DROP COLUMN test_column").Error
	if err != nil {
		logging.Warn().Err(err).Msg("Failed to drop test_column (may not exist or SQLite limitation)")
		// Don't fail rollback if column doesn't exist
		return nil
	}

	logging.Info().Msg("test_column removed successfully")
	return nil
}
