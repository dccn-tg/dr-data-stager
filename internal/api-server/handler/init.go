package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Donders-Institute/dr-data-stager/internal/api-server/config"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/models"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/restapi/operations"
	"github.com/Donders-Institute/dr-data-stager/pkg/tasks"
	"github.com/go-openapi/runtime/middleware"
	"github.com/hibiken/asynq"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

var (
	// PathProject is the top-leve directory in which directories of active projects are located.
	PathProject string = "/project"

	// StagerJobQueueName is the queue name for stager jobs.
	StagerJobQueueName string = "stager"

	// MaxRetries (i.e. total attempts - 1)
	MaxRetries int = 2

	// RetryIntervalSeconds
	RetryIntervalSeconds int = 1
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
func GetJob(ctx context.Context, inspector *asynq.Inspector) func(params operations.GetJobIDParams, principal *models.Principal) middleware.Responder {
	return func(params operations.GetJobIDParams, principal *models.Principal) middleware.Responder {
		id := params.ID

		// retrieve task from the queue
		taskInfo, err := findTask(inspector, id)
		if err != nil {
			log.Errorf("[%s] cannot get task: %s", id, err)
			return operations.NewGetJobIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobQueueError,
				},
			)
		}

		res, err := composeResponseBodyJobInfo(taskInfo)
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

// DeleteJob cancels the job.
func DeleteJob(ctx context.Context, inspector *asynq.Inspector) func(params operations.DeleteJobIDParams, principal *models.Principal) middleware.Responder {

	return func(params operations.DeleteJobIDParams, principal *models.Principal) middleware.Responder {

		id := params.ID

		// retrieve task from the queue
		taskInfo, err := findTask(inspector, id)
		if err != nil {
			log.Errorf("[%s] error retrieve task info: %s", err)
			return operations.NewDeleteJobIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobDataError,
				},
			)
		}

		// unmarshal task payload
		var payload tasks.StagerPayload
		err = json.Unmarshal(taskInfo.Payload, &payload)
		if err != nil {
			log.Errorf("[%s] error unmarshal task payload: %s", err)
			return operations.NewDeleteJobIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobDataError,
				},
			)
		}

		// check stager user
		if payload.StagerUser != *(*string)(principal) {
			return operations.NewDeleteJobIDNotFound().WithPayload(
				fmt.Sprintf("[%s] task doesn't exist or not owned by the authenticated user: %s", id, *principal),
			)
		}

		// try to cancel the task before deleting it
		if err := inspector.CancelProcessing(id); err != nil {
			log.Errorf("[%s] cannot cancel task: %s", id, err)
		}

		if err := inspector.DeleteTask(taskInfo.Queue, id); err != nil {
			log.Errorf("[%s] cannot delete task: %s", id, err)
			return operations.NewDeleteJobIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobQueueError,
				},
			)
		}

		log.Infof("[%s] task deleted", id)

		jinfo, err := composeResponseBodyJobInfo(taskInfo)
		if err != nil {
			log.Errorf("%s", err)
			return operations.NewDeleteJobIDInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobDataError,
				},
			)
		}

		return operations.NewDeleteJobIDOK().WithPayload(jinfo)
	}
}

// NewJob registers the incoming transfer request as a new stager job in the queue.
func NewJob(ctx context.Context, client *asynq.Client) func(params operations.PostJobParams, principal *models.Principal) middleware.Responder {
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

		// set default job timeout (24 hours)
		timeout := params.Data.Timeout
		if timeout <= 0 {
			timeout = 86400
		}

		// set default job timeout no progress (1 hour)
		timeoutNp := params.Data.TimeoutNoprogress
		if timeoutNp <= 0 {
			timeoutNp = 3600
		}

		t, err := tasks.NewStagerTask(
			*params.Data.Title,
			*params.Data.DrUser,
			*params.Data.DstURL,
			*params.Data.SrcURL,
			*params.Data.StagerUser,
			timeout,
			timeoutNp,
		)

		if err != nil {
			return operations.NewPostJobInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobCreateError,
				},
			)
		}

		taskInfo, err := client.EnqueueContext(
			ctx,
			t,
			asynq.Retention(2*24*time.Hour),
			asynq.MaxRetry(3),
		)

		if err != nil {
			log.Errorf("cannot enqueue task: %s", err)
			return operations.NewPostJobInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     JobCreateError,
				},
			)
		}

		log.Infof("[%s] task submitted", taskInfo.ID)

		res, err := composeResponseBodyJobInfo(taskInfo)
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

// findTask looks into different queues to retrieve the TaskInfo, or return an error if not found.
func findTask(inspector *asynq.Inspector, id string) (*asynq.TaskInfo, error) {

	for q := range tasks.StagerQueues {
		if t, err := inspector.GetTaskInfo(q, id); err == nil {
			return t, err
		}
	}

	return nil, asynq.ErrTaskNotFound
}

func composeResponseBodyJobInfo(task *asynq.TaskInfo) (*models.ResponseBodyJobInfo, error) {

	var j tasks.StagerPayload
	if err := json.Unmarshal(task.Payload, &j); err != nil {
		return nil, err
	}

	// job status
	jStatus := task.State.String()

	// job result
	var jResult tasks.StagerTaskResult
	if len(task.Result) != 0 {
		if err := json.Unmarshal(task.Result, &jResult); err != nil {
			log.Warnf("cannot unmarshal payload result: %s", err)
		}
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
				Total:     &jResult.Progress.Total,
				Processed: &jResult.Progress.Processed,
			},
			Error: &task.LastErr,
		},
	}, nil
}
