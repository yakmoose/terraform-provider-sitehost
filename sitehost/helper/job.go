package helper

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/sitehostnz/gosh/pkg/api"
	"github.com/sitehostnz/gosh/pkg/api/job"
	"github.com/sitehostnz/gosh/pkg/models"
	"time"
)

const (
	// JobStatusPending is the status for a pending job.
	JobStatusPending = "Pending"
	// JobStatusCompleted is the status for a completed job.
	JobStatusCompleted = "Completed"
	// JobStatusFailed is the status for a failed job.
	JobStatusFailed = "Failed"
	// JobRequestDelay is the time wait to send a new request to check the job status.
	JobRequestDelay = 1 * time.Second
	// JobRequestTimeout is the time to wait before timeout.
	JobRequestTimeout = 60 * time.Minute
	// JobRequestMinTimeout is the minimum time to wait before refreshes.
	JobRequestMinTimeout = 3 * time.Second
	// JobRequestNotFoundChecks is the number of times to allow not found.
	JobRequestNotFoundChecks = 60
)

// WaitForJob is a function to check the Job status in a refresh function.
func WaitForJob(client *api.Client, aJob models.Job) error {
	var (
		pending   = JobStatusPending
		target    = JobStatusCompleted
		ctx       = context.Background()
		refreshFn = func() (result any, state string, err error) {
			client := job.New(client)

			j, err := client.Get(ctx, job.GetRequest{
				ID:   aJob.ID,
				Type: aJob.Type,
			})
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

	_, err := (&retry.StateChangeConf{
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
