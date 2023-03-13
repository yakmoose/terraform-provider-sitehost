// Package helper provides the functions to work with SiteHost API.
package helper

import (
	"context"
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sitehostnz/gosh/pkg/api"
	"github.com/sitehostnz/gosh/pkg/api/job"
)

const (
	// JobStatusPending is the status for a pending job.
	JobStatusPending = "Pending"
	// JobStatusCompleted is the status for a completed job.
	JobStatusCompleted = "Completed"
	// JobStatusFailed is the status for a failed job.
	JobStatusFailed = "Failed"
	// JobRequestDelay is the time wait to send a new request to check the job status.
	JobRequestDelay = 10 * time.Second
	// JobRequestTimeout is the time to wait before timeout.
	JobRequestTimeout = 60 * time.Minute
	// JobRequestMinTimeout is the minimum time to wait before refreshes.
	JobRequestMinTimeout = 3 * time.Second
	// JobRequestNotFoundChecks is the number of times to allow not found.
	JobRequestNotFoundChecks = 60
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

// WaitForAction is a function to check the Job status in a refresh function.
func WaitForAction(client *api.Client, request job.GetRequest) error {
	var (
		pending   = JobStatusPending
		target    = JobStatusCompleted
		ctx       = context.Background()
		refreshFn = func() (result any, state string, err error) {
			svc := job.New(client)

			j, err := svc.Get(ctx, request)
			if err != nil {
				return nil, "", err
			}

			if !j.Status {
				return nil, "", errors.New("An error has occurred with a message: " + j.Msg)
			}

			switch j.Return.State {
			case JobStatusFailed:
				return j, JobStatusFailed, nil
			case target:
				return j, target, nil
			default:
				return j, pending, nil
			}
		}
	)

	_, err := (&resource.StateChangeConf{
		Pending:        []string{pending},
		Refresh:        refreshFn,
		Target:         []string{target},
		Delay:          JobRequestDelay,
		Timeout:        JobRequestTimeout,
		MinTimeout:     JobRequestMinTimeout,
		NotFoundChecks: JobRequestNotFoundChecks,
	}).WaitForStateContext(ctx)

	return err
}
