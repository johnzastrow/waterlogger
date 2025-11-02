# Waterlogger Build, Run, and Debug Guide

**Last Updated:** 2025-11-01
**Version:** 1.4.0

Complete guide for building, running, debugging, and troubleshooting Waterlogger across all platforms.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Building from Source](#building-from-source)
- [Running Locally](#running-locally)
- [Running in Production](#running-in-production)
- [Debugging](#debugging)
- [Troubleshooting](#troubleshooting)
- [Command Line Options](#command-line-options)

---

## Prerequisites

### All Platforms

- **Go 1.21 or later** - Download from https://go.dev/dl/
- **Git** (optional, for cloning the repository)
- **CGO support** (required for SQLite) - see platform-specific instructions

### Platform-Specific Requirements

#### Windows

- **MSYS2 with MinGW-w64 GCC** (required for CGO/SQLite)
  - Download: https://www.msys2.org/
  - ⚠️ **DO NOT use Cygwin** - it cannot build native Windows Go programs

**Quick Setup:**
```cmd
# 1. Install MSYS2 from https://www.msys2.org/
# 2. In MSYS2 terminal, install GCC:
pacman -S mingw-w64-x86_64-gcc

# 3. Add C:\msys64\mingw64\bin to Windows PATH
# 4. Restart terminal and verify
gcc --version
```

See [BUILD_REQUIREMENTS.md](BUILD_REQUIREMENTS.md) for detailed Windows setup instructions.

#### Linux

**Ubuntu/Debian:**
```bash
# Install build tools
sudo apt update
sudo apt install build-essential golang-go

# Or if Go not available via apt
# Download from https://go.dev/dl/
```

**Fedora/RHEL:**
```bash
sudo dnf groupinstall "Development Tools"
sudo dnf install golang
```

#### macOS

```bash
# Install Xcode Command Line Tools (includes C compiler for CGO)
xcode-select --install

# Install Go (if not already installed)
brew install go
```

#### Windows Subsystem for Linux (WSL2) - Alternative

If you have WSL2 enabled, you can build the Linux version directly without dealing with Windows C compilers:

```bash
# In WSL2 terminal
cd /mnt/c/Users/YourUsername/path/to/waterlogger
sudo apt update
sudo apt install build-essential golang-go
./build.sh
```

---

## Building from Source

### Quick Build (All Platforms)

The easiest way to build is using the provided build scripts:

**Linux/macOS:**
```bash
cd waterlogger
chmod +x build.sh
./build.sh
```

**Windows:**
```cmd
cd waterlogger
build.bat
```

The build script automatically:
- Enables CGO (required for SQLite)
- Sets build timestamp for deployment tracking
- Produces executable: `./waterlogger` (Linux/macOS) or `waterlogger.exe` (Windows)

### Manual Build (If Scripts Don't Work)

**Linux/macOS:**
```bash
export CGO_ENABLED=1
go mod download
go mod tidy
go build -ldflags "-X main.BuildDate=$(date +%Y-%m-%d) -X main.BuildTime=$(date +%H:%M:%S)" -o waterlogger ./cmd/waterlogger
```

**Windows (Command Prompt):**
```cmd
set CGO_ENABLED=1
go mod download
go mod tidy
go build -ldflags "-X main.BuildDate=%date% -X main.BuildTime=%time%" -o waterlogger.exe .\cmd\waterlogger
```

### Build Verification

After building, verify CGO is enabled:

```bash
# Linux/macOS
./waterlogger -version

# Windows
waterlogger.exe -version
```

Should output version information like: `Waterlogger v1.4.0`

If you see a SQLite error about CGO, it wasn't enabled during build. Check [Troubleshooting](#troubleshooting) section.

### Cross-Platform Building

⚠️ **Important:** Cross-compiling with CGO is complex due to SQLite's C dependencies.

**Recommended approach:** Build on the target platform using provided build scripts.

**For Linux → Windows (requires mingw-w64 cross-compiler):**
```bash
sudo apt install gcc-mingw-w64
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
  go build -o waterlogger.exe ./cmd/waterlogger
```

**For Windows → Linux:** Not recommended. Build on Linux instead.

---

## Running Locally

### First Time Setup

1. **Download dependencies:**
   ```bash
   go mod download
   ```

2. **Start the application:**
   ```bash
   ./waterlogger  # Linux/macOS
   waterlogger.exe  # Windows
   ```

3. **Open web browser:**
   - Navigate to: `http://localhost:2342`

4. **Complete setup wizard:**
   - Create an administrator account
   - Configure database settings (SQLite is default, no setup needed)
   - Set server preferences

### Development Mode

Run with debug logging and Gin debug mode:

**Linux/macOS:**
```bash
export GIN_MODE=debug
./waterlogger
```

**Windows:**
```cmd
set GIN_MODE=debug
waterlogger.exe
```

### Development with Configuration

Create a `config.yaml` for development:

```bash
cp config.example.yaml config.yaml
# Edit config.yaml as needed
```

Run with custom config:
```bash
./waterlogger -config config.yaml
```

### Using a Different Port

**Via config file:**
```yaml
server:
  port: 3000  # Change from default 2342
  host: "localhost"
```

**Via environment (overrides config):**
```bash
# Linux/macOS
PORT=3000 ./waterlogger

# Windows
set PORT=3000
waterlogger.exe
```

Then access at: `http://localhost:3000`

---

## Running in Production

### Single Binary Deployment

**Step 1: Build for target platform**
```bash
./build.sh  # On Linux/macOS target
# or
build.bat   # On Windows target
```

**Step 2: Copy files to production location**

**Linux:**
```bash
sudo cp waterlogger /opt/waterlogger/
sudo cp config.yaml /opt/waterlogger/
sudo chmod 600 /opt/waterlogger/config.yaml
```

**Windows:**
```cmd
copy waterlogger.exe "C:\Program Files\Waterlogger\"
copy config.yaml "C:\Program Files\Waterlogger\"
```

**Step 3: Run application**

**Linux:**
```bash
./waterlogger -config /opt/waterlogger/config.yaml
```

**Windows:**
```cmd
waterlogger.exe -config "C:\Program Files\Waterlogger\config.yaml"
```

### Automated Deployment Scripts

**Linux (Recommended):**
```bash
sudo ./deploy-linux.sh
```

Creates production setup with:
- Directory structure (`/opt/waterlogger`)
- Dedicated user with proper permissions
- systemd service file
- Security hardening
- Optional production logging configuration

**Windows (Recommended):**
```cmd
REM Run as Administrator
deploy-windows.bat
```

Creates production setup with:
- Directory structure (`C:\Program Files\Waterlogger`)
- Windows service with auto-start
- Automatic restart on failure
- Optional production logging configuration

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed deployment instructions.

### Production Configuration

Recommended `config.yaml` for production:

```yaml
server:
  port: 2342
  host: "0.0.0.0"          # Or specific IP for security

database:
  type: "sqlite"            # or "mariadb"
  sqlite:
    path: "/opt/waterlogger/waterlogger.db"

app:
  name: "Waterlogger"
  version: "1.4.0"
  secret_key: "change-this-to-a-secure-random-string"

logging:
  level: "info"             # Not "debug" in production
  format: "json"            # Structured JSON format
  output: "file"            # File output only
  file_path: "/var/log/waterlogger/waterlogger.log"
  max_size: 100             # MB
  max_backups: 7            # Keep 7 backups
  max_age: 30               # Days
  compress: true            # Compress old logs
```

### Production Mode (Gin)

Set Gin mode for better performance:

**Via config:**
```yaml
app:
  name: "production"  # Enables Gin release mode
```

**Via environment:**
```bash
# Linux/macOS
export GIN_MODE=release
./waterlogger

# Windows
set GIN_MODE=release
waterlogger.exe
```

---

## Debugging

### Debug Logging

**Enable debug logging in config.yaml:**
```yaml
logging:
  level: "debug"
  format: "console"
  output: "both"
```

**Or via environment:**
```bash
# Linux/macOS
export GIN_MODE=debug
./waterlogger

# Windows
set GIN_MODE=debug
waterlogger.exe
```

### View Logs

**Console (during development):**
Logs appear directly in terminal output.

**File logs (production):**
```bash
# Linux
tail -f /var/log/waterlogger/waterlogger.log

# Windows
type "C:\Program Files\Waterlogger\logs\waterlogger.log"
```

**Systemd service logs (Linux):**
```bash
sudo journalctl -u waterlogger -f      # Follow live logs
sudo journalctl -u waterlogger -n 100  # Last 100 entries
```

**JSON log parsing (production):**
```bash
# Filter by level
cat logs/waterlogger.log | jq 'select(.level=="error")'

# Filter by component
cat logs/waterlogger.log | jq 'select(.component=="database")'

# Pretty print
cat logs/waterlogger.log | jq .
```

### Testing Endpoints Locally

**Using curl:**
```bash
# Test login
curl -X POST http://localhost:2342/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Get pools
curl http://localhost:2342/api/pools

# Get samples
curl http://localhost:2342/api/samples
```

**Using browser developer tools:**
1. Open application at `http://localhost:2342`
2. Open DevTools (F12 or Right-click → Inspect)
3. Go to Network tab
4. Perform operations and watch API calls
5. Check Console tab for JavaScript errors

### Database Debugging

**SQLite:**
```bash
# Connect to database directly
sqlite3 waterlogger.db

# View tables
.tables

# View schema
.schema users

# Query data
SELECT * FROM users;

# Exit
.exit
```

**MariaDB:**
```bash
# Connect to database
mysql -u waterlogger -p waterlogger

# View tables
SHOW TABLES;

# View schema
DESCRIBE users;

# Query data
SELECT * FROM users;

# Exit
exit
```

### Common Debug Scenarios

#### Port Already in Use

```bash
# Find what's using the port
sudo lsof -i :2342              # Linux/macOS
netstat -ano | findstr :2342    # Windows

# Kill the process
sudo kill -9 <PID>              # Linux/macOS
taskkill /PID <PID> /F          # Windows
```

#### Database Connection Issues

```bash
# Test database connection
# Linux
sqlite3 waterlogger.db ".tables"

# Check database file permissions
ls -la waterlogger.db

# Check logs for detailed error
tail -f logs/waterlogger.log
```

#### Template/Asset Loading Issues

```bash
# Ensure correct working directory
pwd  # Should be project root

# Verify templates exist
ls -la web/templates/

# Verify static files exist
ls -la web/static/
```

---

## Troubleshooting

### Build Issues

#### "CGO_ENABLED=0" Error

**Problem:**
```
Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work
```

**Solution:** CGO was not enabled during build.

```bash
# Verify CGO is available
go env CGO_ENABLED  # Should output: 1

# Re-build with CGO enabled
export CGO_ENABLED=1  # Linux/macOS
set CGO_ENABLED=1     # Windows

# Then build
./build.sh  # or build.bat
```

See [BUILD_REQUIREMENTS.md](BUILD_REQUIREMENTS.md#issue-1-cgo_enabled0-error) for detailed instructions.

#### "gcc: command not found" (Windows)

**Problem:** GCC not installed or not in PATH

**Solution:**
1. Install MSYS2: https://www.msys2.org/
2. In MSYS2 terminal:
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```
3. Add `C:\msys64\mingw64\bin` to Windows PATH
4. Restart command prompt and verify: `gcc --version`

See [BUILD_REQUIREMENTS.md](BUILD_REQUIREMENTS.md#windows-build-requirements) for detailed instructions.

#### "cannot find -lgcc_s" (Windows)

**Problem:** Incomplete MSYS2 MinGW installation

**Solution:**
1. In MSYS2 terminal:
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```
2. Verify `C:\msys64\mingw64\bin` is in Windows PATH
3. Restart command prompt

#### Cygwin Compiler Error

**Problem:**
```
#error "don't use the cygwin compiler to build native Windows programs"
```

**Cause:** Cygwin detected instead of MinGW

**Solution:** Use MSYS2 MinGW instead:
1. Uninstall or remove Cygwin from PATH
2. Install MSYS2 (see above)
3. Build using `build.bat`

### Runtime Issues

#### Port Already in Use

**Problem:** Port 2342 is already in use

**Solutions:**
1. Change port in `config.yaml`:
   ```yaml
   server:
     port: 3000
   ```

2. Or find and kill the process:
   ```bash
   # Linux/macOS
   sudo lsof -i :2342
   sudo kill -9 <PID>

   # Windows
   netstat -ano | findstr :2342
   taskkill /PID <PID> /F
   ```

#### Database Connection Failed

**Problem:** Cannot connect to SQLite or MariaDB

**Solutions for SQLite:**
```bash
# Check file exists and is writable
ls -la waterlogger.db
chmod 666 waterlogger.db

# Check disk space
df -h .

# Check logs for details
tail logs/waterlogger.log
```

**Solutions for MariaDB:**
```bash
# Verify MariaDB is running
sudo systemctl status mariadb

# Test connection manually
mysql -u waterlogger -p waterlogger

# Check config.yaml credentials
cat config.yaml | grep -A 5 "mariadb:"

# Check logs
tail logs/waterlogger.log
```

#### Template Not Found

**Problem:** Template loading error

**Solutions:**
```bash
# Ensure correct working directory
pwd  # Should be project root

# Verify templates exist
ls -la web/templates/

# Run from correct directory
cd /path/to/waterlogger
./waterlogger
```

#### Logs Not Appearing

**Problem:** No logs being written

**Solutions:**
```bash
# Check config.yaml logging section
cat config.yaml | grep -A 10 "logging:"

# Verify log directory exists
mkdir -p logs/

# Check permissions
chmod 755 logs/

# Verify log level allows message
# Messages above configured level won't show

# Check output setting
# Change to "both" for console + file
```

### Service Issues (Linux)

#### Service Won't Start

```bash
# Check service status
sudo systemctl status waterlogger

# View detailed logs
sudo journalctl -u waterlogger -f

# Test manual startup
sudo -u waterlogger /opt/waterlogger/waterlogger \
  -config /opt/waterlogger/config.yaml

# Check file permissions
sudo chown -R waterlogger:waterlogger /opt/waterlogger
sudo chmod 755 /opt/waterlogger/waterlogger
```

#### Service Crashes on Startup

```bash
# View crash logs
sudo journalctl -u waterlogger --since "5 minutes ago"

# Check database connectivity
sudo -u waterlogger sqlite3 /opt/waterlogger/waterlogger.db ".tables"

# Verify config file is readable
sudo -u waterlogger cat /opt/waterlogger/config.yaml
```

### Performance Issues

#### Slow Startup

**Check database:**
```bash
# Verify database file isn't huge
du -h waterlogger.db

# Check for slow migrations
tail -50 logs/waterlogger.log | grep -i migration
```

**Check logging:**
```yaml
# Set to info or warn in production
logging:
  level: "info"  # Not "debug"
```

#### High Memory Usage

```bash
# Check process memory
ps aux | grep waterlogger

# For large datasets, consider:
# - Archiving old samples
# - Using MariaDB instead of SQLite
# - Increasing available memory
```

#### Slow API Responses

```bash
# Check logs for slow queries
tail logs/waterlogger.log | grep -i "slow"

# Review database performance
# - Check for missing indices
# - Review query patterns
# - Consider caching
```

---

## Command Line Options

### Available Options

```
waterlogger [options]

Options:
  -config string           Path to configuration file (default: config.yaml)
  -version                 Show version information
  -help                    Show help message
  -reset-password string   Reset password for specified username
  -migrate-to-mariadb      Migrate data from SQLite to MariaDB
  -migrate-to-sqlite       Migrate data from MariaDB to SQLite
  -export string           Export database data to backup file
  -import string           Import database data from backup file
```

### Usage Examples

**Show version:**
```bash
./waterlogger -version
```

**Use custom config file:**
```bash
./waterlogger -config /etc/waterlogger/config.yaml
```

**Reset a user's password:**
```bash
./waterlogger -reset-password username
# Prompts for new password

# Non-interactive:
echo "newpassword" | ./waterlogger -reset-password username
```

**Migrate SQLite to MariaDB:**
```bash
# Update config.yaml to MariaDB settings first
./waterlogger -migrate-to-mariadb
```

**Export database backup:**
```bash
./waterlogger -export waterlogger_backup.json
```

**Import database backup:**
```bash
./waterlogger -import waterlogger_backup.json
```

**Show help:**
```bash
./waterlogger -help
```

---

## Additional Resources

- **Building:** See [BUILD_REQUIREMENTS.md](BUILD_REQUIREMENTS.md) for detailed platform setup
- **Deployment:** See [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment
- **Logging:** See [LOGGING.md](LOGGING.md) for logging configuration and monitoring
- **Configuration:** See `config.example.yaml` for all available settings
- **API:** See [API.md](API.md) for REST API documentation
- **Migrations:** See [MIGRATIONS.md](MIGRATIONS.md) for database schema migrations

---

## Getting Help

If you encounter issues:

1. **Check logs first:** Most problems appear in logs
   ```bash
   tail -f logs/waterlogger.log
   sudo journalctl -u waterlogger -f  # Linux systemd
   ```

2. **Verify prerequisites:** Ensure Go, GCC/C compiler installed
   ```bash
   go version
   gcc --version
   ```

3. **Try debug mode:** Enable detailed logging
   ```bash
   export GIN_MODE=debug
   ./waterlogger
   ```

4. **Report issues:** https://github.com/johnzastrow/waterlogger/issues
   - Include: OS, Go version, GCC version, full error message, build command used

---

**Last Updated:** 2025-11-01
**Version:** 1.4.0
