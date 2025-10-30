# Waterlogger TODO List

**Last Updated:** 2025-10-27
**Version:** 1.5.0
**Status:** Core functionality complete, logging system complete, schema migrations complete, backup import complete, automated deployment complete, production features incomplete

---

## ✅ COMPLETED FEATURES

### Core Application (v1.0)
- [x] User authentication with bcrypt password hashing
- [x] Multi-user management (CRUD operations)
- [x] Session-based authentication with HTTP cookies
- [x] SQLite database support with GORM ORM
- [x] MariaDB database support with GORM ORM
- [x] Database abstraction layer (configurable SQLite/MariaDB)
- [x] Database migration utilities (SQLite ↔ MariaDB)
- [x] Configuration management via YAML (`config.yaml`)
- [x] Setup wizard for first-time configuration
- [x] Pool management (CRUD operations)
- [x] Test kit tracking (CRUD operations)
- [x] Water sample recording (CRUD operations)
- [x] Water chemistry measurements with all parameters
- [x] LSI (Langelier Saturation Index) calculations
- [x] RSI (Ryznar Stability Index) calculations
- [x] Mid-range defaults for missing parameters (TDS=300, CH=250, TA=100)
- [x] Unit conversion utilities (Imperial ↔ Metric)
- [x] User preferences for unit system selection
- [x] Responsive web interface (11 HTML templates)
- [x] Alpine.js integration for reactive UI
- [x] Professional favicon integration
- [x] Build timestamp display on all pages
- [x] Cross-platform compilation (Linux, Windows, macOS)
- [x] Command-line password reset utility
- [x] JSON backup export functionality
- [x] JSON backup import/restore functionality

### Advanced Features (v1.2)
- [x] Pool volume calculator (5 shapes supported)
  - Rectangular, Round, Oval, Kidney, Irregular
- [x] Chemical adjustment recommendation system
- [x] Support for 11 chemical types
- [x] Water balance analysis with LSI/RSI
- [x] Adjustment history tracking
- [x] Command-line data export utility
- [x] Command-line data import utility
- [x] Database migration commands (SQLite ↔ MariaDB)
- [x] Audit trail (created_by, updated_by timestamps)

### Logging System (v1.3)
- [x] Structured logging with zerolog
- [x] Log rotation with lumberjack
- [x] Multiple log levels (debug, info, warn, error, fatal)
- [x] Multiple output destinations (stdout, file, both)
- [x] Request ID middleware for request tracing
- [x] Security audit logging
- [x] Custom GORM logger with error categorization
- [x] Slow query detection
- [x] HTTP request logging with comprehensive details
- [x] Logging configuration in config.yaml
- [x] LOGGING.md documentation

### Database Management (v1.4)
- [x] Schema migration system with version tracking
- [x] SchemaMigration model for tracking applied migrations
- [x] Migration interface with Up() and Down() methods
- [x] Automatic migration execution on startup
- [x] Transaction-wrapped migrations for atomicity
- [x] Command-line migration tools (`-migration-status`, `-migration-rollback`)
- [x] Cross-database schema migration support (SQLite & MariaDB)
- [x] MIGRATIONS.md comprehensive documentation
- [x] Enhanced Settings UI with system information
- [x] Database schema migrations history display
- [x] One-click database backup from web UI
- [x] New API endpoints (`/api/settings/migrations`, `/api/settings/backup`)
- [x] Foreign key validation in data migration

### Data Migration Between Database Types (v1.0+)
- [x] Database migration implementation (`internal/database/migration.go`)
- [x] MigrateSQLiteToMariaDB function with schema and data transfer
- [x] MigrateMariaDBToSQLite function with schema and data transfer
- [x] Data type conversion handling between database types
- [x] Comprehensive error handling and logging
- [x] Command-line flags (`-migrate-to-mariadb`, `-migrate-to-sqlite`)
- [x] JSON backup export functionality (ExportData)
- [x] JSON backup import/restore functionality (ImportData)
- [x] Foreign key relationship validation
- [x] Tested with sample databases
- [x] Documentation in README.md and BUILD_REQUIREMENTS.md

### Web-Based Backup Import (v1.5)
- [x] File upload interface in Settings page
- [x] JSON file validation (.json only)
- [x] New API endpoint (`POST /api/settings/import`)
- [x] Multipart form-data file upload handling
- [x] Temporary file management with auto-cleanup
- [x] Real-time import progress with loading states
- [x] Success/error messages with detailed feedback
- [x] Automatic page refresh after import
- [x] Backwards compatibility with older backups
- [x] Auto-schema migration during import
- [x] Foreign key validation and fixing
- [x] Graceful handling of missing fields
- [x] CSS styling for message boxes (success, warning, info)
- [x] User warnings about data being added (not replaced)

### Automated Deployment (v1.4-1.5)
- [x] Linux deployment script (`deploy-linux.sh`)
- [x] Windows deployment script (`deploy-windows.bat`)
- [x] Enhanced systemd service file with security hardening
- [x] Windows service with auto-restart configuration
- [x] Production configuration prompts
- [x] Directory structure creation (`logs/`, `backups/`)
- [x] User/permission management
- [x] Comprehensive deployment documentation

---

## 🔴 HIGH PRIORITY - Critical Gaps

October 27, 2025 Testing:

1. ~~The version tag that appears in the bottom-right corner of the UI should be updated to  the current version shown in the CHANGELOG.md (1.4.0). However, right now it shows "v | Built on 2025-10-27" and it should show "v1.4.0 | Built on 2025-10-27".~~ **RESOLVED** - Updated all 10 HTML templates with correct changelog anchor links (#140---2025-10-27)
2. Finishing checking the documentation to reflect the current build approaches, particularly for Windows with MSYS2 and MinGW-w64 GCC and on linux.
3. ~~Check the README instructions for setting up systemd on Linux. The example service file has issues.~~ **REVIEWED** - Service file is functional but could use improvements:
   - Missing instruction to create `logs/` directory before first run
   - Should recommend `format: json` and `output: file` for production logging config
   - Missing mention of `backups/` directory (auto-created but worth documenting)
   - Could add security hardening directives (ReadWritePaths, NoNewPrivileges, etc.) 


### 1. Excel Export is Broken
**Status:** PARTIAL - Currently exports CSV text, not true .xlsx format
**Location:** `internal/handlers/handlers.go:770-850`
**Issue:** Missing `excelize` library, no multi-worksheet support

**Tasks:**
- [ ] Add `excelize` library to `go.mod`
- [ ] Replace CSV export with true .xlsx format
- [ ] Implement multi-worksheet structure:
  - [ ] Users worksheet
  - [ ] Pools worksheet
  - [ ] Kits worksheet
  - [ ] Samples worksheet
  - [ ] Measurements worksheet
  - [ ] Indices worksheet
  - [ ] Adjustments worksheet
- [ ] Add proper cell formatting and column headers
- [ ] Add date formatting consistency
- [ ] Fix filename to use `WL[timestamp].xlsx` format
- [ ] Test with large datasets (memory optimization)

**Reference:** CLAUDE.md Line 144-149

---

### 2. PDF Export Not Implemented
**Status:** NOT IMPLEMENTED - Documented as complete but no code exists
**Location:** Should be in `internal/handlers/handlers.go`
**Issue:** No PDF generation code found

**Tasks:**
- [ ] Choose PDF generation approach (jsPDF browser-based OR Go library)
- [ ] Implement adjustment report PDF export
  - [ ] Include starting/target water conditions
  - [ ] Include LSI/RSI calculations
  - [ ] Include chemical recommendations with dosages
  - [ ] Include chemical safety guidelines
  - [ ] Include water balance explanations
- [ ] Add PDF export endpoint (`POST /api/adjustments/:id/pdf`)
- [ ] Add PDF download button to adjustments UI
- [ ] Test browser-based PDF generation (print to PDF)
- [ ] Ensure mobile compatibility for PDF generation

**Reference:** CLAUDE.md Line 127, 154

---

### 3. Charts Not Implemented
**Status:** NOT IMPLEMENTED - No Chart.js library, no templates, no backend
**Location:** Should be in `web/templates/samples.html` or `dashboard.html`
**Issue:** Charts feature is completely missing

**Tasks:**
- [ ] Add Chart.js library to `web/static/js/` (CDN or local)
- [ ] Implement chart data API endpoint (`GET /api/charts/data`)
- [ ] Add date range filtering logic (default: last 30 days)
- [ ] Add parameter selection logic
- [ ] Exclude TDS, CYA, SAL by default (per requirements)
- [ ] Implement multi-parameter overlays for correlation
- [ ] Add chart rendering JavaScript in samples/dashboard template
- [ ] Add canvas elements to HTML templates
- [ ] Test chart zoom capabilities
- [ ] Add loading states for chart data
- [ ] Optimize queries for large datasets
- [ ] Add chart export functionality (PNG download)

**Reference:** CLAUDE.md Line 30-33, 141-143

---

### 4. Service Deployment - COMPLETED ✅
**Status:** COMPLETE - Automated deployment scripts created
**Location:** `deploy-linux.sh`, `deploy-windows.bat`, README.md
**What was implemented:**
- ✅ Linux automated deployment script (`deploy-linux.sh`)
  - Creates directory structure (`/opt/waterlogger`, `logs/`, `backups/`)
  - Creates dedicated user
  - Configures for production (optional)
  - Creates systemd service with security hardening
  - Auto-enables service
- ✅ Windows automated deployment script (`deploy-windows.bat`)
  - Creates directory structure (`C:\Program Files\Waterlogger`, `logs/`, `backups/`)
  - Configures for production (optional)
  - Creates Windows service with auto-restart
  - Service management commands
- ✅ Comprehensive systemd service file with:
  - Security hardening (NoNewPrivileges, ProtectSystem, etc.)
  - Proper restart policies
  - Journald integration
  - Working directory and paths
- ✅ Windows service configuration with:
  - Automatic restart on failure
  - Proper paths and config
  - Service management commands
- ✅ Complete documentation in README.md
- ✅ Both automated and manual deployment instructions

**Reference:** README.md Lines 133-377

---

### 5. Static Assets Not Embedded
**Status:** NOT IMPLEMENTED - Not a true single executable
**Location:** `cmd/waterlogger/main.go`, `web/` directory
**Issue:** Requires external files at runtime (templates, static files)

**Tasks:**
- [ ] Add Go 1.16+ `embed` package imports
- [ ] Embed static files from `web/static/`
  - [ ] CSS files
  - [ ] JavaScript files
  - [ ] Images (favicon, etc.)
- [ ] Embed HTML templates from `web/templates/`
- [ ] Update Gin router to use embedded filesystem
  - [ ] Replace `router.Static()` with embedded FS
  - [ ] Replace `router.LoadHTMLGlob()` with embedded templates
- [ ] Test single executable deployment
- [ ] Update build scripts to verify embedding
- [ ] Update documentation for deployment

**Reference:** CLAUDE.md Line 7, 201

---

## 🟡 MEDIUM PRIORITY - Feature Gaps

### 6. Setup Wizard Incomplete
**Status:** PARTIAL - Works but doesn't validate connectivity
**Location:** `internal/handlers/handlers.go:370-450`
**Issue:** Accepts MariaDB credentials without testing connection

**Tasks:**
- [ ] Add MariaDB connectivity validation before saving config
- [ ] Show clear error messages for connection failures
- [ ] Add database creation helper (optional)
- [ ] Add user creation helper for MariaDB (optional)
- [ ] Test with invalid credentials
- [ ] Test with valid credentials
- [ ] Add retry logic for connection testing
- [ ] Update setup wizard UI with validation feedback

**Reference:** CLAUDE.md Line 226-233

---

### 7. Dashboard Analytics Not Populated
**Status:** PARTIAL - Template exists but widgets are empty
**Location:** `web/templates/dashboard.html`, `internal/handlers/handlers.go`
**Issue:** Widgets don't display real data

**Tasks:**
- [ ] Implement recent samples query (last 10 across all pools)
- [ ] Implement recent adjustments query (last 10 across all pools)
- [ ] Implement pool status calculations (water quality indicators)
- [ ] Add water quality status logic (LSI/RSI thresholds)
- [ ] Add color-coded indicators (green=balanced, yellow=caution, red=action needed)
- [ ] Populate dashboard widgets with real data
- [ ] Add error handling for empty data states
- [ ] Test with multiple pools
- [ ] Test with no data

**Reference:** CLAUDE.md Line 127-131

---

### 8. Markdown Export Incomplete
**Status:** PARTIAL - Code appears truncated
**Location:** `internal/handlers/handlers.go:947`
**Issue:** Export may not include all documented sections

**Tasks:**
- [ ] Review markdown export code for completeness
- [ ] Verify all sections are included:
  - [ ] Users table
  - [ ] Pools table
  - [ ] Kits table
  - [ ] Samples table (sorted by date)
  - [ ] Measurements table
  - [ ] Indices table (LSI/RSI)
  - [ ] Adjustments table with full details
- [ ] Fix filename to use `WL[timestamp].md` format
- [ ] Test with large datasets
- [ ] Verify proper markdown table formatting
- [ ] Add hierarchical structure with headings

**Reference:** CLAUDE.md Line 147-152

---

### 9. Export Filename Format
**Status:** NOT IMPLEMENTED - Uses generic names
**Location:** All export handlers in `internal/handlers/handlers.go`
**Issue:** Filenames don't match `WL[timestamp]` specification

**Tasks:**
- [ ] Implement timestamp formatting function (YYYYMMDD_HHMMSS)
- [ ] Update Excel export filename: `WL20250126_143022.xlsx`
- [ ] Update Markdown export filename: `WL20250126_143022.md`
- [ ] Update JSON backup filename: `WL20250126_143022.json`
- [ ] Update PDF export filename: `WL20250126_143022.pdf`
- [ ] Test filename generation
- [ ] Verify browser download behavior

**Reference:** CLAUDE.md Line 147, 189

---

### 10. Unit Display in UI
**Status:** PARTIAL - Backend ready, frontend incomplete
**Location:** `web/templates/*.html`, `internal/chemistry/calculations.go`
**Issue:** UI doesn't show dual units (e.g., "75°F (24°C)")

**Tasks:**
- [ ] Update measurement display templates to show both units
- [ ] Add helper functions for dual-unit formatting
- [ ] Update pool volume display (gallons + liters)
- [ ] Update temperature display (°F + °C)
- [ ] Update all chemistry parameters with dual units
- [ ] Test with Imperial preference
- [ ] Test with Metric preference
- [ ] Ensure consistent formatting across all pages

**Reference:** CLAUDE.md Line 187-188

---

## 🟢 LOW PRIORITY - Code Quality & Improvements

### 11. Code Cleanup
**Location:** Various files
**Issue:** Technical debt and quality improvements needed

**Tasks:**
- [ ] Remove debug `fmt.Println()` statements (e.g., `handlers.go:331-333`)
- [ ] Split `handlers.go` (1500 lines) into smaller packages:
  - [ ] `handlers/users.go`
  - [ ] `handlers/pools.go`
  - [ ] `handlers/kits.go`
  - [ ] `handlers/samples.go`
  - [ ] `handlers/adjustments.go`
  - [ ] `handlers/exports.go`
- [ ] Improve session validation logic (currently too simple)
- [ ] Add input sanitization for user inputs
- [ ] Add transaction rollback on cascading deletes
- [ ] Optimize database queries with proper indexes
- [ ] Add database query logging (optional)

---

### 12. Testing Infrastructure
**Location:** `testing/` directory (currently only shell scripts)
**Issue:** No comprehensive test suite

**Tasks:**
- [ ] Add unit tests for chemistry calculations
- [ ] Add unit tests for volume calculator
- [ ] Add unit tests for adjustment calculator
- [ ] Add integration tests for API endpoints
- [ ] Add database migration tests
- [ ] Add authentication tests
- [ ] Set up CI/CD pipeline (GitHub Actions)
- [ ] Add code coverage reporting
- [ ] Add benchmark tests for performance
- [ ] Test with large datasets (1000+ samples)

---

### 13. Documentation Improvements
**Location:** Various files
**Issue:** Missing or outdated documentation

**Tasks:**
- [ ] Update `README.md` with complete feature list
- [ ] Add API documentation (Swagger/OpenAPI)
- [ ] Add deployment guide (systemd setup)
- [ ] Add database migration guide
- [ ] Add troubleshooting guide
- [ ] Add development setup guide
- [ ] Document configuration options
- [ ] Add screenshots to documentation
- [ ] Create user manual (Markdown)
- [ ] Create video tutorials (optional)

---

### 14. Configuration Inconsistencies
**Location:** `config.yaml`, `config.example.yaml`
**Issue:** Default port mismatch

**Tasks:**
- [ ] Fix default port consistency:
  - Current `config.yaml`: 2342
  - `config.example.yaml`: 2341
  - CLAUDE.md spec: 2342
- [ ] Update example config to match spec (port 2342)
- [ ] Verify all documentation references correct port
- [ ] Add validation for port range (1024-65535)

**Reference:** CLAUDE.md Line 193, 265

---

### 15. Error Handling Improvements
**Location:** Various handlers
**Issue:** Inconsistent error handling and user feedback

**Tasks:**
- [ ] Standardize error response format (JSON structure)
- [ ] Add user-friendly error messages
- [ ] Add error logging with severity levels
- [ ] Add error tracking (Sentry integration optional)
- [ ] Improve validation error messages
- [ ] Add graceful degradation for missing data
- [ ] Add timeout handling for long operations
- [ ] Add proper HTTP status codes throughout

---

## 📦 FUTURE ENHANCEMENTS (Post v1.2)

### 16. Advanced Features (Nice to Have)
- [ ] Email notifications for water quality alerts
- [ ] Mobile app (React Native or Flutter)
- [ ] Multi-language support (i18n)
- [ ] Dark mode theme
- [ ] API rate limiting
- [ ] User roles and permissions (admin, viewer, editor)
- [ ] Automatic backup scheduling
- [ ] Cloud sync capabilities
- [ ] Weather integration (temperature correlation)
- [ ] Chemical inventory tracking
- [ ] Cost tracking for chemicals
- [ ] Maintenance scheduling and reminders
- [ ] Integration with smart pool sensors
- [ ] Historical comparison reports
- [ ] AI-powered recommendations
- [ ] Mobile-optimized PWA features

---

## 🎯 RECOMMENDED IMPLEMENTATION ORDER

### Phase 1: Fix Critical Exports (High Priority)
**Estimated Time:** 2-3 days
1. Implement proper Excel export with excelize
2. Implement PDF export for adjustments
3. Fix export filenames to use `WL[timestamp]` format

**Why First:** These are documented as complete but broken - biggest gap between docs and reality.

---

### Phase 2: Complete UI Features (High Priority)
**Estimated Time:** 2-3 days
4. Wire up charts to backend data with date range filtering
5. Populate dashboard analytics widgets with real data
6. Fix markdown export completion
7. Implement dual-unit display in UI

**Why Second:** Completes the core user experience and makes the app fully functional.

---

### Phase 3: Production Deployment (High Priority)
**Estimated Time:** 2-3 days
8. Generate systemd/Windows service files
9. Embed static assets using Go's `embed` package
10. Add MariaDB connection validation to setup wizard

**Why Third:** Makes the app production-ready and truly single-executable.

---

### Phase 4: Code Quality (Low Priority)
**Estimated Time:** 3-5 days
11. Refactor handlers.go into smaller packages
12. Add comprehensive test suite
13. Remove debug statements and improve error handling
14. Fix configuration inconsistencies

**Why Last:** Technical debt that doesn't affect functionality but improves maintainability.

---

## 📊 COMPLETION STATUS

### Overall Progress
- **Core Features:** ~85% complete
- **Advanced Features:** ~70% complete
- **Production Ready:** ~60% complete
- **Code Quality:** ~70% complete

### Feature Breakdown
| Category | Complete | In Progress | Not Started | Total |
|----------|----------|-------------|-------------|-------|
| Authentication | 5 | 0 | 0 | 5 |
| Database | 7 | 1 | 0 | 8 |
| CRUD Operations | 4 | 0 | 0 | 4 |
| Water Chemistry | 5 | 0 | 0 | 5 |
| Calculations | 3 | 0 | 0 | 3 |
| Export | 1 | 2 | 1 | 4 |
| UI/UX | 7 | 3 | 0 | 10 |
| Deployment | 2 | 1 | 2 | 5 |
| Testing | 0 | 0 | 10 | 10 |

**Total:** 34/54 (63% complete)

---

## 🐛 KNOWN BUGS & ISSUES

### Critical
- [ ] **Excel export returns CSV instead of .xlsx** - Users expecting Excel format
- [ ] **Charts show empty containers** - No data visualization working
- [ ] **Setup wizard doesn't validate MariaDB** - Can save broken config

### Medium
- [ ] Session validation too simple (returns input as-is)
- [ ] No transaction rollback on cascading deletes
- [ ] Export handlers missing proper error handling

### Minor
- [ ] Debug println statements in production code
- [ ] Configuration port mismatch between example and actual
- [ ] Missing input sanitization in some endpoints

---

## 📝 NOTES

### Development Environment
- **Go Version:** 1.21+
- **Database:** SQLite 3.x OR MariaDB 10.x
- **Frontend:** Alpine.js 3.x, Chart.js 4.x
- **Build Tool:** Make + bash scripts
- **Platform:** Ubuntu Linux + Windows

### Key Dependencies to Add
- `github.com/xuri/excelize/v2` - For Excel export
- `golang.org/x/sys/windows/svc` - For Windows service support
- Consider: PDF generation library (if not using browser-based)

### Testing Notes
- Test with SQLite AND MariaDB configurations
- Test with large datasets (1000+ samples)
- Test on both Linux and Windows
- Test mobile responsive design
- Test all export formats with real data

### Documentation References
- Main requirements: `requirements.md`
- Development notes: `CLAUDE.md`
- Build instructions: `build.sh`, `Makefile`
- Configuration: `config.yaml`, `config.example.yaml`

---

**Last Review Date:** 2025-10-26
**Next Review Date:** After Phase 1 completion
**Maintainer:** Development Team
