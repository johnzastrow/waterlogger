# Waterlogger Build & Deployment Guide

## Quick Start

### Building the Application

```bash
# Build for current platform
go build -o waterlogger cmd/waterlogger/main.go

# Or use the Makefile
make build
```

### Cross-Platform Building

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o waterlogger-linux cmd/waterlogger/main.go

# Windows  
GOOS=windows GOARCH=amd64 go build -o waterlogger.exe cmd/waterlogger/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o waterlogger-mac cmd/waterlogger/main.go

# Or build all platforms
make build-all
```


### jcz
```
go build -o waterlogger ./cmd/waterlogger
 export GIN_MODE=release && ./waterlogger > /tmp/server.log 2>&1 &
ps -aux | grep 'waterlogger'
```
## Testing the Application

### 1. Start the Application
```bash
./waterlogger
```

### 2. Open Web Browser
Navigate to: `http://localhost:2342`

### 3. Complete Setup Wizard
- Create admin user account
- Configure database settings  
- Set preferences

### 4. Test Password Reset
```bash
# Reset password for any user
./waterlogger -reset-password username

# Example: Reset admin password to "admin"
echo "admin" | ./waterlogger -reset-password admin
```

## Key Features Verified

✅ **Simplified Authentication**
- Password complexity requirements removed
- Only requires non-empty passwords
- Secure bcrypt hashing

✅ **Port Configuration**  
- Default port changed to 2342
- Configurable via config.yaml

✅ **Command Line Tools**
- Password reset utility
- Database migration (SQLite ↔ MariaDB)
- Data export/import
- Version and help commands

✅ **Water Chemistry**
- Automatic LSI/RSI calculations
- Unit conversion (Imperial/Metric)
- Mid-range defaults for missing parameters

✅ **Web Interface**
- Dashboard with overview
- Sample management with forms
- Pool/kit management
- Charts and export functionality

✅ **Database Support**
- SQLite (default, no setup required)
- MariaDB (configurable via setup wizard)

## Deployment Options

### Option 1: Standalone Binary
```bash
# Copy binary and config
cp waterlogger /usr/local/bin/
cp config.yaml /etc/waterlogger/

# Run directly
waterlogger -config /etc/waterlogger/config.yaml
```

### Option 2: Docker
```bash
# Build Docker image
docker build -t waterlogger:1.0.0 .

# Run with Docker Compose
docker-compose up -d
```

### Option 3: System Service

**Linux (systemd):**
```bash
# Copy service file
sudo cp docs/systemd/waterlogger.service /etc/systemd/system/
sudo systemctl enable waterlogger
sudo systemctl start waterlogger
```

**Windows:**
```cmd
# Create Windows service
sc create Waterlogger binpath="C:\waterlogger\waterlogger.exe" start=auto
sc start Waterlogger
```

## Configuration

### Basic config.yaml
```yaml
server:
  port: 2342
  host: "0.0.0.0"

database:
  type: "sqlite"  # or "mariadb"
  sqlite:
    path: "waterlogger.db"

app:
  name: "Waterlogger"
  version: "1.0.0"
  secret_key: "change-this-secret-key"
```

## Troubleshooting

### Common Issues

1. **Port already in use**
   - Change port in config.yaml
   - Or set via environment: `PORT=3000`

2. **Database connection failed**
   - Check file permissions for SQLite
   - Verify MariaDB connection details

3. **Forgot admin password**
   ```bash
   ./waterlogger -reset-password admin
   ```

4. **Template not found errors**
   - Ensure `web/templates/` directory is present
   - Run from correct working directory

### Logs and Debugging

```bash
# Run with debug logging
GIN_MODE=debug ./waterlogger

# Check system service logs (Linux)
sudo journalctl -u waterlogger -f

# Check Windows service logs
eventvwr.msc → Windows Logs → Application
```

## Development

### Requirements
- Go 1.21 or later
- SQLite or MariaDB (optional)

### Development workflow
```bash
# Install dependencies
go mod tidy

# Run tests
go test ./...

# Format code
go fmt ./...

# Run with hot reload (if air installed)
make dev

# Build and test
make build && ./waterlogger
```

## Security Notes

- Change default secret key in production
- Use HTTPS in production environments
- Regularly backup database files
- Monitor application logs for security events

## Support

- **Documentation**: See README.md for comprehensive guide
- **Issues**: Report at GitHub repository
- **Configuration**: All settings in config.yaml are documented