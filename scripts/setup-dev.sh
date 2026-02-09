#!/bin/bash
#
# Development environment setup script
#

set -euo pipefail

echo "Setting up sslcheckdomain development environment..."

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Found Go version: ${GO_VERSION}"

# Check Make installation
if ! command -v make &> /dev/null; then
    echo "Warning: Make is not installed. Some build commands may not work."
else
    echo "✓ Found Make"
fi

# Download dependencies
echo ""
echo "Downloading Go dependencies..."
go mod download
go mod verify
echo "✓ Dependencies downloaded"

# Install development tools (optional)
echo ""
echo "Installing development tools..."

# golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    echo "✓ golangci-lint installed"
else
    echo "✓ golangci-lint already installed"
fi

# goimports
if ! command -v goimports &> /dev/null; then
    echo "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
    echo "✓ goimports installed"
else
    echo "✓ goimports already installed"
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo ""
    echo "Creating .env file from template..."
    cp .env.example .env
    echo "✓ .env file created. Please edit it with your credentials."
else
    echo "✓ .env file already exists"
fi

# Run tests
echo ""
echo "Running tests..."
if make test; then
    echo "✓ All tests passed"
else
    echo "⚠️  Some tests failed. Please review the output above."
fi

# Build the project
echo ""
echo "Building project..."
if make build; then
    echo "✓ Build successful: bin/sslcheckdomain"
else
    echo "✗ Build failed. Please review the output above."
    exit 1
fi

echo ""
echo "========================================="
echo "✓ Development environment setup complete!"
echo "========================================="
echo ""
echo "Next steps:"
echo "  1. Edit .env with your Cloudflare API token"
echo "  2. Run: make run"
echo "  3. See README.md for more information"
echo ""
