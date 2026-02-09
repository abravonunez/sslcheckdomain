package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"sslcheckdomain/internal/checker"
	"sslcheckdomain/internal/config"
	"sslcheckdomain/internal/output"
	"sslcheckdomain/internal/provider"
	"sslcheckdomain/internal/provider/cloudflare"
	"sslcheckdomain/pkg/models"
)

var (
	// Version is set during build
	Version   = "dev"
	BuildTime = "unknown"

	// CLI flags
	providerFlag    string
	zoneFlag        string
	expiringInFlag  int
	thresholdFlag   int
	outputFlag      string
	concurrentFlag  int
	verboseFlag     bool
	timeoutFlag     int
	versionFlag     bool
	testDomainFlag  string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(3)
	}
}

var rootCmd = &cobra.Command{
	Use:   "sslcheckdomain [domain1 domain2 ...]",
	Short: "Check SSL certificate expiration for multiple domains",
	Long: `sslcheckdomain is a CLI tool for monitoring SSL certificate expiration
across multiple domains managed in DNS providers (Cloudflare, Route53, etc.).

It automatically discovers domains from your DNS provider and checks their
SSL certificate expiration status, displaying results sorted by expiration date.`,
	Example: `  # Test a single domain (no provider needed)
  sslcheckdomain --test example.com
  sslcheckdomain -d example.com

  # Check all domains in your Cloudflare account
  sslcheckdomain

  # Check specific domains
  sslcheckdomain example.com api.example.com

  # Show only certificates expiring in 7 days
  sslcheckdomain --expiring-in 7

  # Output as JSON
  sslcheckdomain --output json

  # Check specific zone
  sslcheckdomain --zone example.com`,
	RunE: run,
}

func init() {
	rootCmd.Flags().StringVarP(&providerFlag, "provider", "p", "", "DNS provider (cloudflare, route53)")
	rootCmd.Flags().StringVarP(&zoneFlag, "zone", "z", "", "Filter by specific zone/domain")
	rootCmd.Flags().IntVarP(&expiringInFlag, "expiring-in", "e", 0, "Show only certs expiring in N days (0 = show all)")
	rootCmd.Flags().IntVarP(&thresholdFlag, "threshold", "t", 0, "Warning threshold in days (default from config)")
	rootCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output format (table, json, prometheus)")
	rootCmd.Flags().IntVarP(&concurrentFlag, "concurrent", "c", 0, "Number of concurrent checks (default from config)")
	rootCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().IntVar(&timeoutFlag, "timeout", 0, "HTTP timeout in seconds (default from config)")
	rootCmd.Flags().BoolVar(&versionFlag, "version", false, "Show version information")
	rootCmd.Flags().StringVarP(&testDomainFlag, "test", "d", "", "Test a single domain (bypasses provider lookup)")
}

func run(cmd *cobra.Command, args []string) error {
	// Show version if requested
	if versionFlag {
		fmt.Printf("sslcheckdomain version %s (built %s)\n", Version, BuildTime)
		return nil
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override config with CLI flags if provided
	if providerFlag != "" {
		cfg.Provider = providerFlag
	}
	if zoneFlag != "" {
		cfg.Zone = zoneFlag
	}
	if expiringInFlag > 0 {
		cfg.ExpiringIn = expiringInFlag
	}
	if thresholdFlag > 0 {
		cfg.Threshold = thresholdFlag
	}
	if outputFlag != "" {
		cfg.Output = outputFlag
	}
	if concurrentFlag > 0 {
		cfg.Concurrent = concurrentFlag
	}
	if timeoutFlag > 0 {
		cfg.Timeout = timeoutFlag
	}
	cfg.Verbose = verboseFlag
	cfg.Domains = args

	// Validate configuration (skip provider validation if using --test flag)
	if testDomainFlag == "" {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}
	} else {
		// Only validate non-provider settings when using --test
		if cfg.Timeout <= 0 {
			return fmt.Errorf("timeout must be greater than 0")
		}
		if cfg.Concurrent <= 0 {
			return fmt.Errorf("concurrent must be greater than 0")
		}
		if cfg.Threshold < 0 {
			return fmt.Errorf("threshold must be non-negative")
		}
		validOutputs := map[string]bool{
			"table":      true,
			"json":       true,
			"prometheus": true,
		}
		if !validOutputs[cfg.Output] {
			return fmt.Errorf("invalid output format: %s (valid: table, json, prometheus)", cfg.Output)
		}
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Configuration loaded:\n")
		fmt.Fprintf(os.Stderr, "  Provider: %s\n", cfg.Provider)
		fmt.Fprintf(os.Stderr, "  Timeout: %ds\n", cfg.Timeout)
		fmt.Fprintf(os.Stderr, "  Concurrent: %d\n", cfg.Concurrent)
		fmt.Fprintf(os.Stderr, "  Threshold: %d days\n", cfg.Threshold)
		fmt.Fprintf(os.Stderr, "  Output: %s\n", cfg.Output)
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Create context
	ctx := context.Background()

	// Get domains to check
	var domains []string
	if !cfg.Verbose && testDomainFlag == "" {
		// Show spinner only if not in verbose mode and not testing a single domain
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Fetching domains from provider..."
		s.Start()
		domains, err = getDomains(ctx, cfg)
		s.Stop()
	} else {
		domains, err = getDomains(ctx, cfg)
	}

	if err != nil {
		return fmt.Errorf("failed to get domains: %w", err)
	}

	if len(domains) == 0 {
		return fmt.Errorf("no domains to check")
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Found %d domains to check\n", len(domains))
	}

	// Check SSL certificates
	sslChecker := checker.New(time.Duration(cfg.Timeout)*time.Second, cfg.Concurrent)

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Checking SSL certificates...\n")
	}

	var certificates []models.Certificate
	if !cfg.Verbose {
		// Show spinner only if not in verbose mode
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" Checking SSL certificates for %d domains...", len(domains))
		s.Start()
		certificates, err = sslChecker.CheckDomains(ctx, domains, cfg.Threshold)
		s.Stop()
	} else {
		certificates, err = sslChecker.CheckDomains(ctx, domains, cfg.Threshold)
	}

	if err != nil {
		return fmt.Errorf("failed to check certificates: %w", err)
	}

	// Filter by expiring-in if specified
	if cfg.ExpiringIn > 0 {
		filtered := make([]models.Certificate, 0)
		for _, cert := range certificates {
			if cert.DaysLeft <= cfg.ExpiringIn {
				filtered = append(filtered, cert)
			}
		}
		certificates = filtered
	}

	// Sort by days left (ascending)
	sort.Slice(certificates, func(i, j int) bool {
		return certificates[i].DaysLeft < certificates[j].DaysLeft
	})

	// Create report
	report := createReport(certificates)

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Format and display output
	formatter, err := output.GetFormatter(cfg.Output)
	if err != nil {
		return fmt.Errorf("failed to create formatter: %w", err)
	}

	if err := formatter.Format(report); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Set exit code based on results
	exitCode := getExitCode(report)
	os.Exit(exitCode)

	return nil
}

func getDomains(ctx context.Context, cfg *config.Config) ([]string, error) {
	// If test domain flag is provided, use it (highest priority)
	if testDomainFlag != "" {
		return []string{testDomainFlag}, nil
	}

	// If specific domains provided via CLI, use those
	if len(cfg.Domains) > 0 {
		return cfg.Domains, nil
	}

	// Otherwise, fetch from DNS provider
	var dnsProvider provider.DNSProvider
	var err error

	switch cfg.Provider {
	case "cloudflare":
		dnsProvider, err = cloudflare.New(cfg.CloudflareToken)
		if err != nil {
			return nil, fmt.Errorf("failed to create cloudflare provider: %w", err)
		}
	case "route53":
		return nil, fmt.Errorf("route53 provider not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	// Get domains
	var domains []string
	if cfg.Zone != "" {
		domains, err = dnsProvider.GetDomainsByZone(ctx, cfg.Zone)
	} else {
		domains, err = dnsProvider.GetDomains(ctx)
	}

	if err != nil {
		return nil, err
	}

	return domains, nil
}

func createReport(certificates []models.Certificate) *models.CertificateReport {
	report := &models.CertificateReport{
		Timestamp:    time.Now(),
		TotalDomains: len(certificates),
		Certificates: certificates,
		Summary: models.ReportSummary{
			Expired: 0,
			Warning: 0,
			OK:      0,
			Error:   0,
		},
	}

	for _, cert := range certificates {
		switch cert.Status {
		case models.StatusExpired:
			report.Summary.Expired++
		case models.StatusWarning:
			report.Summary.Warning++
		case models.StatusOK:
			report.Summary.OK++
		case models.StatusError:
			report.Summary.Error++
		}
	}

	return report
}

func getExitCode(report *models.CertificateReport) int {
	if report.Summary.Expired > 0 {
		return 2 // Critical: one or more certificates expired
	}
	if report.Summary.Warning > 0 {
		return 1 // Warning: certificates expiring soon
	}
	if report.Summary.Error > 0 {
		return 3 // Error: API failure or network issues
	}
	return 0 // All OK
}
