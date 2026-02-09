# Quick Start Guide

Get started with sslcheckdomain in 5 minutes.

## Prerequisites

- Go 1.21+ installed
- Cloudflare API token (or AWS credentials for Route53)

## Step 1: Get Cloudflare API Token

1. Go to https://dash.cloudflare.com/profile/api-tokens
2. Click "Create Token"
3. Use "Read all resources" template or create custom with:
   - Permissions: `Zone:Read`, `DNS:Read`
   - Zone Resources: `Include All zones`
4. Copy the token

## Step 2: Setup

```bash
# Run the setup script
./scripts/setup-dev.sh

# OR manually:
go mod download
make build
```

## Step 3: Configure

```bash
# Edit .env file
vim .env

# Add your token:
CLOUDFLARE_API_TOKEN=your-token-here
```

## Step 4: Run

```bash
# Check all domains
./bin/sslcheckdomain

# Check specific domains
./bin/sslcheckdomain example.com api.example.com

# Show only expiring soon
./bin/sslcheckdomain --expiring-in 7

# JSON output
./bin/sslcheckdomain --output json
```

## Example Output

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SSL Certificate Expiration Report                      ║
╠════════════════════════════════╦════════════╦═══════════╦═════════════════════╣
║ Domain                         ║ Status     ║ Days Left ║ Expires             ║
╠════════════════════════════════╬════════════╬═══════════╬═════════════════════╣
║ critical.example.com           ║ ⚠️  WARN   ║       3   ║ 2025-12-21 14:22:15 ║
║ soon.example.com               ║ ⚠️  WARN   ║      15   ║ 2026-01-02 08:15:30 ║
║ example.com                    ║ ✅ OK      ║      45   ║ 2026-02-01 12:00:00 ║
╚════════════════════════════════╩════════════╩═══════════╩═════════════════════╝
```

## Common Use Cases

### Daily Monitoring (Cron)

```bash
# Add to crontab
0 9 * * * /usr/local/bin/sslcheckdomain --output json > /var/log/ssl-check.json
```

### Prometheus Metrics

```bash
# Export metrics every 5 minutes
*/5 * * * * /usr/local/bin/sslcheckdomain --output prometheus > /var/lib/node_exporter/ssl.prom
```

### Alert on Expiring Certificates

```bash
#!/bin/bash
RESULT=$(sslcheckdomain --expiring-in 7 --output json)
EXPIRING=$(echo $RESULT | jq '.summary.expired + .summary.warning')

if [ $EXPIRING -gt 0 ]; then
    echo "Alert: $EXPIRING certificates expiring soon!"
    # Send notification (Slack, PagerDuty, etc.)
fi
```

## Troubleshooting

### "CLOUDFLARE_API_TOKEN is required"
- Make sure .env file exists and contains your token
- Or export: `export CLOUDFLARE_API_TOKEN=your-token`

### "failed to connect"
- Check domain is accessible on port 443
- Verify DNS resolution
- Check firewall rules

### "timeout"
- Increase timeout: `--timeout 30`
- Check network connectivity
- Reduce concurrent checks: `--concurrent 5`

## Next Steps

- Read the full [README.md](README.md) for all features
- Set up automated monitoring (see examples/)
- Integrate with your alerting system
- Configure for multiple DNS providers

## Help

```bash
# Show all options
./bin/sslcheckdomain --help

# Show version
./bin/sslcheckdomain --version
```

For more information, see [README.md](README.md) or [CONTRIBUTING.md](CONTRIBUTING.md).
