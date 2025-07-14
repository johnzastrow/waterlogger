# Waterlogger

A comprehensive web application for managing pool and hot tub water chemistry parameters, with built-in calculations, data visualization, and export capabilities.

## Features

- **Multi-User Support**: Multiple users can manage water testing data
- **Pool Management**: Track multiple pools and hot tubs with detailed specifications
- **Water Chemistry**: Record comprehensive water parameter measurements
- **Automatic Calculations**: LSI (Langelier Saturation Index) and RSI (Ryznar Stability Index) calculations
- **Data Visualization**: Interactive charts showing parameter trends over time
- **Export Functionality**: Export data to Excel and Markdown formats
- **Responsive Design**: Mobile-friendly interface with modern UI
- **Database Flexibility**: Support for SQLite and MariaDB databases
- **Cross-Platform**: Single executable for Windows and Linux

## Screenshots

![Dashboard](docs/images/dashboard.png)
*Main dashboard showing recent samples and pool status*

![Setup Wizard](docs/images/setup.png)
*First-run setup wizard for configuration*

## Quick Start

### Prerequisites

- Go 1.21 or later (for building from source)
- SQLite (included) or MariaDB (optional)

### Installation

#### Option 1: Download Pre-built Binary

1. Download the latest release for your platform from [Releases](https://github.com/your-org/waterlogger/releases)
2. Extract the binary to your desired location
3. Run the application

#### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/your-org/waterlogger.git
cd waterlogger

# Build the application
go build -o waterlogger cmd/waterlogger/main.go

# Run the application
./waterlogger
```

### First Run

1. Start the application
2. Open your web browser to `http://localhost:2341`
3. Complete the setup wizard to:
   - Create an administrator account
   - Configure database settings
   - Set server preferences

## Building and Running

### Windows

#### Prerequisites
- Go 1.21 or later
- Git (optional, for cloning)

#### Build Steps
```cmd
# Clone the repository (if not already done)
git clone https://github.com/your-org/waterlogger.git
cd waterlogger

# Download dependencies
go mod tidy

# Build for Windows
go build -o waterlogger.exe cmd/waterlogger/main.go

# Run the application
waterlogger.exe
```

#### Running as Windows Service
1. Copy the executable to your preferred location (e.g., `C:\Program Files\Waterlogger\`)
2. Create a Windows service using `sc create` or a service manager
3. Configure the service to run at startup

Example service creation:
```cmd
sc create Waterlogger binpath="C:\Program Files\Waterlogger\waterlogger.exe" start=auto
sc description Waterlogger "Pool and Hot Tub Water Management System"
sc start Waterlogger
```

### Linux

#### Prerequisites
- Go 1.21 or later
- Git (optional, for cloning)

#### Build Steps
```bash
# Clone the repository (if not already done)
git clone https://github.com/your-org/waterlogger.git
cd waterlogger

# Download dependencies
go mod tidy

# Build for Linux
go build -o waterlogger cmd/waterlogger/main.go

# Make executable
chmod +x waterlogger

# Run the application
./waterlogger
```

#### Running as Linux Service (systemd)
1. Copy the executable to `/usr/local/bin/waterlogger`
2. Create a systemd service file:

```bash
sudo tee /etc/systemd/system/waterlogger.service > /dev/null <<EOF
[Unit]
Description=Waterlogger - Pool and Hot Tub Water Management System
After=network.target

[Service]
Type=simple
User=waterlogger
WorkingDirectory=/var/lib/waterlogger
ExecStart=/usr/local/bin/waterlogger
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
```

3. Enable and start the service:
```bash
# Create user and directory
sudo useradd -r -s /bin/false waterlogger
sudo mkdir -p /var/lib/waterlogger
sudo chown waterlogger:waterlogger /var/lib/waterlogger

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable waterlogger
sudo systemctl start waterlogger

# Check status
sudo systemctl status waterlogger
```

### Cross-Platform Building

Build for multiple platforms:

```bash
# Build for Windows (from any platform)
GOOS=windows GOARCH=amd64 go build -o waterlogger.exe cmd/waterlogger/main.go

# Build for Linux (from any platform)
GOOS=linux GOARCH=amd64 go build -o waterlogger cmd/waterlogger/main.go

# Build for macOS (from any platform)
GOOS=darwin GOARCH=amd64 go build -o waterlogger-mac cmd/waterlogger/main.go
```

## Configuration

### Configuration File

The application uses a YAML configuration file (`config.yaml`) with the following structure:

```yaml
server:
  port: 2341
  host: "localhost"

database:
  type: "sqlite" # sqlite or mariadb
  sqlite:
    path: "waterlogger.db"
  mariadb:
    host: "localhost"
    port: 3306
    username: "waterlogger"
    password: "password"
    database: "waterlogger"

app:
  name: "Waterlogger"
  version: "1.0.0"
  secret_key: "your-secret-key-change-this"
```

### Command Line Options

```bash
waterlogger [options]

Options:
  -config string    Path to configuration file (default: config.yaml)
  -version          Show version information
  -help             Show help message
```

### Database Setup

#### SQLite (Default)
- No additional setup required
- Database file is created automatically
- Recommended for single-user or small deployments

#### MariaDB
1. Install MariaDB server
2. Create database and user:
```sql
CREATE DATABASE waterlogger;
CREATE USER 'waterlogger'@'localhost' IDENTIFIED BY 'your-password';
GRANT ALL PRIVILEGES ON waterlogger.* TO 'waterlogger'@'localhost';
FLUSH PRIVILEGES;
```
3. Update configuration file with connection details
4. Restart the application

## Usage

### Water Parameters

The application tracks the following water chemistry parameters:

- **FC (Free Chlorine)**: 1.0-4.0 ppm - Available chlorine for sanitization
- **TC (Total Chlorine)**: Should match FC - Total chlorine including combined chlorine
- **pH**: 7.4-7.6 - Acidity/alkalinity level
- **TA (Total Alkalinity)**: 80-120 ppm - pH buffering capacity
- **CH (Calcium Hardness)**: 200-400 ppm - Dissolved calcium concentration
- **CYA (Cyanuric Acid)**: 30-50 ppm - Chlorine stabilizer (optional)
- **Temperature**: Water temperature in °F
- **Salinity**: 2,700-3,400 ppm - For saltwater pools (optional)
- **TDS (Total Dissolved Solids)**: Total dissolved substances (optional)

### Calculated Indices

- **LSI (Langelier Saturation Index)**: Indicates water balance (-0.3 to +0.3 ideal)
- **RSI (Ryznar Stability Index)**: Predicts scaling tendency (6.0-7.0 ideal)

### Data Export

Export your data in two formats:

1. **Excel Export**: Multi-worksheet file with separate sheets for each data type
2. **Markdown Export**: Structured text report with tables and summaries

Files are named with format: `WL[timestamp].xlsx` or `WL[timestamp].md`

## API Documentation

### REST Endpoints

#### Authentication
- `POST /api/login` - User login
- `POST /api/logout` - User logout

#### Pools
- `GET /api/pools` - List all pools
- `POST /api/pools` - Create new pool
- `PUT /api/pools/:id` - Update pool
- `DELETE /api/pools/:id` - Delete pool

#### Samples
- `GET /api/samples` - List all samples
- `POST /api/samples` - Create new sample
- `PUT /api/samples/:id` - Update sample
- `DELETE /api/samples/:id` - Delete sample

#### Export
- `GET /api/export/excel` - Export data to Excel
- `GET /api/export/markdown` - Export data to Markdown

## Development

### Project Structure

```
waterlogger/
├── cmd/waterlogger/          # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database abstraction layer
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # HTTP middleware
│   ├── models/              # Data models
│   ├── services/            # Business logic
│   └── chemistry/           # Water chemistry calculations
├── web/
│   ├── static/              # Static assets (CSS, JS)
│   └── templates/           # HTML templates
├── migrations/              # Database migrations
├── docs/                    # Documentation
└── config.yaml              # Configuration file
```

### Testing

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags integration ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go best practices and conventions
- Use `gofmt` for code formatting
- Include unit tests for new features
- Update documentation for API changes

## Troubleshooting

### Common Issues

#### Port Already in Use
If port 2341 is already in use, modify the configuration file:
```yaml
server:
  port: 3000  # Change to available port
```

#### Database Connection Issues
1. **SQLite**: Check file permissions and available disk space
2. **MariaDB**: Verify connection details and database server status

#### Template Loading Issues
Ensure the `web/templates` directory is present and accessible from the working directory.

### Log Files

Application logs are written to stdout by default. For service deployments, logs are typically captured by the service manager (systemd on Linux, Event Log on Windows).

### Performance Tuning

- **SQLite**: Use WAL mode for better concurrent access
- **MariaDB**: Configure connection pooling and query optimization
- **Memory**: Monitor memory usage for large datasets

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: Report bugs and feature requests on [GitHub Issues](https://github.com/your-org/waterlogger/issues)
- **Discussions**: Join the community on [GitHub Discussions](https://github.com/your-org/waterlogger/discussions)
- **Documentation**: Visit the [Wiki](https://github.com/your-org/waterlogger/wiki) for detailed guides

## Acknowledgments

- Water chemistry calculations based on research from [WaterPy](https://github.com/johnzastrow/WaterPy)
- UI framework: [Alpine.js](https://alpinejs.dev/)
- Database ORM: [GORM](https://gorm.io/)
- Web framework: [Gin](https://gin-gonic.com/)

## Changelog

### Version 1.0.0
- Initial release
- Core water chemistry tracking
- Multi-user support
- Export functionality
- Setup wizard
- Cross-platform support

---

**Waterlogger** - Making pool and hot tub water management simple and efficient.