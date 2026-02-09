package output

import (
	"fmt"

	"sslcheckdomain/pkg/models"
)

// PrometheusFormatter formats certificate report as Prometheus metrics
type PrometheusFormatter struct{}

// NewPrometheusFormatter creates a new Prometheus formatter
func NewPrometheusFormatter() *PrometheusFormatter {
	return &PrometheusFormatter{}
}

// Format formats the certificate report as Prometheus metrics
func (f *PrometheusFormatter) Format(report *models.CertificateReport) error {
	// Certificate expiry days metric
	fmt.Println("# HELP ssl_certificate_expiry_days Days until SSL certificate expiration")
	fmt.Println("# TYPE ssl_certificate_expiry_days gauge")

	for _, cert := range report.Certificates {
		if cert.Error == nil {
			fmt.Printf("ssl_certificate_expiry_days{domain=\"%s\",issuer=\"%s\",status=\"%s\"} %d\n",
				cert.Domain,
				cert.Issuer,
				cert.Status,
				cert.DaysLeft,
			)
		}
	}

	fmt.Println()

	// Certificate status metric
	fmt.Println("# HELP ssl_certificate_status SSL certificate status (0=expired, 1=warning, 2=ok, 3=error)")
	fmt.Println("# TYPE ssl_certificate_status gauge")

	for _, cert := range report.Certificates {
		statusValue := f.statusToValue(cert.Status)
		fmt.Printf("ssl_certificate_status{domain=\"%s\",issuer=\"%s\"} %d\n",
			cert.Domain,
			cert.Issuer,
			statusValue,
		)
	}

	fmt.Println()

	// Summary metrics
	fmt.Println("# HELP ssl_certificates_total Total number of certificates checked")
	fmt.Println("# TYPE ssl_certificates_total gauge")
	fmt.Printf("ssl_certificates_total %d\n", report.TotalDomains)

	fmt.Println()

	fmt.Println("# HELP ssl_certificates_expired Number of expired certificates")
	fmt.Println("# TYPE ssl_certificates_expired gauge")
	fmt.Printf("ssl_certificates_expired %d\n", report.Summary.Expired)

	fmt.Println()

	fmt.Println("# HELP ssl_certificates_warning Number of certificates with warnings")
	fmt.Println("# TYPE ssl_certificates_warning gauge")
	fmt.Printf("ssl_certificates_warning %d\n", report.Summary.Warning)

	fmt.Println()

	fmt.Println("# HELP ssl_certificates_ok Number of OK certificates")
	fmt.Println("# TYPE ssl_certificates_ok gauge")
	fmt.Printf("ssl_certificates_ok %d\n", report.Summary.OK)

	fmt.Println()

	fmt.Println("# HELP ssl_certificates_error Number of certificates with errors")
	fmt.Println("# TYPE ssl_certificates_error gauge")
	fmt.Printf("ssl_certificates_error %d\n", report.Summary.Error)

	return nil
}

// statusToValue converts status to numeric value
func (f *PrometheusFormatter) statusToValue(status models.CertificateStatus) int {
	switch status {
	case models.StatusExpired:
		return 0
	case models.StatusWarning:
		return 1
	case models.StatusOK:
		return 2
	case models.StatusError:
		return 3
	default:
		return 3
	}
}
