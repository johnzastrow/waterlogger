# Contributing to Waterlogger

Thank you for your interest in contributing to Waterlogger! This document provides guidelines and information for contributors.

## Code of Conduct

By participating in this project, you are expected to uphold our Code of Conduct. Please report unacceptable behavior to the project maintainers.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, please include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **Environment details** (OS, Go version, database type)
- **Screenshots** if applicable
- **Log output** if relevant

### Suggesting Features

Feature requests are welcome! Please provide:

- **Clear description** of the feature
- **Use case** and motivation
- **Proposed implementation** (if you have ideas)
- **Alternative solutions** considered

### Pull Requests

1. **Fork** the repository
2. **Create a feature branch** from `main`
3. **Make your changes** with clear, descriptive commits
4. **Add tests** for new functionality
5. **Update documentation** as needed
6. **Ensure tests pass** and code follows style guidelines
7. **Submit a pull request**

#### Pull Request Guidelines

- Use a clear and descriptive title
- Describe what changes you made and why
- Reference any related issues
- Keep changes focused and atomic
- Add tests for new features
- Update documentation for API changes

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- SQLite (included with Go)
- MariaDB (optional, for testing)

### Local Development

```bash
# Clone your fork
git clone https://github.com/your-username/waterlogger.git
cd waterlogger

# Download dependencies
go mod tidy

# Run tests
go test ./...

# Build the application
go build -o waterlogger cmd/waterlogger/main.go

# Run the application
./waterlogger
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/chemistry/

# Run integration tests
go test -tags integration ./...
```

### Code Style

- Follow Go best practices and conventions
- Use `gofmt` for code formatting
- Use `go vet` for static analysis
- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines

### Commit Messages

Use clear and descriptive commit messages:

```
feat: add water chemistry calculation for LSI index
fix: correct temperature conversion in chemistry module
docs: update API documentation for pool endpoints
refactor: simplify database connection handling
test: add unit tests for user authentication
```

Prefixes:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Build process or auxiliary tool changes

## Project Structure

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
└── tests/                   # Integration tests
```

## Areas for Contribution

### High Priority

- **Sample Management Interface**: Complete the measurements input forms
- **Charts and Visualization**: Implement Chart.js integration
- **Export Functionality**: Excel and Markdown export features
- **Unit Tests**: Increase test coverage
- **Documentation**: API documentation and user guides

### Medium Priority

- **Database Migration Utility**: Tool for switching between database types
- **User Interface Improvements**: Enhanced mobile responsiveness
- **Performance Optimization**: Database query optimization
- **Internationalization**: Multi-language support

### Low Priority

- **Additional Database Support**: PostgreSQL support
- **Advanced Charts**: More chart types and customization
- **Email Notifications**: Alert system for parameter thresholds
- **Mobile App**: Native mobile application

## Testing Guidelines

### Unit Tests

- Test individual functions and methods
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Aim for high coverage of critical paths

Example:
```go
func TestCalculateLSI(t *testing.T) {
    tests := []struct {
        name     string
        tempC    float64
        ph       float64
        tds      float64
        ca       float64
        hco3     float64
        expected float64
    }{
        {
            name:     "typical pool water",
            tempC:    25.0,
            ph:       7.5,
            tds:      300.0,
            ca:       250.0,
            hco3:     100.0,
            expected: -0.2,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateLSI(tt.tempC, tt.ph, tt.tds, tt.ca, tt.hco3)
            if math.Abs(result-tt.expected) > 0.1 {
                t.Errorf("CalculateLSI() = %v, expected %v", result, tt.expected)
            }
        })
    }
}
```

### Integration Tests

- Test complete workflows
- Use test database
- Test API endpoints
- Test database operations

### Frontend Tests

- Test JavaScript functionality
- Test form validation
- Test responsive design
- Test accessibility

## Documentation

### Code Documentation

- Use Go doc comments for public functions
- Include examples in documentation
- Document complex algorithms
- Explain design decisions

### API Documentation

- Document all endpoints
- Include request/response examples
- Document error codes
- Update OpenAPI specification

### User Documentation

- Keep README up to date
- Write clear setup instructions
- Document configuration options
- Provide troubleshooting guides

## Release Process

1. **Version Bump**: Update version in relevant files
2. **Changelog**: Update CHANGELOG.md with new features and fixes
3. **Testing**: Ensure all tests pass
4. **Documentation**: Update documentation
5. **Tag Release**: Create git tag
6. **Build Artifacts**: Build binaries for all platforms
7. **GitHub Release**: Create release with binaries and changelog

## Getting Help

- **Questions**: Open a GitHub Discussion
- **Issues**: Report bugs via GitHub Issues
- **Chat**: Join our community chat (if available)
- **Documentation**: Check the Wiki for detailed guides

## Recognition

Contributors will be recognized in:
- README.md acknowledgments
- Release notes
- GitHub contributors page
- Annual contributor report

Thank you for contributing to Waterlogger!