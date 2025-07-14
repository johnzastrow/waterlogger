package database

import (
	"fmt"
	"log"

	"waterlogger/internal/config"
	"waterlogger/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch cfg.Database.Type {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.Database.SQLite.Path), gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
		}
	case "mariadb":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.MariaDB.Username,
			cfg.Database.MariaDB.Password,
			cfg.Database.MariaDB.Host,
			cfg.Database.MariaDB.Port,
			cfg.Database.MariaDB.Database,
		)
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MariaDB database: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(
		&models.User{},
		&models.UserPreferences{},
		&models.Pool{},
		&models.Kit{},
		&models.Sample{},
		&models.Measurements{},
		&models.Indices{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate schema: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CreateDefaultAdminUser creates a default admin user if no users exist
func (db *DB) CreateDefaultAdminUser() error {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		log.Println("No users found, creating default admin user")
		// This will be replaced by setup wizard
		return nil
	}

	return nil
}

// Migration utilities
func (db *DB) ExportData() (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// Export users
	var users []models.User
	if err := db.Preload("Preferences").Find(&users).Error; err != nil {
		return nil, err
	}
	data["users"] = users

	// Export pools
	var pools []models.Pool
	if err := db.Find(&pools).Error; err != nil {
		return nil, err
	}
	data["pools"] = pools

	// Export kits
	var kits []models.Kit
	if err := db.Find(&kits).Error; err != nil {
		return nil, err
	}
	data["kits"] = kits

	// Export samples with measurements and indices
	var samples []models.Sample
	if err := db.Preload("Measurements").Preload("Indices").Find(&samples).Error; err != nil {
		return nil, err
	}
	data["samples"] = samples

	return data, nil
}

func (db *DB) ImportData(data map[string]interface{}) error {
	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Import users
	if users, ok := data["users"].([]models.User); ok {
		for _, user := range users {
			if err := tx.Create(&user).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Import pools
	if pools, ok := data["pools"].([]models.Pool); ok {
		for _, pool := range pools {
			if err := tx.Create(&pool).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Import kits
	if kits, ok := data["kits"].([]models.Kit); ok {
		for _, kit := range kits {
			if err := tx.Create(&kit).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Import samples
	if samples, ok := data["samples"].([]models.Sample); ok {
		for _, sample := range samples {
			if err := tx.Create(&sample).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (db *DB) BackupSQLite(backupPath string) error {
	if db.Dialector.Name() != "sqlite" {
		return fmt.Errorf("backup only supported for SQLite")
	}

	// For SQLite backup, we need to get the database file path
	// This is a simplified approach - in production, use SQLite backup API
	return fmt.Errorf("backup functionality not implemented yet")
}