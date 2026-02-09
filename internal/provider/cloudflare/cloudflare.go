package cloudflare

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

// Provider implements the DNSProvider interface for Cloudflare
type Provider struct {
	client *cloudflare.API
}

// New creates a new Cloudflare provider
func New(apiToken string) (*Provider, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("cloudflare API token is required")
	}

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudflare client: %w", err)
	}

	return &Provider{
		client: api,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "cloudflare"
}

// GetDomains retrieves all domains from Cloudflare
func (p *Provider) GetDomains(ctx context.Context) ([]string, error) {
	zones, err := p.client.ListZones(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list zones: %w", err)
	}

	domains := make([]string, 0, len(zones))
	for _, zone := range zones {
		domains = append(domains, zone.Name)

		// Also get subdomains from DNS records
		subdomains, err := p.getSubdomains(ctx, zone.ID, zone.Name)
		if err != nil {
			// Log error but continue with other zones
			continue
		}
		domains = append(domains, subdomains...)
	}

	return domains, nil
}

// GetDomainsByZone retrieves domains filtered by zone
func (p *Provider) GetDomainsByZone(ctx context.Context, zoneName string) ([]string, error) {
	// Find the zone ID first
	zones, err := p.client.ListZones(ctx, zoneName)
	if err != nil {
		return nil, fmt.Errorf("failed to list zones: %w", err)
	}

	if len(zones) == 0 {
		return nil, fmt.Errorf("zone not found: %s", zoneName)
	}

	zone := zones[0]
	domains := []string{zone.Name}

	// Get subdomains
	subdomains, err := p.getSubdomains(ctx, zone.ID, zone.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get subdomains: %w", err)
	}

	domains = append(domains, subdomains...)
	return domains, nil
}

// getSubdomains retrieves all subdomains for a given zone
func (p *Provider) getSubdomains(ctx context.Context, zoneID, zoneName string) ([]string, error) {
	records, _, err := p.client.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to list DNS records: %w", err)
	}

	// Use map to deduplicate domains
	domainSet := make(map[string]bool)

	for _, record := range records {
		// Only consider A, AAAA, and CNAME records that point to external resources
		if record.Type == "A" || record.Type == "AAAA" || record.Type == "CNAME" {
			// Skip if it's the zone apex
			if record.Name == zoneName {
				continue
			}

			// Only include if it's a subdomain and not a wildcard
			if strings.HasSuffix(record.Name, "."+zoneName) && !strings.Contains(record.Name, "*") {
				domainSet[record.Name] = true
			}
		}
	}

	// Convert map to slice
	subdomains := make([]string, 0, len(domainSet))
	for domain := range domainSet {
		subdomains = append(subdomains, domain)
	}

	return subdomains, nil
}
