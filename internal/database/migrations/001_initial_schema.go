package migrations

import (
	"gorm.io/gorm"
	"waterlogger/internal/logging"
	"waterlogger/internal/models"
)

// Migration001InitialSchema marks the existing schema as the baseline
type Migration001InitialSchema struct{}

func (m *Migration001InitialSchema) Version() string {
	return "001"
}

func (m *Migration001InitialSchema) Name() string {
	return "initial_schema"
}

// Up ensures all tables exist (using GORM AutoMigrate)
func (m *Migration001InitialSchema) Up(db *gorm.DB) error {
	logging.Info().Msg("Running initial schema migration")

	// Create all tables if they don't exist
	err := db.AutoMigrate(
		&models.User{},
		&models.UserPreferences{},
		&models.Pool{},
		&models.Kit{},
		&models.Sample{},
		&models.Measurements{},
		&models.Indices{},
		&models.Adjustment{},
	)

	if err != nil {
		return err
	}

	logging.Info().Msg("Initial schema migration completed")
	return nil
}

// Down is a no-op for initial schema (we don't want to drop everything)
func (m *Migration001InitialSchema) Down(db *gorm.DB) error {
	logging.Warn().Msg("Rollback of initial schema is not supported")
	return nil
}
