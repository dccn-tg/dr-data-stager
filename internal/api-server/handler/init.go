package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Donders-Institute/dr-data-stager/internal/api-server/config"
	"github.com/Donders-Institute/dr-data-stager/internal/job"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/models"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/thoas/bokchoy"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

var (
	// PathProject is the top-leve directory in which directories of active projects are located.
	PathProject string = "/project"

	// StagerJobQueueName is the queue name for stager jobs.
	StagerJobQueueName string = "stager"

	// MaxRetries (i.e. total attempts - 1)
	MaxRetries int = 4
)

// Error code definitions.
var (
	// NotImplementedError indicates the implementation of the handler is not implemented yet.
	NotImplementedError int64 = 100

	// JobQueueError indicates failure adding/retrieving job from the stager job queue.
	JobQueueError int64 = 101

	// JobDataError indicates invalid job data format causing (de-)serialization failure.
	JobDataError int64 = 102

	// JobCreateError indicates failure creating/registering new job to the stager job queue.
	JobCreateError int64 = 103
)

// Common payload for the ResponseBody500.
var responseNotImplemented = models.ResponseBody500{
	ErrorMessage: "not implemented",
	ExitCode:     NotImplementedError,
}

// GetPing returns dummy string for health check, including the authentication.
func GetPing(cfg config.Configuration) func(params operations.GetPingParams) middleware.Responder {
	return func(params operations.GetPingParams) middleware.Responder {
		return operations.NewGetPingOK().WithPayload("pong")
	}
}

// GetJob retrieves job information.
func GetJob(ctx context.Context, bok *bokchoy.Bokchoy) func(params operations.GetJobIDParams, principal *models.Principal) middleware.Responder {
	return func(params operations.GetJobIDParams, principal *models.Principal) middleware.Responder {
		id := params.ID

		// retrieve task from the queue
		tbok, err := bok.Queue(StagerJobQueueName).Get(ctx, id)

		if err != nil {
			log.Errorf("%s", err)
			return operations.NewGetJobIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobQueueError,
				},
			)
		}

		res, err := composeResponseBodyJobInfo(tbok)
		if err != nil {
			log.Errorf("%s", err)
			return operations.NewGetJobIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobDataError,
				},
			)
		}

		// check if job owner (stager user) is the same as principal.
		// Return 404 not found error if no match.
		if *res.Data.StagerUser != *(*string)(principal) {
			return operations.NewGetJobIDNotFound().WithPayload(
				fmt.Sprintf("job %s doesn't exist or not owned by the authenticated user: %s", id, *principal),
			)
		}

		return operations.NewGetJobIDOK().WithPayload(res)
	}
}

// NewJob registers the incoming transfer request as a new stager job in the queue.
func NewJob(ctx context.Context, bok *bokchoy.Bokchoy) func(params operations.PostJobParams, principal *models.Principal) middleware.Responder {
	return func(params operations.PostJobParams, principal *models.Principal) middleware.Responder {

		// check if job owner matches the authenticated principal
		if *params.Data.StagerUser != *(*string)(principal) {
			return operations.NewPostJobInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: fmt.Sprintf("job not owned by the authenticated user: %s", *principal),
					ExitCode:     JobCreateError,
				},
			)
		}

		j := job.Stager{
			Title:             *params.Data.Title,
			DrUser:            *params.Data.DrUser,
			StagerUser:        *params.Data.StagerUser,
			SrcURL:            *params.Data.SrcURL,
			DstURL:            *params.Data.DstURL,
			Timeout:           params.Data.Timeout,
			TimeoutNoprogress: params.Data.TimeoutNoprogress,
		}

		// set default job timeout (24 hours)
		if j.Timeout <= 0 {
			j.Timeout = 86400
		}

		// set default job timeout no progress (1 hour)
		if j.TimeoutNoprogress <= 0 {
			j.TimeoutNoprogress = 3600
		}

		// publish data to queue
		// 1. a big default job timeout of 7 days is set as the actual timeout
		//    logic is handled by the worker using the `Timeout` attribute of
		//    the job data.
		// 2. after the job is finished, the retention time of the result is 3 days.
		// 3. the default job max retries is defined by `MaxRetries`
		tbok, err := bok.Queue(StagerJobQueueName).Publish(ctx, &j,
			bokchoy.WithMaxRetries(MaxRetries),
			bokchoy.WithTimeout(7*24*time.Hour),
			bokchoy.WithTTL(3*24*time.Hour))

		if err != nil {
			return operations.NewPostJobInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobCreateError,
				},
			)
		}

		res, err := composeResponseBodyJobInfo(tbok)
		if err != nil {
			log.Errorf("%s", err)
			return operations.NewPostJobInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobDataError,
				},
			)
		}

		return operations.NewPostJobOK().WithPayload(res)
	}
}

func composeResponseBodyJobInfo(task *bokchoy.Task) (*models.ResponseBodyJobInfo, error) {
	// convert Payload to models.JobData structure
	res, err := json.Marshal(task.Payload)
	if err != nil {
		return nil, err
	}

	var j job.Stager
	err = json.Unmarshal(res, &j)
	if err != nil {
		return nil, err
	}

	// job status
	jStatus := task.StatusDisplay()

	// job progress
	jProgress, err := resultToProgress(task.Result)
	if err != nil {
		return nil, err
	}

	// job error
	jErr := ""
	if task.Error != nil {
		jErr = fmt.Sprintf("%s", task.Error)
	}

	// job identifier
	jid := models.JobID(task.ID)

	return &models.ResponseBodyJobInfo{
		ID: &jid,
		Data: &models.JobData{
			Title:             &j.Title,
			DrUser:            &j.DrUser,
			StagerUser:        &j.StagerUser,
			SrcURL:            &j.SrcURL,
			DstURL:            &j.DstURL,
			Timeout:           j.Timeout,
			TimeoutNoprogress: j.TimeoutNoprogress,
		},
		Status: &models.JobStatus{
			Status: &jStatus,
			Progress: &models.JobProgress{
				Total:     &jProgress.Total,
				Processed: &jProgress.Processed,
			},
			Error: &jErr,
		},
	}, nil
}

// resultToProgress convert any task result interface into job.Progress
// structure.
func resultToProgress(rslt interface{}) (*job.Progress, error) {
	if rslt == nil {
		return &job.Progress{
			Total:     0,
			Processed: 0,
		}, nil
	}

	res, err := json.Marshal(rslt)
	if err != nil {
		return nil, err
	}

	var p job.Progress
	err = json.Unmarshal(res, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
