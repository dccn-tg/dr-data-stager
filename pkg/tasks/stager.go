package tasks

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hibiken/asynq"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

// A list of task types.
const (
	TypeStager = "stager"
)

// Queues for different task types, with their associated task priority
var (
	StagerQueues = map[string]int{
		"critical": 6,
		"default":  3,
		"low":      1,
	}
)

// StagerPayload defines the data structure of the stager file transfer payload.
type StagerPayload struct {

	// creation time of the payload
	CreatedAt int64 `json:"createdAt,omitempty"`

	// short description about the job
	Title string `json:"title"`

	// username of the DR data-access account
	DrUser string `json:"drUser"`

	// password of the DR data-access account
	DrPass string `json:"drPass,omitempty"`

	// path or DR namespace (prefixed with irods:) of the destination endpoint
	DstURL string `json:"dstURL"`

	// path or DR namespace (prefixed with irods:) of the source endpoint
	SrcURL string `json:"srcURL"`

	// username of stager's local account
	StagerUser string `json:"stagerUser"`

	// email of stager's local account
	StagerUserEmail string `json:"stagerUserEmail,omitempty"`

	// allowed duration in seconds for entire transfer job (0 for no timeout)
	Timeout int64 `json:"timeout,omitempty"`

	// allowed duration in seconds for no further transfer progress (0 for no timeout)
	TimeoutNoprogress int64 `json:"timeout_noprogress,omitempty"`
}

// NewStagerTask wraps payload data into a `asynq.Task` ready for enqueuing.
func NewStagerTask(Title, DrUser, DrPass, DstURL, SrcURL, StagerUser, StagerUserEmail string, Timeout, TimeoutNoprogress int64) (*asynq.Task, error) {
	payload, err := json.Marshal(StagerPayload{
		CreatedAt:         time.Now().Unix(),
		Title:             Title,
		DrUser:            DrUser,
		DrPass:            DrPass,
		DstURL:            DstURL,
		SrcURL:            SrcURL,
		StagerUser:        StagerUser,
		StagerUserEmail:   StagerUserEmail,
		Timeout:           Timeout,
		TimeoutNoprogress: TimeoutNoprogress,
	})
	if err != nil {
		return nil, err
	}
	// task options are default settings which can be overridden at enqueue time.
	return asynq.NewTask(
		TypeStager,
		payload,
		asynq.MaxRetry(2),
		asynq.Timeout(time.Duration(Timeout)*time.Second),
	), nil
}

// Stager implements asynq.Handler interface.
type Stager struct {
	config config.Configuration
}

func (stager *Stager) ProcessTask(ctx context.Context, t *asynq.Task) error {

	updateRslt := func(rslt *StagerTaskResult) {
		if d, err := json.Marshal(rslt); err == nil {
			t.ResultWriter().Write(d)
		}
	}

	tid, ok := asynq.GetTaskID(ctx)
	if !ok {
		return fmt.Errorf("cannot get task id from the task context")
	}

	var p StagerPayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("cannot unmarshal stager payload: %v: %w", err, asynq.SkipRetry)
	}
	log.Debugf("[%s] payload: %+v", tid, p)

	timer := time.NewTimer(time.Duration(p.TimeoutNoprogress) * time.Second)

	cout, cerr, cmd, err := runSyncAs(ctx, p)
	if err != nil {
		log.Errorf("[%s] %s", tid, err)
		return err
	}

	// updata task progress
	done := make(chan error, 1)
	go func() {
		rslt := new(StagerTaskResult)

		percent := 0

		for progress := range cout {
			// stop timer
			if !timer.Stop() {
				<-timer.C
			}

			// increase the counter by 1, and update the queue data
			rslt.Progress.Total = progress.Total
			rslt.Progress.Processed = progress.Success + progress.Failure
			rslt.Progress.Failed = progress.Failure

			npercent := int(100 * (progress.Success + progress.Failure) / progress.Total)

			log.Debugf("[%s] %d/%d (%d%%) processed", tid, rslt.Progress.Processed, rslt.Progress.Total, npercent)

			if npercent > percent {
				updateRslt(rslt)
				percent = npercent
			}

			// reset timer
			timer.Reset(time.Duration(p.TimeoutNoprogress) * time.Second)
		}

		// wait for command to stop
		done <- cmd.Wait()
	}()

	lastErr := ""
	// control loop for:
	// - catch last message on stderr
	// - cmd process has been finished
	// - kill cmd process on timeout and process terminiation by context
	for {
		select {

		case errStr, more := <-cerr:
			if more {
				lastErr = errStr
			}
		case e := <-done:

			if e != nil {
				err := fmt.Errorf("s-isync failed: %s - %s", e, lastErr)
				log.Errorf("[%s] %s", tid, err)
				return err
			}
			return nil

		case <-timer.C:
			// receive times up signal for `timeoutNoprogress`
			err := fmt.Errorf("no progress more than %d seconds", p.TimeoutNoprogress)
			log.Errorf("[%s] %s", tid, err)

			// send kill to the cmd's process
			if err := cmd.Process.Kill(); err != nil {
				log.Errorf("[%s] fail to terminate s-isync: %s", tid, err)
			}

			return err

		case <-ctx.Done():
			// receive abort signal from parent context
			err := fmt.Errorf("aborted by context")
			log.Debugf("[%s] %s", tid, err)

			// send kill to the cmd's process
			if err := cmd.Process.Kill(); err != nil {
				log.Errorf("[%s] fail to terminate s-isync: %s", tid, err)
			}

			return ctx.Err()
		}
	}
}

// progress stores total number of processed files.
type progress struct {
	Total   int64
	Success int64
	Failure int64
}

// runSyncAs runs `s-isync` as the `stagerUser` in a go routine.
func runSyncAs(ctx context.Context, payload StagerPayload) (chan progress, chan string, *exec.Cmd, error) {

	tid, ok := asynq.GetTaskID(ctx)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid context: missing asynq task id")
	}

	u, err := user.Lookup(payload.StagerUser)
	if err != nil {
		return nil, nil, nil, err
	}

	uid, _ := strconv.ParseInt(u.Uid, 10, 32)
	gid, _ := strconv.ParseInt(u.Gid, 10, 32)

	cmd := exec.Command(
		"/opt/stager/s-isync",
		"-v",
		"-c", "/etc/stager/worker.yml",
		"-l", fmt.Sprintf("/tmp/s-isync-%s.log", tid),
		"--task", tid,
		"--druser", payload.DrUser,
		"--drpass", payload.DrPass,
		payload.SrcURL,
		payload.DstURL,
	)

	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, cmd, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, cmd, err
	}

	log.Debugf("[%s] run s-isync as %s (%d:%d)\n", tid, payload.StagerUser, uid, gid)
	if err := cmd.Start(); err != nil {
		return nil, nil, cmd, err
	}

	cout := make(chan progress, 1)
	// go routine to read and process the stdout of s-irsync (progress information)
	go func() {
		defer close(cout)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			data := strings.Split(line, ",")

			if len(data) != 3 {
				log.Errorf("unexpected progress output: %s", string(line))
				continue
			}

			t, _ := strconv.ParseInt(data[0], 10, 64)
			s, _ := strconv.ParseInt(data[1], 10, 64)
			f, _ := strconv.ParseInt(data[2], 10, 64)

			cout <- progress{
				Total:   t,
				Success: s,
				Failure: f,
			}
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

	return cout, cerr, cmd, nil
}

func NewStager(config config.Configuration) *Stager {
	return &Stager{
		config: config,
	}
}

// StagerTaskResult
type StagerTaskResult struct {
	Progress struct {
		Total     int64 `json:"total"`
		Processed int64 `json:"processes"`
		Failed    int64 `json:"failed"`
	} `json:"progress"`
}
