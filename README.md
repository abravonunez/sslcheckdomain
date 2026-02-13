# SSL Check Domain

A robust CLI tool written in Go for monitoring SSL certificate expiration across multiple domains managed in DNS providers (Cloudflare, Route53, etc.).

## Overview

`sslcheckdomain` is designed for SRE and DevOps teams managing multiple domains. It automatically discovers domains from your DNS provider and checks their SSL certificate expiration status, displaying results in an easy-to-read table sorted by expiration date.

## Features

- **Multi-Provider Support**: Cloudflare, AWS Route53 (extensible architecture)
- **Automatic Domain Discovery**: Fetches all domains from your DNS provider
- **Concurrent Checks**: Fast parallel SSL certificate verification
- **Smart Sorting**: Displays results sorted by expiration date (earliest first)
- **Rich Output**: Colored terminal output with table formatting
- **Configurable**: Environment variables and CLI flags support
- **Portable**: Single binary distribution for Linux, macOS, and Windows
- **SRE-Friendly**: Exit codes, metrics export, and integration-ready

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap abravonunez/tap
brew install sslcheckdomain
```

### Linux Packages

**Debian/Ubuntu (.deb):**
```bash
wget https://github.com/abravonunez/sslcheckdomain/releases/latest/download/sslcheckdomain_<version>_amd64.deb
sudo dpkg -i sslcheckdomain_<version>_amd64.deb
```

**RedHat/Fedora (.rpm):**
```bash
wget https://github.com/abravonunez/sslcheckdomain/releases/latest/download/sslcheckdomain_<version>_amd64.rpm
sudo rpm -i sslcheckdomain_<version>_amd64.rpm
```

### Binary Release

Download the latest release for your platform:

```bash
# macOS (Apple Silicon)
curl -L https://github.com/abravonunez/sslcheckdomain/releases/latest/download/sslcheckdomain-darwin-arm64 -o sslcheckdomain
chmod +x sslcheckdomain
sudo mv sslcheckdomain /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/yourusername/sslcheckdomain/releases/latest/download/sslcheckdomain-linux-amd64 -o sslcheckdomain
chmod +x sslcheckdomain
sudo mv sslcheckdomain /usr/local/bin/

# Windows
# Download from releases page
```

### Build from Source

```bash
git clone https://github.com/yourusername/sslcheckdomain.git
cd sslcheckdomain
make build
```

## Quick Start

### 1. Configure Environment Variables

```bash
# Cloudflare
export CLOUDFLARE_API_TOKEN="your-api-token-here"
export CLOUDFLARE_ACCOUNT_ID="your-account-id"  # Optional

# OR use .env file
cp .env.example .env
# Edit .env with your credentials
```

### 2. Run Basic Check

```bash
# Check all domains in your Cloudflare account
sslcheckdomain

# Check specific domains
sslcheckdomain example.com api.example.com

# Filter by zone/domain
sslcheckdomain --zone example.com

# Show only expiring soon (default: 30 days)
sslcheckdomain --expiring-in 30
```

## Usage

```
sslcheckdomain [flags] [domain1 domain2 ...]

Flags:
  -p, --provider string       DNS provider (cloudflare, route53) (default "cloudflare")
  -z, --zone string          Filter by specific zone/domain
  -e, --expiring-in int      Show only certs expiring in N days (default: show all)
  -t, --threshold int        Warning threshold in days (default: 30)
  -o, --output string        Output format (table, json, prometheus) (default "table")
  -c, --concurrent int       Number of concurrent checks (default: 10)
  -v, --verbose             Verbose output
      --timeout int         HTTP timeout in seconds (default: 10)
      --version             Show version information
  -h, --help                Show help
```

### Examples

```bash
# Check all domains, show only those expiring in 7 days
sslcheckdomain --expiring-in 7

# Check with JSON output (for automation)
sslcheckdomain --output json

# Export metrics in Prometheus format
sslcheckdomain --output prometheus > metrics.txt

# Verbose mode with custom concurrency
sslcheckdomain -v --concurrent 20

# Check specific zone only
sslcheckdomain --zone mycompany.com
```

## Configuration

### Environment Variables

```bash
# Cloudflare
CLOUDFLARE_API_TOKEN=        # Required: Your Cloudflare API Token
CLOUDFLARE_ACCOUNT_ID=       # Optional: Account ID for multi-account setups
CLOUDFLARE_EMAIL=            # Optional: Email (for legacy API key auth)

# AWS Route53
AWS_ACCESS_KEY_ID=           # AWS credentials
AWS_SECRET_ACCESS_KEY=
AWS_REGION=                  # Default: us-east-1

# General
SSL_CHECK_TIMEOUT=10         # HTTP timeout in seconds
SSL_CHECK_CONCURRENT=10      # Number of concurrent checks
SSL_CHECK_THRESHOLD=30       # Warning threshold in days
```

### Configuration File (Optional)

Create `~/.sslcheckdomain.yaml`:

```yaml
provider: cloudflare
timeout: 10
concurrent: 15
threshold: 30
output: table
```

## Output Examples

### Table Format (Default)

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SSL Certificate Expiration Report                      ║
╠════════════════════════════════╦════════════╦═══════════╦═════════════════════╣
║ Domain                         ║ Status     ║ Days Left ║ Expires             ║
╠════════════════════════════════╬════════════╬═══════════╬═════════════════════╣
║ expired.example.com            ║ ❌ EXPIRED ║      -5   ║ 2025-12-13 10:30:00 ║
║ critical.example.com           ║ ⚠️  WARN   ║       3   ║ 2025-12-21 14:22:15 ║
║ soon.example.com               ║ ⚠️  WARN   ║      15   ║ 2026-01-02 08:15:30 ║
║ example.com                    ║ ✅ OK      ║      45   ║ 2026-02-01 12:00:00 ║
║ api.example.com                ║ ✅ OK      ║      78   ║ 2026-03-06 16:45:22 ║
╚════════════════════════════════╩════════════╩═══════════╩═════════════════════╝

Summary: 5 domains checked, 1 expired, 2 warnings, 2 ok
```

### JSON Format

```json
{
  "timestamp": "2025-12-18T10:00:00Z",
  "total_domains": 5,
  "summary": {
    "expired": 1,
    "warning": 2,
    "ok": 2
  },
  "domains": [
    {
      "domain": "expired.example.com",
      "status": "expired",
      "days_left": -5,
      "expires_at": "2025-12-13T10:30:00Z",
      "issuer": "Let's Encrypt"
    }
  ]
}
```

### Prometheus Format

```
# HELP ssl_certificate_expiry_days Days until SSL certificate expiration
# TYPE ssl_certificate_expiry_days gauge
ssl_certificate_expiry_days{domain="example.com",issuer="Let's Encrypt"} 45
ssl_certificate_expiry_days{domain="api.example.com",issuer="Let's Encrypt"} 78
ssl_certificate_expiry_days{domain="critical.example.com",issuer="DigiCert"} 3

# HELP ssl_certificate_status SSL certificate status (0=expired, 1=warning, 2=ok)
# TYPE ssl_certificate_status gauge
ssl_certificate_status{domain="example.com"} 2
ssl_certificate_status{domain="critical.example.com"} 1
```

## SRE Integration

### Monitoring & Alerting

```bash
# Cron job example (daily check)
0 9 * * * /usr/local/bin/sslcheckdomain --output json > /var/log/ssl-check.json

# Alert on critical certificates
sslcheckdomain --expiring-in 7 --output json | jq -e '.summary.expired + .summary.warning == 0' || notify-slack

# Prometheus metrics export
*/5 * * * * /usr/local/bin/sslcheckdomain --output prometheus > /var/lib/node_exporter/ssl_certs.prom
```

### Exit Codes

- `0`: All certificates OK
- `1`: Warning threshold reached (some certificates expiring soon)
- `2`: Critical (one or more certificates expired)
- `3`: Error (API failure, network issues, etc.)

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Check SSL Certificates
  run: |
    ./sslcheckdomain --expiring-in 14 --output json
  env:
    CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
```

## API Token Setup

### Cloudflare API Token

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com/profile/api-tokens)
2. Click "Create Token"
3. Use "Read all resources" template or create custom with:
   - Permissions: `Zone:Read`, `DNS:Read`
   - Zone Resources: `Include All zones`
4. Copy the token

### AWS Route53 (IAM)

Create IAM user with policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "route53:ListHostedZones",
        "route53:ListResourceRecordSets"
      ],
      "Resource": "*"
    }
  ]
}
```

## Development

### Prerequisites

- Go 1.21+
- Make

### Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint

# Install locally
make install
```

### Project Structure

```
sslcheckdomain/
├── cmd/
│   └── sslcheckdomain/
│       └── main.go              # CLI entrypoint
├── internal/
│   ├── provider/
│   │   ├── cloudflare/          # Cloudflare DNS client
│   │   ├── route53/             # AWS Route53 client
│   │   └── provider.go          # Provider interface
│   ├── checker/
│   │   └── ssl.go               # SSL certificate checker
│   ├── output/
│   │   ├── table.go             # Table formatter
│   │   ├── json.go              # JSON formatter
│   │   └── prometheus.go        # Prometheus exporter
│   └── config/
│       └── config.go            # Configuration management
├── pkg/
│   └── models/
│       └── certificate.go       # Data models
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── CLAUDE.md                    # Project instructions
├── .env.example
├── .gitignore
└── LICENSE
```

## Roadmap

- [x] Cloudflare provider support
- [ ] AWS Route53 provider support
- [ ] Google Cloud DNS provider support
- [ ] Azure DNS provider support
- [ ] Slack/Discord webhook notifications
- [ ] PagerDuty integration
- [ ] HTTP API mode
- [ ] TLS/SSL detailed analysis (cipher suites, vulnerabilities)
- [ ] Historical tracking database
- [ ] Web UI dashboard

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Support

- Issues: [GitHub Issues](https://github.com/yourusername/sslcheckdomain/issues)
- Documentation: [Wiki](https://github.com/yourusername/sslcheckdomain/wiki)

## Credits

Built with:
- [Cloudflare Go SDK](https://github.com/cloudflare/cloudflare-go)
- [Cobra CLI](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
- [PrettyTable](https://github.com/jedib0t/go-pretty)

---

**Sources:**
- [Cloudflare SSL/TLS Documentation](https://developers.cloudflare.com/ssl/)
- [sslcheck GitHub](https://github.com/i04n/sslcheck)
- [checkssl.org](https://www.checkssl.org/)
- [testssl.sh](https://testssl.sh/)
