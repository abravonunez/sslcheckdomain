# Contributing to sslcheckdomain

Thank you for your interest in contributing to sslcheckdomain! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful and constructive in all interactions.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Make
- Git

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/yourusername/sslcheckdomain.git
cd sslcheckdomain

# Install dependencies
make mod-download

# Run tests
make test

# Build
make build
```

## Development Workflow

1. **Fork the repository** and create your branch from `main`

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Write clean, readable code
   - Follow Go best practices and idioms
   - Add tests for new functionality
   - Update documentation as needed

4. **Run tests and linting**
   ```bash
   make check
   ```

5. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

6. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request**

## Commit Message Format

Follow conventional commits:

```
<type>: <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat: add route53 provider support
fix: handle timeout errors in SSL checker
docs: update README with new examples
```

## Code Style

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` and `goimports` for formatting
- Run `golangci-lint` and fix all issues
- Write meaningful comments for exported functions
- Keep functions small and focused

### Testing

- Write unit tests for new functionality
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

Example:
```go
func TestCertificateStatus(t *testing.T) {
    tests := []struct {
        name      string
        daysLeft  int
        threshold int
        want      models.CertificateStatus
    }{
        {"expired", -5, 30, models.StatusExpired},
        {"warning", 15, 30, models.StatusWarning},
        {"ok", 45, 30, models.StatusOK},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Adding a New DNS Provider

1. Create a new package under `internal/provider/yourprovider/`
2. Implement the `DNSProvider` interface
3. Add configuration options to `internal/config/config.go`
4. Update main.go to register the provider
5. Add tests
6. Update documentation

Example structure:
```go
package yourprovider

type Provider struct {
    client *YourClient
}

func New(config Config) (*Provider, error) {
    // Implementation
}

func (p *Provider) GetDomains(ctx context.Context) ([]string, error) {
    // Implementation
}

func (p *Provider) GetDomainsByZone(ctx context.Context, zone string) ([]string, error) {
    // Implementation
}

func (p *Provider) Name() string {
    return "yourprovider"
}
```

## Adding a New Output Format

1. Create a new file under `internal/output/yourformat.go`
2. Implement the `Formatter` interface
3. Register in `output.go`
4. Add tests
5. Update documentation

## Project Structure

```
sslcheckdomain/
├── cmd/
│   └── sslcheckdomain/     # Main application
├── internal/               # Private application code
│   ├── checker/           # SSL certificate checker
│   ├── config/            # Configuration management
│   ├── output/            # Output formatters
│   └── provider/          # DNS provider implementations
├── pkg/                   # Public libraries
│   └── models/            # Data models
├── examples/              # Example scripts and configs
├── .github/               # GitHub workflows
└── Makefile              # Build automation
```

## Testing

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

### Run Specific Tests
```bash
go test -v ./internal/checker/...
```

## Building

### Local Build
```bash
make build
```

### Build for All Platforms
```bash
make build-all
```

### Create Release
```bash
make release
```

## Documentation

- Update README.md for user-facing changes
- Update CLAUDE.md for project context
- Add code comments for complex logic
- Update examples/ for new features

## Pull Request Process

1. Ensure all tests pass
2. Update documentation
3. Add yourself to contributors list
4. Request review from maintainers
5. Address review feedback
6. Squash commits if requested

## Questions?

- Open an issue for bugs or feature requests
- Join discussions in GitHub Discussions
- Contact maintainers via email (if urgent)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
