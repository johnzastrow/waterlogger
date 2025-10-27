# Waterlogger TODO List

**Last Updated:** 2025-10-26
**Version:** 1.2
**Status:** Core functionality complete, production features incomplete

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
- [x] Database migration commands
- [x] Audit trail (created_by, updated_by timestamps)

---

## 🔴 HIGH PRIORITY - Critical Gaps

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

### 3. Charts Not Working
**Status:** PARTIAL - Templates exist but no data binding
**Location:** `web/templates/dashboard.html`, need backend handlers
**Issue:** Chart.js included but not connected to backend data

**Tasks:**
- [ ] Implement chart data API endpoint (`GET /api/charts/data`)
- [ ] Add date range filtering logic (default: last 30 days)
- [ ] Add parameter selection logic
- [ ] Exclude TDS, CYA, SAL by default (per requirements)
- [ ] Implement multi-parameter overlays for correlation
- [ ] Add chart rendering JavaScript in dashboard template
- [ ] Test chart zoom capabilities
- [ ] Add loading states for chart data
- [ ] Optimize queries for large datasets
- [ ] Add chart export functionality (PNG download)

**Reference:** CLAUDE.md Line 30-33, 141-143

---

### 4. Service File Generation Missing
**Status:** NOT IMPLEMENTED - Critical for production deployment
**Location:** Should be in `internal/service/` (new package)
**Issue:** No systemd or Windows service files exist

**Tasks:**
- [ ] Create `internal/service/` package
- [ ] Implement systemd service file generation for Linux
  - [ ] Template: `/etc/systemd/system/waterlogger.service`
  - [ ] Auto-detect installation path
  - [ ] Set proper user/group permissions
  - [ ] Add restart policies
- [ ] Implement Windows service file generation
  - [ ] Use `golang.org/x/sys/windows/svc` package
  - [ ] Create service installer/uninstaller
- [ ] Add command-line flags:
  - [ ] `./waterlogger -install-service`
  - [ ] `./waterlogger -uninstall-service`
  - [ ] `./waterlogger -generate-service-file`
- [ ] Add service management to setup wizard
- [ ] Test on Ubuntu Linux
- [ ] Test on Windows Server
- [ ] Update documentation with service installation steps

**Reference:** CLAUDE.md Line 67, 200

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
