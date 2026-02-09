package output

import (
	"encoding/json"
	"fmt"

	"sslcheckdomain/pkg/models"
)

// JSONFormatter formats certificate report as JSON
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format formats the certificate report as JSON
func (f *JSONFormatter) Format(report *models.CertificateReport) error {
	output, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(output))
	return nil
}
