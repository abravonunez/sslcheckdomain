#!/bin/bash
#
# Example cron job for SSL certificate monitoring
# Add to crontab: 0 9 * * * /path/to/cron-job.sh
#

set -euo pipefail

# Configuration
SSLCHECK_BIN="/usr/local/bin/sslcheckdomain"
LOG_DIR="/var/log/sslcheck"
ALERT_WEBHOOK="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
EXPIRING_THRESHOLD=7

# Create log directory if it doesn't exist
mkdir -p "${LOG_DIR}"

# Set log file with timestamp
LOG_FILE="${LOG_DIR}/ssl-check-$(date +%Y%m%d).log"
JSON_FILE="${LOG_DIR}/ssl-check-$(date +%Y%m%d).json"

# Run check and save output
echo "=== SSL Certificate Check - $(date) ===" | tee -a "${LOG_FILE}"

# Check certificates
${SSLCHECK_BIN} --output json > "${JSON_FILE}" 2>> "${LOG_FILE}"

# Parse results
EXPIRED=$(jq '.summary.expired' "${JSON_FILE}")
WARNING=$(jq '.summary.warning' "${JSON_FILE}")

# Alert if there are issues
if [ "${EXPIRED}" -gt 0 ] || [ "${WARNING}" -gt 0 ]; then
    MESSAGE="🚨 SSL Certificate Alert: ${EXPIRED} expired, ${WARNING} expiring soon"

    # Send to Slack (if webhook configured)
    if [ -n "${ALERT_WEBHOOK}" ]; then
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"${MESSAGE}\"}" \
            "${ALERT_WEBHOOK}"
    fi

    # Log alert
    echo "ALERT: ${MESSAGE}" | tee -a "${LOG_FILE}"

    # Exit with error code
    exit 1
fi

echo "✓ All certificates OK" | tee -a "${LOG_FILE}"
exit 0
