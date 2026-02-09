package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	// Provider settings
	Provider         string
	CloudflareToken  string
	CloudflareEmail  string
	CloudflareAccountID string

	// AWS settings
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string

	// Application settings
	Timeout    int
	Concurrent int
	Threshold  int
	Output     string
	Verbose    bool

	// Filter settings
	Zone        string
	ExpiringIn  int
	Domains     []string
}

// Load loads configuration from environment variables and config file
func Load() (*Config, error) {
	// Set defaults
	viper.SetDefault("timeout", 10)
	viper.SetDefault("concurrent", 10)
	viper.SetDefault("threshold", 30)
	viper.SetDefault("output", "table")
	viper.SetDefault("provider", "cloudflare")
	viper.SetDefault("aws_region", "us-east-1")

	// Bind environment variables
	viper.SetEnvPrefix("SSL_CHECK")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Also bind Cloudflare and AWS specific env vars
	viper.BindEnv("cloudflare_token", "CLOUDFLARE_API_TOKEN")
	viper.BindEnv("cloudflare_email", "CLOUDFLARE_EMAIL")
	viper.BindEnv("cloudflare_account_id", "CLOUDFLARE_ACCOUNT_ID")
	viper.BindEnv("aws_access_key_id", "AWS_ACCESS_KEY_ID")
	viper.BindEnv("aws_secret_access_key", "AWS_SECRET_ACCESS_KEY")
	viper.BindEnv("aws_region", "AWS_REGION")

	// Try to load config file from multiple locations
	viper.SetConfigName("sslcheckdomain")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("$HOME")

	// Read config file if it exists (ignore error if not found)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		if err := viper.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("error reading .env file: %w", err)
		}
	}

	cfg := &Config{
		Provider:            viper.GetString("provider"),
		CloudflareToken:     viper.GetString("cloudflare_token"),
		CloudflareEmail:     viper.GetString("cloudflare_email"),
		CloudflareAccountID: viper.GetString("cloudflare_account_id"),
		AWSAccessKeyID:      viper.GetString("aws_access_key_id"),
		AWSSecretAccessKey:  viper.GetString("aws_secret_access_key"),
		AWSRegion:           viper.GetString("aws_region"),
		Timeout:             viper.GetInt("timeout"),
		Concurrent:          viper.GetInt("concurrent"),
		Threshold:           viper.GetInt("threshold"),
		Output:              viper.GetString("output"),
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	switch c.Provider {
	case "cloudflare":
		if c.CloudflareToken == "" {
			return fmt.Errorf("CLOUDFLARE_API_TOKEN is required for cloudflare provider")
		}
	case "route53":
		if c.AWSAccessKeyID == "" || c.AWSSecretAccessKey == "" {
			return fmt.Errorf("AWS credentials are required for route53 provider")
		}
	default:
		return fmt.Errorf("unsupported provider: %s (supported: cloudflare, route53)", c.Provider)
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	if c.Concurrent <= 0 {
		return fmt.Errorf("concurrent must be greater than 0")
	}

	if c.Threshold < 0 {
		return fmt.Errorf("threshold must be non-negative")
	}

	validOutputs := map[string]bool{
		"table":      true,
		"json":       true,
		"prometheus": true,
	}

	if !validOutputs[c.Output] {
		return fmt.Errorf("invalid output format: %s (valid: table, json, prometheus)", c.Output)
	}

	return nil
}
