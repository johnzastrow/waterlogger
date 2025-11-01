# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- TBD

### Changed
- TBD

### Fixed
- TBD

## [1.5.0] - 2025-10-27

### Added - Database Backup Import & Automated Deployment

#### Database Backup Import Feature
- **Web-based backup import** from Settings page
  - File upload interface for JSON backup files
  - Validation of file type (.json only)
  - Real-time import progress with loading states
  - Success/error messages with detailed feedback
  - Automatic refresh of settings and migrations after import
- **Backwards compatibility handling**
  - Imports data from older system versions gracefully
  - Auto-updates database schema if needed
  - Foreign key validation and automatic fixing
  - Skips missing or incompatible fields without failing
  - Handles partial imports when some data is unavailable
- **User experience enhancements**
  - Warning about data being added (not replaced)
  - Clear file selection feedback
  - Automatic cleanup of uploaded files
  - 10-second auto-dismiss of success messages
- **New API endpoint**: `POST /api/settings/import` - Import database from uploaded backup file
- **CSS styling** for message boxes (success, warning, info)

#### Automated Deployment Scripts
- **Linux deployment script** (`deploy-linux.sh`) for one-command systemd service installation
  - Creates installation directory structure (`/opt/waterlogger`, `logs/`, `backups/`)
  - Creates dedicated `waterlogger` user with proper permissions
  - Optional production configuration (json logging, file output)
  - Generates systemd service file with security hardening
  - Auto-enables service to start on boot
- **Windows deployment script** (`deploy-windows.bat`) for one-command Windows service installation
  - Creates installation directory structure (`C:\Program Files\Waterlogger`, `logs/`, `backups/`)
  - Optional production configuration
  - Creates Windows service with automatic restart on failure
  - Includes service management commands
- **Enhanced systemd service file** with:
  - Security hardening (NoNewPrivileges, ProtectSystem, ProtectHome, PrivateTmp)
  - Proper working directory and read-write paths
  - Journald integration with SyslogIdentifier
  - Restart policy with 10-second delay
  - Documentation link
- **Comprehensive deployment documentation** in README.md
  - Quick deployment (automated) and manual deployment options
  - Production configuration recommendations
  - Log viewing instructions
  - Service management commands

### Changed
- **Settings page** redesigned with "Backup Database" and "Restore from Backup" sections
- **Database Management section** now includes both backup creation and import capabilities
- **Deployment process** significantly simplified with automated scripts

### Technical Details
- Import uses existing `database.ImportData()` function with enhanced error handling
- Temporary file handling with automatic cleanup
- File upload via multipart/form-data
- Frontend uses FormData API for file upload
- Imports data in correct order respecting foreign key constraints

## [1.4.0] - 2025-10-27

### Added - Database Schema Migration System & Enhanced Settings UI

#### Schema Migration System
- **Version-tracked schema migrations** for managing database structure changes
- **SchemaMigration model** to track applied migrations with version, name, and timestamp
- **Migration interface** with Up() and Down() methods for applying and reverting changes
- **Automatic migration execution** on application startup
- **Migration registry** for centralized migration management
- **Transaction-wrapped migrations** for atomicity (all-or-nothing execution)
- **Command-line migration tools**:
  - `-migration-status` - View all applied migrations
  - `-migration-rollback` - Revert the last migration
- **Cross-database support** - Works with both SQLite and MariaDB
- **Migration documentation** (MIGRATIONS.md) with comprehensive examples and best practices

#### Enhanced Settings UI
- **System Information section** displaying:
  - Application version
  - Build date and time
  - Database type (SQLite/MariaDB)
  - Current schema version
  - Total migrations applied
  - Server host and port
- **Database Schema Migrations section** with:
  - Complete migration history table
  - Refresh capability for latest status
  - Applied date/time for each migration
  - Command-line instructions for advanced operations
- **Database Management section** with:
  - Current database connection details (sanitized)
  - One-click database backup creation
  - Auto-generated timestamped backup files (WL<timestamp>.json)
  - Database migration instructions (SQLite ↔ MariaDB)
  - Safety warnings and documentation links
- **User-friendly interface** with loading states, success messages, and error handling

#### API Endpoints
- `GET /api/settings/migrations` - Retrieve migration history
- `POST /api/settings/backup` - Create database backup on demand
- Enhanced `GET /api/settings` - Returns comprehensive system information

#### Automated Deployment Scripts
- **Linux deployment script** (`deploy-linux.sh`) for automated systemd service installation
  - Creates installation directory structure (`/opt/waterlogger`, `logs/`, `backups/`)
  - Creates dedicated `waterlogger` user with proper permissions
  - Optional production configuration (json logging, file output)
  - Generates systemd service file with security hardening
  - Auto-enables service to start on boot
- **Windows deployment script** (`deploy-windows.bat`) for automated Windows service installation
  - Creates installation directory structure (`C:\Program Files\Waterlogger`, `logs/`, `backups/`)
  - Optional production configuration
  - Creates Windows service with automatic restart on failure
  - Includes service management commands
- **Enhanced systemd service file** with:
  - Security hardening (NoNewPrivileges, ProtectSystem, ProtectHome, PrivateTmp)
  - Proper working directory and read-write paths
  - Journald integration with SyslogIdentifier
  - Restart policy with 10-second delay
  - Documentation link
- **Comprehensive deployment documentation** in README.md
  - Quick deployment (automated) and manual deployment options
  - Production configuration recommendations
  - Log viewing instructions
  - Service management commands

### Changed
- **GetSettings API** now includes schema version, migration count, and database details
- **Settings page** completely redesigned with organized sections
- **Database initialization** now uses migration system instead of simple AutoMigrate
- **Foreign key validation** enhanced in data migration to handle all relationships

### Technical Details
- Migration system uses GORM for database operations
- Migrations tracked in `schema_migrations` table
- Backups stored in dedicated `backups/` directory
- All migration operations logged with structured logging
- Migration runner supports pending migration detection
- Rollback validation prevents reverting other developers' changes

## [1.2.0] - 2025-10-25

### Added - Pool Volume Calculator, Chemical Adjustments, and PDF Export

#### Pool Volume Calculator
- **Comprehensive volume calculation system** supporting multiple pool shapes:
  - Rectangular pools with variable depth
  - Round/circular pools
  - Oval pools
  - Kidney-shaped pools
  - L-shaped pools with dual sections
- **Advanced features:**
  - Steps configuration (width, length, depth)
  - Attached spa calculations
  - Dynamic volume updates
  - Visual pool shape preview

#### Chemical Adjustment System
- **Professional-grade water balance calculations** with:
  - Starting and target water conditions
  - LSI (Langelier Saturation Index) calculations
  - RSI (Ryznar Stability Index) calculations
  - Water balance recommendations
- **11 different pool chemicals supported:**
  - Muriatic Acid
  - Sodium Bisulfate
  - Soda Ash
  - Borax
  - Sodium Bicarbonate
  - Calcium Chloride
  - Bleach
  - Trichlor
  - Dichlor
  - Cal-Hypo
  - Salt
- **Chemical dosing calculations** based on:
  - Pool volume
  - Starting parameters
  - Target parameters
  - Water balance indices

#### Adjustment History Tracking
- **Complete adjustment records** including:
  - Starting water conditions
  - Target water conditions
  - Chemical additions (in fluid ounces or pounds)
  - Water balance before/after
  - User notes
  - Timestamp tracking
- **Dashboard display** showing last 10 adjustments across all pools
- **Pool-filtered history** on adjustments management page
- **Search and filtering** capabilities

#### PDF Export
- **Browser-based PDF generation** for adjustment details
- **Professional report formatting** with:
  - Adjustment summary
  - Chemical additions list
  - Water balance explanations
  - Safety guidelines for chemicals
  - Printable format
  - Print-friendly styling

#### Enhanced Features
- **Dashboard Analytics** with recent adjustments widget
- **Water Balance Analysis** with color-coded LSI/RSI indicators
- **JSON Backup Export** with all tables and relationships
- **Professional Favicon** integration across all pages
- **Enhanced Markdown Export** including adjustment sections
- **Improvement fixes** for build timestamp display and PDF generation

### Changed
- Updated dashboard to include adjustment history
- Enhanced markdown export structure with adjustments section
- Improved water balance visualization with indicators

### Technical Details
- Adjustments model with comprehensive field support
- LSI/RSI calculation functions
- Chemical dosing algorithms
- JSON backup includes all new data types
- PDF generation using browser print API

## [1.3.0] - 2025-10-26

### Added - Comprehensive Logging System
- **Structured logging** with zerolog for high-performance JSON and console output
- **Log rotation** with lumberjack (configurable max size, backups, age, compression)
- **Multiple log levels**: debug, info, warn, error, fatal
- **Multiple output destinations**: stdout, file, or both simultaneously
- **Request ID middleware** for complete request tracing across all operations
- **Security audit logging** for sensitive operations (login, user management, data exports, etc.)
- **Custom GORM logger** with comprehensive database error tracking:
  - Duplicate key errors detection
  - Foreign key constraint violations
  - NOT NULL, UNIQUE, and CHECK constraint violations
  - Record not found errors
  - Slow query detection (threshold: 200ms)
  - Query execution time tracking
  - SQL query logging with sensitive data sanitization
- **HTTP request logging** with method, path, status, latency, user agent, and client IP
- **Configuration section** in config.yaml for logging settings
- **Logging documentation** (LOGGING.md) with usage examples and best practices

### Changed
- Replaced standard library `log` package with structured zerolog
- Updated all log statements throughout codebase to use new structured logger
- Changed Gin middleware to custom logging middleware for better control
- Database initialization now uses custom GORM logger for detailed error tracking

### Removed
- Debug `fmt.Printf` statements from handlers (replaced with proper debug logging)
- Standard library log imports (replaced with zerolog)

### Technical Details
- Zero-allocation logging for high performance
- Automatic log file compression for rotated logs
- Request tracing with unique request IDs
- Component-based logging for better organization
- Error categorization for database operations
- Audit trail for compliance and security monitoring

## [1.0.0] - 2024-07-14

### Added
- Initial release of Waterlogger
- Multi-user authentication system with secure password hashing
- Pool and hot tub management with full CRUD operations
- Water chemistry parameter tracking (FC, TC, pH, TA, CH, CYA, Temperature, Salinity, TDS)
- Automatic LSI (Langelier Saturation Index) and RSI (Ryznar Stability Index) calculations
- Setup wizard for first-run configuration
- Database support for SQLite and MariaDB
- Responsive web interface with Alpine.js integration
- REST API endpoints for all data operations
- Configuration management with YAML files
- Audit trail for all database operations
- Cross-platform support (Windows and Linux)
- Command-line interface with help and version commands

### Technical Features
- GORM-based database abstraction layer
- Gin web framework for HTTP handling
- bcrypt password hashing for security
- Session-based authentication
- Auto-migration of database schemas
- Comprehensive input validation
- Error handling and logging
- Static file serving for CSS and JavaScript
- HTML template rendering with Go templates

### Water Chemistry
- Free Chlorine (FC) measurement and tracking
- Total Chlorine (TC) measurement and tracking
- pH level monitoring with ideal range indicators
- Total Alkalinity (TA) measurement
- Calcium Hardness (CH) tracking
- Cyanuric Acid (CYA) optional measurement
- Temperature recording in Fahrenheit
- Salinity tracking for saltwater pools
- Total Dissolved Solids (TDS) measurement
- Water appearance and maintenance notes
- Automatic calculation of water balance indices
- Mid-range defaults for missing parameters with comment tracking

### User Interface
- Dark navy navigation bar
- Modern, clean design with responsive layout
- Dashboard with recent samples and pool status
- Pool management with card-based layout
- Form validation with required field indicators
- Modal dialogs for data entry
- Loading states and error messages
- Mobile-friendly responsive design
- Hover tooltips for parameter descriptions

### Database Schema
- Users table with username, email, and password
- UserPreferences for unit system selection
- Pools table with name, volume, type, and system description
- Kits table for test equipment tracking
- Samples table linking pools, users, and kits
- Measurements table for all water chemistry parameters
- Indices table for calculated LSI and RSI values
- Audit fields (created_at, updated_at, created_by, updated_by) on all tables

### Configuration
- Default port 2341 for web server
- Configurable database type (SQLite or MariaDB)
- YAML-based configuration file
- Environment-specific settings
- Database connection parameters
- Application metadata (name, version, secret key)

### Security Features
- Password complexity requirements
- Secure session management
- SQL injection prevention through ORM
- XSS protection in templates
- CSRF protection considerations
- Secure cookie handling

[Unreleased]: https://github.com/your-org/waterlogger/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/your-org/waterlogger/releases/tag/v1.0.0