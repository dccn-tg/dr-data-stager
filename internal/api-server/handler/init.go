package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Donders-Institute/dr-data-stager/internal/api-server/config"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/models"
	"github.com/Donders-Institute/dr-data-stager/pkg/swagger/server/restapi/operations"
	"github.com/Donders-Institute/dr-data-stager/pkg/tasks"
	"github.com/go-openapi/runtime/middleware"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

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

	// FileSystemError indicates failure on filesystem operation
	FileSystemError int64 = 104
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

// GetJobs retrieves all jobs owned by the `principal`.
func GetJobs(ctx context.Context, inspector *asynq.Inspector) func(params operations.GetJobsParams, principal *models.Principal) middleware.Responder {
	return func(params operations.GetJobsParams, principal *models.Principal) middleware.Responder {
		return operations.NewGetJobsOK().WithPayload(
			&models.ResponseBodyJobs{
				Jobs: getTasksOwnedBy(inspector, string(*principal)),
			},
		)
	}
}

// GetJob retrieves job information.
func GetJob(ctx context.Context, inspector *asynq.Inspector) func(params operations.GetJobIDParams, principal *models.Principal) middleware.Responder {
	return func(params operations.GetJobIDParams, principal *models.Principal) middleware.Responder {
		id := params.ID

		// retrieve task from the queue
		taskInfo, err := findTask(inspector, id)
		if errors.Is(err, asynq.ErrTaskNotFound) {
			return operations.NewGetJobIDNotFound().WithPayload(
				fmt.Sprintf("job %s doesn't exist", id),
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
		// Return 404 not found if it doesn't match
		if *res.Data.StagerUser != *(*string)(principal) {
			return operations.NewGetJobIDNotFound().WithPayload(
				fmt.Sprintf("job %s not owned by the authenticated user: %s", id, *principal),
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
		if errors.Is(err, asynq.ErrTaskNotFound) {
			return operations.NewDeleteJobIDNotFound().WithPayload(
				fmt.Sprintf("job %s doesn't exist", id),
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

func ListDir(ctx context.Context) func(params operations.GetDirParams, principal *models.Principal) middleware.Responder {
	return func(params operations.GetDirParams, principal *models.Principal) middleware.Responder {

		cout, cerr, cmd, err := runCmdAs(
			string(*principal),
			"ls",
			"-l",
			"-G",
			"-g",
			`--time-style=+%s`,
			*params.Dir.Path,
		)

		if err != nil {
			log.Errorf("%s", err)
			return operations.NewGetDirInternalServerError().WithPayload(
				&models.ResponseBody500{
					ErrorMessage: err.Error(),
					ExitCode:     FileSystemError,
				},
			)
		}

		// catching stderr output
		done := make(chan error, 1)
		go func() {
			defer close(done)
			lastErr := ""
			for lerr := range cerr {
				lastErr = lerr
			}
			done <- fmt.Errorf("%s", lastErr)
		}()

		// an empty slice for entries
		entries := []*models.DirEntry{}

		// process output line by line
		for lout := range cout {

			fields := strings.Fields(lout)

			if len(fields) < 5 {
				continue
			}

			// file type
			ftype := models.DirEntryTypeUnknown
			switch string(fields[0])[0] {
			case 'd':
				ftype = models.DirEntryTypeDir
			case '-':
				ftype = models.DirEntryTypeRegular
			case 'l':
				ftype = models.DirEntryTypeSymlink
			}

			// file size
			fsize := int64(0)
			s, err := strconv.Atoi(fields[2])
			if err != nil {
				log.Errorf("unknown file size: %s\n", fields[2])
			} else {
				fsize = int64(s)
			}

			_, fname, _ := strings.Cut(lout, fields[3])
			fname = strings.TrimSpace(fname)

			// file name
			entries = append(entries, &models.DirEntry{
				Name: &fname,
				Size: &fsize,
				Type: &ftype,
			})
		}

		// wait until the command is finished
		lastErr := <-done
		err = cmd.Wait()

		if err != nil {
			// only log commandline execution error on server.
			// always return 200 OK to the API client.
			log.Errorf("%s: %s\n", err.Error(), lastErr)
		}

		return operations.NewGetDirOK().WithPayload(
			&models.ResponseDirEntries{
				Entries: entries,
			},
		)
	}
}

// NewJobs registers the incoming transfer request as multiple stager jobs in the queue.
func NewJobs(ctx context.Context, client *asynq.Client, rdb *redis.Client) func(params operations.PostJobsParams, principal *models.Principal) middleware.Responder {
	return func(params operations.PostJobsParams, principal *models.Principal) middleware.Responder {

		submitted := []*models.JobInfo{}

		for _, jdata := range params.Data.Jobs {
			if *jdata.StagerUser != *(*string)(principal) {
				log.Errorf("job not owned by the authenticated user: %s\n", *principal)
				continue
			}

			taskInfo, err := enqueueStagerTask(ctx, client, jdata, rdb)
			if err != nil {
				log.Errorf("cannot enqueue task: %s", err)
				continue
			}

			log.Infof("[%s] task submitted", taskInfo.ID)

			jinfo, err := composeResponseBodyJobInfo(taskInfo)
			if err != nil {
				log.Errorf("[%s] fail to wrap up job info: %s\n", err)
				continue
			}
			submitted = append(submitted, jinfo)
		}

		if len(submitted) != len(params.Data.Jobs) {
			if len(submitted) == 0 {
				return operations.NewPostJobsInternalServerError().WithPayload(
					&models.ResponseBody500{
						ErrorMessage: "fail to create stager jobs",
						ExitCode:     JobCreateError,
					},
				)
			}

			return operations.NewPostJobsMultiStatus().WithPayload(
				&models.ResponseBodyJobs{
					Jobs: submitted,
				},
			)
		}

		return operations.NewPostJobsOK().WithPayload(
			&models.ResponseBodyJobs{
				Jobs: submitted,
			},
		)
	}
}

// NewJob registers the incoming transfer request as a new stager job in the queue.
func NewJob(ctx context.Context, client *asynq.Client, rdb *redis.Client) func(params operations.PostJobParams, principal *models.Principal) middleware.Responder {
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

		taskInfo, err := enqueueStagerTask(ctx, client, params.Data, rdb)
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

// runCmdAs spawns a new process and run the `cmd` with `args` as the `username`.
func runCmdAs(username string, cmd string, args ...string) (chan string, chan string, *exec.Cmd, error) {

	c := exec.Command(cmd, args...)

	if username != "" {
		u, err := user.Lookup(username)
		if err != nil {
			return nil, nil, nil, err
		}

		uid, _ := strconv.ParseInt(u.Uid, 10, 32)
		gid, _ := strconv.ParseInt(u.Gid, 10, 32)

		c.SysProcAttr = &syscall.SysProcAttr{}
		c.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}

	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, nil, c, err
	}

	stderr, err := c.StderrPipe()
	if err != nil {
		return nil, nil, c, err
	}

	//log.Debugf("run command %s as %s (%d:%d)\n", cmd, username, uid, gid)
	if err := c.Start(); err != nil {
		return nil, nil, c, err
	}

	cout := make(chan string, 1)
	// go routine to read the cmd's stdout line-by-line
	go func() {
		defer close(cout)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			cout <- strings.TrimSpace(scanner.Text())
		}
	}()

	cerr := make(chan string, 1)
	// go routine to read and process the stderr
	go func() {
		defer close(cerr)
		// read stderr
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			cerr <- strings.TrimSpace(scanner.Text())
		}
	}()

	return cout, cerr, c, nil
}

// findTask looks into different queues to retrieve the TaskInfo, or return `asynq.ErrTaskNotFound` if not found.
func findTask(inspector *asynq.Inspector, id string) (*asynq.TaskInfo, error) {

	for q := range tasks.StagerQueues {
		if t, err := inspector.GetTaskInfo(q, id); err == nil {
			return t, err
		}
	}

	return nil, asynq.ErrTaskNotFound
}

// getTasksInStateOwnedBy retrieves all tasks "belong" to the `username` by list all
// tasks with `stagerUser` prefix on the task ID.
func getTasksInStateOwnedBy(inspector *asynq.Inspector, state asynq.TaskState, username string) []*models.JobInfo {

	ctasks := make(chan *asynq.TaskInfo)

	// generic fetch function loops over queues and pages
	fetch := func(fetcher func(q string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error)) {

		defer close(ctasks)

		tidPrefix := fmt.Sprintf("%s.", username)
		for queue := range tasks.StagerQueues {
			pn := 0
			for {
				pn++
				qtasks, err := fetcher(queue, asynq.PageSize(3), asynq.Page(pn))
				if err != nil || len(qtasks) == 0 {
					break
				}
				for _, t := range qtasks {
					if !strings.HasPrefix(t.ID, tidPrefix) {
						continue
					}
					ctasks <- t
				}
			}
		}
	}

	switch state {
	case asynq.TaskStateScheduled:
		go fetch(inspector.ListScheduledTasks)
	case asynq.TaskStatePending:
		go fetch(inspector.ListPendingTasks)
	case asynq.TaskStateActive:
		go fetch(inspector.ListActiveTasks)
	case asynq.TaskStateRetry:
		go fetch(inspector.ListRetryTasks)
	case asynq.TaskStateCompleted:
		go fetch(inspector.ListCompletedTasks)
	case asynq.TaskStateArchived:
		go fetch(inspector.ListArchivedTasks)
	}

	tasks := []*models.JobInfo{}

	for t := range ctasks {
		info, err := composeResponseBodyJobInfo(t)
		if err != nil {
			log.Errorf("%s\n", err)
			continue
		}
		tasks = append(tasks, info)
	}

	return tasks
}

// getTasksOwnedBy retrieves all tasks "belong" to the `username` by list all
// tasks with `stagerUser` prefix on the task ID.
func getTasksOwnedBy(inspector *asynq.Inspector, username string) []*models.JobInfo {

	tasks := []*models.JobInfo{}

	for _, s := range []asynq.TaskState{
		asynq.TaskStatePending,
		asynq.TaskStateActive,
		asynq.TaskStateRetry,
		asynq.TaskStateCompleted,
		asynq.TaskStateArchived,
	} {
		tasks = append(tasks, getTasksInStateOwnedBy(inspector, s, username)...)
	}

	return tasks
}

// enqueueStagerTask creates a new stager task in the asynq queue.
func enqueueStagerTask(ctx context.Context, client *asynq.Client, job *models.JobData, rdb *redis.Client) (*asynq.TaskInfo, error) {
	// set default job timeout (24 hours)
	timeout := job.Timeout
	if timeout <= 0 {
		timeout = 86400
	}

	// set default job timeout no progress (1 hour)
	timeoutNp := job.TimeoutNoprogress
	if timeoutNp <= 0 {
		timeoutNp = 3600
	}

	t, err := tasks.NewStagerTask(
		*job.Title,
		*job.DrUser,
		job.DrPass,
		*job.DstURL,
		*job.SrcURL,
		*job.StagerUser,
		job.StagerUserEmail.String(),
		timeout,
		timeoutNp,
	)

	if err != nil {
		return nil, err
	}

	tid, err := rdb.Incr(ctx, "stager:lastTaskId").Result()
	if err != nil {
		return nil, fmt.Errorf("cannot get next task id: %s", err)
	}

	return client.EnqueueContext(
		ctx,
		t,
		asynq.TaskID(fmt.Sprintf("%s.%d", *job.StagerUser, tid)),
		asynq.Retention(2*24*time.Hour),
		asynq.MaxRetry(4),
		asynq.Timeout(time.Duration(timeout)*time.Second), // this set the hard timeout
	)
}

// composeResponseBodyJobInfo wraps the data structure of `asynq.TaskInfo` into `models.JobInfo`.
func composeResponseBodyJobInfo(task *asynq.TaskInfo) (*models.JobInfo, error) {

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

	createdAt := j.CreatedAt
	nextProcessedAt := task.NextProcessAt.Unix()
	lastFailedAt := task.LastFailedAt.Unix()
	completedAt := task.CompletedAt.Unix()

	attempts := int64(task.Retried + 1)

	return &models.JobInfo{
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
		Timestamps: &models.JobTimestamps{
			CreatedAt:     &createdAt,
			NextProcessAt: &nextProcessedAt,
			LastFailedAt:  &lastFailedAt,
			CompletedAt:   &completedAt,
		},
		Status: &models.JobStatus{
			Status: &jStatus,
			Progress: &models.JobProgress{
				Total:     &jResult.Progress.Total,
				Processed: &jResult.Progress.Processed,
				Failed:    &jResult.Progress.Failed,
			},
			Error:    &task.LastErr,
			Attempts: &attempts,
		},
	}, nil
}
