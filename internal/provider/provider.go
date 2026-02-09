package provider

import (
	"context"
	"fmt"
)

// DNSProvider is the interface that all DNS providers must implement
type DNSProvider interface {
	// GetDomains retrieves all domains from the DNS provider
	GetDomains(ctx context.Context) ([]string, error)

	// GetDomainsByZone retrieves domains filtered by zone/parent domain
	GetDomainsByZone(ctx context.Context, zone string) ([]string, error)

	// Name returns the provider name
	Name() string
}

// ProviderFactory creates a DNS provider based on the provider type
type ProviderFactory struct {
	providers map[string]func() (DNSProvider, error)
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[string]func() (DNSProvider, error)),
	}
}

// Register registers a provider constructor
func (f *ProviderFactory) Register(name string, constructor func() (DNSProvider, error)) {
	f.providers[name] = constructor
}

// Create creates a provider instance
func (f *ProviderFactory) Create(name string) (DNSProvider, error) {
	constructor, ok := f.providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return constructor()
}

// AvailableProviders returns list of available providers
func (f *ProviderFactory) AvailableProviders() []string {
	providers := make([]string, 0, len(f.providers))
	for name := range f.providers {
		providers = append(providers, name)
	}
	return providers
}
