package checker

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"sslcheckdomain/pkg/models"
)

// SSLChecker checks SSL certificates for domains
type SSLChecker struct {
	timeout    time.Duration
	concurrent int
	verbose    bool
	logger     io.Writer
}

// New creates a new SSL checker
func New(timeout time.Duration, concurrent int) *SSLChecker {
	return &SSLChecker{
		timeout:    timeout,
		concurrent: concurrent,
		verbose:    false,
		logger:     nil,
	}
}

// WithVerbose enables verbose logging
func (c *SSLChecker) WithVerbose(verbose bool, logger io.Writer) *SSLChecker {
	c.verbose = verbose
	c.logger = logger
	return c
}

func (c *SSLChecker) logf(format string, args ...any) {
	if c.verbose && c.logger != nil {
		fmt.Fprintf(c.logger, format, args...)
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
	startTime := time.Now()
	cert := models.Certificate{
		Domain: domain,
	}

	c.logf("\n[DEBUG] Checking domain: %s\n", domain)

	// Create context with timeout
	checkCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Step 1: DNS Resolution
	c.logf("  [1/4] DNS Resolution...\n")
	resolver := &net.Resolver{}
	ips, err := resolver.LookupHost(ctx, domain)
	if err != nil {
		c.logf("  ❌ DNS resolution failed: %v\n", err)
		cert.Error = fmt.Errorf("DNS resolution failed: %w", err)
		cert.DetermineStatus(threshold)
		c.logf("  ⏱️  Total time: %v\n", time.Since(startTime))
		return cert
	}
	c.logf("  ✅ Resolved to IPs: %v\n", ips)

	// Step 2: TCP Connection
	c.logf("  [2/4] TCP Connection to %s:443...\n", domain)
	dialer := &net.Dialer{
		Timeout: c.timeout,
	}

	targetAddr := domain + ":443"
	c.logf("  → Connecting to: %s\n", targetAddr)

	tcpConnStart := time.Now()
	conn, err := tls.DialWithDialer(dialer, "tcp", targetAddr, &tls.Config{
		ServerName: domain,
		MinVersion: tls.VersionTLS12,
	})
	tcpDuration := time.Since(tcpConnStart)

	if err != nil {
		c.logf("  ❌ Connection failed after %v: %v\n", tcpDuration, err)
		cert.Error = fmt.Errorf("failed to connect: %w", err)
		cert.DetermineStatus(threshold)
		c.logf("  ⏱️  Total time: %v\n", time.Since(startTime))
		return cert
	}
	defer conn.Close()

	c.logf("  ✅ TCP+TLS handshake successful (%v)\n", tcpDuration)
	c.logf("  → Remote address: %s\n", conn.RemoteAddr())
	c.logf("  → Local address: %s\n", conn.LocalAddr())

	// Check if context was cancelled
	select {
	case <-checkCtx.Done():
		c.logf("  ❌ Context timeout\n")
		cert.Error = fmt.Errorf("check timeout")
		cert.DetermineStatus(threshold)
		c.logf("  ⏱️  Total time: %v\n", time.Since(startTime))
		return cert
	default:
	}

	// Step 3: TLS Information
	c.logf("  [3/4] TLS Connection State...\n")
	connState := conn.ConnectionState()
	c.logf("  → TLS Version: %s\n", tlsVersionString(connState.Version))
	c.logf("  → Cipher Suite: %s\n", tls.CipherSuiteName(connState.CipherSuite))
	c.logf("  → Server Name: %s\n", connState.ServerName)
	c.logf("  → Negotiated Protocol: %s\n", connState.NegotiatedProtocol)

	// Step 4: Certificate Information
	c.logf("  [4/4] Certificate Validation...\n")
	if len(connState.PeerCertificates) == 0 {
		c.logf("  ❌ No certificates found in chain\n")
		cert.Error = fmt.Errorf("no certificate found")
		cert.DetermineStatus(threshold)
		c.logf("  ⏱️  Total time: %v\n", time.Since(startTime))
		return cert
	}

	c.logf("  → Certificate chain length: %d\n", len(connState.PeerCertificates))
	peerCert := connState.PeerCertificates[0]

	cert.ExpiresAt = peerCert.NotAfter
	cert.IssuedAt = peerCert.NotBefore
	cert.Issuer = peerCert.Issuer.CommonName
	cert.Subject = peerCert.Subject.CommonName
	cert.SerialNumber = peerCert.SerialNumber.String()

	c.logf("  ✅ Certificate details:\n")
	c.logf("     - Subject: %s\n", cert.Subject)
	c.logf("     - Issuer: %s\n", cert.Issuer)
	c.logf("     - Serial: %s\n", cert.SerialNumber)
	c.logf("     - Valid from: %s\n", cert.IssuedAt.Format("2006-01-02 15:04:05"))
	c.logf("     - Valid until: %s\n", cert.ExpiresAt.Format("2006-01-02 15:04:05"))
	c.logf("     - DNS Names: %v\n", peerCert.DNSNames)

	// Determine status
	cert.DetermineStatus(threshold)

	totalDuration := time.Since(startTime)
	c.logf("  ✅ Check completed successfully\n")
	c.logf("  ⏱️  Total time: %v\n", totalDuration)

	return cert
}

// tlsVersionString converts TLS version number to string
func tlsVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (0x%04x)", version)
	}
}

// CheckDomain checks SSL certificate for a single domain (public method)
func (c *SSLChecker) CheckDomain(ctx context.Context, domain string, threshold int) models.Certificate {
	return c.checkDomain(ctx, domain, threshold)
}
