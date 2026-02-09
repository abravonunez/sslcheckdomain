package models

import (
	"time"
)

// CertificateStatus represents the status of a certificate
type CertificateStatus string

const (
	StatusExpired CertificateStatus = "expired"
	StatusWarning CertificateStatus = "warning"
	StatusOK      CertificateStatus = "ok"
	StatusError   CertificateStatus = "error"
)

// Certificate represents SSL certificate information
type Certificate struct {
	Domain      string            `json:"domain"`
	Status      CertificateStatus `json:"status"`
	ExpiresAt   time.Time         `json:"expires_at"`
	IssuedAt    time.Time         `json:"issued_at"`
	Issuer      string            `json:"issuer"`
	Subject     string            `json:"subject"`
	DaysLeft    int               `json:"days_left"`
	SerialNumber string           `json:"serial_number"`
	Error       error             `json:"error,omitempty"`
}

// CertificateReport represents a collection of certificate checks
type CertificateReport struct {
	Timestamp    time.Time      `json:"timestamp"`
	TotalDomains int            `json:"total_domains"`
	Summary      ReportSummary  `json:"summary"`
	Certificates []Certificate  `json:"certificates"`
}

// ReportSummary provides aggregated statistics
type ReportSummary struct {
	Expired int `json:"expired"`
	Warning int `json:"warning"`
	OK      int `json:"ok"`
	Error   int `json:"error"`
}

// DaysUntilExpiration calculates days left until expiration
func (c *Certificate) DaysUntilExpiration() int {
	return int(time.Until(c.ExpiresAt).Hours() / 24)
}

// DetermineStatus determines the status based on days left and threshold
func (c *Certificate) DetermineStatus(warningThreshold int) {
	if c.Error != nil {
		c.Status = StatusError
		return
	}

	c.DaysLeft = c.DaysUntilExpiration()

	switch {
	case c.DaysLeft < 0:
		c.Status = StatusExpired
	case c.DaysLeft <= warningThreshold:
		c.Status = StatusWarning
	default:
		c.Status = StatusOK
	}
}

// IsHealthy returns true if certificate is OK
func (c *Certificate) IsHealthy() bool {
	return c.Status == StatusOK
}

// NeedsAttention returns true if certificate needs attention
func (c *Certificate) NeedsAttention() bool {
	return c.Status == StatusExpired || c.Status == StatusWarning
}
