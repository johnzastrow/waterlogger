# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Sample management interface with measurement input forms
- Chart.js integration for data visualization
- Excel and Markdown export functionality
- Database migration utility for switching between SQLite and MariaDB
- Unit conversion system with dual-unit display
- Enhanced mobile responsiveness

### Changed
- TBD

### Fixed
- TBD

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