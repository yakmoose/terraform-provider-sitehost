// Package helper provides the functions to work with SiteHost API.
package helper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/sitehostnz/gosh/pkg/api"
	"log"
	"net/url"
)

// Config is a wrapper to save the configuration connection from terraform.
type Config struct {
	APIKey           string
	ClientID         string
	APIEndpoint      string
	TerraformVersion string
}

// CombinedConfig is a struct with API wrapper and the Config.
type CombinedConfig struct {
	Client *api.Client
	Config *Config
}

// Client returns a new CombinedConfig instance.
func (c *Config) Client() (*CombinedConfig, diag.Diagnostics) {
	client := api.NewClient(c.APIKey, c.ClientID)

	client.UserAgent = "Terraform/" + c.TerraformVersion

	if c.APIEndpoint != "" {
		apiURL, err := url.Parse(c.APIEndpoint)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		client.BaseURL = apiURL
		log.Printf("[INFO] SiteHost Client configured for URL: %s", client.BaseURL.String())
	}

	return &CombinedConfig{
		Client: client,
		Config: c,
	}, nil
}
