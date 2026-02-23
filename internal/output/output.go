package output

import (
	"fmt"

	"sslcheckdomain/pkg/models"
)

// Formatter is the interface for output formatters
type Formatter interface {
	Format(report *models.CertificateReport) error
}

// GetFormatter returns the appropriate formatter based on format string
func GetFormatter(format string, themeName string) (Formatter, error) {
	switch format {
	case "table":
		return NewTableFormatterWithTheme(themeName), nil
	case "json":
		return NewJSONFormatter(), nil
	case "prometheus":
		return NewPrometheusFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}
