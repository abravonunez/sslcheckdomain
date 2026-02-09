package checker

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"sslcheckdomain/pkg/models"
)

// SSLChecker checks SSL certificates for domains
type SSLChecker struct {
	timeout    time.Duration
	concurrent int
}

// New creates a new SSL checker
func New(timeout time.Duration, concurrent int) *SSLChecker {
	return &SSLChecker{
		timeout:    timeout,
		concurrent: concurrent,
	}
}

// CheckDomains checks SSL certificates for multiple domains concurrently
func (c *SSLChecker) CheckDomains(ctx context.Context, domains []string, threshold int) ([]models.Certificate, error) {
	if len(domains) == 0 {
		return nil, fmt.Errorf("no domains to check")
	}

	// Create channels for work distribution
	jobs := make(chan string, len(domains))
	results := make(chan models.Certificate, len(domains))

	// Create worker pool
	var wg sync.WaitGroup
	for i := 0; i < c.concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range jobs {
				cert := c.checkDomain(ctx, domain, threshold)
				results <- cert
			}
		}()
	}

	// Send jobs
	for _, domain := range domains {
		jobs <- domain
	}
	close(jobs)

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	certificates := make([]models.Certificate, 0, len(domains))
	for cert := range results {
		certificates = append(certificates, cert)
	}

	return certificates, nil
}

// checkDomain checks SSL certificate for a single domain
func (c *SSLChecker) checkDomain(ctx context.Context, domain string, threshold int) models.Certificate {
	cert := models.Certificate{
		Domain: domain,
	}

	// Create context with timeout
	checkCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Connect to the domain
	dialer := &net.Dialer{
		Timeout: c.timeout,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", domain+":443", &tls.Config{
		ServerName: domain,
		MinVersion: tls.VersionTLS12,
	})

	if err != nil {
		cert.Error = fmt.Errorf("failed to connect: %w", err)
		cert.DetermineStatus(threshold)
		return cert
	}
	defer conn.Close()

	// Check if context was cancelled
	select {
	case <-checkCtx.Done():
		cert.Error = fmt.Errorf("check timeout")
		cert.DetermineStatus(threshold)
		return cert
	default:
	}

	// Get certificate information
	if len(conn.ConnectionState().PeerCertificates) == 0 {
		cert.Error = fmt.Errorf("no certificate found")
		cert.DetermineStatus(threshold)
		return cert
	}

	peerCert := conn.ConnectionState().PeerCertificates[0]

	cert.ExpiresAt = peerCert.NotAfter
	cert.IssuedAt = peerCert.NotBefore
	cert.Issuer = peerCert.Issuer.CommonName
	cert.Subject = peerCert.Subject.CommonName
	cert.SerialNumber = peerCert.SerialNumber.String()

	// Determine status
	cert.DetermineStatus(threshold)

	return cert
}

// CheckDomain checks SSL certificate for a single domain (public method)
func (c *SSLChecker) CheckDomain(ctx context.Context, domain string, threshold int) models.Certificate {
	return c.checkDomain(ctx, domain, threshold)
}
