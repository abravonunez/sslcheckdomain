#!/bin/bash
#
# Example script to export SSL metrics for Prometheus Node Exporter
# Run this script via cron: */5 * * * * /path/to/prometheus-exporter.sh
#

set -euo pipefail

SSLCHECK_BIN="/usr/local/bin/sslcheckdomain"
METRICS_DIR="/var/lib/node_exporter/textfile_collector"
METRICS_FILE="${METRICS_DIR}/ssl_certificates.prom"
TEMP_FILE="${METRICS_FILE}.tmp"

# Ensure metrics directory exists
mkdir -p "${METRICS_DIR}"

# Generate metrics
${SSLCHECK_BIN} --output prometheus > "${TEMP_FILE}"

# Atomic move to prevent partial reads
mv "${TEMP_FILE}" "${METRICS_FILE}"

# Set proper permissions
chmod 644 "${METRICS_FILE}"

echo "Metrics exported to ${METRICS_FILE}"
