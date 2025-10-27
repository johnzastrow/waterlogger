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