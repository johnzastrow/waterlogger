# App Development Discussion Notes

---

## ⚠️ IMPORTANT: VERSION MANAGEMENT

**Before making ANY code changes, read:** [`VERSION_MANAGEMENT.md`](VERSION_MANAGEMENT.md)

**Current Version:** 1.4.0

**Repository:** https://github.com/johnzastrow/waterlogger

**Key Rule:** ALWAYS increment version when making code changes. See VERSION_MANAGEMENT.md for complete instructions.

---

## Initial Discussion

Based on requirements.md analysis - building a web application called "Waterlogger" for tracking pool and hot tub water parameters with calculations, charting, and export capabilities.

## App Concept

**Waterlogger** - A comprehensive water quality management system for pools and hot tubs that:
- Tracks water chemistry parameters over time
- Calculates water balance indices (LSI, RSI)
- Provides data visualization through line charts
- Supports multi-user environments
- Exports data to Excel and Markdown formats
- Runs as a single executable on Ubuntu Linux

## Technical Requirements

- **Platform**: Ubuntu Linux AND Windows, single executable
- **Database**: SQLite OR MariaDB (user configurable)
- **Authentication**: Basic user auth with single role + setup wizard
- **Web Interface**: Modern responsive design for mobile/desktop
- **Configuration**: File-based config including port settings (default :2342)
- **Deployment**: Single executable + service configuration
- **Data Migration**: Version tracking with database migrations
- **Units**: User-selectable display units (Imperial/Metric) with proper conversions

## Database Schema

### USERS Table
- id (PRIMARY KEY)
- username (UNIQUE, NOT NULL)
- email (UNIQUE, NOT NULL) 
- password (NOT NULL)
- created_at, updated_at (NOT NULL)
- created_by, updated_by (NOT NULL)

### POOLS Table
- id (PRIMARY KEY)
- name (UNIQUE, NOT NULL)
- volume_gallons (DECIMAL) - Calculated or manually entered volume
- type (ENUM: 'pool', 'hot_tub')
- system_description (TEXT)
- shape (ENUM: 'rectangular', 'round', 'oval', 'kidney', 'l_shaped') - Pool shape for volume calculation
- length, width, diameter (DECIMAL) - Pool dimensions in feet
- shallow_depth, deep_depth (DECIMAL) - Pool depths in feet
- has_steps, steps_width, steps_length, steps_depth (BOOLEAN/DECIMAL) - Step configuration
- has_spa, spa_length, spa_width, spa_depth (DECIMAL) - Attached spa configuration
- created_at, updated_at (NOT NULL)
- created_by, updated_by (NOT NULL)

### KITS Table
- id (PRIMARY KEY)
- name (NOT NULL)
- description (TEXT)
- purchased_date (DATE)
- replenished_date (DATE)
- created_at, updated_at (NOT NULL)
- created_by, updated_by (NOT NULL)

### SAMPLES Table
- id (PRIMARY KEY)
- pool_id (FOREIGN KEY to POOLS, NOT NULL)
- sample_datetime (DATETIME, NOT NULL)
- user_id (FOREIGN KEY to USERS, NOT NULL)
- kit_id (FOREIGN KEY to KITS, NOT NULL)
- created_at, updated_at (NOT NULL)
- created_by, updated_by (NOT NULL)

### MEASUREMENTS Table
- id (PRIMARY KEY)
- sample_id (FOREIGN KEY to SAMPLES, NOT NULL)
- fc (DECIMAL, NOT NULL) - Free Chlorine (ppm)
- tc (DECIMAL, NOT NULL) - Total Chlorine (ppm)
- ph (DECIMAL, NOT NULL) - pH (0-14 scale)
- ta (DECIMAL, NOT NULL) - Total Alkalinity (ppm)
- ch (DECIMAL, NOT NULL) - Calcium Hardness (ppm)
- cya (DECIMAL) - Cyanuric Acid (ppm)
- temperature (DECIMAL, NOT NULL) - Temperature (°F)
- salinity (DECIMAL) - Salinity (ppm)
- tds (DECIMAL) - Total Dissolved Solids (mg/l)
- appearance (TEXT) - Water appearance notes
- maintenance (TEXT) - Maintenance notes
- created_at, updated_at (NOT NULL)
- created_by, updated_by (NOT NULL)

### INDICES Table
- id (PRIMARY KEY)
- sample_id (FOREIGN KEY to SAMPLES, NOT NULL)
- lsi (DECIMAL) - Langelier Saturation Index
- rsi (DECIMAL) - Ryznar Stability Index
- comment (TEXT) - Notes about estimation/missing parameters
- created_at, updated_at (NOT NULL)
- created_by, updated_by (NOT NULL)

### USER_PREFERENCES Table
- id (PRIMARY KEY)
- user_id (FOREIGN KEY to USERS, NOT NULL)
- unit_system (ENUM: 'imperial', 'metric', NOT NULL, DEFAULT 'imperial')
- created_at, updated_at (NOT NULL)

### ADJUSTMENTS Table
- id (PRIMARY KEY)
- pool_id (FOREIGN KEY to POOLS, NOT NULL)
- sample_id (FOREIGN KEY to SAMPLES) - Optional reference to originating sample
- starting_fc, starting_ph, starting_ta, starting_ch, starting_temp (DECIMAL, NOT NULL) - Starting water conditions
- starting_cya, starting_salt, starting_tds (DECIMAL) - Optional starting parameters
- target_fc, target_ph, target_ta, target_ch, target_temp (DECIMAL, NOT NULL) - Target water conditions
- target_cya, target_salt, target_tds (DECIMAL) - Optional target parameters
- starting_lsi, starting_rsi, target_lsi, target_rsi (DECIMAL) - Water balance indices
- add_muriatic_acid, add_sodium_bisulfate, add_soda_ash, add_borax (DECIMAL) - Chemical additions in fluid ounces
- add_sodium_bicarbonate, add_calcium_chloride, add_bleach, add_trichlor (DECIMAL) - Chemical additions in fluid ounces
- add_dichlor, add_cal_hypo, add_salt (DECIMAL) - Chemical additions in fluid ounces/pounds
- notes (TEXT) - User notes about the adjustment
- created_at, updated_at (NOT NULL)
- created_by, updated_by (NOT NULL)

## Features Required

### Core CRUD Operations
- Users: Create, Read, Update, Delete user accounts
- Pools: Manage multiple pools with specifications and volume calculations
- Test Kits: Track testing equipment and supplies
- Samples: Log water testing sessions
- Measurements: Record all water chemistry parameters
- Indices: Auto-calculate and display water balance indices
- Adjustments: Chemical dosing recommendations and application tracking

### Data Visualization
- Line charts showing parameter trends over time
- Multi-parameter overlays for correlation analysis
- Date range filtering and zoom capabilities

### Export Functionality
- **Excel Export**: Separate worksheets for each entity (users, pools, kits, samples, measurements, indices, adjustments)
- **Markdown Reports**: Structured reports with data tables sorted by date, including comprehensive adjustment details
- **JSON Backup**: Complete database backup with all tables and relationships for data migration
- **PDF Generation**: Browser-based PDF export for adjustment details with professional formatting

### Advanced Features (Version 1.2+)
- **Pool Volume Calculator**: Comprehensive volume calculation system supporting multiple pool shapes (rectangular, round, oval, kidney, L-shaped) with varying depths, steps, and attached spas
- **Chemical Adjustment System**: Professional water balance calculations with LSI/RSI indices and precise chemical dosing recommendations for 11 different pool chemicals
- **Adjustment History**: Complete tracking of chemical adjustments with before/after conditions, chemical additions, and user notes
- **Water Balance Analysis**: Real-time LSI/RSI calculations with color-coded indicators for optimal water balance
- **Dashboard Analytics**: Quick overview of recent samples, water quality status, and recent adjustments across all pools

### User Interface Requirements
- Dark navy navigation background
- Full field names with units displayed
- Hover tooltips with detailed descriptions
- Responsive design for mobile and desktop
- Modern, clean appearance
- Professional favicon integration
- Color-coded water quality indicators
- Mobile-optimized forms and displays

## Implementation Plan

### Phase 1: Foundation & Setup
1. **Technology Stack Selection**
   - Backend: Go with Gin framework (single executable requirement)
   - Database: SQLite OR MariaDB with GORM (configurable)
   - Frontend: HTML/CSS/JavaScript with Alpine.js for reactivity
   - Charts: Chart.js for data visualization
   - Build: Cross-platform binary compilation (Linux + Windows)

2. **Project Structure Setup**
   ```
   waterlogger/
   ├── cmd/waterlogger/main.go
   ├── internal/
   │   ├── config/
   │   ├── database/
   │   ├── handlers/
   │   ├── models/
   │   ├── services/
   │   └── middleware/
   ├── web/
   │   ├── static/
   │   └── templates/
   ├── migrations/
   └── config.yaml
   ```

3. **Database Setup**
   - Implement database models with GORM
   - Database abstraction layer (SQLite/MariaDB)
   - Create migration system
   - Set up connection pooling and transactions

### Phase 2: Core Backend Development
1. **Authentication System**
   - Setup wizard for first user creation
   - User registration/login with password complexity requirements
   - Session management
   - Password hashing (bcrypt)

2. **API Endpoints**
   - RESTful APIs for all entities
   - Input validation and sanitization
   - Error handling middleware
   - Command-line utilities (password reset, database migration, data export/import)

3. **Water Chemistry Calculations**
   - Implement LSI and RSI calculation functions
   - Unit conversion utilities (Imperial/Metric with user preferences)
   - Auto-calculation triggers on measurement creation
   - Mid-range defaults for missing parameters with comment tracking

### Phase 3: Frontend Development
1. **UI Components**
   - Responsive layout with navigation
   - Data grids with CRUD operations
   - Form validation with required field indicators (red asterisks)
   - Unit system selection in user preferences

2. **Data Visualization**
   - Chart.js integration for line charts (exclude TDS, CYA, SAL by default)
   - Interactive date range selection (default: last 30 days)
   - Parameter selection and filtering

3. **User Experience**
   - Tooltips with field descriptions
   - Loading states and error messages
   - Mobile-friendly responsive design

### Phase 4: Advanced Features
1. **Export Functionality**
   - Excel file generation with multiple worksheets (all pools, all data)
   - Markdown report generation with formatting
   - File naming: WL[timestamp] format
   - Download management

2. **Configuration & Deployment**
   - Config file parsing (YAML) with database type selection
   - Database migration utility (SQLite ↔ MariaDB)
   - Cross-platform executable compilation (Linux + Windows)
   - Service file creation (systemd for Linux, service for Windows)
   - Installation documentation

## Architecture Decisions (from requirements.md)

### Data Management
- **Database Options**: SQLite OR MariaDB (user configurable)
- **GORM ORM**: Provides type safety and migration support with database abstraction
- **Audit Trail**: All tables include created/updated timestamps and user tracking
- **Data Integrity**: Foreign key constraints and validation rules
- **Unit System**: User preferences for Imperial/Metric display with proper conversions

### UI Layout & Design
- **Server-Side Rendering**: HTML templates with minimal JavaScript
- **Alpine.js**: Lightweight reactivity for dynamic UI elements
- **CSS Grid/Flexbox**: Modern responsive layouts
- **Component-Based**: Reusable template components

### Export Functionality
- **Excel**: Using excelize library for multi-worksheet generation
- **Markdown**: Template-based generation with table formatting
- **Streaming**: Large dataset handling without memory issues

### Error Handling & UX
- **Graceful Degradation**: Progressive enhancement approach
- **Input Validation**: Client and server-side validation
- **User Feedback**: Clear error messages and success indicators
- **Logging**: Structured logging for debugging and monitoring

### Packaging & Distribution
- **Cross-Platform Binary**: Go's static compilation for Linux + Windows deployment
- **Embedded Assets**: Static files embedded in executable
- **Configuration**: External YAML file for runtime settings (database type, port :2342)
- **Service Integration**: Systemd (Linux) and Windows Service support
- **Setup Wizard**: Comprehensive first-run wizard (admin user + database + config)
- **Database Migration**: Bidirectional utility for SQLite ↔ MariaDB data transfer

## Additional Implementation Details to Consider

### Code Structure & Organization
- Clean Architecture with domain separation
- Dependency injection for testability
- Interface-based design for modularity
- Comprehensive error handling

### Database Details
- Connection pooling optimization
- Transaction management for data consistency
- Index optimization for query performance
- Backup and recovery procedures

### UI Component Specifics
- Form builders for consistent CRUD interfaces
- Data grid components with sorting/filtering
- Chart configuration and customization
- Mobile-first responsive breakpoints

### Export File Details
- **Excel Format**: 
  - Separate sheets: Users, Pools, Kits, Samples, Measurements, Indices
  - Data formatting and column headers
  - Date formatting consistency
- **Markdown Format**:
  - Hierarchical structure with headings
  - Data tables with proper alignment
  - Chronological sorting (oldest first)

### Error Scenarios
- Database connection failures
- Invalid input data handling
- Export generation errors
- Chart rendering failures
- Session timeout management

### User Experience
- Intuitive navigation flow
- Consistent data entry patterns
- Visual feedback for actions
- Help documentation integration
- Performance optimization for large datasets

## Development Notes

### Environment Setup ✓
- Go development environment
- SQLite database tools
- Frontend build pipeline
- Testing framework setup

### Implementation Complete ✓

#### Version 1.0 Core Features
- ✅ User authentication and management system
- ✅ Pool management with CRUD operations
- ✅ Test kit tracking and management
- ✅ Water sample recording and tracking
- ✅ Water chemistry measurements
- ✅ LSI/RSI water balance calculations
- ✅ Data visualization with Chart.js
- ✅ Excel export functionality
- ✅ Markdown export functionality
- ✅ SQLite database with GORM ORM
- ✅ Responsive web interface
- ✅ Build system with timestamp injection

#### Version 1.2 Advanced Features
- ✅ Pool volume calculator with multiple shapes
- ✅ Chemical adjustment recommendation system
- ✅ Water balance analysis (LSI/RSI) with color indicators
- ✅ Adjustment history tracking and display
- ✅ PDF export for adjustment details
- ✅ Enhanced dashboard with recent adjustments
- ✅ Comprehensive backup export (JSON)
- ✅ Professional favicon integration
- ✅ Mobile-optimized responsive design

### Issues Fixed During Development
- **User Management**: Implemented comprehensive user CRUD operations with modal-based interface
- **Build Timestamps**: Added build date/time injection and display across all pages (fixed "Built on unknown at unknown" issue)
- **Navigation**: Updated all template files with user management links
- **Authentication**: Enhanced password validation and security measures
- **Responsive Design**: Improved mobile compatibility across all pages
- **PDF Generation**: Fixed jsPDF library loading issues by switching to browser-based print functionality
- **Build Script**: Resolved Windows line ending issues causing bash execution errors
- **Favicon Integration**: Professional branding implementation across all HTML templates

### Recent Updates (Latest Sessions)

#### Version 1.2 Implementation (Previous Sessions)
- **User Management System**: Complete implementation with create, edit, delete functionality
- **Build Timestamp Integration**: All pages now display build date/time in bottom-right corner
- **CSS Enhancements**: Added styling for user management interface and build timestamp
- **Template Updates**: Updated all 10 HTML templates with user management navigation and build timestamps
- **Build Process**: Created automated build script with timestamp injection

#### Version 1.2 Features
- **Pool Volume Calculator**: Comprehensive volume calculation system with support for rectangular, round, oval, kidney, and L-shaped pools with varying depths
- **Chemical Adjustment System**: Professional-grade water balance calculations with LSI/RSI indices and chemical dosing recommendations
- **PDF Export for Adjustments**: Browser-based PDF generation with comprehensive adjustment details, chemical safety guidelines, and water balance explanations
- **Favicon Integration**: Professional branding with favicon.ico across all pages
- **Adjustment History Tracking**: Past adjustments display on both adjustments screen (pool-filtered) and dashboard (last 10 across all pools)
- **Enhanced Markdown Export**: Complete adjustments section with starting/target conditions, water balance indices, and chemical recommendations
- **Full Backup Export Update**: All new data tables included in JSON backup exports (user preferences, measurements, indices, adjustments)

#### Version 1.3 Features
- **Comprehensive Logging System**: Structured logging with zerolog for high-performance JSON and console output
- **Log Rotation**: Automatic log rotation with configurable size, backups, and compression
- **Request Tracing**: Unique request IDs for complete request tracing across all operations
- **Security Audit Logging**: Comprehensive audit trail for sensitive operations
- **Custom GORM Logger**: Database operation logging with error categorization and slow query detection
- **Version Display Fix**: Corrected version display in UI footer
- **Documentation Updates**: Added LOGGING.md with comprehensive logging documentation
- **Build System Enhancements**: MSYS2 build instructions for Windows

#### Version 1.5 Features (Current Session)
- **Web-Based Backup Import**: Complete backup restore functionality from Settings page
  - File upload interface for JSON backups
  - Validation of file type (.json only)
  - Real-time import progress with loading states
  - Success/error messages with detailed feedback
  - Automatic refresh of settings and migrations after import
  - **Backwards compatibility**: Imports data from older versions gracefully
    - Auto-updates database schema if needed
    - Foreign key validation and automatic fixing
    - Skips missing or incompatible fields without failing
    - Handles partial imports when some data is unavailable
  - **New API endpoint**: `POST /api/settings/import` - Import database from uploaded backup
  - **CSS styling** for message boxes (success, warning, info boxes)
  - Warning about data being added (not replaced)
  - Automatic cleanup of uploaded temporary files

#### Version 1.4 Features
- **Schema Migration System**: Version-tracked database schema changes with automatic migration on startup
  - `SchemaMigration` model tracks version, name, and timestamp
  - Migration interface with Up() and Down() methods
  - Transaction-wrapped migrations for atomicity
  - Command-line tools: `-migration-status` and `-migration-rollback`
  - MIGRATIONS.md documentation with comprehensive examples
- **Enhanced Settings UI**: Completely redesigned settings page with:
  - **System Information section**: App version, build date/time, database type, schema version, migrations count, server address
  - **Database Schema Migrations section**: Complete migration history table with refresh capability
  - **Database Management section**: Current database details, one-click backup creation, migration instructions
- **New API Endpoints**:
  - `GET /api/settings/migrations` - Retrieve migration history
  - `POST /api/settings/backup` - Create database backup on demand
  - Enhanced `GET /api/settings` - Returns comprehensive system information
- **Database Backup System**: One-click JSON backups from web UI with auto-generated timestamped filenames (WL<timestamp>.json)
- **Migration Testing**: Comprehensive testing of schema migrations on both SQLite and MariaDB
- **Foreign Key Validation**: Enhanced data migration code to validate and fix all foreign key relationships
- **Automated Deployment Scripts**:
  - **Linux**: `deploy-linux.sh` - One-command deployment with systemd service creation
    - Creates directory structure (`/opt/waterlogger`, `logs/`, `backups/`)
    - Creates dedicated user with proper permissions
    - Optional production configuration (json logging, file output)
    - Generates systemd service file with security hardening
    - Auto-enables service to start on boot
  - **Windows**: `deploy-windows.bat` - One-command deployment with Windows service creation
    - Creates directory structure (`C:\Program Files\Waterlogger`, `logs/`, `backups/`)
    - Optional production configuration
    - Creates Windows service with automatic restart on failure
    - Includes service management commands
  - **Enhanced service files** with security hardening, proper restart policies, and complete documentation
  - **Simplified deployment**: Single command goes from binary to production-ready service

## Notes

### Water Chemistry Calculation Details
- **Unit Conversions**: Dynamic based on user preference (Imperial ↔ Metric)
- **Required Parameters**: Temperature, pH, TDS, Calcium Hardness, Total Alkalinity
- **Missing Parameter Handling**: Use mid-range defaults with comment tracking
- **Formulas**: LSI = pH - pHs, RSI = 2×pHs - pH
- **Implementation**: Direct port of Python functions with proper error handling

### User Clarifications Implemented

#### Initial Requirements
- **Setup**: Setup wizard for first admin user
- **Security**: Modern password complexity requirements with display
- **Calculations**: Mid-range defaults for missing parameters with comment field
- **Units**: User-selectable Imperial/Metric with proper conversion handling
- **Charts**: Default exclude TDS, CYA, SAL; 30-day default range
- **Exports**: All data, all pools; WL[timestamp] naming convention
- **Platform**: Cross-platform support (Linux + Windows)
- **Database**: Configurable SQLite OR MariaDB
- **Port**: Default :2342
- **UI**: Required fields marked with red asterisks

#### Final Clarifications
- **Database Config**: MariaDB connection details (host, port, username, password, database name)
- **Database Migration**: Bidirectional migration utility between SQLite ↔ MariaDB
- **Mid-Range Defaults**: TDS=300mg/L, Calcium Hardness=250ppm, Total Alkalinity=100ppm
- **Unit Display**: Show both units (e.g., "75°F (24°C)")
- **Chart Parameters**: TDS, CYA, SAL permanently excluded from charts
- **Export Format**: Filename WL20240714_143022.xlsx format
- **Markdown Export**: Include calculated indices (LSI/RSI) as separate section
- **Setup Wizard**: Configure database type, connection details, and all configuration options
- **Password Security**: ~~Basic complexity constraints (length, special characters, etc.)~~ **REMOVED** - Now accepts any non-empty password for simplicity
- **Port Configuration**: Changed from :2341 to :2342 as the default port
- **Password Reset**: Command-line utility for resetting user passwords (`./waterlogger -reset-password username`)
