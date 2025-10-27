# Database Schema Migrations

Waterlogger uses a versioned migration system to manage database schema changes. This ensures safe, traceable, and reversible database upgrades across both SQLite and MariaDB.

## Overview

The migration system provides:
- **Schema version tracking** - Know exactly which migrations have been applied
- **Automatic migrations** - Pending migrations run automatically on startup
- **Rollback capability** - Safely revert migrations if needed
- **Cross-database support** - Works with both SQLite and MariaDB
- **Migration history** - Complete audit trail of schema changes

## Migration Basics

### What Gets Migrated

The migration system handles **schema changes**:
- ✅ Adding/removing tables
- ✅ Adding/removing columns
- ✅ Changing column types
- ✅ Adding/removing indexes
- ✅ Modifying foreign keys
- ✅ Data transformations

### What Doesn't Get Migrated

The migration system does NOT handle **data migrations** (SQLite ↔ MariaDB). For that, use:
```bash
# Migrate data from SQLite to MariaDB
./waterlogger -config config.yaml -migrate-to-mariadb

# Migrate data from MariaDB to SQLite
./waterlogger -config config.yaml -migrate-to-sqlite
```

## Command-Line Interface

### Check Migration Status

See which migrations have been applied:

```bash
./waterlogger -config config.yaml -migration-status
```

Output:
```
Database Migration Status:
==========================
[✓] 001 - initial_schema
[✓] 002 - add_test_column_to_pools
[ ] 003 - add_user_roles
```

### Rollback Last Migration

Revert the most recently applied migration:

```bash
./waterlogger -config config.yaml -migration-rollback
```

**Warning**: Only use rollback if you understand the implications. Data loss may occur.

### Apply Pending Migrations

Migrations run automatically on application startup. Just start the application normally:

```bash
./waterlogger -config config.yaml
```

The system will:
1. Check for pending migrations
2. Apply them in order
3. Record each successful migration
4. Continue with normal startup

## Creating New Migrations

### Step 1: Create Migration File

Create a new file in `internal/database/migrations/` with the naming pattern:
```
NNN_description.go
```

Where:
- `NNN` = Zero-padded version number (e.g., 003, 004, 005)
- `description` = Snake_case description (e.g., add_user_roles, remove_old_table)

### Step 2: Implement Migration Interface

```go
package migrations

import (
	"gorm.io/gorm"
	"waterlogger/internal/logging"
)

// Migration003AddUserRoles adds a roles column to users table
type Migration003AddUserRoles struct{}

func (m *Migration003AddUserRoles) Version() string {
	return "003"
}

func (m *Migration003AddUserRoles) Name() string {
	return "add_user_roles"
}

// Up applies the migration
func (m *Migration003AddUserRoles) Up(db *gorm.DB) error {
	logging.Info().Msg("Adding roles column to users table")

	// Add column - works for both SQLite and MariaDB
	err := db.Exec("ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'user'").Error
	if err != nil {
		return err
	}

	logging.Info().Msg("Roles column added successfully")
	return nil
}

// Down reverts the migration
func (m *Migration003AddUserRoles) Down(db *gorm.DB) error {
	logging.Info().Msg("Removing roles column from users table")

	err := db.Exec("ALTER TABLE users DROP COLUMN role").Error
	if err != nil {
		logging.Warn().Err(err).Msg("Failed to drop role column")
		return nil // Don't fail rollback if column doesn't exist
	}

	logging.Info().Msg("Roles column removed successfully")
	return nil
}
```

### Step 3: Register Migration

Add your migration to `internal/database/migrations/registry.go`:

```go
func GetMigrationRunner(db *gorm.DB) *MigrationRunner {
	runner := NewMigrationRunner(db)

	// Register migrations in order
	runner.Register(&Migration001InitialSchema{})
	runner.Register(&Migration002AddTestColumn{})
	runner.Register(&Migration003AddUserRoles{})  // <- Add your migration here

	logging.Debug().Int("count", len(runner.migrations)).Msg("Registered migrations")

	return runner
}
```

### Step 4: Test Your Migration

```bash
# Build the application
./build.sh

# Check migration status (should show your migration as pending)
./waterlogger -config config.yaml -migration-status

# Start application (will apply your migration)
./waterlogger -config config.yaml

# Verify migration was applied
./waterlogger -config config.yaml -migration-status

# Test rollback (optional)
./waterlogger -config config.yaml -migration-rollback
```

## Migration Best Practices

### 1. Make Migrations Idempotent

Migrations should be safe to run multiple times:

```go
// GOOD: Check if column exists before adding
func (m *MyMigration) Up(db *gorm.DB) error {
	// Try to add column
	err := db.Exec("ALTER TABLE pools ADD COLUMN new_field VARCHAR(255)").Error
	if err != nil {
		// Check if column already exists
		var count int64
		checkErr := db.Raw("SELECT COUNT(*) FROM pragma_table_info('pools') WHERE name='new_field'").Scan(&count).Error
		if checkErr == nil && count > 0 {
			logging.Info().Msg("Column already exists, skipping")
			return nil
		}
		return err
	}
	return nil
}

// BAD: Assumes column doesn't exist
func (m *MyMigration) Up(db *gorm.DB) error {
	return db.Exec("ALTER TABLE pools ADD COLUMN new_field VARCHAR(255)").Error
}
```

### 2. Handle Both SQLite and MariaDB

Test migrations on both database types:

```go
func (m *MyMigration) Up(db *gorm.DB) error {
	// This works for both SQLite and MariaDB
	return db.Exec("ALTER TABLE pools ADD COLUMN test_field VARCHAR(255)").Error
}
```

### 3. Use Transactions

Migrations automatically run in transactions. If any part fails, the entire migration rolls back.

### 4. Add Logging

Help with debugging:

```go
func (m *MyMigration) Up(db *gorm.DB) error {
	logging.Info().Msg("Starting complex migration")

	// Step 1
	logging.Debug().Msg("Creating new table")
	if err := db.Exec("CREATE TABLE...").Error; err != nil {
		return err
	}

	// Step 2
	logging.Debug().Msg("Migrating data")
	if err := db.Exec("INSERT INTO...").Error; err != nil {
		return err
	}

	logging.Info().Msg("Complex migration completed")
	return nil
}
```

### 5. Don't Modify Existing Migrations

Once a migration is deployed to production:
- ❌ **Never modify it** - This breaks version tracking
- ✅ **Create a new migration** instead

### 6. Test Rollback

Always implement and test the `Down()` method:

```go
func (m *MyMigration) Down(db *gorm.DB) error {
	// Reverse what Up() did
	logging.Info().Msg("Rolling back my_migration")

	err := db.Exec("DROP TABLE new_table").Error
	if err != nil {
		logging.Warn().Err(err).Msg("Failed to drop table (may not exist)")
		return nil // Don't fail rollback if table doesn't exist
	}

	return nil
}
```

## Migration Examples

### Adding a Column

```go
func (m *Migration) Up(db *gorm.DB) error {
	return db.Exec("ALTER TABLE pools ADD COLUMN owner_id INTEGER").Error
}

func (m *Migration) Down(db *gorm.DB) error {
	return db.Exec("ALTER TABLE pools DROP COLUMN owner_id").Error
}
```

### Renaming a Column

```go
func (m *Migration) Up(db *gorm.DB) error {
	// SQLite: Create new column, copy data, drop old column
	// MariaDB: Use CHANGE COLUMN
	return db.Exec(`
		ALTER TABLE pools ADD COLUMN volume_gallons_new DECIMAL(10,2);
		UPDATE pools SET volume_gallons_new = volume_gallons;
		ALTER TABLE pools DROP COLUMN volume_gallons;
		ALTER TABLE pools RENAME COLUMN volume_gallons_new TO volume_gallons;
	`).Error
}
```

### Adding an Index

```go
func (m *Migration) Up(db *gorm.DB) error {
	return db.Exec("CREATE INDEX idx_samples_pool_date ON samples(pool_id, sample_date_time)").Error
}

func (m *Migration) Down(db *gorm.DB) error {
	return db.Exec("DROP INDEX idx_samples_pool_date").Error
}
```

### Data Transformation

```go
func (m *Migration) Up(db *gorm.DB) error {
	// Add new column
	if err := db.Exec("ALTER TABLE measurements ADD COLUMN temp_celsius DECIMAL(5,2)").Error; err != nil {
		return err
	}

	// Transform existing data: Fahrenheit to Celsius
	if err := db.Exec("UPDATE measurements SET temp_celsius = (temperature - 32) * 5.0/9.0").Error; err != nil {
		return err
	}

	return nil
}
```

## Troubleshooting

### Migration Fails Midway

The transaction will automatically rollback. The migration won't be recorded as applied, so it will retry next startup.

### Need to Force a Migration

If you need to mark a migration as applied without running it:

```sql
INSERT INTO schema_migrations (version, name, applied_at)
VALUES ('003', 'my_migration', CURRENT_TIMESTAMP);
```

**Use with caution!** Only do this if you've manually applied the changes.

### Migration Applied but Database Doesn't Match

Check if the migration actually succeeded:
```bash
# Check database schema manually
sqlite3 waterlogger.db ".schema"

# Or for MariaDB
mysql -u user -p database -e "SHOW CREATE TABLE tablename"
```

### Rollback Fails

If rollback fails, you may need to:
1. Fix the `Down()` method
2. Manually revert the schema changes
3. Delete the migration record from `schema_migrations`

## Schema Migrations Table

The system tracks migrations in the `schema_migrations` table:

```sql
CREATE TABLE schema_migrations (
    id          INTEGER PRIMARY KEY,
    version     VARCHAR(255) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    applied_at  DATETIME NOT NULL
);
```

**Do not manually modify this table** unless troubleshooting a failed migration.

## Continuous Integration

In CI/CD pipelines:

```bash
# Run pending migrations
./waterlogger -config config.yaml -migration-status

# Application startup will auto-apply pending migrations
./waterlogger -config config.yaml
```

## Production Deployment

1. **Backup database** before deploying
2. **Review pending migrations**: `./waterlogger -migration-status`
3. **Test migrations** on a copy of production data
4. **Deploy new version** (migrations run automatically on startup)
5. **Verify migrations applied**: Check logs and migration status
6. **Monitor application** for issues

## Getting Help

- Check logs in `logs/` directory for detailed migration output
- Review `internal/database/migrations/` for existing migration examples
- See migration status: `./waterlogger -migration-status`
- Review this documentation

## See Also

- [Database Migration (Data)](/mnt/c/Users/br8kw/Github/waterlogger/internal/database/migration.go) - For migrating data between SQLite and MariaDB
- [VERSION_MANAGEMENT.md](/mnt/c/Users/br8kw/Github/waterlogger/VERSION_MANAGEMENT.md) - Application versioning
- [DEPLOYMENT.md](/mnt/c/Users/br8kw/Github/waterlogger/docs/DEPLOYMENT.md) - Deployment procedures
