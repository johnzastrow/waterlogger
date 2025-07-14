# App Development Discussion Notes

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
- **Configuration**: File-based config including port settings (default :2341)
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
- volume_gallons (DECIMAL)
- type (ENUM: 'pool', 'hot_tub')
- system_description (TEXT)
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

## Features Required

### Core CRUD Operations
- Users: Create, Read, Update, Delete user accounts
- Pools: Manage multiple pools with specifications
- Test Kits: Track testing equipment and supplies
- Samples: Log water testing sessions
- Measurements: Record all water chemistry parameters
- Indices: Auto-calculate and display water balance indices

### Data Visualization
- Line charts showing parameter trends over time
- Multi-parameter overlays for correlation analysis
- Date range filtering and zoom capabilities

### Export Functionality
- **Excel Export**: Separate worksheets for each entity
- **Markdown Reports**: Structured reports with data tables sorted by date

### User Interface Requirements
- Dark navy navigation background
- Full field names with units displayed
- Hover tooltips with detailed descriptions
- Responsive design for mobile and desktop
- Modern, clean appearance

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
- **Configuration**: External YAML file for runtime settings (database type, port :2341)
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
- [To be updated during development]

### Issues Fixed During Development
- [To be documented as encountered]

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
- **Port**: Default :2341
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
- **Password Security**: Basic complexity constraints (length, special characters, etc.)
